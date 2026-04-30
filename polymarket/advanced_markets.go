package polymarket

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

// AnnualizedReturnMarketsParams controls market scanning and ranking.
type AnnualizedReturnMarketsParams struct {
	// Limit caps final rows returned. Default 20.
	Limit int
	// MaxPages caps Gamma /events/keyset pages scanned. Default 3.
	MaxPages int
	// EventsPageLimit is the page size for each Gamma keyset request (1–500). Default 100.
	EventsPageLimit int
	// TagSlug filters events/markets via Gamma tag_slug (optional).
	TagSlug string
	// MinBestAsk keeps only outcomes whose best ask is >= this threshold. Default 0.5.
	MinBestAsk float64
	// Now overrides current time for deterministic callers/tests.
	Now time.Time
}

// MarketAnnualizedReturn is a ranked market/outcome opportunity snapshot.
type MarketAnnualizedReturn struct {
	ConditionID      string        `json:"condition_id"`
	Question         string        `json:"question"`
	TokenID          string        `json:"token_id"`
	Outcome          string        `json:"outcome"`
	BestAsk          float64       `json:"best_ask"`
	BestBid          float64       `json:"best_bid,omitempty"`
	SettlementTime   time.Time     `json:"settlement_time"`
	HoldingPeriod    time.Duration `json:"holding_period"`
	PnLPerShare      float64       `json:"pnl_per_share"`
	ROI              float64       `json:"roi"`
	AnnualizedReturn float64       `json:"annualized_return"`
}

type marketTokenRef struct {
	TokenID string
	Outcome string
}

type marketScanRow struct {
	ConditionID string
	Question    string
	EndDate     time.Time
	Tokens      []marketTokenRef
}

// GetMarketsByAnnualizedReturn scans markets from the Gamma /events/keyset feed and returns rows
// sorted by annualized return descending.
//
// Strategy:
//  1. Page Gamma GET /events/keyset (optional tag_slug via [AnnualizedReturnMarketsParams.TagSlug]).
//  2. Flatten each event's markets; decode condition id, settlement time, and CLOB token ids.
//  3. For each market, fetch token order books and pick the outcome with highest best ask >= min threshold.
//  4. Compute ROI=(1-ask)/ask and annualized return=(1+ROI)^(1/years)-1 until settlement.
func (c *Client) GetMarketsByAnnualizedReturn(params *AnnualizedReturnMarketsParams) ([]MarketAnnualizedReturn, error) {
	if c == nil {
		return nil, fmt.Errorf("nil client")
	}
	cfg := normalizeAnnualizedParams(params)
	now := cfg.Now

	rows := make([]MarketAnnualizedReturn, 0, cfg.Limit)
	seen := make(map[string]struct{})
	closed := false
	after := ""
	for page := 0; page < cfg.MaxPages; page++ {
		evp, err := c.GetEventsKeyset(&EventsKeysetParams{
			Limit:       cfg.EventsPageLimit,
			AfterCursor: after,
			TagSlug:     cfg.TagSlug,
			Closed:      &closed,
		})
		if err != nil {
			return nil, err
		}
		for _, ev := range evp.Events {
			for _, raw := range ev.Markets {
				m, ok := decodeGammaMarketToScanRow(raw)
				if !ok || !m.EndDate.After(now) || len(m.Tokens) == 0 {
					continue
				}
				if _, ok := seen[m.ConditionID]; ok {
					continue
				}
				seen[m.ConditionID] = struct{}{}

				best, ok := c.bestOutcomeByAsk(m, cfg.MinBestAsk, now)
				if !ok {
					continue
				}
				rows = append(rows, best)
			}
		}
		if evp.NextCursor == "" {
			break
		}
		after = evp.NextCursor
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].AnnualizedReturn > rows[j].AnnualizedReturn
	})
	if len(rows) > cfg.Limit {
		rows = rows[:cfg.Limit]
	}
	return rows, nil
}

