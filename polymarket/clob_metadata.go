package polymarket

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// GetServerTime returns CLOB server time in unix seconds (GET /time).
func (c *Client) GetServerTime() (int64, error) {
	body, err := c.clobRequest("GET", PathTime, nil, nil, nil)
	if err != nil {
		return 0, err
	}
	var n json.Number
	if err := json.Unmarshal(body, &n); err == nil {
		i, err := n.Int64()
		if err != nil {
			return 0, fmt.Errorf("server time: %w", err)
		}
		return i, nil
	}
	var f float64
	if err := json.Unmarshal(body, &f); err != nil {
		return 0, fmt.Errorf("decode server time: %w", err)
	}
	return int64(f), nil
}

func (c *Client) authTimestampSeconds() (int64, error) {
	if c.useServerTime {
		return c.GetServerTime()
	}
	return time.Now().Unix(), nil
}

type tickSizeAPIResponse struct {
	MinimumTickSize json.RawMessage `json:"minimum_tick_size"`
	Error           string          `json:"error"`
}

func parseTickSizeWire(raw json.RawMessage) (string, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return "", fmt.Errorf("missing minimum_tick_size")
	}
	if raw[0] == '"' {
		var s string
		if err := json.Unmarshal(raw, &s); err != nil {
			return "", err
		}
		return s, nil
	}
	var f float64
	if err := json.Unmarshal(raw, &f); err != nil {
		return "", err
	}
	return strconv.FormatFloat(f, 'f', -1, 64), nil
}

// GetTickSize returns minimum tick size for tokenID (cached).
func (c *Client) GetTickSize(tokenID string) (string, error) {
	c.ensureMetadataCaches()
	if tokenID == "" {
		return "", fmt.Errorf("tokenID required")
	}
	if ts, ok := c.tickSizes[tokenID]; ok {
		if at, ok2 := c.tickSizeAt[tokenID]; ok2 && time.Since(at) < c.tickSizeTTL {
			return ts, nil
		}
	}
	if conditionID, ok := c.tokenConditionMap[tokenID]; ok {
		if _, err := c.GetClobMarketInfo(conditionID); err != nil {
			return "", err
		}
		if ts, ok := c.tickSizes[tokenID]; ok {
			return ts, nil
		}
	}
	q := url.Values{}
	q.Set("token_id", tokenID)
	body, err := c.clobRequest("GET", PathTickSize, q, nil, nil)
	if err != nil {
		return "", err
	}
	var r tickSizeAPIResponse
	if err := json.Unmarshal(body, &r); err != nil {
		return "", fmt.Errorf("tick-size json: %w", err)
	}
	if r.Error != "" {
		return "", fmt.Errorf("%s", r.Error)
	}
	tick, err := parseTickSizeWire(r.MinimumTickSize)
	if err != nil {
		return "", err
	}
	if tick == "" {
		return "", fmt.Errorf("empty minimum_tick_size")
	}
	c.tickSizes[tokenID] = tick
	c.tickSizeAt[tokenID] = time.Now()
	return tick, nil
}

// ClearTickSizeCache drops tick size cache for one token or all when tokenID is empty.
func (c *Client) ClearTickSizeCache(tokenID string) {
	if tokenID != "" {
		delete(c.tickSizes, tokenID)
		delete(c.tickSizeAt, tokenID)
		return
	}
	clear(c.tickSizes)
	clear(c.tickSizeAt)
}

// parseNegRiskCLOBResponse reads GET /neg-risk JSON. Wrong neg_risk picks the wrong EIP-712
// verifyingContract (CTF vs neg-risk exchange), which surfaces as "invalid signature" on POST /order.
func parseNegRiskCLOBResponse(body []byte) (negRisk bool, apiErr string, err error) {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(body, &m); err != nil {
		return false, "", fmt.Errorf("neg-risk json: %w", err)
	}
	if raw, ok := m["error"]; ok && string(raw) != "null" && len(raw) > 0 {
		var es string
		_ = json.Unmarshal(raw, &es)
		if es != "" {
			return false, es, nil
		}
	}
	for _, key := range []string{"neg_risk", "negRisk"} {
		raw, ok := m[key]
		if !ok {
			continue
		}
		var b bool
		if err := json.Unmarshal(raw, &b); err == nil {
			return b, "", nil
		}
		var f float64
		if err := json.Unmarshal(raw, &f); err == nil {
			return f != 0, "", nil
		}
		var s string
		if err := json.Unmarshal(raw, &s); err == nil {
			switch s {
			case "true", "1":
				return true, "", nil
			case "false", "0":
				return false, "", nil
			}
		}
	}
	return false, "", nil
}

