package polymarket

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/polymarket/go-order-utils/pkg/model"
)

// CreateOrder builds and signs a limit order without posting it.
func (c *Client) CreateOrder(req OrderRequest) (*SignedOrderV2, error) {
	return c.buildSignedLimitOrder(req)
}

// CreateAndPostOrder builds, signs, and posts a limit order.
func (c *Client) CreateAndPostOrder(req OrderRequest, orderType string, postOnly, deferExec bool) (*OrderResponse, error) {
	if orderType == "" {
		orderType = OrderTypeGTC
	}
	signed, err := c.CreateOrder(req)
	if err != nil {
		return nil, err
	}
	return c.postSignedOrder(signed, orderType, deferExec, postOnly)
}

// validateOrderSignatureConfig enforces clob-client / CLOB rules: POLY_PROXY and POLY_GNOSIS_SAFE
// orders must use a separate on-chain maker (funder); otherwise the API often returns "invalid signature".
func (c *Client) validateOrderSignatureConfig() error {
	switch c.signatureType {
	case model.EOA:
		return nil
	case model.POLY_PROXY, model.POLY_GNOSIS_SAFE:
		if c.funderAddress == (common.Address{}) {
			return fmt.Errorf("signatureType %d requires non-zero funderAddress (maker); use WithPolymarketSafeMaker, WithPolymarketProxyMaker, or WithFunderAddress", c.signatureType)
		}
		return nil
	default:
		return nil
	}
}

func (c *Client) buildSignedLimitOrder(req OrderRequest) (*SignedOrderV2, error) {
	if err := c.validateOrderSignatureConfig(); err != nil {
		return nil, err
	}
	tick, err := c.GetTickSize(req.TokenID)
	if err != nil {
		return nil, fmt.Errorf("tick size: %w", err)
	}
	cfg, ok := roundingConfig[tick]
	if !ok {
		return nil, fmt.Errorf("unsupported tick size %q (expected 0.1, 0.01, 0.001, 0.0001)", tick)
	}
	okPrice, err := priceValid(req.Price, tick)
	if err != nil {
		return nil, err
	}
	minTick, _ := strconv.ParseFloat(tick, 64)
	if !okPrice {
		return nil, fmt.Errorf("invalid price (%g), min: %g - max: %g", req.Price, minTick, 1-minTick)
	}
	amounts, err := buildOrderCreationArgs(req, cfg)
	if err != nil {
		return nil, err
	}
	return c.newSignedOrderV2(
		req.TokenID,
		amounts.MakerAmount,
		amounts.TakerAmount,
		req.Side,
		c.signatureType,
		req.Metadata,
		req.BuilderCode,
		req.Expiration,
	)
}

type postOrderEnvelope struct {
	DeferExec bool           `json:"deferExec"`
	Order     postOrderInner `json:"order"`
	Owner     string         `json:"owner"`
	OrderType string         `json:"orderType"`
	PostOnly  *bool          `json:"postOnly,omitempty"`
}

// newOrderWireFromSigned mirrors clob-client src/utilities.ts orderToJson (POST /order body).
type postOrderInner struct {
	Salt          json.Number `json:"salt"`
	Maker         string      `json:"maker"`
	Signer        string      `json:"signer"`
	TokenID       string      `json:"tokenId"`
	MakerAmount   string      `json:"makerAmount"`
	TakerAmount   string      `json:"takerAmount"`
	Side          string      `json:"side"`
	SignatureType int         `json:"signatureType"`
	Timestamp     string      `json:"timestamp"`
	Metadata      string      `json:"metadata"`
	Builder       string      `json:"builder"`
	Expiration    string      `json:"expiration"`
	Signature     string      `json:"signature"`
}

func newOrderWireFromSigned(s *SignedOrderV2, owner, orderType string, deferExec, postOnly bool) *postOrderEnvelope {
	env := &postOrderEnvelope{
		DeferExec: deferExec,
		Order: postOrderInner{
			Salt:          json.Number(s.Salt.String()),
			Maker:         s.Maker.Hex(),
			Signer:        s.Signer.Hex(),
			TokenID:       s.TokenID.String(),
			MakerAmount:   s.MakerAmount.String(),
			TakerAmount:   s.TakerAmount.String(),
			Side:          string(s.Side),
			SignatureType: int(s.SignatureType),
			Timestamp:     s.Timestamp.String(),
			Metadata:      s.Metadata.Hex(),
			Builder:       s.Builder.Hex(),
			Expiration:    s.Expiration.String(),
			Signature:     encodeOrderSignatureHex(s.Signature),
		},
		Owner:     owner,
		OrderType: orderType,
	}
	if postOnly {
		t := true
		env.PostOnly = &t
	}
	return env
}

func (c *Client) postSignedOrder(signed *SignedOrderV2, orderType string, deferExec, postOnly bool) (*OrderResponse, error) {
	if c.apiKeyCredentials == nil {
		return nil, fmt.Errorf("API key not set, please call SetAPIKeyCredentials first")
	}
	if postOnly && (orderType == OrderTypeFOK || orderType == OrderTypeFAK) {
		return nil, fmt.Errorf("postOnly is not supported for FOK/FAK orders")
	}
	payload := newOrderWireFromSigned(signed, c.apiKeyCredentials.ApiKey, orderType, deferExec, postOnly)
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal order: %w", err)
	}

	path := PathPostOrder
	headers, err := c.l2Headers("POST", path, string(jsonData), true)
	if err != nil {
		return nil, err
	}

	body, err := c.clobRequest("POST", path, nil, headers, jsonData)
	if err != nil {
		return nil, fmt.Errorf("submit order: %w", err)
	}

	var orderResp OrderResponse
	if err := json.Unmarshal(body, &orderResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if !orderResp.Success {
		return &orderResp, fmt.Errorf("order submission failed: %s", orderResp.ErrorMsg)
	}
	return &orderResp, nil
}

// GetOrder retrieves a single open order by ID (order hash) from the CLOB API (GET /data/order/{id}).
func (c *Client) GetOrder(orderID string, _ int64) (*OrderRef, error) {
	if c.apiKeyCredentials == nil {
		return nil, fmt.Errorf("API key not set, please call SetAPIKeyCredentials first")
	}

	path := PathDataOrderPrefix + orderID
	headers, err := c.l2Headers("GET", path, "", true)
	if err != nil {
		return nil, fmt.Errorf("l2 headers: %w", err)
	}

	body, err := c.clobRequest("GET", path, nil, headers, nil)
	if err != nil {
		return nil, fmt.Errorf("get order: %w", err)
	}

	var orderResp GetOrderResponse
	if err := json.Unmarshal(body, &orderResp); err != nil {
		return nil, fmt.Errorf("decode order response: %w", err)
	}
	return &OrderRef{ID: orderResp.ID, Status: orderResp.Status}, nil
}
