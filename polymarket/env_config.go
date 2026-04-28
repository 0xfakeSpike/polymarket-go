package polymarket

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/polymarket/go-order-utils/pkg/model"
)

// EnvConfig holds polymarket runtime configuration loaded from environment variables.
type EnvConfig struct {
	PrivateKeyHex string

	APIKey     string
	APISecret  string
	APIPass    string
	SkipL2Auth bool

	GammaBaseURL string
	DataBaseURL  string
	CLOBHost     string

	GeoBlockToken string
	UseServerTime bool
	ThrowOnError  bool
	RetryPost     bool
	ForceNegRisk  bool

	SignatureType model.SignatureType
	FunderAddress common.Address

	Timeout        time.Duration
	TickSizeTTL    time.Duration
	PolygonRPCURLs []string
}

// LoadEnvConfig reads all polymarket SDK settings from environment variables.
func LoadEnvConfig() (EnvConfig, error) {
	cfg := EnvConfig{
		PrivateKeyHex: strings.TrimSpace(os.Getenv("POLYMARKET_PRIVATE_KEY")),
		APIKey:        strings.TrimSpace(os.Getenv("POLYMARKET_API_KEY")),
		APISecret:     strings.TrimSpace(os.Getenv("POLYMARKET_API_SECRET")),
		APIPass:       strings.TrimSpace(os.Getenv("POLYMARKET_API_PASSPHRASE")),
		GammaBaseURL:  strings.TrimSpace(os.Getenv("POLYMARKET_GAMMA_BASE_URL")),
		DataBaseURL:   strings.TrimSpace(os.Getenv("POLYMARKET_DATA_API_BASE_URL")),
		CLOBHost:      strings.TrimSpace(os.Getenv("POLYMARKET_CLOB_HOST")),
		GeoBlockToken: strings.TrimSpace(os.Getenv("POLYMARKET_GEO_BLOCK_TOKEN")),
		SignatureType: model.POLY_GNOSIS_SAFE,
		Timeout:       DefaultTimeout,
		TickSizeTTL:   defaultTickCache,
	}

	var err error
	if cfg.SkipL2Auth, err = envBool("POLYMARKET_SKIP_L2_BOOTSTRAP", false); err != nil {
		return EnvConfig{}, err
	}
	if cfg.UseServerTime, err = envBool("POLYMARKET_USE_SERVER_TIME", false); err != nil {
		return EnvConfig{}, err
	}
	if cfg.ThrowOnError, err = envBool("POLYMARKET_THROW_ON_ERROR", false); err != nil {
		return EnvConfig{}, err
	}
	if cfg.RetryPost, err = envBool("POLYMARKET_RETRY_POST_ON_ERROR", false); err != nil {
		return EnvConfig{}, err
	}
	if cfg.ForceNegRisk, err = envBool("POLYMARKET_FORCE_NEG_RISK_SIGNING", false); err != nil {
		return EnvConfig{}, err
	}

	if raw := strings.TrimSpace(os.Getenv("POLYMARKET_TIMEOUT")); raw != "" {
		d, parseErr := time.ParseDuration(raw)
		if parseErr != nil {
			return EnvConfig{}, fmt.Errorf("POLYMARKET_TIMEOUT: %w", parseErr)
		}
		cfg.Timeout = d
	}
	if raw := strings.TrimSpace(os.Getenv("POLYMARKET_TICK_SIZE_TTL")); raw != "" {
		d, parseErr := time.ParseDuration(raw)
		if parseErr != nil {
			return EnvConfig{}, fmt.Errorf("POLYMARKET_TICK_SIZE_TTL: %w", parseErr)
		}
		cfg.TickSizeTTL = d
	}

	if raw := strings.TrimSpace(os.Getenv("POLYMARKET_SIGNATURE_TYPE")); raw != "" {
		t, parseErr := parseSignatureType(raw)
		if parseErr != nil {
			return EnvConfig{}, parseErr
		}
		cfg.SignatureType = t
	}

	if raw := strings.TrimSpace(os.Getenv("POLYMARKET_FUNDER_ADDRESS")); raw != "" {
		if !common.IsHexAddress(raw) {
			return EnvConfig{}, fmt.Errorf("POLYMARKET_FUNDER_ADDRESS: invalid hex address %q", raw)
		}
		cfg.FunderAddress = common.HexToAddress(raw)
	}

	if raw := strings.TrimSpace(os.Getenv("POLYMARKET_RPC_URLS")); raw != "" {
		parts := strings.Split(raw, ",")
		cfg.PolygonRPCURLs = make([]string, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				cfg.PolygonRPCURLs = append(cfg.PolygonRPCURLs, p)
			}
		}
	}

	if cfg.PrivateKeyHex == "" && (cfg.APIKey != "" || cfg.APISecret != "" || cfg.APIPass != "") {
		return EnvConfig{}, fmt.Errorf("api key credentials require POLYMARKET_PRIVATE_KEY")
	}
	if (cfg.APIKey == "") != (cfg.APISecret == "") || (cfg.APIKey == "") != (cfg.APIPass == "") {
		return EnvConfig{}, fmt.Errorf("POLYMARKET_API_KEY / POLYMARKET_API_SECRET / POLYMARKET_API_PASSPHRASE must be all set or all empty")
	}

	return cfg, nil
}

