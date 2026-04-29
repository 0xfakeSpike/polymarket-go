package polymarket

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestCLOBServer(t *testing.T, handler http.HandlerFunc) (*Client, func()) {
	t.Helper()
	srv := httptest.NewServer(handler)
	c, err := NewPublicClient(WithCLOBHost(srv.URL))
	if err != nil {
		srv.Close()
		t.Fatal(err)
	}
	return c, srv.Close
}

func TestGetOKUsesOKEndpoint(t *testing.T) {
	c, closeServer := newTestCLOBServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != PathOK {
			t.Fatalf("path = %s, want %s", r.URL.Path, PathOK)
		}
		_, _ = w.Write([]byte(`{"ok":true}`))
	})
	defer closeServer()

	if _, err := c.GetOK(); err != nil {
		t.Fatal(err)
	}
}

func TestGetClobMarketInfoCachesTokenMetadata(t *testing.T) {
	c, closeServer := newTestCLOBServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != PathCLOBMarketPrefix+"cond-1" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{
			"c":"cond-1",
			"t":[{"t":"token-yes","o":"YES"},{"t":"token-no","o":"NO"}],
			"mts":0.01,
			"nr":true,
			"fd":{"r":0.02,"e":2,"to":true}
		}`))
	})
	defer closeServer()

	info, err := c.GetClobMarketInfo("cond-1")
	if err != nil {
		t.Fatal(err)
	}
	if info.ConditionID != "cond-1" {
		t.Fatalf("condition id = %s", info.ConditionID)
	}
	if got := c.tokenConditionMap["token-yes"]; got != "cond-1" {
		t.Fatalf("token condition cache = %s", got)
	}
	if got := c.tickSizes["token-yes"]; got != "0.01" {
		t.Fatalf("tick cache = %s", got)
	}
	if got := c.negRiskCache["token-yes"]; !got {
		t.Fatal("expected neg risk cache")
	}
	exp, err := c.GetFeeExponent("token-yes")
	if err != nil {
		t.Fatal(err)
	}
	if exp != 2 {
		t.Fatalf("fee exponent = %v", exp)
	}
}

func TestGetBuilderTradesUsesPublicBuilderCodeQuery(t *testing.T) {
	c, closeServer := newTestCLOBServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != PathBuilderTrades {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("builder_code"); got != "builder-1" {
			t.Fatalf("builder_code = %s", got)
		}
		if got := r.Header.Get("POLY_SIGNATURE"); got != "" {
			t.Fatalf("unexpected auth header %s", got)
		}
		_, _ = w.Write([]byte(`{"data":[],"next_cursor":"LTE=","limit":0,"count":0}`))
	})
	defer closeServer()

	page, err := c.GetBuilderTrades(&BuilderTradeParams{BuilderCode: "builder-1"}, "")
	if err != nil {
		t.Fatal(err)
	}
	if page.NextCursor != EndCursor {
		t.Fatalf("next cursor = %s", page.NextCursor)
	}
	b, err := json.Marshal(page)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), `"trades"`) || strings.Contains(string(b), `"data"`) {
		t.Fatalf("unexpected public JSON: %s", string(b))
	}
}

func TestOrderBookSummaryHashUsesStringLevels(t *testing.T) {
	book := &Book{
		Market:         "m",
		AssetID:        "a",
		Timestamp:      "1",
		Bids:           []LimitOrder{{Price: 0.4, Size: 100}},
		Asks:           []LimitOrder{{Price: 0.6, Size: 100}},
		MinOrderSize:   "1",
		TickSize:       "0.01",
		NegRisk:        false,
		Hash:           "old",
		LastTradePrice: "0.5",
	}
	got, err := OrderBookSummaryHash(book)
	if err != nil {
		t.Fatal(err)
	}
	wire := `{"market":"m","asset_id":"a","timestamp":"1","bids":[{"price":"0.4","size":"100"}],"asks":[{"price":"0.6","size":"100"}],"min_order_size":"1","tick_size":"0.01","neg_risk":false,"hash":"","last_trade_price":"0.5"}`
	want, err := sha1Hex(wire)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("hash = %s, want %s", got, want)
	}
}

func sha1Hex(s string) (string, error) {
	sum := sha1.Sum([]byte(s))
	return hex.EncodeToString(sum[:]), nil
}