func normalizeAnnualizedParams(p *AnnualizedReturnMarketsParams) AnnualizedReturnMarketsParams {
	out := AnnualizedReturnMarketsParams{
		Limit:           20,
		MaxPages:        3,
		EventsPageLimit: 100,
		MinBestAsk:      0.5,
		Now:             time.Now(),
	}
	if p == nil {
		return out
	}
	if p.Limit > 0 {
		out.Limit = p.Limit
	}
	if p.MaxPages > 0 {
		out.MaxPages = p.MaxPages
	}
	if p.EventsPageLimit > 0 {
		out.EventsPageLimit = p.EventsPageLimit
		if out.EventsPageLimit > 500 {
			out.EventsPageLimit = 500
		}
	}
	if strings.TrimSpace(p.TagSlug) != "" {
		out.TagSlug = strings.TrimSpace(p.TagSlug)
	}
	if p.MinBestAsk > 0 {
		out.MinBestAsk = p.MinBestAsk
	}
	if !p.Now.IsZero() {
		out.Now = p.Now
	}
	return out
}

func (c *Client) bestOutcomeByAsk(m marketScanRow, minAsk float64, now time.Time) (MarketAnnualizedReturn, bool) {
	var best MarketAnnualizedReturn
	ok := false
	for _, t := range m.Tokens {
		book, err := c.GetOrderBook(t.TokenID)
		if err != nil {
			continue
		}
		asks := book.AsksData()
		if len(asks) == 0 {
			continue
		}
		bestAsk := asks[0].Price
		if bestAsk < minAsk || bestAsk <= 0 {
			continue
		}
		bestBid := 0.0
		bids := book.BidsData()
		if len(bids) > 0 {
			bestBid = bids[0].Price
		}

		returnRow, rowOK := buildAnnualizedRow(m, t, bestAsk, bestBid, now)
		if !rowOK {
			continue
		}
		if !ok || returnRow.BestAsk > best.BestAsk {
			best = returnRow
			ok = true
		}
	}
	return best, ok
}

func buildAnnualizedRow(m marketScanRow, t marketTokenRef, bestAsk, bestBid float64, now time.Time) (MarketAnnualizedReturn, bool) {
	holding := m.EndDate.Sub(now)
	if holding <= 0 {
		return MarketAnnualizedReturn{}, false
	}
	roi := (1 - bestAsk) / bestAsk
	annualized, ok := annualizedReturnFromROI(roi, holding)
	if !ok {
		return MarketAnnualizedReturn{}, false
	}
	return MarketAnnualizedReturn{
		ConditionID:      m.ConditionID,
		Question:         m.Question,
		TokenID:          t.TokenID,
		Outcome:          t.Outcome,
		BestAsk:          bestAsk,
		BestBid:          bestBid,
		SettlementTime:   m.EndDate,
		HoldingPeriod:    holding,
		PnLPerShare:      1 - bestAsk,
		ROI:              roi,
		AnnualizedReturn: annualized,
	}, true
}

func annualizedReturnFromROI(roi float64, holding time.Duration) (float64, bool) {
	if holding <= 0 || roi <= -1 {
		return 0, false
	}
	years := holding.Seconds() / (365.25 * 24 * 3600)
	if years <= 0 {
		return 0, false
	}
	return math.Pow(1+roi, 1/years) - 1, true
}

func decodeGammaMarketToScanRow(item map[string]any) (marketScanRow, bool) {
	if item == nil {
		return marketScanRow{}, false
	}
	conditionID := firstString(item, "conditionId", "condition_id", "id")
	if conditionID == "" {
		return marketScanRow{}, false
	}
	endISO := firstString(item, "endDateIso", "end_date_iso", "endDate", "end_date", "umaEndDateIso", "uma_end_date_iso")
	endDate, ok := parseMarketEndTime(endISO)
	if !ok {
		return marketScanRow{}, false
	}
	tokens := decodeGammaMarketTokens(item)
	if len(tokens) == 0 {
		return marketScanRow{}, false
	}
	return marketScanRow{
		ConditionID: conditionID,
		Question:    firstString(item, "question", "title", "groupItemTitle", "group_item_title"),
		EndDate:     endDate,
		Tokens:      tokens,
	}, true
}

func decodeGammaMarketTokens(item map[string]any) []marketTokenRef {
	ids := firstFlexibleStringSlice(item, "clobTokenIds", "clob_token_ids", "assets_ids", "asset_ids")
	outcomes := firstFlexibleStringSlice(item, "outcomes", "shortOutcomes", "short_outcomes")
	if len(ids) == 0 {
		return nil
	}
	out := make([]marketTokenRef, 0, len(ids))
	for i, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		t := marketTokenRef{TokenID: id}
		if i < len(outcomes) {
			t.Outcome = outcomes[i]
		}
		out = append(out, t)
	}
	return out
}

