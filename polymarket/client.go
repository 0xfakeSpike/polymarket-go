package polymarket

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/polymarket/go-order-utils/pkg/model"
)

const (
	CLOBHost         = "https://clob.polymarket.com"
	DefaultTimeout   = 30 * time.Second
	defaultTickCache = 5 * time.Minute
)

var defaultPolygonRPCURLs = []string{"https://api.zan.top/polygon-mainnet", "https://polygon-rpc.com"}

// Client provides Polymarket CLOB APIs.
type Client struct {
	httpClient *http.Client

	clobHost string

	ethClient       *ethclient.Client
	chainID         *big.Int
	chainIDOverride *big.Int
	polygonRPCURLs  []string

	privateKey  *ecdsa.PrivateKey
	fromAddress common.Address

	apiKeyCredentials *APIKeyCredentials

	geoBlockToken    string
	useServerTime    bool
	retryPostOnError bool
	throwOnError     bool
	signatureType    model.SignatureType
	// funderAddress is the order maker when using proxy/safe flows; zero means maker == signer (EOA).
	funderAddress common.Address
	builderSigner BuilderSigner

	tickSizes   map[string]string
	tickSizeAt  map[string]time.Time
	tickSizeTTL time.Duration

	negRiskCache map[string]bool
	feeRateCache map[string]int
	feeInfoCache map[string]FeeInfo

	tokenConditionMap map[string]string
	builderFeeRates   map[string]BuilderFeeRate
	cachedVersion     *int

	// skipL2APIKeyBootstrap skips CreateOrDeriveAPIKey inside NewClient (public-only or manual creds).
	skipL2APIKeyBootstrap bool

	// forceNegRiskExchange signs all orders against the neg-risk CTF exchange (matches a fixed
	// verifyingContract workflow). When false, the contract is chosen from GET /neg-risk per token.
	forceNegRiskExchange bool
}

// NewClient builds a client from a hex-encoded secp256k1 private key (with or without 0x).
// By default it also obtains CLOB L2 API credentials via CreateOrDeriveAPIKey (network call).
// Skip that with [WithAPIKeyCredentials] (non-nil) or [WithSkipL2APIKeyBootstrap].
func NewClient(privateKeyHex string, opts ...ClientOption) (*Client, error) {
	if len(privateKeyHex) >= 2 && privateKeyHex[0:2] == "0x" {
		privateKeyHex = privateKeyHex[2:]
	}
	pk, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("private key: %w", err)
	}
	from := crypto.PubkeyToAddress(pk.PublicKey)

	c := &Client{
		privateKey:        pk,
		fromAddress:       from,
		clobHost:          CLOBHost,
		polygonRPCURLs:    defaultPolygonRPCURLs,
		signatureType:     model.POLY_GNOSIS_SAFE,
		tickSizes:         make(map[string]string),
		tickSizeAt:        make(map[string]time.Time),
		tickSizeTTL:       defaultTickCache,
		negRiskCache:      make(map[string]bool),
		feeRateCache:      make(map[string]int),
		feeInfoCache:      make(map[string]FeeInfo),
		tokenConditionMap: make(map[string]string),
		builderFeeRates:   make(map[string]BuilderFeeRate),
		httpClient:        &http.Client{Transport: &http.Transport{Proxy: http.ProxyFromEnvironment}, Timeout: DefaultTimeout},
	}
	for _, o := range opts {
		o(c)
	}

	if c.ethClient == nil {
		var dialErr error
		for _, u := range c.polygonRPCURLs {
			c.ethClient, dialErr = ethclient.Dial(u)
			if dialErr == nil {
				break
			}
		}
		if c.ethClient == nil {
			return nil, fmt.Errorf("dial ethereum rpc: %w", dialErr)
		}
	}

	var cid *big.Int
	if c.chainIDOverride != nil {
		cid = new(big.Int).Set(c.chainIDOverride)
	} else {
		cid, err = c.ethClient.NetworkID(context.Background())
		if err != nil {
			return nil, fmt.Errorf("chain id: %w", err)
		}
	}
	c.chainID = cid

	if !c.skipL2APIKeyBootstrap && c.apiKeyCredentials == nil {
		creds, err := c.CreateOrDeriveAPIKey()
		if err != nil {
			return nil, fmt.Errorf("clob l2 api key: %w", err)
		}
		if creds == nil || creds.ApiKey == "" || creds.Secret == "" || creds.Passphrase == "" {
			return nil, fmt.Errorf("clob l2 api key: empty credentials from server")
		}
		c.apiKeyCredentials = creds
	}

	return c, nil
}

