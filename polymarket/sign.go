package polymarket

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

const clobAuthMessage = "This message attests that I control the given wallet"

// signL2 builds the Polymarket CLOB L2 HMAC (timestamp + method + path + body), matching clob-client/src/signing/hmac.ts.
func signL2(method, path, body, timestamp, secretB64 string) (string, error) {
	secretBytes, err := base64.URLEncoding.DecodeString(secretB64)
	if err != nil {
		secretBytes, err = base64.StdEncoding.DecodeString(secretB64)
		if err != nil {
			secretBytes, err = base64.RawURLEncoding.DecodeString(secretB64)
			if err != nil {
				return "", fmt.Errorf("decode api secret: %w", err)
			}
		}
	}

	message := timestamp + method + path
	if body != "" {
		message += body
	}

	h := hmac.New(sha256.New, secretBytes)
	h.Write([]byte(message))
	sig := base64.StdEncoding.EncodeToString(h.Sum(nil))
	sig = replaceB64ToURLSafe(sig)
	return sig, nil
}

func replaceB64ToURLSafe(s string) string {
	b := []byte(s)
	for i, ch := range b {
		switch ch {
		case '+':
			b[i] = '-'
		case '/':
			b[i] = '_'
		}
	}
	return string(b)
}

// signClobAuthEIP712 signs L1 auth using the same typed data as clob-client/src/signing/eip712.ts (viem-compatible).
func (c *Client) signClobAuthEIP712(ts, nonce int64) (string, error) {
	tsStr := strconv.FormatInt(ts, 10)
	types := apitypes.Types{
		"EIP712Domain": {
			{Name: "name", Type: "string"},
			{Name: "version", Type: "string"},
			{Name: "chainId", Type: "uint256"},
		},
		"ClobAuth": {
			{Name: "address", Type: "address"},
			{Name: "timestamp", Type: "string"},
			{Name: "nonce", Type: "uint256"},
			{Name: "message", Type: "string"},
		},
	}
	domain := apitypes.TypedDataDomain{
		Name:    "ClobAuthDomain",
		Version: "1",
		ChainId: (*math.HexOrDecimal256)(c.chainID),
	}
	msg := apitypes.TypedDataMessage{
		"address":   c.fromAddress.Hex(),
		"timestamp": tsStr,
		"nonce":     new(big.Int).SetInt64(nonce),
		"message":   clobAuthMessage,
	}
	typedData := apitypes.TypedData{
		Types:       types,
		PrimaryType: "ClobAuth",
		Domain:      domain,
		Message:     msg,
	}
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return "", fmt.Errorf("eip712 domain: %w", err)
	}
	structHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return "", fmt.Errorf("eip712 clob auth: %w", err)
	}
	rawData := append([]byte{0x19, 0x01}, append(domainSeparator, structHash...)...)
	challengeHash := crypto.Keccak256Hash(rawData)

	sig, err := crypto.Sign(challengeHash.Bytes(), c.privateKey)
	if err != nil {
		return "", fmt.Errorf("sign clob auth: %w", err)
	}
	if sig[64] < 27 {
		sig[64] += 27
	}
	return "0x" + hex.EncodeToString(sig), nil
}

func (c *Client) buildL1AuthHeaders(nonce int64) (map[string]string, error) {
	ts, err := c.authTimestampSeconds()
	if err != nil {
		return nil, err
	}
	sig, err := c.signClobAuthEIP712(ts, nonce)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"POLY_ADDRESS":   c.fromAddress.Hex(),
		"POLY_SIGNATURE": sig,
		"POLY_TIMESTAMP": strconv.FormatInt(ts, 10),
		"POLY_NONCE":     strconv.FormatInt(nonce, 10),
	}, nil
}

func (c *Client) buildL2AuthHeaders(method, requestPath, body string) (map[string]string, error) {
	if c.apiKeyCredentials == nil {
		return nil, fmt.Errorf("no API key credentials set")
	}
	ts, err := c.authTimestampSeconds()
	if err != nil {
		return nil, err
	}
	tsStr := strconv.FormatInt(ts, 10)
	sig, err := signL2(method, requestPath, body, tsStr, c.apiKeyCredentials.Secret)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"POLY_ADDRESS":    c.fromAddress.Hex(),
		"POLY_SIGNATURE":  sig,
		"POLY_TIMESTAMP":  tsStr,
		"POLY_API_KEY":    c.apiKeyCredentials.ApiKey,
		"POLY_PASSPHRASE": c.apiKeyCredentials.Passphrase,
	}, nil
}
