package polymarket

import (
	"encoding/json"
	"fmt"
)

// GetComments retrieves a list of comments from the Polymarket API
func (c *Client) GetComments(params *CommentsParams) ([]Comment, error) {
	body, err := c.GammaGET("/comments", params.Values())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comments: %w", err)
	}

	var comments []Comment
	if err := json.Unmarshal(body, &comments); err != nil {
		return nil, fmt.Errorf("failed to parse comments response: %w", err)
	}

	return comments, nil
}

// GetMarketComments retrieves comments for a specific market
func (c *Client) GetMarketComments(marketID int, params *CommentsParams) ([]Comment, error) {
	if params == nil {
		params = &CommentsParams{}
	}

	// Set market-specific filters
	params.ParentEntityType = "market"
	params.ParentEntityID = &marketID

	return c.GetComments(params)
}

// GetEventComments retrieves comments for a specific event
func (c *Client) GetEventComments(eventID int, params *CommentsParams) ([]Comment, error) {
	if params == nil {
		params = &CommentsParams{}
	}

	// Set event-specific filters
	params.ParentEntityType = "Event"
	params.ParentEntityID = &eventID

	return c.GetComments(params)
}

// GetSeriesComments retrieves comments for a specific series
func (c *Client) GetSeriesComments(seriesID int, params *CommentsParams) ([]Comment, error) {
	if params == nil {
		params = &CommentsParams{}
	}

	// Set series-specific filters
	params.ParentEntityType = "Series"
	params.ParentEntityID = &seriesID

	return c.GetComments(params)
}
