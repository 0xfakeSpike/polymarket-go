package polymarket

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const bytes32Zero = "0x0000000000000000000000000000000000000000000000000000000000000000"

func tradeParamsValues(p *TradeParams, nextCursor string) url.Values {
	q := url.Values{}
	if nextCursor != "" {
		q.Set("next_cursor", nextCursor)
	}
	if p == nil {
		return q
	}
	if p.ID != "" {
		q.Set("id", p.ID)
	}
	if p.MakerAddress != "" {
		q.Set("maker_address", p.MakerAddress)
	}
	if p.Market != "" {
		q.Set("market", p.Market)
	}
	if p.AssetID != "" {
		q.Set("asset_id", p.AssetID)
	}
	if p.Before != "" {
		q.Set("before", p.Before)
	}
	if p.After != "" {
		q.Set("after", p.After)
	}
	return q
}

func builderTradeParamsValues(p *BuilderTradeParams, nextCursor string) url.Values {
	var base *TradeParams
	if p != nil {
		base = &p.TradeParams
	}
	q := tradeParamsValues(base, nextCursor)
	if p != nil && p.BuilderCode != "" {
		q.Set("builder_code", p.BuilderCode)
	}
	return q
}

func openOrderParamsValues(p *OpenOrderParams, nextCursor string) url.Values {
	q := url.Values{}
	if nextCursor != "" {
		q.Set("next_cursor", nextCursor)
	}
	if p == nil {
		return q
	}
	if p.ID != "" {
		q.Set("id", p.ID)
	}
	if p.Market != "" {
		q.Set("market", p.Market)
	}
	if p.AssetID != "" {
		q.Set("asset_id", p.AssetID)
	}
	return q
}

// GetTrades fetches all trades matching params by walking next_cursor until end (unless onlyFirstPage).
func (c *Client) GetTrades(params *TradeParams, onlyFirstPage bool, nextCursor string) ([]Trade, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	if nextCursor == "" {
		nextCursor = InitialCursor
	}
	var all []Trade
	for nextCursor != EndCursor && (nextCursor == InitialCursor || !onlyFirstPage) {
		path := PathDataTrades
		bodyStr := ""
		h, err := c.l2Headers("GET", path, bodyStr)
		if err != nil {
			return nil, err
		}
		q := tradeParamsValues(params, nextCursor)
		data, err := c.clobRequest("GET", path, q, h, nil)
		if err != nil {
			return nil, err
		}
		var page TradesPage
		if err := json.Unmarshal(data, &page); err != nil {
			return nil, err
		}
		all = append(all, page.Trades...)
		nextCursor = page.NextCursor
		if onlyFirstPage {
			break
		}
	}
	return all, nil
}

// GetTradesPaginated returns a single page of trades.
func (c *Client) GetTradesPaginated(params *TradeParams, nextCursor string) (*TradesPage, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	if nextCursor == "" {
		nextCursor = InitialCursor
	}
	path := PathDataTrades
	h, err := c.l2Headers("GET", path, "")
	if err != nil {
		return nil, err
	}
	q := tradeParamsValues(params, nextCursor)
	data, err := c.clobRequest("GET", path, q, h, nil)
	if err != nil {
		return nil, err
	}
	var page TradesPage
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, err
	}
	return &page, nil
}

// GetBuilderTrades returns builder-attributed trades.
func (c *Client) GetBuilderTrades(params *BuilderTradeParams, nextCursor string) (*BuilderTradesPage, error) {
	if params == nil || params.BuilderCode == "" || params.BuilderCode == bytes32Zero {
		return nil, fmt.Errorf("builderCode is required and cannot be zero")
	}
	if nextCursor == "" {
		nextCursor = InitialCursor
	}
	path := PathBuilderTrades
	q := builderTradeParamsValues(params, nextCursor)
	data, err := c.clobRequest("GET", path, q, nil, nil)
	if err != nil {
		return nil, err
	}
	var page BuilderTradesPage
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, err
	}
	return &page, nil
}

// GetPreMigrationOrders walks all pages of pre-migration orders unless onlyFirstPage.
func (c *Client) GetPreMigrationOrders(onlyFirstPage bool, nextCursor string) ([]OpenOrder, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	if nextCursor == "" {
		nextCursor = InitialCursor
	}
	var all []OpenOrder
	path := PathPreMigrationOrders
	for nextCursor != EndCursor && (nextCursor == InitialCursor || !onlyFirstPage) {
		h, err := c.l2Headers("GET", path, "")
		if err != nil {
			return nil, err
		}
		q := url.Values{}
		q.Set("next_cursor", nextCursor)
		data, err := c.clobRequest("GET", path, q, h, nil)
		if err != nil {
			return nil, err
		}
		var page OpenOrdersPage
		if err := json.Unmarshal(data, &page); err != nil {
			return nil, err
		}
		all = append(all, page.Orders...)
		nextCursor = page.NextCursor
		if onlyFirstPage {
			break
		}
	}
	return all, nil
}

