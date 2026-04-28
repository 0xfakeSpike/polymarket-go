package polymarket

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// GetValidatedEventMarkets returns markets for an event and ensures each has at least two CLOB token IDs.
func (c *Client) GetValidatedEventMarkets(eventID string) ([]Market, error) {
	markets, err := c.GetEventMarkets(eventID)
	if err != nil {
		return nil, fmt.Errorf("get markets for event %s: %w", eventID, err)
	}
	if err := validateMarketsClobTokenIDs(markets); err != nil {
		return nil, err
	}
	return markets, nil
}

// ListMarkets queries the Polymarket /markets endpoint with optional filters.
func (c *Client) ListMarkets(params *MarketsParams) ([]Market, error) {
	body, err := c.GammaGET("/markets", params.Values())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch markets: %w", err)
	}
	var markets []Market
	if err := json.Unmarshal(body, &markets); err != nil {
		return nil, fmt.Errorf("failed to parse markets response: %w", err)
	}
	return markets, nil
}

// GetMarketsKeyset retrieves markets using the keyset pagination endpoint from gamma API.
func (c *Client) GetMarketsKeyset(params *MarketsKeysetParams) (*MarketsKeysetResponse, error) {
	query := params.Values()
	fullURL := c.baseURL + "/markets/keyset"
	if len(query) > 0 {
		fullURL += "?" + query.Encode()
	}
	log.Printf("polymarket GetMarketsKeyset url: %s", fullURL)

	body, err := c.GammaGET("/markets/keyset", query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch markets keyset: %w", err)
	}

	var resp MarketsKeysetResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse markets keyset response: %w", err)
	}
	return &resp, nil
}

// GetMarketDetail fetches a single market by ID.
func (c *Client) GetMarketDetail(marketID string) (*Market, error) {
	endpoint := fmt.Sprintf("/markets/%s", marketID)
	body, err := c.GammaGET(endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch market %s: %w", marketID, err)
	}
	var market Market
	if err := json.Unmarshal(body, &market); err != nil {
		return nil, fmt.Errorf("failed to parse market response: %w", err)
	}
	return &market, nil
}

// GetMarketBySlug retrieves a specific market by its slug
func (c *Client) GetMarketBySlug(slug string) (*Market, error) {
	markets, err := c.ListMarkets(&MarketsParams{Slug: slug, Limit: 1})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch market by slug %s: %w", slug, err)
	}
	if len(markets) == 0 {
		return nil, fmt.Errorf("market with slug %s not found", slug)
	}
	return &markets[0], nil
}

// GetFavoredSidePNLFromOrderBook fetches live books for market ClobTokenIds and computes
// PnL for the side whose best ask is above 0.5.
func (c *Client) GetFavoredSidePNLFromOrderBook(market *Market, now time.Time) (*FavoredSidePNL, error) {
	if c == nil {
		return nil, fmt.Errorf("nil client")
	}
	if market == nil {
		return nil, fmt.Errorf("nil market")
	}
	if len(market.ClobTokenIds) == 0 {
		return nil, fmt.Errorf("market %s: missing CLOB token IDs", market.ID)
	}

	books := make(map[string]*Book, len(market.ClobTokenIds))
	for _, tokenID := range market.ClobTokenIds {
		book, err := c.GetOrderBook(tokenID)
		if err != nil {
			return nil, fmt.Errorf("get order book token %s: %w", tokenID, err)
		}
		if len(book.Bids) == 0 || len(book.Asks) == 0 {
			return nil, fmt.Errorf("token %s has empty order book", tokenID)
		}
		books[tokenID] = book
	}

	return market.FavoredSidePNLFromOrderBooks(now, books)
}

// ParseOutcomes parses the outcomes JSON string and returns a slice of outcome names
func ParseOutcomes(outcomes string) ([]string, error) {
	if outcomes == "" {
		return []string{}, nil
	}
	var outcomeNames []string
	if err := json.Unmarshal([]byte(outcomes), &outcomeNames); err != nil {
		return nil, fmt.Errorf("failed to parse outcomes: %w", err)
	}
	return outcomeNames, nil
}

func validateMarketsClobTokenIDs(markets []Market) error {
	for _, m := range markets {
		if len(m.ClobTokenIds) < 2 {
			return fmt.Errorf("market %s: missing CLOB token IDs", m.ID)
		}
	}
	return nil
}
