package polymarket

import (
	"encoding/json"
	"fmt"
	"strings"
)

// GetAPIKeys lists API keys for the authenticated user (L2).
func (c *Client) GetAPIKeys() (*APIKeysResponse, error) {
	path := PathGetAPIKeys
	headers, err := c.l2Headers("GET", path, "")
	if err != nil {
		return nil, err
	}
	body, err := c.clobRequest("GET", path, nil, headers, nil)
	if err != nil {
		return nil, fmt.Errorf("get API keys: %w", err)
	}
	var resp APIKeysResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal API keys: %w", err)
	}
	return &resp, nil
}

// CreateAPIKey creates a new API key (L1).
func (c *Client) CreateAPIKey() (*APIKeyCredentials, error) {
	if err := c.requireL1(); err != nil {
		return nil, err
	}
	headers, err := c.buildL1AuthHeaders(0)
	if err != nil {
		return nil, fmt.Errorf("l1 headers: %w", err)
	}
	path := PathCreateAPIKey
	body, err := c.clobRequest("POST", path, nil, headers, nil)
	if err != nil {
		return nil, fmt.Errorf("create API key: %w, body: %s", err, string(body))
	}
	var credentials APIKeyCredentials
	if err := json.Unmarshal(body, &credentials); err != nil {
		return nil, fmt.Errorf("unmarshal API key: %w", err)
	}
	return &credentials, nil
}

// DeriveAPIKey derives an existing API key (L1).
func (c *Client) DeriveAPIKey(nonce int64) (*APIKeyCredentials, error) {
	if err := c.requireL1(); err != nil {
		return nil, err
	}
	headers, err := c.buildL1AuthHeaders(nonce)
	if err != nil {
		return nil, fmt.Errorf("l1 headers: %w", err)
	}
	path := PathDeriveAPIKey
	body, err := c.clobRequest("GET", path, nil, headers, nil)
	if err != nil {
		return nil, fmt.Errorf("derive API key: %w", err)
	}
	var credentials APIKeyCredentials
	if err := json.Unmarshal(body, &credentials); err != nil {
		return nil, fmt.Errorf("unmarshal API key: %w", err)
	}
	return &credentials, nil
}

// CreateOrDeriveAPIKey matches @polymarket/clob-client createOrDeriveApiKey:
// create first; if the response has no api key, derive.
//
// The CLOB also returns HTTP 400 {"error":"Could not create api key"} when a key
// already exists for this wallet — in that case we fall back to DeriveAPIKey.
func (c *Client) CreateOrDeriveAPIKey() (*APIKeyCredentials, error) {
	creds, err := c.CreateAPIKey()
	if err == nil {
		if creds != nil && creds.ApiKey != "" {
			return creds, nil
		}
		return c.DeriveAPIKey(0)
	}
	if shouldDeriveAfterCreateAPIKeyFailure(err) {
		derived, derr := c.DeriveAPIKey(0)
		if derr != nil {
			return nil, fmt.Errorf("create API key: %w; derive API key: %w", err, derr)
		}
		return derived, nil
	}
	return nil, err
}

func shouldDeriveAfterCreateAPIKeyFailure(err error) bool {
	if err == nil {
		return false
	}
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "could not create api key")
}

// DeleteAPIKey deletes the current L2 API key.
func (c *Client) DeleteAPIKey() (json.RawMessage, error) {
	path := PathDeleteAPIKey
	headers, err := c.l2Headers("DELETE", path, "")
	if err != nil {
		return nil, err
	}
	body, err := c.clobRequest("DELETE", path, nil, headers, nil)
	if err != nil {
		return nil, fmt.Errorf("delete API key: %w", err)
	}
	return body, nil
}

func (c *Client) GetClosedOnlyMode() (*ClosedOnlyModeStatus, error) {
	path := PathClosedOnly
	headers, err := c.l2Headers("GET", path, "")
	if err != nil {
		return nil, err
	}
	body, err := c.clobRequest("GET", path, nil, headers, nil)
	if err != nil {
		return nil, fmt.Errorf("closed-only status: %w", err)
	}
	var status ClosedOnlyModeStatus
	if err := json.Unmarshal(body, &status); err != nil {
		return nil, fmt.Errorf("unmarshal closed-only: %w", err)
	}
	return &status, nil
}
