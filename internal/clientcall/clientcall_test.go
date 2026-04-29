package clientcall

import (
	"encoding/json"
	"testing"

	"github.com/0xfakeSpike/polymarket-go/polymarket"
)

func TestListClientMethods_containsGetOK(t *testing.T) {
	names := ListClientMethods()
	found := false
	for _, n := range names {
		if n == "GetOK" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected GetOK in method list")
	}
}

func TestInvoke_ChainID_noNetwork(t *testing.T) {
	c, err := polymarket.NewPublicClient()
	if err != nil {
		t.Fatal(err)
	}
	out, err := Invoke(c, "ChainID", []byte("[]"))
	if err != nil {
		t.Fatal(err)
	}
	if out == nil {
		t.Fatal("expected chain id")
	}
	_, err = json.Marshal(out)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
}
