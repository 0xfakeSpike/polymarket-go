package polymarket

import "testing"

func TestMarketsKeysetParamsValues(t *testing.T) {
	closed := false
	includeTag := true
	p := &MarketsKeysetParams{
		Limit:       20,
		Order:       "volume_num,liquidity_num",
		Ascending:   false,
		AfterCursor: "abc",
		Locale:      "zh",
		Slug:        []string{"foo"},
		Closed:      &closed,
		IncludeTag:  &includeTag,
	}

	v := p.Values()
	if got := v.Get("limit"); got != "20" {
		t.Fatalf("limit mismatch: %s", got)
	}
	if got := v.Get("order"); got != "volume_num,liquidity_num" {
		t.Fatalf("order mismatch: %s", got)
	}
	if got := v.Get("ascending"); got != "false" {
		t.Fatalf("ascending mismatch: %s", got)
	}
	if got := v.Get("after_cursor"); got != "abc" {
		t.Fatalf("after_cursor mismatch: %s", got)
	}
	if got := v.Get("locale"); got != "zh" {
		t.Fatalf("locale mismatch: %s", got)
	}
	if got := v.Get("closed"); got != "false" {
		t.Fatalf("closed mismatch: %s", got)
	}
	if got := v.Get("include_tag"); got != "true" {
		t.Fatalf("include_tag mismatch: %s", got)
	}
}
