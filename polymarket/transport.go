package polymarket

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const clobUserAgent = "@polymarket/clob-client"

// clobRequest executes a CLOB HTTP call: path is e.g. "/order" (no host). Query is merged with geo_block_token.
func (c *Client) clobRequest(method, path string, query url.Values, headers map[string]string, body []byte) ([]byte, error) {
	data, _, err := c.clobRequestStatus(method, path, query, headers, body)
	return data, err
}

func (c *Client) clobRequestStatus(method, path string, query url.Values, headers map[string]string, body []byte) ([]byte, int, error) {
	u := c.clobHost + path
	q := c.mergeGeoQuery(query)
	if len(q) > 0 {
		u += "?" + q.Encode()
	}

	do := func() (*http.Response, []byte, error) {
		var r io.Reader
		if body != nil {
			r = bytes.NewReader(body)
		}
		req, err := http.NewRequest(method, u, r)
		if err != nil {
			return nil, nil, err
		}
		req.Header.Set("User-Agent", clobUserAgent)
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Connection", "keep-alive")
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, nil, err
		}
		b, err := readHTTPBody(resp)
		if err != nil {
			return resp, nil, err
		}
		return resp, b, nil
	}

	resp, data, err := do()
	if err != nil && c.retryPostOnError && method == http.MethodPost {
		time.Sleep(30 * time.Millisecond)
		resp, data, err = do()
	}
	if err != nil {
		return nil, 0, err
	}
	status := resp.StatusCode
	if status < 200 || status >= 300 {
		return data, status, fmt.Errorf("http %d: %s", status, string(data))
	}
	if c.throwOnError {
		if err := maybeThrowAPIError(status, data); err != nil {
			return data, status, err
		}
	}
	return data, status, nil
}

func maybeThrowAPIError(status int, data []byte) error {
	if len(data) == 0 {
		return nil
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(data, &m); err != nil {
		return nil
	}
	raw, ok := m["error"]
	if !ok || len(raw) == 0 || string(raw) == "null" {
		return nil
	}
	var msg string
	if err := json.Unmarshal(raw, &msg); err != nil {
		return &ApiError{Message: string(raw), Status: status, Data: json.RawMessage(data)}
	}
	if msg == "" {
		return nil
	}
	return &ApiError{Message: msg, Status: status, Data: json.RawMessage(data)}
}

func readHTTPBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	reader := io.Reader(resp.Body)
	contentEncoding := strings.ToLower(strings.TrimSpace(resp.Header.Get("Content-Encoding")))
	if strings.Contains(contentEncoding, "gzip") {
		gz, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("gzip decode response: %w", err)
		}
		defer gz.Close()
		reader = gz
	}

	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	return b, nil
}

// mergeGeoQuery appends geo_block_token when configured.
func (c *Client) mergeGeoQuery(q url.Values) url.Values {
	if q == nil {
		q = url.Values{}
	}
	if c.geoBlockToken != "" {
		q.Set("geo_block_token", c.geoBlockToken)
	}
	return q
}
