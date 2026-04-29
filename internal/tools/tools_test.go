package tools

import (
	"encoding/json"
	"testing"

	"github.com/0xfakeSpike/polymarket-go"
)

func TestList_sortedAndIncludesCoreTools(t *testing.T) {
	got := List()
	if len(got) == 0 {
		t.Fatal("expected tools")
	}
	for i := 1; i < len(got); i++ {
		if got[i-1].Name > got[i].Name {
			t.Fatalf("tools are not sorted: %q before %q", got[i-1].Name, got[i].Name)
		}
	}

	names := map[string]bool{}
	for _, tool := range got {
		names[tool.Name] = true
	}
	for _, want := range []string{"client_call", "get_orderbook", "methods", "rank_markets_by_annualized_return", "search_events"} {
		if !names[want] {
			t.Fatalf("missing tool %q", want)
		}
	}
}

func TestCall_clientCallNoNetwork(t *testing.T) {
	c, err := polymarket.NewPublicClient()
	if err != nil {
		t.Fatal(err)
	}
	params := json.RawMessage(`{"method":"ChainID","args":[]}`)
	got, err := Call(c, "client_call", params)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil {
		t.Fatal("expected result")
	}
}

func TestCall_validatesParams(t *testing.T) {
	c, err := polymarket.NewPublicClient()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Call(c, "search_events", nil); err == nil {
		t.Fatal("expected missing query error")
	}
	if _, err := Call(c, "missing_tool", nil); err == nil {
		t.Fatal("expected unknown tool error")
	}
}

func TestEventMatchesKeyword(t *testing.T) {
	ev := polymarket.Event{
		Title:       "Will Iran close strait this year?",
		Description: "Geopolitical risk event",
	}
	if !eventMatchesKeyword(ev, "iran") {
		t.Fatal("expected keyword match in title")
	}
	if eventMatchesKeyword(ev, "bitcoin") {
		t.Fatal("did not expect unrelated keyword to match")
	}
}

func TestRankMarketsParamsDecodeMinAnnualizedReturn(t *testing.T) {
	var p rankMarketsParams
	raw := json.RawMessage(`{"min_annualized_return":0.3}`)
	if err := decodeParams(raw, &p); err != nil {
		t.Fatal(err)
	}
	if p.MinAnnualizedReturn == nil {
		t.Fatal("expected min_annualized_return to be set")
	}
	if *p.MinAnnualizedReturn != 0.3 {
		t.Fatalf("unexpected min_annualized_return: got %v", *p.MinAnnualizedReturn)
	}
}
