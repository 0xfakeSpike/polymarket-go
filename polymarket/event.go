package polymarket

import (
	"encoding/json"
	"fmt"
	"log"
)

// GetEvents retrieves a list of events from the Polymarket API.
func (c *Client) GetEvents(params *EventsParams) ([]Event, error) {
	body, err := c.GammaGET("/events", params.Values())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch events: %w", err)
	}

	var events []Event
	if err := json.Unmarshal(body, &events); err != nil {
		return nil, fmt.Errorf("failed to parse events response: %w", err)
	}

	return events, nil
}

// GetEvent retrieves a specific event by its ID with optional parameters.
func (c *Client) GetEvent(eventID string, params *GetEventParams) (*Event, error) {
	endpoint := fmt.Sprintf("/events/%s", eventID)

	body, err := c.GammaGET(endpoint, params.Values())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch event %s: %w", eventID, err)
	}

	var event Event
	if err := json.Unmarshal(body, &event); err != nil {
		return nil, fmt.Errorf("failed to parse event response: %w", err)
	}

	return &event, nil
}

// GetEventBySlug retrieves a specific event by its slug
func (c *Client) GetEventBySlug(slug string) ([]Market, error) {
	params := &EventsParams{
		Slug:  []string{slug},
		Limit: 1,
	}

	events, err := c.GetEvents(params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch event by slug %s: %w", slug, err)
	}

	if len(events) == 0 {
		return nil, fmt.Errorf("event with slug %s not found", slug)
	}
	event := events[0]
	activeMarkets := []Market{}
	for _, market := range event.Markets {
		if !market.Closed {
			activeMarkets = append(activeMarkets, market)
		}
	}
	if len(activeMarkets) == 0 {
		return nil, fmt.Errorf("event has no active markets")
	}

	return activeMarkets, nil
}

// GetEventMarkets retrieves all markets for a specific event
func (c *Client) GetEventMarkets(eventID string) ([]Market, error) {
	params := &MarketsParams{
		EventID: eventID,
	}

	return c.ListMarkets(params)
}

// GetBiggestMovers retrieves the biggest movers from Polymarket
func (c *Client) GetBiggestMovers() ([]BiggestMover, error) {
	body, err := c.GammaGET("/api/biggest-movers", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch biggest movers: %w", err)
	}

	var movers []BiggestMover
	if err := json.Unmarshal(body, &movers); err != nil {
		return nil, fmt.Errorf("failed to parse biggest movers response: %w", err)
	}

	return movers, nil
}

// GetEventsPagination retrieves events using the pagination endpoint from gamma API
func (c *Client) GetEventsPagination(params *EventsPaginationParams) ([]Event, error) {
	body, err := c.GammaGET("/events/pagination", params.Values())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch events pagination: %w", err)
	}

	var eventPaginationResponse EventPaginationResponse
	if err := json.Unmarshal(body, &eventPaginationResponse); err != nil {
		return nil, fmt.Errorf("failed to parse events pagination response: %w", err)
	}

	return eventPaginationResponse.Events, nil
}

// GetEventsKeyset retrieves events using cursor-based keyset pagination (GET /events/keyset).
// Pass next_cursor from the response as after_cursor on the next page.
func (c *Client) GetEventsKeyset(params *EventsKeysetParams) (*EventsKeysetResponse, error) {
	if params == nil {
		params = &EventsKeysetParams{}
	}
	query := params.Values()
	fullURL := c.baseURL + "/events/keyset"
	if len(query) > 0 {
		fullURL += "?" + query.Encode()
	}
	log.Printf("polymarket GetEventsKeyset url: %s", fullURL)

	body, err := c.GammaGET("/events/keyset", query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch events keyset: %w", err)
	}

	var resp EventsKeysetResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse events keyset response: %w", err)
	}
	return &resp, nil
}
