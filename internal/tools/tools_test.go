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
	for _, want := range []string{"client_call", "get_orderbook", "methods"} {
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
	if _, err := Call(c, "get_orderbook", nil); err == nil {
		t.Fatal("expected missing token_id error")
	}
	if _, err := Call(c, "missing_tool", nil); err == nil {
		t.Fatal("expected unknown tool error")
	}
}
