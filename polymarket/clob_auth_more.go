package polymarket

import (
	"encoding/json"
)

// CreateReadonlyAPIKey creates a readonly API key (L2).
func (c *Client) CreateReadonlyAPIKey() (json.RawMessage, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	path := PathCreateReadonlyKey
	h, err := c.l2Headers("POST", path, "", false)
	if err != nil {
		return nil, err
	}
	return c.clobRequest("POST", path, nil, h, nil)
}

// GetReadonlyAPIKeys lists readonly keys (L2).
func (c *Client) GetReadonlyAPIKeys() ([]string, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	path := PathGetReadonlyKeys
	h, err := c.l2Headers("GET", path, "", false)
	if err != nil {
		return nil, err
	}
	data, err := c.clobRequest("GET", path, nil, h, nil)
	if err != nil {
		return nil, err
	}
	var keys []string
	if err := json.Unmarshal(data, &keys); err != nil {
		return nil, err
	}
	return keys, nil
}

// DeleteReadonlyAPIKey deletes a readonly key (L2, body: {"key": ...}).
func (c *Client) DeleteReadonlyAPIKey(key string) error {
	if err := c.requireL2(); err != nil {
		return err
	}
	payload := map[string]string{"key": key}
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	path := PathDeleteReadonlyKey
	h, err := c.l2Headers("DELETE", path, string(b), false)
	if err != nil {
		return err
	}
	_, err = c.clobRequest("DELETE", path, nil, h, b)
	return err
}

// CreateBuilderAPIKey creates a builder API key (L2).
func (c *Client) CreateBuilderAPIKey() (json.RawMessage, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	path := PathCreateBuilderKey
	h, err := c.l2Headers("POST", path, "", false)
	if err != nil {
		return nil, err
	}
	return c.clobRequest("POST", path, nil, h, nil)
}

// GetBuilderAPIKeys lists builder API keys (L2).
func (c *Client) GetBuilderAPIKeys() (json.RawMessage, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	path := PathGetBuilderKeys
	h, err := c.l2Headers("GET", path, "", false)
	if err != nil {
		return nil, err
	}
	return c.clobRequest("GET", path, nil, h, nil)
}

// RevokeBuilderAPIKey revokes builder API key (builder headers only).
func (c *Client) RevokeBuilderAPIKey() error {
	if err := c.requireBuilder(); err != nil {
		return err
	}
	path := PathRevokeBuilderKey
	h, err := c.builderHeadersOnly("DELETE", path, "")
	if err != nil {
		return err
	}
	_, err = c.clobRequest("DELETE", path, nil, h, nil)
	return err
}
