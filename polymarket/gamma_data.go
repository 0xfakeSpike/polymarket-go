package polymarket

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// hostHTTP performs one HTTP request against a full API host (Gamma or Data API).
// path must start with "/" (e.g. "/markets"). query may be nil.
func (c *Client) hostHTTP(baseURL, method, path string, query url.Values, body []byte) ([]byte, error) {
	if path == "" || path[0] != '/' {
		return nil, fmt.Errorf("path must start with /, got %q", path)
	}
	u := baseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}
	var r io.Reader
	if body != nil {
		r = bytes.NewReader(body)
	}
	req, err := http.NewRequest(method, u, r)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", clobUserAgent)
	req.Header.Set("Accept", "*/*")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	out, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%s %s: http %d: %s", method, path, resp.StatusCode, string(out))
	}
	return out, nil
}

// GammaGET requests the Polymarket Gamma API (see DefaultBaseURL / WithGammaBaseURL).
func (c *Client) GammaGET(path string, query url.Values) ([]byte, error) {
	return c.hostHTTP(c.baseURL, http.MethodGet, path, query, nil)
}

// GammaPOST sends JSON to the Gamma API.
func (c *Client) GammaPOST(path string, query url.Values, body []byte) ([]byte, error) {
	return c.hostHTTP(c.baseURL, http.MethodPost, path, query, body)
}

// DataGET requests the Polymarket Data API (see DataAPIBaseURL / WithDataAPIBaseURL).
func (c *Client) DataGET(path string, query url.Values) ([]byte, error) {
	return c.hostHTTP(c.dataAPIBaseURL, http.MethodGet, path, query, nil)
}

// DataPOST sends JSON to the Data API.
func (c *Client) DataPOST(path string, query url.Values, body []byte) ([]byte, error) {
	return c.hostHTTP(c.dataAPIBaseURL, http.MethodPost, path, query, body)
}
