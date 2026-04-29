package polymarket

import (
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/polymarket/go-order-utils/pkg/model"
)

// ClientOption configures Client construction.
type ClientOption func(*Client)

// WithHTTPClient sets the HTTP client used for all APIs.
func WithHTTPClient(h *http.Client) ClientOption {
	return func(c *Client) { c.httpClient = h }
}

// WithCLOBHost sets the CLOB API base URL (no trailing slash).
func WithCLOBHost(host string) ClientOption {
	return func(c *Client) { c.clobHost = host }
}

// WithEthereumClient sets an existing eth client (otherwise one is dialed).
func WithEthereumClient(ec *ethclient.Client) ClientOption {
	return func(c *Client) { c.ethClient = ec }
}

// WithChainID overrides chain id (otherwise from eth client NetworkID).
func WithChainID(id *big.Int) ClientOption {
	return func(c *Client) { c.chainIDOverride = id }
}

// WithPolygonRPCURLs sets RPC endpoints tried in order for dialing.
func WithPolygonRPCURLs(urls []string) ClientOption {
	return func(c *Client) { c.polygonRPCURLs = urls }
}

// WithGeoBlockToken sets geo_block_token query param on CLOB calls.
func WithGeoBlockToken(token string) ClientOption {
	return func(c *Client) { c.geoBlockToken = token }
}

// WithUseServerTime uses GET /time for L1/L2 auth timestamps.
func WithUseServerTime(v bool) ClientOption {
	return func(c *Client) { c.useServerTime = v }
}

// WithRetryPostOnError retries failed POST once after 30ms (matches clob-client).
func WithRetryPostOnError(v bool) ClientOption {
	return func(c *Client) { c.retryPostOnError = v }
}

// WithThrowOnError turns JSON `{error: ...}` payloads on 2xx into ApiError returns.
func WithThrowOnError(v bool) ClientOption {
	return func(c *Client) { c.throwOnError = v }
}

// WithTickSizeTTL sets tick size cache duration.
func WithTickSizeTTL(d time.Duration) ClientOption {
	return func(c *Client) { c.tickSizeTTL = d }
}

// WithSignatureType sets CLOB order signature type (EOA / proxy / safe).
func WithSignatureType(t model.SignatureType) ClientOption {
	return func(c *Client) { c.signatureType = t }
}

// WithFunderAddress sets the on-chain maker (funder) for signed orders, matching clob-client's
// optional funderAddress passed into OrderBuilder. When unset (zero address), maker is the signer EOA.
func WithFunderAddress(addr common.Address) ClientOption {
	return func(c *Client) { c.funderAddress = addr }
}

// WithPolymarketSafeMaker configures orders like polymarket.com when the UI uses a Gnosis-Safe-style
// maker (HTTP body order.maker) with signatureType 2 ([model.POLY_GNOSIS_SAFE]) and EOA signer.
// Pass the same maker address shown on your profile / in a captured POST /order payload.
func WithPolymarketSafeMaker(maker common.Address) ClientOption {
	return func(c *Client) {
		c.funderAddress = maker
		c.signatureType = model.POLY_GNOSIS_SAFE
	}
}

// WithPolymarketProxyMaker is the POLY_PROXY ([model.POLY_PROXY]) analogue of [WithPolymarketSafeMaker].
func WithPolymarketProxyMaker(maker common.Address) ClientOption {
	return func(c *Client) {
		c.funderAddress = maker
		c.signatureType = model.POLY_PROXY
	}
}

// WithAPIKeyCredentials sets L2 credentials and skips automatic key bootstrap in [NewClient].
func WithAPIKeyCredentials(cred *APIKeyCredentials) ClientOption {
	return func(c *Client) {
		if cred != nil {
			c.apiKeyCredentials = cred
		}
	}
}

// WithSkipL2APIKeyBootstrap skips the default [CreateOrDeriveAPIKey] call in [NewClient].
// Trading endpoints will still require [SetAPIKeyCredentials] later.
func WithSkipL2APIKeyBootstrap() ClientOption {
	return func(c *Client) { c.skipL2APIKeyBootstrap = true }
}

// WithForceNegRiskSigning signs every CLOB order with the v2 neg-risk exchange verifying contract
// (`0xe2222d279d744050d28e00520010520000310F59`), ignoring GET /neg-risk per token.
// Use this when your flow always targets neg-risk markets.
func WithForceNegRiskSigning() ClientOption {
	return func(c *Client) { c.forceNegRiskExchange = true }
}