// GetOpenOrders walks all pages of open orders unless onlyFirstPage.
func (c *Client) GetOpenOrders(params *OpenOrderParams, onlyFirstPage bool, nextCursor string) ([]OpenOrder, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	if nextCursor == "" {
		nextCursor = InitialCursor
	}
	var all []OpenOrder
	path := PathDataOrders
	for nextCursor != EndCursor && (nextCursor == InitialCursor || !onlyFirstPage) {
		h, err := c.l2Headers("GET", path, "")
		if err != nil {
			return nil, err
		}
		q := openOrderParamsValues(params, nextCursor)
		data, err := c.clobRequest("GET", path, q, h, nil)
		if err != nil {
			return nil, err
		}
		var page struct {
			Data       []OpenOrder `json:"data"`
			NextCursor string      `json:"next_cursor"`
		}
		if err := json.Unmarshal(data, &page); err != nil {
			return nil, err
		}
		all = append(all, page.Data...)
		nextCursor = page.NextCursor
		if onlyFirstPage {
			break
		}
	}
	return all, nil
}

// PostOrder posts a signed order with execution type (GTC, GTD, FOK, FAK).
func (c *Client) PostOrder(signed *SignedOrderV2, orderType string, postOnly, deferExec bool) (*OrderResponse, error) {
	return c.postSignedOrder(signed, orderType, deferExec, postOnly)
}

// PostOrderBatchItem is one entry for PostOrders.
type PostOrderBatchItem struct {
	Order     *SignedOrderV2
	OrderType string
}

// PostOrders posts multiple signed orders (POST /orders).
func (c *Client) PostOrders(orders []PostOrderBatchItem, postOnly, deferExec bool) (json.RawMessage, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	outs := make([]*postOrderEnvelope, 0, len(orders))
	if postOnly {
		for _, o := range orders {
			if o.OrderType == OrderTypeFOK || o.OrderType == OrderTypeFAK {
				return nil, fmt.Errorf("postOnly is not supported for FOK/FAK orders")
			}
		}
	}
	for _, o := range orders {
		outs = append(outs, newOrderWireFromSigned(o.Order, c.apiKeyCredentials.ApiKey, o.OrderType, deferExec, postOnly))
	}
	body, err := json.Marshal(outs)
	if err != nil {
		return nil, err
	}
	path := PathPostOrders
	h, err := c.l2Headers("POST", path, string(body))
	if err != nil {
		return nil, err
	}
	return c.clobRequest("POST", path, nil, h, body)
}

// CancelOrder deletes a single order by id payload.
func (c *Client) CancelOrder(payload OrderPayload) (json.RawMessage, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	path := PathCancelOrder
	h, err := c.l2Headers("DELETE", path, string(b))
	if err != nil {
		return nil, err
	}
	return c.clobRequest("DELETE", path, nil, h, b)
}

// CancelOrders deletes multiple orders by order hashes.
func (c *Client) CancelOrders(orderHashes []string) (json.RawMessage, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	b, err := json.Marshal(orderHashes)
	if err != nil {
		return nil, err
	}
	path := PathCancelOrders
	h, err := c.l2Headers("DELETE", path, string(b))
	if err != nil {
		return nil, err
	}
	return c.clobRequest("DELETE", path, nil, h, b)
}

// CancelAll cancels all open orders.
func (c *Client) CancelAll() (json.RawMessage, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	path := PathCancelAll
	h, err := c.l2Headers("DELETE", path, "")
	if err != nil {
		return nil, err
	}
	return c.clobRequest("DELETE", path, nil, h, nil)
}

// CancelMarketOrders cancels by market filter.
func (c *Client) CancelMarketOrders(payload OrderMarketCancelParams) (json.RawMessage, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	path := PathCancelMarket
	h, err := c.l2Headers("DELETE", path, string(b))
	if err != nil {
		return nil, err
	}
	return c.clobRequest("DELETE", path, nil, h, b)
}

// IsOrderScoring GET /order-scoring
func (c *Client) IsOrderScoring(params *OrderScoringParams) (OrderScoring, error) {
	if err := c.requireL2(); err != nil {
		return OrderScoring{}, err
	}
	path := PathOrderScoring
	h, err := c.l2Headers("GET", path, "")
	if err != nil {
		return OrderScoring{}, err
	}
	q := url.Values{}
	if params != nil && params.OrderID != "" {
		q.Set("order_id", params.OrderID)
	}
	data, err := c.clobRequest("GET", path, q, h, nil)
	if err != nil {
		return OrderScoring{}, err
	}
	var out OrderScoring
	if err := json.Unmarshal(data, &out); err != nil {
		return OrderScoring{}, err
	}
	return out, nil
}

// AreOrdersScoring POST /orders-scoring with body = JSON array of order ids.
func (c *Client) AreOrdersScoring(params *OrdersScoringParams) (OrdersScoring, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	var ids []string
	if params != nil {
		ids = params.OrderIDs
	}
	b, err := json.Marshal(ids)
	if err != nil {
		return nil, err
	}
	path := PathOrdersScoring
	h, err := c.l2Headers("POST", path, string(b))
	if err != nil {
		return nil, err
	}
	data, err := c.clobRequest("POST", path, nil, h, b)
	if err != nil {
		return nil, err
	}
	var out OrdersScoring
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	return out, nil
}