// GetNegRisk returns whether the market for tokenID uses the neg-risk exchange contract (cached).
func (c *Client) GetNegRisk(tokenID string) (bool, error) {
	c.ensureMetadataCaches()
	if v, ok := c.negRiskCache[tokenID]; ok {
		return v, nil
	}
	if conditionID, ok := c.tokenConditionMap[tokenID]; ok {
		if _, err := c.GetClobMarketInfo(conditionID); err != nil {
			return false, err
		}
		if v, ok := c.negRiskCache[tokenID]; ok {
			return v, nil
		}
	}
	q := url.Values{}
	q.Set("token_id", tokenID)
	body, err := c.clobRequest("GET", PathNegRisk, q, nil, nil)
	if err != nil {
		return false, err
	}
	neg, apiErr, err := parseNegRiskCLOBResponse(body)
	if err != nil {
		return false, err
	}
	if apiErr != "" {
		return false, fmt.Errorf("%s", apiErr)
	}
	c.negRiskCache[tokenID] = neg
	return neg, nil
}

// parseFeeRateCLOBResponse reads GET /fee-rate JSON. Fee is part of the signed Order struct; a
// mismatch vs the CLOB's expectation yields invalid signature.
func parseFeeRateCLOBResponse(body []byte) (bps int, apiErr string, err error) {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(body, &m); err != nil {
		return 0, "", fmt.Errorf("fee-rate json: %w", err)
	}
	if raw, ok := m["error"]; ok && string(raw) != "null" && len(raw) > 0 {
		var es string
		_ = json.Unmarshal(raw, &es)
		if es != "" {
			return 0, es, nil
		}
	}
	for _, key := range []string{"base_fee", "baseFee"} {
		raw, ok := m[key]
		if !ok {
			continue
		}
		var n int
		if err := json.Unmarshal(raw, &n); err == nil {
			return n, "", nil
		}
		var f float64
		if err := json.Unmarshal(raw, &f); err == nil {
			return int(f), "", nil
		}
		var s string
		if err := json.Unmarshal(raw, &s); err == nil {
			if v, err := strconv.Atoi(s); err == nil {
				return v, "", nil
			}
		}
	}
	return 0, "", nil
}

// GetFeeRateBps returns base fee rate in basis points for tokenID (cached).
func (c *Client) GetFeeRateBps(tokenID string) (int, error) {
	c.ensureMetadataCaches()
	if v, ok := c.feeRateCache[tokenID]; ok {
		return v, nil
	}
	q := url.Values{}
	q.Set("token_id", tokenID)
	body, err := c.clobRequest("GET", PathFeeRate, q, nil, nil)
	if err != nil {
		return 0, err
	}
	bps, apiErr, err := parseFeeRateCLOBResponse(body)
	if err != nil {
		return 0, err
	}
	if apiErr != "" {
		return 0, fmt.Errorf("%s", apiErr)
	}
	c.feeRateCache[tokenID] = bps
	return bps, nil
}

// GetFeeExponent returns fee exponent from CLOB market info.
func (c *Client) GetFeeExponent(tokenID string) (float64, error) {
	c.ensureMetadataCaches()
	if info, ok := c.feeInfoCache[tokenID]; ok {
		return info.Exponent, nil
	}
	if err := c.ensureMarketInfoCached(tokenID); err != nil {
		return 0, err
	}
	return c.feeInfoCache[tokenID].Exponent, nil
}

func (c *Client) ensureMarketInfoCached(tokenID string) error {
	c.ensureMetadataCaches()
	if _, ok := c.feeInfoCache[tokenID]; ok {
		return nil
	}
	conditionID := c.tokenConditionMap[tokenID]
	if conditionID == "" {
		body, err := c.clobRequest("GET", PathMarketByTokenPrefix+tokenID, nil, nil, nil)
		if err != nil {
			return err
		}
		var r struct {
			ConditionID string `json:"condition_id"`
		}
		if err := json.Unmarshal(body, &r); err != nil {
			return err
		}
		if r.ConditionID == "" {
			return fmt.Errorf("failed to resolve condition id for token %s", tokenID)
		}
		conditionID = r.ConditionID
		c.tokenConditionMap[tokenID] = conditionID
	}
	_, err := c.GetClobMarketInfo(conditionID)
	return err
}

func priceValid(price float64, tickSize string) (bool, error) {
	minTick, err := strconv.ParseFloat(tickSize, 64)
	if err != nil {
		return false, err
	}
	return price >= minTick && price <= 1-minTick, nil
}

func (c *Client) rememberTickFromBook(assetID, tick string) {
	if assetID == "" || tick == "" {
		return
	}
	c.tickSizes[assetID] = tick
	c.tickSizeAt[assetID] = time.Now()
}

// ResolveTickSize returns user-provided tick if valid vs market minimum, else minimum tick.
func (c *Client) ResolveTickSize(tokenID string, userTick *string) (string, error) {
	minTick, err := c.GetTickSize(tokenID)
	if err != nil {
		return "", err
	}
	if userTick == nil || *userTick == "" {
		return minTick, nil
	}
	smaller, err := tickSizeSmaller(*userTick, minTick)
	if err != nil {
		return "", err
	}
	if smaller {
		return "", fmt.Errorf("invalid tick size (%s), minimum for the market is %s", *userTick, minTick)
	}
	return *userTick, nil
}

func tickSizeSmaller(a, b string) (bool, error) {
	af, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return false, err
	}
	bf, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return false, err
	}
	return af < bf, nil
}
