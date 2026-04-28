package polymarket

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// GetLiveVolume retrieves live volume data for an event
func (c *Client) GetLiveVolume(eventID int) (*LiveVolume, error) {
	if eventID < 1 {
		return nil, fmt.Errorf("event ID must be >= 1, got %d", eventID)
	}
	
	// Build query parameters
	params := url.Values{}
	params.Add("id", strconv.Itoa(eventID))
	
	// Make request to the data API
	body, err := c.DataGET("/live-volume", params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch live volume for event %d: %w", eventID, err)
	}
	
	// The API returns an array, but we expect only one element for a single event
	var volumes []LiveVolume
	if err := json.Unmarshal(body, &volumes); err != nil {
		return nil, fmt.Errorf("failed to parse live volume response: %w", err)
	}
	
	if len(volumes) == 0 {
		return nil, fmt.Errorf("no volume data found for event %d", eventID)
	}
	
	return &volumes[0], nil
}

// GetLiveVolumeMultiple retrieves live volume data for multiple events
func (c *Client) GetLiveVolumeMultiple(eventIDs []int) ([]LiveVolume, error) {
	if len(eventIDs) == 0 {
		return nil, fmt.Errorf("at least one event ID is required")
	}
	
	// Build query parameters for multiple IDs
	params := url.Values{}
	for _, id := range eventIDs {
		if id < 1 {
			return nil, fmt.Errorf("all event IDs must be >= 1, got %d", id)
		}
		params.Add("id", strconv.Itoa(id))
	}
	
	// Make request to the data API
	body, err := c.DataGET("/live-volume", params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch live volume for events: %w", err)
	}
	
	var volumes []LiveVolume
	if err := json.Unmarshal(body, &volumes); err != nil {
		return nil, fmt.Errorf("failed to parse live volume response: %w", err)
	}
	
	return volumes, nil
}

// GetEventTotalVolume returns just the total volume for an event (convenience method)
func (c *Client) GetEventTotalVolume(eventID int) (float64, error) {
	volume, err := c.GetLiveVolume(eventID)
	if err != nil {
		return 0, err
	}
	
	return volume.Total, nil
}

// GetEventMarketVolumes returns volume data for all markets in an event (convenience method)
func (c *Client) GetEventMarketVolumes(eventID int) ([]MarketVolume, error) {
	volume, err := c.GetLiveVolume(eventID)
	if err != nil {
		return nil, err
	}
	
	return volume.Markets, nil
}