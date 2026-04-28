package polymarket

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
)

// GetOK hits the CLOB health/root endpoint.
func (c *Client) GetOK() (json.RawMessage, error) {
	return c.clobRequest("GET", "/", nil, nil, nil)
}

func (c *Client) getPaginated(path, nextCursor string) (*PaginationPayload, error) {
	if nextCursor == "" {
		nextCursor = InitialCursor
	}
	q := url.Values{}
	q.Set("next_cursor", nextCursor)
	body, err := c.clobRequest("GET", path, q, nil, nil)
	if err != nil {
		return nil, err
	}
	var p PaginationPayload
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// GetSamplingSimplifiedMarkets returns a page of sampling simplified markets.
func (c *Client) GetSamplingSimplifiedMarkets(nextCursor string) (*PaginationPayload, error) {
	return c.getPaginated(PathSamplingSimplifiedMarkets, nextCursor)
}

// GetSamplingMarkets returns a page of sampling markets.
func (c *Client) GetSamplingMarkets(nextCursor string) (*PaginationPayload, error) {
	return c.getPaginated(PathSamplingMarkets, nextCursor)
}

// GetSimplifiedMarkets returns a page of simplified markets.
func (c *Client) GetSimplifiedMarkets(nextCursor string) (*PaginationPayload, error) {
	return c.getPaginated(PathSimplifiedMarkets, nextCursor)
}

// GetMarkets returns a page of markets (CLOB /markets).
func (c *Client) GetMarkets(nextCursor string) (*PaginationPayload, error) {
	return c.getPaginated(PathMarkets, nextCursor)
}

// GetCLOBMarket fetches a single market by condition id from the CLOB API.
func (c *Client) GetCLOBMarket(conditionID string) (json.RawMessage, error) {
	return c.clobRequest("GET", PathMarketPrefix+conditionID, nil, nil, nil)
}

// GetOrderBook returns the order book for a token.
func (c *Client) GetOrderBook(tokenID string) (*Book, error) {
	q := url.Values{}
	q.Set("token_id", tokenID)
	body, err := c.clobRequest("GET", PathOrderBook, q, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get order book: %w", err)
	}
	var book Book
	if err := json.Unmarshal(body, &book); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order book: %w", err)
	}
	sortOrderBookLevels(&book)
	c.rememberTickFromBook(book.AssetID, book.TickSize)
	return &book, nil
}

// sortOrderBookLevels sorts bids high→low and asks low→high (same convention as Book.BidsData / AsksData).
func sortOrderBookLevels(b *Book) {
	sort.Slice(b.Bids, func(i, j int) bool {
		return b.Bids[i].Price > b.Bids[j].Price
	})
	sort.Slice(b.Asks, func(i, j int) bool {
		return b.Asks[i].Price < b.Asks[j].Price
	})
}

// GetOrderBooks returns multiple order books (POST /books).
func (c *Client) GetOrderBooks(params []BookParams) ([]Book, error) {
	b, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	body, err := c.clobRequest("POST", PathOrderBooks, nil, nil, b)
	if err != nil {
		return nil, err
	}
	var books []Book
	if err := json.Unmarshal(body, &books); err != nil {
		return nil, err
	}
	for i := range books {
		sortOrderBookLevels(&books[i])
		c.rememberTickFromBook(books[i].AssetID, books[i].TickSize)
	}
	return books, nil
}

// GetMidpoint returns midpoint for a token.
func (c *Client) GetMidpoint(tokenID string) (json.RawMessage, error) {
	q := url.Values{}
	q.Set("token_id", tokenID)
	return c.clobRequest("GET", PathMidpoint, q, nil, nil)
}

// GetMidpoints batch midpoints.
func (c *Client) GetMidpoints(params []BookParams) (json.RawMessage, error) {
	b, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	return c.clobRequest("POST", PathMidpoints, nil, nil, b)
}

// GetPrice returns best price for token and side ("BUY" / "SELL").
func (c *Client) GetPrice(tokenID, side string) (json.RawMessage, error) {
	q := url.Values{}
	q.Set("token_id", tokenID)
	q.Set("side", side)
	return c.clobRequest("GET", PathPrice, q, nil, nil)
}

// GetPrices batch prices.
func (c *Client) GetPrices(params []BookParams) (json.RawMessage, error) {
	b, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	return c.clobRequest("POST", PathPrices, nil, nil, b)
}

// GetSpread returns spread for a token.
func (c *Client) GetSpread(tokenID string) (json.RawMessage, error) {
	q := url.Values{}
	q.Set("token_id", tokenID)
	return c.clobRequest("GET", PathSpread, q, nil, nil)
}

// GetSpreads batch spreads.
func (c *Client) GetSpreads(params []BookParams) (json.RawMessage, error) {
	b, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	return c.clobRequest("POST", PathSpreads, nil, nil, b)
}

// GetLastTradePrice returns last trade price for token.
func (c *Client) GetLastTradePrice(tokenID string) (json.RawMessage, error) {
	q := url.Values{}
	q.Set("token_id", tokenID)
	return c.clobRequest("GET", PathLastTradePrice, q, nil, nil)
}

// GetLastTradesPrices batch last trade prices.
func (c *Client) GetLastTradesPrices(params []BookParams) (json.RawMessage, error) {
	b, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	return c.clobRequest("POST", PathLastTradesPrices, nil, nil, b)
}

// GetPricesHistory returns historical prices for filters.
func (c *Client) GetPricesHistory(params PriceHistoryFilterParams) ([]MarketPrice, error) {
	q := url.Values{}
	if params.Market != "" {
		q.Set("market", params.Market)
	}
	if params.StartTs != 0 {
		q.Set("startTs", fmt.Sprintf("%d", params.StartTs))
	}
	if params.EndTs != 0 {
		q.Set("endTs", fmt.Sprintf("%d", params.EndTs))
	}
	if params.Fidelity != 0 {
		q.Set("fidelity", fmt.Sprintf("%d", params.Fidelity))
	}
	if params.Interval != "" {
		q.Set("interval", params.Interval)
	}
	body, err := c.clobRequest("GET", PathPricesHistory, q, nil, nil)
	if err != nil {
		return nil, err
	}
	var out []MarketPrice
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetMarketTradesEvents returns live activity events for a condition.
func (c *Client) GetMarketTradesEvents(conditionID string) (json.RawMessage, error) {
	return c.clobRequest("GET", PathLiveActivity+conditionID, nil, nil, nil)
}

// GetCurrentRewards paginates GET /rewards/markets/current.
func (c *Client) GetCurrentRewards(nextCursor string) (*PaginationPayload, error) {
	return c.getPaginated(PathRewardsMarketsCurrent, nextCursor)
}

// GetRawRewardsForMarket paginates /rewards/markets/{conditionId}.
func (c *Client) GetRawRewardsForMarket(conditionID, nextCursor string) (*PaginationPayload, error) {
	if nextCursor == "" {
		nextCursor = InitialCursor
	}
	q := url.Values{}
	q.Set("next_cursor", nextCursor)
	body, err := c.clobRequest("GET", PathRewardsMarketsPrefix+conditionID, q, nil, nil)
	if err != nil {
		return nil, err
	}
	var p PaginationPayload
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, err
	}
	return &p, nil
}
