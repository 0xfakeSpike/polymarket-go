package polymarket

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPublicProfile_mock(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public-profile" || r.URL.Query().Get("address") != "0x1111111111111111111111111111111111111111" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.String())
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"tester","proxyWallet":"0x1111111111111111111111111111111111111111"}`))
	}))
	defer srv.Close()

	c, err := NewPublicClient(WithGammaAPIHost(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	p, err := c.GetPublicProfile("0x1111111111111111111111111111111111111111")
	if err != nil {
		t.Fatal(err)
	}
	if p.Name == nil || *p.Name != "tester" {
		t.Fatalf("name: %+v", p.Name)
	}
}

func TestGetCurrentPositions_mock(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/positions" || r.URL.Query().Get("user") != "0x2222222222222222222222222222222222222222" {
			t.Fatalf("unexpected request: %s", r.URL.String())
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"conditionId":"0xabc","size":1,"title":"x"}]`))
	}))
	defer srv.Close()

	c, err := NewPublicClient(WithDataAPIHost(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	pos, err := c.GetCurrentPositions(&CurrentPositionsParams{User: "0x2222222222222222222222222222222222222222"})
	if err != nil {
		t.Fatal(err)
	}
	if len(pos) != 1 || pos[0].Title != "x" {
		t.Fatalf("got %+v", pos)
	}
}

func TestGetClosedPositions_mock(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/closed-positions" {
			t.Fatalf("path %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	c, err := NewPublicClient(WithDataAPIHost(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetClosedPositions(&ClosedPositionsParams{User: "0x3333333333333333333333333333333333333333"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetUserActivity_mock(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/activity" || r.URL.Query().Get("type") != "TRADE,REDEEM" {
			t.Fatalf("unexpected %s", r.URL.String())
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"type":"TRADE","timestamp":1}]`))
	}))
	defer srv.Close()

	c, err := NewPublicClient(WithDataAPIHost(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	act, err := c.GetUserActivity(&UserActivityParams{
		User:          "0x4444444444444444444444444444444444444444",
		ActivityTypes: []string{"TRADE", "REDEEM"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(act) != 1 || act[0].Type != "TRADE" {
		t.Fatalf("got %+v", act)
	}
}

func TestGetPublicProfile_requiresAddress(t *testing.T) {
	c, err := NewPublicClient()
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetPublicProfile(" ")
	if err == nil {
		t.Fatal("expected error")
	}
}
