package polymarket

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetEventsKeyset_mock(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/events/keyset" {
			t.Fatalf("path %s", r.URL.Path)
		}
		if r.URL.Query().Get("tag_slug") != "crypto" {
			t.Fatalf("tag_slug: %s", r.URL.Query().Get("tag_slug"))
		}
		if r.URL.Query().Get("closed") != "false" {
			t.Fatalf("closed: %s", r.URL.Query().Get("closed"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"events":[{"markets":[]}],"next_cursor":"abc"}`))
	}))
	defer srv.Close()

	c, err := NewPublicClient(WithGammaAPIHost(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	closed := false
	out, err := c.GetEventsKeyset(&EventsKeysetParams{Limit: 10, TagSlug: "crypto", Closed: &closed})
	if err != nil {
		t.Fatal(err)
	}
	if out.NextCursor != "abc" || len(out.Events) != 1 {
		t.Fatalf("response: %+v", out)
	}
}