// NewClientFromEnv creates a polymarket client from environment variables.
func NewClientFromEnv() (*Client, error) {
	cfg, err := LoadEnvConfig()
	if err != nil {
		return nil, err
	}
	return cfg.NewClient()
}

// NewClient builds either an authenticated client (when private key is set) or a public-only client.
func (c EnvConfig) NewClient() (*Client, error) {
	opts := []ClientOption{
		WithUseServerTime(c.UseServerTime),
		WithThrowOnError(c.ThrowOnError),
		WithRetryPostOnError(c.RetryPost),
		WithTickSizeTTL(c.TickSizeTTL),
		WithSignatureType(c.SignatureType),
	}
	if c.Timeout > 0 {
		opts = append(opts, WithHTTPClient(&http.Client{Transport: &http.Transport{Proxy: http.ProxyFromEnvironment}, Timeout: c.Timeout}))
	}
	if c.GammaBaseURL != "" {
		opts = append(opts, WithGammaBaseURL(c.GammaBaseURL))
	}
	if c.DataBaseURL != "" {
		opts = append(opts, WithDataAPIBaseURL(c.DataBaseURL))
	}
	if c.CLOBHost != "" {
		opts = append(opts, WithCLOBHost(c.CLOBHost))
	}
	if c.GeoBlockToken != "" {
		opts = append(opts, WithGeoBlockToken(c.GeoBlockToken))
	}
	if len(c.PolygonRPCURLs) > 0 {
		opts = append(opts, WithPolygonRPCURLs(c.PolygonRPCURLs))
	}
	if c.FunderAddress != (common.Address{}) {
		opts = append(opts, WithFunderAddress(c.FunderAddress))
	}
	if c.ForceNegRisk {
		opts = append(opts, WithForceNegRiskSigning())
	}
	if c.SkipL2Auth {
		opts = append(opts, WithSkipL2APIKeyBootstrap())
	}
	if c.APIKey != "" {
		opts = append(opts, WithAPIKeyCredentials(&APIKeyCredentials{ApiKey: c.APIKey, Secret: c.APISecret, Passphrase: c.APIPass}))
	}

	if c.PrivateKeyHex == "" {
		return NewPublicClient(opts...)
	}
	return NewClient(c.PrivateKeyHex, opts...)
}

func envBool(name string, fallback bool) (bool, error) {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return fallback, nil
	}
	v, err := strconv.ParseBool(raw)
	if err != nil {
		return false, fmt.Errorf("%s: %w", name, err)
	}
	return v, nil
}

func parseSignatureType(raw string) (model.SignatureType, error) {
	raw = strings.TrimSpace(strings.ToLower(raw))
	switch raw {
	case "eoa", "0":
		return model.EOA, nil
	case "poly_proxy", "proxy", "1":
		return model.POLY_PROXY, nil
	case "poly_gnosis_safe", "safe", "gnosis", "2":
		return model.POLY_GNOSIS_SAFE, nil
	default:
		return model.POLY_GNOSIS_SAFE, fmt.Errorf("POLYMARKET_SIGNATURE_TYPE: unsupported value %q", raw)
	}
}
