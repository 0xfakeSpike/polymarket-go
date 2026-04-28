package polymarket

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
)

// OrderBookSummaryHash computes SHA-1 over JSON of the book with hash cleared (matches clob-client utilities.ts).
func OrderBookSummaryHash(book *Book) (string, error) {
	if book == nil {
		return "", fmt.Errorf("nil book")
	}
	cp := *book
	cp.Hash = ""
	b, err := json.Marshal(cp)
	if err != nil {
		return "", err
	}
	sum := sha1.Sum(b)
	return hex.EncodeToString(sum[:]), nil
}

// CalculateBuyMarketPrice walks asks (ascending price) to match a USDC amount (matches TS calculateBuyMarketPrice).
func CalculateBuyMarketPrice(asks []LimitOrder, amountToMatch float64, orderType string) (float64, error) {
	if len(asks) == 0 {
		return 0, fmt.Errorf("no match")
	}
	// TS iterates from end (highest price first) — asks sorted low→high in API; TS sorts asks ascending then iterates reversed
	sum := 0.0
	for i := len(asks) - 1; i >= 0; i-- {
		p := asks[i]
		sum += p.Size * p.Price
		if sum >= amountToMatch {
			return p.Price, nil
		}
	}
	if orderType == OrderTypeFOK {
		return 0, fmt.Errorf("no match")
	}
	return asks[0].Price, nil
}

// CalculateSellMarketPrice walks bids to match share amount.
func CalculateSellMarketPrice(bids []LimitOrder, amountToMatch float64, orderType string) (float64, error) {
	if len(bids) == 0 {
		return 0, fmt.Errorf("no match")
	}
	sum := 0.0
	for i := len(bids) - 1; i >= 0; i-- {
		p := bids[i]
		sum += p.Size
		if sum >= amountToMatch {
			return p.Price, nil
		}
	}
	if orderType == OrderTypeFOK {
		return 0, fmt.Errorf("no match")
	}
	return bids[0].Price, nil
}

// CalculateMarketPrice loads the order book and computes the market price for side/amount (FOK default).
func (c *Client) CalculateMarketPrice(tokenID string, side Side, amount float64, orderType string) (float64, error) {
	book, err := c.GetOrderBook(tokenID)
	if err != nil {
		return 0, err
	}
	if side == SideBuy {
		asks := book.AsksData()
		if len(asks) == 0 {
			return 0, fmt.Errorf("no match")
		}
		return CalculateBuyMarketPrice(asks, amount, orderType)
	}
	bids := book.BidsData()
	if len(bids) == 0 {
		return 0, fmt.Errorf("no match")
	}
	return CalculateSellMarketPrice(bids, amount, orderType)
}

// CreateMarketOrder builds and signs a market order (price optional — filled from book when empty).
func (c *Client) CreateMarketOrder(req MarketOrderRequest, orderType string) (*SignedOrderV2, error) {
	tick, err := c.ResolveTickSize(req.TokenID, nil)
	if err != nil {
		return nil, err
	}
	cfg, ok := roundingConfig[tick]
	if !ok {
		return nil, fmt.Errorf("unsupported tick size %q", tick)
	}
	price := req.Price
	if price == 0 {
		ot := orderType
		if ot == "" {
			ot = OrderTypeFOK
		}
		price, err = c.CalculateMarketPrice(req.TokenID, req.Side, req.Amount, ot)
		if err != nil {
			return nil, err
		}
	}
	if ok, err := priceValid(price, tick); err != nil || !ok {
		if err != nil {
			return nil, err
		}
		minT, _ := strconv.ParseFloat(tick, 64)
		return nil, fmt.Errorf("invalid price (%g), min: %g - max: %g", price, minT, 1-minT)
	}
	if err := c.validateOrderSignatureConfig(); err != nil {
		return nil, err
	}
	amounts, err := buildMarketOrderCreationArgs(req, price, cfg)
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
		nil,
	)
}

// CreateAndPostMarketOrder creates, signs, and posts a market order.
func (c *Client) CreateAndPostMarketOrder(req MarketOrderRequest, orderType string, deferExec bool) (*OrderResponse, error) {
	if orderType == "" {
		orderType = OrderTypeFOK
	}
	signed, err := c.CreateMarketOrder(req, orderType)
	if err != nil {
		return nil, err
	}
	return c.postSignedOrder(signed, orderType, deferExec, false)
}
