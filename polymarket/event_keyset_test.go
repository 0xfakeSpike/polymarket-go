package polymarket

import (
	"testing"
	"time"
)

func TestEventsKeysetParamsValues(t *testing.T) {
	closed := false
	live := true
	week := 3
	ts := time.Date(2026, 1, 2, 15, 4, 5, 0, time.UTC)
	p := &EventsKeysetParams{
		Limit:        50,
		Order:        "volume",
		Ascending:    false,
		AfterCursor:  "cursor-token",
		Locale:       "en",
		ID:           []int{1, 2},
		Slug:         []string{"a-b"},
		Closed:       &closed,
		Live:         &live,
		TagID:        []int{10},
		ExcludeTagID: []int{20},
		SeriesID:     []int{100},
		GameID:       []int{200, 201},
		EventWeek:    &week,
		EventDate:    &ts,
	}

	v := p.Values()
	if got := v.Get("limit"); got != "50" {
		t.Fatalf("limit mismatch: %s", got)
	}
	if got := v.Get("order"); got != "volume" {
		t.Fatalf("order mismatch: %s", got)
	}
	if got := v.Get("ascending"); got != "false" {
		t.Fatalf("ascending mismatch: %s", got)
	}
	if got := v.Get("after_cursor"); got != "cursor-token" {
		t.Fatalf("after_cursor mismatch: %s", got)
	}
	if got := v.Get("locale"); got != "en" {
		t.Fatalf("locale mismatch: %s", got)
	}
	if got := v["id"]; len(got) != 2 || got[0] != "1" || got[1] != "2" {
		t.Fatalf("id mismatch: %v", got)
	}
	if got := v.Get("slug"); got != "a-b" {
		t.Fatalf("slug mismatch: %s", got)
	}
	if got := v.Get("closed"); got != "false" {
		t.Fatalf("closed mismatch: %s", got)
	}
	if got := v.Get("live"); got != "true" {
		t.Fatalf("live mismatch: %s", got)
	}
	if got := v.Get("tag_id"); got != "10" {
		t.Fatalf("tag_id mismatch: %s", got)
	}
	if got := v.Get("exclude_tag_id"); got != "20" {
		t.Fatalf("exclude_tag_id mismatch: %s", got)
	}
	if got := v.Get("series_id"); got != "100" {
		t.Fatalf("series_id mismatch: %s", got)
	}
	if got := v["game_id"]; len(got) != 2 || got[0] != "200" || got[1] != "201" {
		t.Fatalf("game_id mismatch: %v", got)
	}
	if got := v.Get("event_week"); got != "3" {
		t.Fatalf("event_week mismatch: %s", got)
	}
	if got := v.Get("event_date"); got != "2026-01-02T15:04:05Z" {
		t.Fatalf("event_date mismatch: %s", got)
	}
}

func TestEventsKeysetParamsValuesNil(t *testing.T) {
	var p *EventsKeysetParams
	if p.Values() == nil {
		t.Fatal("expected non-nil url.Values")
	}
}
