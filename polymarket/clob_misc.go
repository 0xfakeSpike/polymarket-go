package polymarket

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// GetNotifications returns user notifications (L2 + optional builder on same pattern as TS).
func (c *Client) GetNotifications() (json.RawMessage, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	path := PathNotifications
	h, err := c.l2Headers("GET", path, "")
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("signature_type", fmt.Sprintf("%d", c.signatureType))
	return c.clobRequest("GET", path, q, h, nil)
}

// DropNotifications deletes notifications; ids become comma-separated query param.
func (c *Client) DropNotifications(params *DropNotificationParams) error {
	if err := c.requireL2(); err != nil {
		return err
	}
	path := PathNotifications
	h, err := c.l2Headers("DELETE", path, "")
	if err != nil {
		return err
	}
	q := url.Values{}
	if params != nil && len(params.IDs) > 0 {
		q.Set("ids", strings.Join(params.IDs, ","))
	}
	_, err = c.clobRequest("DELETE", path, q, h, nil)
	return err
}

// GetBalanceAllowance returns balance and allowance for an asset type.
func (c *Client) GetBalanceAllowance(params *BalanceAllowanceParams) (BalanceAllowanceResponse, error) {
	if err := c.requireL2(); err != nil {
		return BalanceAllowanceResponse{}, err
	}
	path := PathBalanceAllowance
	h, err := c.l2Headers("GET", path, "")
	if err != nil {
		return BalanceAllowanceResponse{}, err
	}
	q := url.Values{}
	if params != nil {
		q.Set("asset_type", params.AssetType)
		if params.TokenID != "" {
			q.Set("token_id", params.TokenID)
		}
	}
	q.Set("signature_type", fmt.Sprintf("%d", c.signatureType))
	data, err := c.clobRequest("GET", path, q, h, nil)
	if err != nil {
		return BalanceAllowanceResponse{}, err
	}
	var out BalanceAllowanceResponse
	if err := json.Unmarshal(data, &out); err != nil {
		return BalanceAllowanceResponse{}, err
	}
	return out, nil
}

// UpdateBalanceAllowance triggers allowance refresh (GET with update path in TS).
func (c *Client) UpdateBalanceAllowance(params *BalanceAllowanceParams) error {
	if err := c.requireL2(); err != nil {
		return err
	}
	path := PathBalanceUpdate
	h, err := c.l2Headers("GET", path, "")
	if err != nil {
		return err
	}
	q := url.Values{}
	if params != nil {
		q.Set("asset_type", params.AssetType)
		if params.TokenID != "" {
			q.Set("token_id", params.TokenID)
		}
	}
	q.Set("signature_type", fmt.Sprintf("%d", c.signatureType))
	_, err = c.clobRequest("GET", path, q, h, nil)
	return err
}

// PostHeartbeat sends a heartbeat; pass empty heartbeatID to start a new chain.
func (c *Client) PostHeartbeat(heartbeatID *string) (*HeartbeatResponse, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	var id any
	if heartbeatID == nil {
		id = nil
	} else {
		id = *heartbeatID
	}
	payload := map[string]any{"heartbeat_id": id}
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	path := PathHeartbeats
	h, err := c.l2Headers("POST", path, string(b))
	if err != nil {
		return nil, err
	}
	data, err := c.clobRequest("POST", path, nil, h, b)
	if err != nil {
		return nil, err
	}
	var out HeartbeatResponse
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
