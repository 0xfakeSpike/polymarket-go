package polymarket

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// EventsKeysetParams are query parameters for Gamma GET /events/keyset.
// See https://docs.polymarket.com/api-reference/events/list-events-keyset-pagination
//
// Do not set Offset; the API rejects it. Use AfterCursor from the previous response's NextCursor.
type EventsKeysetParams struct {
	Limit       int    // 1–500; default 20 when unset or out of range
	AfterCursor string `json:"after_cursor,omitempty"`
	TagSlug     string `json:"tag_slug,omitempty"`
	Closed      *bool  `json:"closed,omitempty"`
	Live        *bool  `json:"live,omitempty"`
	Order       string `json:"order,omitempty"`
	Ascending   *bool  `json:"ascending,omitempty"`
}

// KeysetEventsResponse is the JSON body for GET /events/keyset (subset of fields).
type KeysetEventsResponse struct {
	Events     []GammaEvent `json:"events"`
	NextCursor string       `json:"next_cursor,omitempty"`
}

// GammaEvent is the subset of the Gamma Event model needed for market iteration.
type GammaEvent struct {
	Markets []map[string]any `json:"markets"`
}

// GetEventsKeyset lists events from the Gamma API using keyset pagination.
func (c *Client) GetEventsKeyset(params *EventsKeysetParams) (*KeysetEventsResponse, error) {
	if c == nil {
		return nil, fmt.Errorf("nil client")
	}
	q, err := eventsKeysetQueryValues(params)
	if err != nil {
		return nil, err
	}
	var out KeysetEventsResponse
	if err := c.gammaGET("/events/keyset", q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func eventsKeysetQueryValues(p *EventsKeysetParams) (url.Values, error) {
	q := url.Values{}
	limit := 20
	if p != nil && p.Limit > 0 {
		limit = p.Limit
	}
	if limit > 500 {
		limit = 500
	}
	if limit < 1 {
		limit = 1
	}
	q.Set("limit", strconv.Itoa(limit))

	if p == nil {
		return q, nil
	}
	if strings.TrimSpace(p.AfterCursor) != "" {
		q.Set("after_cursor", strings.TrimSpace(p.AfterCursor))
	}
	if strings.TrimSpace(p.TagSlug) != "" {
		q.Set("tag_slug", strings.TrimSpace(p.TagSlug))
	}
	if p.Closed != nil {
		q.Set("closed", strconv.FormatBool(*p.Closed))
	}
	if p.Live != nil {
		q.Set("live", strconv.FormatBool(*p.Live))
	}
	if strings.TrimSpace(p.Order) != "" {
		q.Set("order", strings.TrimSpace(p.Order))
	}
	if p.Ascending != nil {
		q.Set("ascending", strconv.FormatBool(*p.Ascending))
	}
	return q, nil
}
