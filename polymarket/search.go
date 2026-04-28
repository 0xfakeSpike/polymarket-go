package polymarket

import (
	"encoding/json"
	"fmt"

	"github.com/0xfakespike/everything/utils"
)

// https://gamma-api.polymarket.com/search-v2?
// // events_status=active&limit_per_type=20&page=1&q=Match+Winner&sort=volume_24hr&type=events

// https://gamma-api.polymarket.com/search-v2?
// // q=Counter-Strike%3A+SE7ENS+Esport+vs+JUMBO+TEAM+%28BO1%29+-+ESEA+Advanced+Europe+Regular+Season&optimized=true&limit_per_type=6&type=events&search_tags=true&search_profiles=true&cache=true

func (c *Client) Search(params *SearchParams) (*SearchResults, error) {
	if params == nil || params.Q == "" {
		return nil, fmt.Errorf("search query (Q) is required")
	}

	body, err := c.GammaGET("/search-v2", params.Values())
	if err != nil {
		return nil, fmt.Errorf("failed to perform search: %w", err)
	}

	var results SearchResults
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("failed to parse search results: %w", err)
	}

	return &results, nil
}

// SearchEventsWithQuery searches events with the same core query shape as the public site
// (type=events, active, volume_24hr, limit_per_type=20). Use Search with a custom SearchParams for tags/presets.
func (c *Client) SearchEventsWithQuery(query string) ([]Event, error) {
	results, err := c.Search(&SearchParams{
		Q:            query,
		Type:         "events",
		Page:         1,
		LimitPerType: 20,
		EventsStatus: "active",
		Sort:         "volume_24hr",
	})
	if err != nil {
		return nil, err
	}

	return results.Events, nil
}

// SearchProfiles searches specifically for user profiles
func (c *Client) SearchProfiles(query string, params *SearchParams) ([]UserProfile, error) {
	if params == nil {
		params = &SearchParams{}
	}
	params.Q = query
	params.SearchTags = utils.BoolPtr(false)
	params.SearchProfiles = utils.BoolPtr(true)

	results, err := c.Search(params)
	if err != nil {
		return nil, err
	}

	return results.Profiles, nil
}

// SearchTags searches specifically for tags
func (c *Client) SearchTags(query string, params *SearchParams) ([]Tag, error) {
	if params == nil {
		params = &SearchParams{}
	}
	params.Q = query
	params.SearchTags = utils.BoolPtr(true)
	params.SearchProfiles = utils.BoolPtr(false)

	results, err := c.Search(params)
	if err != nil {
		return nil, err
	}

	return results.Tags, nil
}

// SearchByTag searches for events by specific tags
func (c *Client) SearchByTag(query string, tags []string, params *SearchParams) ([]Event, error) {
	if params == nil {
		params = &SearchParams{}
	}
	params.Q = query
	params.EventsTag = tags

	results, err := c.Search(params)
	if err != nil {
		return nil, err
	}

	return results.Events, nil
}
