package polymarket

import (
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/polymarket/go-order-utils/pkg/model"

	sdk "github.com/0xfakeSpike/polymarket-go/polymarket"
)

// Core client and common public types.
type (
	Client            = sdk.Client
	ClientOption      = sdk.ClientOption
	APIKeyCredentials = sdk.APIKeyCredentials

	SearchParams  = sdk.SearchParams
	SearchResults = sdk.SearchResults
	Event         = sdk.Event
	Market        = sdk.Market
	Book          = sdk.Book
	BookParams    = sdk.BookParams

	OrderRequest      = sdk.OrderRequest
	MarketOrderRequest = sdk.MarketOrderRequest
	SignedOrderV2     = sdk.SignedOrderV2
	OrderResponse     = sdk.OrderResponse
)

const (
	DefaultBaseURL = sdk.DefaultBaseURL
	DataAPIBaseURL = sdk.DataAPIBaseURL
	CLOBHost       = sdk.CLOBHost
	DefaultTimeout = sdk.DefaultTimeout
)

func NewClient(privateKeyHex string, opts ...ClientOption) (*Client, error) {
	return sdk.NewClient(privateKeyHex, opts...)
}

func NewPublicClient(opts ...ClientOption) (*Client, error) {
	return sdk.NewPublicClient(opts...)
}

// Option re-exports for convenient root import usage.
func WithHTTPClient(h *http.Client) ClientOption                 { return sdk.WithHTTPClient(h) }
func WithGammaBaseURL(u string) ClientOption                     { return sdk.WithGammaBaseURL(u) }
func WithDataAPIBaseURL(u string) ClientOption                   { return sdk.WithDataAPIBaseURL(u) }
func WithCLOBHost(host string) ClientOption                      { return sdk.WithCLOBHost(host) }
func WithEthereumClient(ec *ethclient.Client) ClientOption       { return sdk.WithEthereumClient(ec) }
func WithChainID(id *big.Int) ClientOption                       { return sdk.WithChainID(id) }
func WithPolygonRPCURLs(urls []string) ClientOption              { return sdk.WithPolygonRPCURLs(urls) }
func WithGeoBlockToken(token string) ClientOption                { return sdk.WithGeoBlockToken(token) }
func WithUseServerTime(v bool) ClientOption                      { return sdk.WithUseServerTime(v) }
func WithRetryPostOnError(v bool) ClientOption                   { return sdk.WithRetryPostOnError(v) }
func WithThrowOnError(v bool) ClientOption                       { return sdk.WithThrowOnError(v) }
func WithTickSizeTTL(d time.Duration) ClientOption               { return sdk.WithTickSizeTTL(d) }
func WithSignatureType(t model.SignatureType) ClientOption       { return sdk.WithSignatureType(t) }
func WithFunderAddress(addr common.Address) ClientOption         { return sdk.WithFunderAddress(addr) }
func WithPolymarketSafeMaker(maker common.Address) ClientOption  { return sdk.WithPolymarketSafeMaker(maker) }
func WithPolymarketProxyMaker(maker common.Address) ClientOption { return sdk.WithPolymarketProxyMaker(maker) }
func WithBuilderSigner(s sdk.BuilderSigner) ClientOption         { return sdk.WithBuilderSigner(s) }
func WithAPIKeyCredentials(cred *APIKeyCredentials) ClientOption { return sdk.WithAPIKeyCredentials(cred) }
func WithSkipL2APIKeyBootstrap() ClientOption                    { return sdk.WithSkipL2APIKeyBootstrap() }
func WithForceNegRiskSigning() ClientOption                      { return sdk.WithForceNegRiskSigning() }