func firstFlexibleStringSlice(m map[string]any, keys ...string) []string {
	for _, k := range keys {
		if v, ok := m[k]; ok {
			if sl := flexibleStringSlice(v); len(sl) > 0 {
				return sl
			}
		}
	}
	return nil
}

// flexibleStringSlice turns a Gamma JSON field into a []string: JSON arrays, Go []any, or JSON-in-string.
func flexibleStringSlice(v any) []string {
	if v == nil {
		return nil
	}
	switch x := v.(type) {
	case string:
		s := strings.TrimSpace(x)
		if s == "" || s == "null" {
			return nil
		}
		if strings.HasPrefix(s, "[") {
			var arr []string
			if err := json.Unmarshal([]byte(s), &arr); err == nil && len(arr) > 0 {
				return arr
			}
			var raw []json.RawMessage
			if err := json.Unmarshal([]byte(s), &raw); err == nil && len(raw) > 0 {
				out := make([]string, 0, len(raw))
				for _, r := range raw {
					var s2 string
					if json.Unmarshal(r, &s2) == nil && strings.TrimSpace(s2) != "" {
						out = append(out, strings.TrimSpace(s2))
						continue
					}
					var n json.Number
					if json.Unmarshal(r, &n) == nil {
						out = append(out, n.String())
					}
				}
				if len(out) > 0 {
					return out
				}
			}
			return nil
		}
		return []string{s}
	case []any:
		return toStringSlice(x)
	default:
		return nil
	}
}

func decodeMarketScanRows(raw json.RawMessage) ([]marketScanRow, error) {
	var wire []map[string]any
	if err := json.Unmarshal(raw, &wire); err != nil {
		return nil, fmt.Errorf("decode markets page data: %w", err)
	}
	out := make([]marketScanRow, 0, len(wire))
	for _, item := range wire {
		row, ok := decodeMarketScanRow(item)
		if !ok {
			continue
		}
		out = append(out, row)
	}
	return out, nil
}

func decodeMarketScanRow(item map[string]any) (marketScanRow, bool) {
	conditionID := firstString(item, "condition_id", "conditionId", "id")
	endISO := firstString(item, "end_date_iso", "endDateIso", "end_date", "endDate")
	endDate, ok := parseMarketEndTime(endISO)
	if !ok || conditionID == "" {
		return marketScanRow{}, false
	}
	tokens := decodeMarketTokens(item)
	if len(tokens) == 0 {
		return marketScanRow{}, false
	}
	return marketScanRow{
		ConditionID: conditionID,
		Question:    firstString(item, "question", "title"),
		EndDate:     endDate,
		Tokens:      tokens,
	}, true
}

func decodeMarketTokens(item map[string]any) []marketTokenRef {
	rawTokens, ok := item["tokens"].([]any)
	if ok && len(rawTokens) > 0 {
		out := make([]marketTokenRef, 0, len(rawTokens))
		for _, rt := range rawTokens {
			m, ok := rt.(map[string]any)
			if !ok {
				continue
			}
			tokenID := firstString(m, "token_id", "tokenId", "t")
			if tokenID == "" {
				continue
			}
			out = append(out, marketTokenRef{
				TokenID: tokenID,
				Outcome: firstString(m, "outcome", "o"),
			})
		}
		if len(out) > 0 {
			return out
		}
	}

	assetIDs := toStringSlice(item["assets_ids"])
	outcomes := toStringSlice(item["outcomes"])
	if len(assetIDs) == 0 {
		assetIDs = toStringSlice(item["clob_token_ids"])
	}
	if len(assetIDs) == 0 {
		return nil
	}
	out := make([]marketTokenRef, 0, len(assetIDs))
	for i, id := range assetIDs {
		t := marketTokenRef{TokenID: id}
		if i < len(outcomes) {
			t.Outcome = outcomes[i]
		}
		out = append(out, t)
	}
	return out
}

func firstString(m map[string]any, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k]; ok {
			if s, ok := v.(string); ok {
				return strings.TrimSpace(s)
			}
		}
	}
	return ""
}

func toStringSlice(v any) []string {
	raw, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(raw))
	for _, it := range raw {
		s, ok := it.(string)
		if !ok || strings.TrimSpace(s) == "" {
			continue
		}
		out = append(out, s)
	}
	return out
}

func parseMarketEndTime(v string) (time.Time, bool) {
	v = strings.TrimSpace(v)
	if v == "" {
		return time.Time{}, false
	}
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05.000000Z",
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, v); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}
