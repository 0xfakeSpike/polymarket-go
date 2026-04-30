package polymarket

import (
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/polymarket/go-order-utils/pkg/model"
)

// NewPublicClient builds a client without private-key auth, useful for read-only CLOB APIs.
func NewPublicClient(opts ...ClientOption) (*Client, error) {
	c := &Client{
		clobHost:              CLOBHost,
		polygonRPCURLs:        defaultPolygonRPCURLs,
		signatureType:         model.POLY_GNOSIS_SAFE,
		tickSizes:             make(map[string]string),
		tickSizeAt:            make(map[string]time.Time),
		tickSizeTTL:           defaultTickCache,
		negRiskCache:          make(map[string]bool),
		feeRateCache:          make(map[string]int),
		feeInfoCache:          make(map[string]FeeInfo),
		tokenConditionMap:     make(map[string]string),
		builderFeeRates:       make(map[string]BuilderFeeRate),
		httpClient:            &http.Client{Transport: &http.Transport{Proxy: http.ProxyFromEnvironment}, Timeout: DefaultTimeout},
		skipL2APIKeyBootstrap: true,
		chainID:               big.NewInt(137),
	}
	for _, o := range opts {
		o(c)
	}
	if c.gammaHost == "" {
		c.gammaHost = GammaAPIHost
	}
	if c.dataHost == "" {
		c.dataHost = DataAPIHost
	}
	if c.httpClient == nil {
		return nil, fmt.Errorf("nil http client")
	}
	if c.chainID == nil {
		c.chainID = big.NewInt(137)
	}
	return c, nil
}
