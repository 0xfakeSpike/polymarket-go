package polymarket

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// GetEarningsForUserForDay walks pagination for /rewards/user.
func (c *Client) GetEarningsForUserForDay(date string) (json.RawMessage, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	path := PathRewardsUser
	h, err := c.l2Headers("GET", path, "", false)
	if err != nil {
		return nil, err
	}
	var combined []json.RawMessage
	next := InitialCursor
	for next != EndCursor {
		q := url.Values{}
		q.Set("date", date)
		q.Set("signature_type", fmt.Sprintf("%d", c.signatureType))
		q.Set("next_cursor", next)
		data, err := c.clobRequest("GET", path, q, h, nil)
		if err != nil {
			return nil, err
		}
		var page struct {
			Data       []json.RawMessage `json:"data"`
			NextCursor string            `json:"next_cursor"`
		}
		if err := json.Unmarshal(data, &page); err != nil {
			return nil, err
		}
		combined = append(combined, page.Data...)
		next = page.NextCursor
	}
	out, err := json.Marshal(combined)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(out), nil
}

// GetTotalEarningsForUserForDay returns GET /rewards/user/total for a date.
func (c *Client) GetTotalEarningsForUserForDay(date string) (json.RawMessage, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	path := PathRewardsUserTotal
	h, err := c.l2Headers("GET", path, "", false)
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("date", date)
	q.Set("signature_type", fmt.Sprintf("%d", c.signatureType))
	return c.clobRequest("GET", path, q, h, nil)
}

// GetUserEarningsAndMarketsConfig walks /rewards/user/markets.
func (c *Client) GetUserEarningsAndMarketsConfig(date, orderBy, position string, noCompetition bool) (json.RawMessage, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	path := PathRewardsUserMarkets
	h, err := c.l2Headers("GET", path, "", false)
	if err != nil {
		return nil, err
	}
	var combined []json.RawMessage
	next := InitialCursor
	for next != EndCursor {
		q := url.Values{}
		q.Set("date", date)
		q.Set("signature_type", fmt.Sprintf("%d", c.signatureType))
		q.Set("next_cursor", next)
		if orderBy != "" {
			q.Set("order_by", orderBy)
		}
		if position != "" {
			q.Set("position", position)
		}
		if noCompetition {
			q.Set("no_competition", "true")
		}
		data, err := c.clobRequest("GET", path, q, h, nil)
		if err != nil {
			return nil, err
		}
		var page struct {
			Data       []json.RawMessage `json:"data"`
			NextCursor string            `json:"next_cursor"`
		}
		if err := json.Unmarshal(data, &page); err != nil {
			return nil, err
		}
		combined = append(combined, page.Data...)
		next = page.NextCursor
	}
	out, err := json.Marshal(combined)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(out), nil
}

// GetRewardPercentages returns liquidity reward percentages for the user.
func (c *Client) GetRewardPercentages() (json.RawMessage, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	path := PathRewardsPercentages
	h, err := c.l2Headers("GET", path, "", false)
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("signature_type", fmt.Sprintf("%d", c.signatureType))
	return c.clobRequest("GET", path, q, h, nil)
}