// SetAPIKeyCredentials sets L2 API credentials.
func (c *Client) SetAPIKeyCredentials(cred *APIKeyCredentials) { c.apiKeyCredentials = cred }

// Credentials returns current L2 credentials.
func (c *Client) Credentials() *APIKeyCredentials { return c.apiKeyCredentials }

// PrivateKey returns the signer private key (read-only use).
func (c *Client) PrivateKey() *ecdsa.PrivateKey { return c.privateKey }

// Address returns the signer Ethereum address.
func (c *Client) Address() common.Address { return c.fromAddress }

func (c *Client) FunderAddress() common.Address {
	return c.funderAddress
}

// SetFunderAddress sets the on-chain maker for new signed orders (same semantics as [WithFunderAddress]).
func (c *Client) SetFunderAddress(addr common.Address) { c.funderAddress = addr }

// SetSignatureType sets EIP712 signature type for new orders (EOA / proxy / Gnosis safe).
func (c *Client) SetSignatureType(t model.SignatureType) { c.signatureType = t }

// SetForceNegRiskSigning forces EIP-712 verifyingContract to the neg-risk exchange for every order
// (see [WithForceNegRiskSigning]). Non–neg-risk markets may reject orders when this is enabled.
func (c *Client) SetForceNegRiskSigning(v bool) { c.forceNegRiskExchange = v }

// ChainID returns the configured chain id.
func (c *Client) ChainID() *big.Int { return new(big.Int).Set(c.chainID) }

// CLOBHost returns the CLOB base URL.
func (c *Client) Host() string { return c.clobHost }

func (c *Client) ensureMetadataCaches() {
	if c.tickSizes == nil {
		c.tickSizes = make(map[string]string)
	}
	if c.tickSizeAt == nil {
		c.tickSizeAt = make(map[string]time.Time)
	}
	if c.negRiskCache == nil {
		c.negRiskCache = make(map[string]bool)
	}
	if c.feeRateCache == nil {
		c.feeRateCache = make(map[string]int)
	}
	if c.feeInfoCache == nil {
		c.feeInfoCache = make(map[string]FeeInfo)
	}
	if c.tokenConditionMap == nil {
		c.tokenConditionMap = make(map[string]string)
	}
	if c.builderFeeRates == nil {
		c.builderFeeRates = make(map[string]BuilderFeeRate)
	}
}

func (c *Client) requireL1() error {
	if c.privateKey == nil {
		return ErrL1AuthRequired
	}
	return nil
}

func (c *Client) requireL2() error {
	if err := c.requireL1(); err != nil {
		return err
	}
	if c.apiKeyCredentials == nil {
		return ErrL2AuthRequired
	}
	return nil
}

func (c *Client) requireBuilder() error {
	if c.builderSigner == nil {
		return ErrBuilderAuthMissing
	}
	return nil
}

// l2Headers merges optional builder headers when useBuilder is true (matches clob-client builder flow).
func (c *Client) builderHeadersOnly(method, path, body string) (map[string]string, error) {
	if c.builderSigner == nil {
		return nil, ErrBuilderAuthMissing
	}
	h, err := c.builderSigner.SignBuilder(method, path, body)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBuilderAuthFailed, err)
	}
	if len(h) == 0 {
		return nil, ErrBuilderAuthFailed
	}
	return h, nil
}

func (c *Client) l2Headers(method, path, body string, useBuilder bool) (map[string]string, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	h, err := c.buildL2AuthHeaders(method, path, body)
	if err != nil {
		return nil, err
	}
	if !useBuilder || c.builderSigner == nil {
		return h, nil
	}
	bh, err := c.builderSigner.SignBuilder(method, path, body)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBuilderAuthFailed, err)
	}
	if len(bh) == 0 {
		return h, nil
	}
	for k, v := range bh {
		h[k] = v
	}
	return h, nil
}
