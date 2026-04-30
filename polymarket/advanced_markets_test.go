package polymarket

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDecodeGammaMarketToScanRow_stringArrays(t *testing.T) {
	raw := `{
		"conditionId": "0xabc",
		"question": "Will it rain?",
		"endDateIso": "2030-01-15T00:00:00Z",
		"clobTokenIds": "[\"111\",\"222\"]",
		"outcomes": "[\"Yes\",\"No\"]"
	}`
	var m map[string]any
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		t.Fatal(err)
	}
	row, ok := decodeGammaMarketToScanRow(m)
	if !ok {
		t.Fatal("expected ok")
	}
	if row.ConditionID != "0xabc" || row.Question != "Will it rain?" {
		t.Fatalf("meta: %+v", row)
	}
	if len(row.Tokens) != 2 || row.Tokens[0].TokenID != "111" || row.Tokens[0].Outcome != "Yes" {
		t.Fatalf("tokens: %+v", row.Tokens)
	}
	if !row.EndDate.After(time.Date(2029, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("end: %v", row.EndDate)
	}
}

func TestDecodeGammaMarketToScanRow_nativeArrays(t *testing.T) {
	m := map[string]any{
		"conditionId":  "0xdef",
		"question":     "Q",
		"endDate":      "2030-06-01T12:00:00Z",
		"clobTokenIds": []any{"aa", "bb"},
		"outcomes":     []any{"A", "B"},
	}
	row, ok := decodeGammaMarketToScanRow(m)
	if !ok {
		t.Fatal("expected ok")
	}
	if len(row.Tokens) != 2 || row.Tokens[1].Outcome != "B" {
		t.Fatalf("tokens: %+v", row.Tokens)
	}
}

func TestDecodeGammaMarketToScanRow_dateOnlyEndDateIso(t *testing.T) {
	m := map[string]any{
		"conditionId":  "0xdate",
		"question":     "Date only",
		"endDateIso":   "2025-12-31",
		"clobTokenIds": "[\"x\",\"y\"]",
		"outcomes":     "[\"Yes\",\"No\"]",
	}
	row, ok := decodeGammaMarketToScanRow(m)
	if !ok {
		t.Fatal("expected ok for date-only endDateIso")
	}
	if got := row.EndDate.UTC().Format("2006-01-02"); got != "2025-12-31" {
		t.Fatalf("unexpected parsed date: %s", got)
	}
}
