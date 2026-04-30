package polymarket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// GetPublicProfile returns Gamma GET /public-profile for a wallet (proxy or EOA).
// Docs: https://docs.polymarket.com/api-reference/profiles/get-public-profile-by-wallet-address
func (c *Client) GetPublicProfile(address string) (*PublicProfileResponse, error) {
	address = strings.TrimSpace(address)
	if address == "" {
		return nil, fmt.Errorf("address is required")
	}
	q := url.Values{}
	q.Set("address", address)
	var out PublicProfileResponse
	if err := c.gammaGET("/public-profile", q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetCurrentPositions returns Data API GET /positions for a user profile address.
// Docs: https://docs.polymarket.com/api-reference/core/get-current-positions-for-a-user
func (c *Client) GetCurrentPositions(params *CurrentPositionsParams) ([]Position, error) {
	q, err := currentPositionsQuery(params)
	if err != nil {
		return nil, err
	}
	var out []Position
	if err := c.dataGET("/positions", q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetClosedPositions returns Data API GET /closed-positions.
// Docs: https://docs.polymarket.com/api-reference/core/get-closed-positions-for-a-user
func (c *Client) GetClosedPositions(params *ClosedPositionsParams) ([]ClosedPosition, error) {
	q, err := closedPositionsQuery(params)
	if err != nil {
		return nil, err
	}
	var out []ClosedPosition
	if err := c.dataGET("/closed-positions", q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetUserActivity returns Data API GET /activity.
// Docs: https://docs.polymarket.com/api-reference/core/get-user-activity
func (c *Client) GetUserActivity(params *UserActivityParams) ([]Activity, error) {
	q, err := userActivityQuery(params)
	if err != nil {
		return nil, err
	}
	var out []Activity
	if err := c.dataGET("/activity", q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) gammaGET(path string, q url.Values, dest any) error {
	return c.jsonGET(c.gammaHost, path, q, dest)
}

func (c *Client) dataGET(path string, q url.Values, dest any) error {
	return c.jsonGET(c.dataHost, path, q, dest)
}

func (c *Client) jsonGET(baseURL, path string, q url.Values, dest any) error {
	u := strings.TrimSuffix(baseURL, "/") + path
	if len(q) > 0 {
		u += "?" + q.Encode()
	}
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", clobUserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	body, err := readHTTPBody(resp)
	if err != nil {
		return err
	}
	status := resp.StatusCode
	if status < 200 || status >= 300 {
		return fmt.Errorf("http %d: %s", status, string(body))
	}
	if c.throwOnError {
		if err := maybeThrowAPIError(status, body); err != nil {
			return err
		}
	}
	if dest != nil {
		if err := json.Unmarshal(body, dest); err != nil {
			return fmt.Errorf("decode json: %w", err)
		}
	}
	return nil
}

func currentPositionsQuery(p *CurrentPositionsParams) (url.Values, error) {
	if p == nil || strings.TrimSpace(p.User) == "" {
		return nil, fmt.Errorf("user is required")
	}
	q := url.Values{}
	q.Set("user", strings.TrimSpace(p.User))
	encodeCSV(q, "market", p.Market)
	encodeIntCSV(q, "eventId", p.EventID)
	if p.SizeThreshold != nil {
		q.Set("sizeThreshold", strconv.FormatFloat(*p.SizeThreshold, 'f', -1, 64))
	}
	if p.Redeemable != nil {
		q.Set("redeemable", strconv.FormatBool(*p.Redeemable))
	}
	if p.Mergeable != nil {
		q.Set("mergeable", strconv.FormatBool(*p.Mergeable))
	}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	if p.Offset != nil {
		q.Set("offset", strconv.Itoa(*p.Offset))
	}
	if p.SortBy != "" {
		q.Set("sortBy", p.SortBy)
	}
	if p.SortDirection != "" {
		q.Set("sortDirection", p.SortDirection)
	}
	if p.Title != "" {
		q.Set("title", p.Title)
	}
	return q, nil
}

func closedPositionsQuery(p *ClosedPositionsParams) (url.Values, error) {
	if p == nil || strings.TrimSpace(p.User) == "" {
		return nil, fmt.Errorf("user is required")
	}
	q := url.Values{}
	q.Set("user", strings.TrimSpace(p.User))
	encodeCSV(q, "market", p.Market)
	if p.Title != "" {
		q.Set("title", p.Title)
	}
	encodeIntCSV(q, "eventId", p.EventID)
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	if p.Offset != nil {
		q.Set("offset", strconv.Itoa(*p.Offset))
	}
	if p.SortBy != "" {
		q.Set("sortBy", p.SortBy)
	}
	if p.SortDirection != "" {
		q.Set("sortDirection", p.SortDirection)
	}
	return q, nil
}

func userActivityQuery(p *UserActivityParams) (url.Values, error) {
	if p == nil || strings.TrimSpace(p.User) == "" {
		return nil, fmt.Errorf("user is required")
	}
	q := url.Values{}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	if p.Offset != nil {
		q.Set("offset", strconv.Itoa(*p.Offset))
	}
	q.Set("user", strings.TrimSpace(p.User))
	encodeCSV(q, "market", p.Market)
	encodeIntCSV(q, "eventId", p.EventID)
	if len(p.ActivityTypes) > 0 {
		q.Set("type", strings.Join(p.ActivityTypes, ","))
	}
	if p.Start != nil {
		q.Set("start", strconv.FormatInt(*p.Start, 10))
	}
	if p.End != nil {
		q.Set("end", strconv.FormatInt(*p.End, 10))
	}
	if p.SortBy != "" {
		q.Set("sortBy", p.SortBy)
	}
	if p.SortDirection != "" {
		q.Set("sortDirection", p.SortDirection)
	}
	if p.Side != "" {
		q.Set("side", p.Side)
	}
	return q, nil
}

func encodeCSV(q url.Values, key string, vals []string) {
	if len(vals) == 0 {
		return
	}
	q.Set(key, strings.Join(vals, ","))
}

func encodeIntCSV(q url.Values, key string, vals []int) {
	if len(vals) == 0 {
		return
	}
	parts := make([]string, len(vals))
	for i, v := range vals {
		parts[i] = strconv.Itoa(v)
	}
	q.Set(key, strings.Join(parts, ","))
}
