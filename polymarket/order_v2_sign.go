package polymarket

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethmath "github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/polymarket/go-order-utils/pkg/model"
	"github.com/polymarket/go-order-utils/pkg/utils"
)

const (
	zeroAddressHex  = "0x0000000000000000000000000000000000000000"
	bytes32ZeroHex  = "0x0000000000000000000000000000000000000000000000000000000000000000"
	v2DomainName    = "Polymarket CTF Exchange"
	v2DomainVersion = "2"
)

type SignedOrderV2 struct {
	Salt          *big.Int
	Maker         common.Address
	Signer        common.Address
	TokenID       *big.Int
	MakerAmount   *big.Int
	TakerAmount   *big.Int
	Side          Side
	SignatureType model.SignatureType
	Timestamp     *big.Int
	Metadata      common.Hash
	Builder       common.Hash
	Expiration    *big.Int
	Signature     []byte
}

func parseBase10Int(name, value string) (*big.Int, error) {
	v, ok := new(big.Int).SetString(value, 10)
	if !ok {
		return nil, fmt.Errorf("invalid %s: %q", name, value)
	}
	return v, nil
}

func normalizeBytes32Hex(v string) common.Hash {
	if v == "" {
		return common.HexToHash(bytes32ZeroHex)
	}
	return common.HexToHash(v)
}

func (c *Client) resolveV2VerifyingContract(tokenID string) (common.Address, error) {
	negRisk := c.forceNegRiskExchange
	if !negRisk {
		var err error
		negRisk, err = c.GetNegRisk(tokenID)
		if err != nil {
			return common.Address{}, fmt.Errorf("resolve neg-risk market: %w", err)
		}
	}
	switch c.chainID.Int64() {
	case 137, 80002:
		if negRisk {
			return common.HexToAddress("0xe2222d279d744050d28e00520010520000310F59"), nil
		}
		return common.HexToAddress("0xE111180000d2663C0091e4f400237545B87B996B"), nil
	default:
		return common.Address{}, fmt.Errorf("unsupported chain id for v2 order signing: %s", c.chainID.String())
	}
}

func (c *Client) signOrderV2(order *SignedOrderV2, verifyingContract common.Address) error {
	types := apitypes.Types{
		"EIP712Domain": {
			{Name: "name", Type: "string"},
			{Name: "version", Type: "string"},
			{Name: "chainId", Type: "uint256"},
			{Name: "verifyingContract", Type: "address"},
		},
		"Order": {
			{Name: "salt", Type: "uint256"},
			{Name: "maker", Type: "address"},
			{Name: "signer", Type: "address"},
			{Name: "tokenId", Type: "uint256"},
			{Name: "makerAmount", Type: "uint256"},
			{Name: "takerAmount", Type: "uint256"},
			{Name: "side", Type: "uint8"},
			{Name: "signatureType", Type: "uint8"},
			{Name: "timestamp", Type: "uint256"},
			{Name: "metadata", Type: "bytes32"},
			{Name: "builder", Type: "bytes32"},
		},
	}

	side := uint8(1)
	if order.Side == SideBuy {
		side = 0
	}
	msg := apitypes.TypedDataMessage{
		"salt":          order.Salt,
		"maker":         order.Maker.Hex(),
		"signer":        order.Signer.Hex(),
		"tokenId":       order.TokenID,
		"makerAmount":   order.MakerAmount,
		"takerAmount":   order.TakerAmount,
		"side":          side,
		"signatureType": uint8(order.SignatureType),
		"timestamp":     order.Timestamp,
		"metadata":      order.Metadata.Hex(),
		"builder":       order.Builder.Hex(),
	}
	typedData := apitypes.TypedData{
		Types:       types,
		PrimaryType: "Order",
		Domain: apitypes.TypedDataDomain{
			Name:              v2DomainName,
			Version:           v2DomainVersion,
			ChainId:           (*ethmath.HexOrDecimal256)(c.chainID),
			VerifyingContract: verifyingContract.Hex(),
		},
		Message: msg,
	}

	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return fmt.Errorf("eip712 v2 domain: %w", err)
	}
	structHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return fmt.Errorf("eip712 v2 order: %w", err)
	}
	rawData := append([]byte{0x19, 0x01}, append(domainSeparator, structHash...)...)
	challengeHash := crypto.Keccak256Hash(rawData)
	sig, err := crypto.Sign(challengeHash.Bytes(), c.privateKey)
	if err != nil {
		return fmt.Errorf("sign v2 order: %w", err)
	}
	order.Signature = sig
	return nil
}

func (c *Client) newSignedOrderV2(
	tokenID, makerAmount, takerAmount string,
	side Side,
	signatureType model.SignatureType,
	metadataHex, builderHex string,
	expiration *int64,
) (*SignedOrderV2, error) {
	token, err := parseBase10Int("tokenId", tokenID)
	if err != nil {
		return nil, err
	}
	makerAmt, err := parseBase10Int("makerAmount", makerAmount)
	if err != nil {
		return nil, err
	}
	takerAmt, err := parseBase10Int("takerAmount", takerAmount)
	if err != nil {
		return nil, err
	}
	exp := big.NewInt(0)
	if expiration != nil {
		exp = big.NewInt(*expiration)
	}

	signer := c.fromAddress
	maker := signer
	if c.funderAddress != (common.Address{}) {
		maker = c.funderAddress
	}
	if side != SideBuy && side != SideSell {
		return nil, fmt.Errorf("invalid order side: %q", side)
	}
	order := &SignedOrderV2{
		Salt:          big.NewInt(utils.GenerateRandomSalt()),
		Maker:         maker,
		Signer:        signer,
		TokenID:       token,
		MakerAmount:   makerAmt,
		TakerAmount:   takerAmt,
		Side:          side,
		SignatureType: signatureType,
		Timestamp:     big.NewInt(time.Now().UnixMilli()),
		Metadata:      normalizeBytes32Hex(metadataHex),
		Builder:       normalizeBytes32Hex(builderHex),
		Expiration:    exp,
	}
	verifyingContract, err := c.resolveV2VerifyingContract(tokenID)
	if err != nil {
		return nil, err
	}
	if err := c.signOrderV2(order, verifyingContract); err != nil {
		return nil, err
	}
	return order, nil
}

func encodeOrderSignatureHex(sig []byte) string {
	return "0x" + hex.EncodeToString(sig)
}
