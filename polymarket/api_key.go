package polymarket

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// GetAPIKeys lists API keys for the authenticated user (L2).
func (c *Client) GetAPIKeys() ([]APIKeyCredentials, error) {
	if c.apiKeyCredentials == nil {
		return nil, fmt.Errorf("no API key credentials set")
	}
	path := PathGetAPIKeys
	headers, err := c.buildL2AuthHeaders("GET", path, "")
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
	return resp.APIKeys, nil
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

// DeleteAPIKey deletes an API key using explicit credentials (L2).
func (c *Client) DeleteAPIKey(apiKey, passphrase, secret string) error {
	if apiKey == "" || passphrase == "" || secret == "" {
		return fmt.Errorf("API key, passphrase, and secret are required")
	}
	path := PathDeleteAPIKey
	ts, err := c.authTimestampSeconds()
	if err != nil {
		return err
	}
	tsStr := strconv.FormatInt(ts, 10)
	sig, err := signL2("DELETE", path, "", tsStr, secret)
	if err != nil {
		return fmt.Errorf("l2 signature: %w", err)
	}
	headers := map[string]string{
		"POLY_ADDRESS":    c.fromAddress.Hex(),
		"POLY_SIGNATURE":  sig,
		"POLY_TIMESTAMP":  tsStr,
		"POLY_API_KEY":    apiKey,
		"POLY_PASSPHRASE": passphrase,
	}
	_, err = c.clobRequest("DELETE", path, nil, headers, nil)
	if err != nil {
		return fmt.Errorf("delete API key: %w", err)
	}
	return nil
}

func (c *Client) GetClosedOnlyModeStatus() (*ClosedOnlyModeStatus, error) {
	if c.apiKeyCredentials == nil {
		return nil, fmt.Errorf("no API key credentials set")
	}
	path := PathClosedOnly
	headers, err := c.buildL2AuthHeaders("GET", path, "")
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

func (c *Client) GetAccessStatus() (*AccessStatus, error) {
	if c.apiKeyCredentials == nil {
		return nil, fmt.Errorf("no API key credentials set")
	}
	path := PathAccessStatus
	headers, err := c.buildL2AuthHeaders("GET", path, "")
	if err != nil {
		return nil, err
	}
	body, err := c.clobRequest("GET", path, nil, headers, nil)
	if err != nil {
		return nil, fmt.Errorf("access status: %w", err)
	}
	var status AccessStatus
	if err := json.Unmarshal(body, &status); err != nil {
		return nil, fmt.Errorf("unmarshal access status: %w", err)
	}
	return &status, nil
}

// CreateAndSetAPIKey creates a new API key and sets it on the client.
func (c *Client) CreateAndSetAPIKey() (*APIKeyCredentials, error) {
	credentials, err := c.CreateAPIKey()
	if err != nil {
		return nil, err
	}
	c.SetAPIKeyCredentials(credentials)
	return credentials, nil
}

// DeleteAPIKeyWithCredentials deletes the current API key using stored credentials.
func (c *Client) DeleteAPIKeyWithCredentials() error {
	if c.apiKeyCredentials == nil {
		return fmt.Errorf("no API key credentials set")
	}
	err := c.DeleteAPIKey(c.apiKeyCredentials.ApiKey, c.apiKeyCredentials.Passphrase, c.apiKeyCredentials.Secret)
	if err != nil {
		return err
	}
	c.apiKeyCredentials = nil
	return nil
}
