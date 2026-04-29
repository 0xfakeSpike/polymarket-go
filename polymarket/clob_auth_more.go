package polymarket

import (
	"encoding/json"
)

// CreateReadonlyAPIKey creates a readonly API key (L2).
func (c *Client) CreateReadonlyAPIKey() (*ReadonlyAPIKeyResponse, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	path := PathCreateReadonlyKey
	h, err := c.l2Headers("POST", path, "")
	if err != nil {
		return nil, err
	}
	data, err := c.clobRequest("POST", path, nil, h, nil)
	if err != nil {
		return nil, err
	}
	var out ReadonlyAPIKeyResponse
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetReadonlyAPIKeys lists readonly keys (L2).
func (c *Client) GetReadonlyAPIKeys() ([]string, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	path := PathGetReadonlyKeys
	h, err := c.l2Headers("GET", path, "")
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
	h, err := c.l2Headers("DELETE", path, string(b))
	if err != nil {
		return err
	}
	_, err = c.clobRequest("DELETE", path, nil, h, b)
	return err
}

// CreateBuilderAPIKey creates a builder API key (L2).
func (c *Client) CreateBuilderAPIKey() (*BuilderAPIKey, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	path := PathCreateBuilderKey
	h, err := c.l2Headers("POST", path, "")
	if err != nil {
		return nil, err
	}
	data, err := c.clobRequest("POST", path, nil, h, nil)
	if err != nil {
		return nil, err
	}
	var out BuilderAPIKey
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetBuilderAPIKeys lists builder API keys (L2).
func (c *Client) GetBuilderAPIKeys() ([]BuilderAPIKeyResponse, error) {
	if err := c.requireL2(); err != nil {
		return nil, err
	}
	path := PathGetBuilderKeys
	h, err := c.l2Headers("GET", path, "")
	if err != nil {
		return nil, err
	}
	data, err := c.clobRequest("GET", path, nil, h, nil)
	if err != nil {
		return nil, err
	}
	var out []BuilderAPIKeyResponse
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// RevokeBuilderAPIKey revokes the current builder API key (L2).
func (c *Client) RevokeBuilderAPIKey() (json.RawMessage, error) {
	path := PathRevokeBuilderKey
	h, err := c.l2Headers("DELETE", path, "")
	if err != nil {
		return nil, err
	}
	return c.clobRequest("DELETE", path, nil, h, nil)
}
