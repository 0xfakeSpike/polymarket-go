package polymarket

import (
	"math"
	"testing"
	"time"
)

func TestFavoredSidePNLFromOrderBooks(t *testing.T) {
	now := time.Date(2026, 4, 22, 12, 0, 0, 0, time.UTC)
	end := now.Add(24 * time.Hour)

	m := &Market{
		ID:           "m1",
		Outcomes:     []string{"YES", "NO"},
		ClobTokenIds: []string{"yes-token", "no-token"},
		EndDate:      &end,
	}
	books := map[string]*Book{
		"yes-token": {
			Asks: []LimitOrder{{Price: 0.62, Size: 100}},
			Bids: []LimitOrder{{Price: 0.60, Size: 100}},
		},
		"no-token": {
			Asks: []LimitOrder{{Price: 0.40, Size: 100}},
			Bids: []LimitOrder{{Price: 0.38, Size: 100}},
		},
	}

	got, err := m.FavoredSidePNLFromOrderBooks(now, books)
	if err != nil {
		t.Fatalf("FavoredSidePNLFromOrderBooks returned error: %v", err)
	}
	if got.Outcome != "YES" {
		t.Fatalf("expected YES, got %s", got.Outcome)
	}
	if got.TokenID != "yes-token" {
		t.Fatalf("expected yes-token, got %s", got.TokenID)
	}
	if got.Price != 0.62 || got.BestAsk != 0.62 || got.BestBid != 0.60 {
		t.Fatalf("unexpected prices: %+v", got)
	}
	if got.PnLPerShare != 0.38 {
		t.Fatalf("expected pnl 0.38, got %.8f", got.PnLPerShare)
	}
	if got.AnnualizedReturn == nil {
		t.Fatal("expected non-nil annualized return")
	}
	years := got.HoldingPeriod.Seconds() / (365.25 * 24 * 3600)
	wantAR := math.Pow(1+got.ROI, 1/years) - 1
	if math.Abs(*got.AnnualizedReturn-wantAR) > 1e-9 {
		t.Fatalf("annualized return mismatch: got %g want %g", *got.AnnualizedReturn, wantAR)
	}
}

func TestFavoredSidePNLFromOrderBooksPickNoSide(t *testing.T) {
	now := time.Date(2026, 4, 22, 12, 0, 0, 0, time.UTC)
	end := now.Add(24 * time.Hour)
	m := &Market{
		Outcomes:     []string{"YES", "NO"},
		ClobTokenIds: []string{"yes-token", "no-token"},
		EndDate:      &end,
	}
	books := map[string]*Book{
		"yes-token": {Asks: []LimitOrder{{Price: 0.45, Size: 100}}, Bids: []LimitOrder{{Price: 0.44, Size: 50}}},
		"no-token":  {Asks: []LimitOrder{{Price: 0.58, Size: 100}}, Bids: []LimitOrder{{Price: 0.57, Size: 60}}},
	}

	got, err := m.FavoredSidePNLFromOrderBooks(now, books)
	if err != nil {
		t.Fatalf("FavoredSidePNLFromOrderBooks returned error: %v", err)
	}
	if got.Outcome != "NO" || got.TokenID != "no-token" {
		t.Fatalf("expected NO/no-token, got outcome=%s token=%s", got.Outcome, got.TokenID)
	}
	if got.AnnualizedReturn == nil {
		t.Fatal("expected non-nil annualized return")
	}
}

func TestFavoredSidePNLFromOrderBooks_NoAskAboveHalf(t *testing.T) {
	now := time.Date(2026, 4, 22, 12, 0, 0, 0, time.UTC)
	end := now.Add(24 * time.Hour)
	m := &Market{
		Outcomes:     []string{"YES", "NO"},
		ClobTokenIds: []string{"yes-token", "no-token"},
		EndDate:      &end,
	}
	books := map[string]*Book{
		"yes-token": {Asks: []LimitOrder{{Price: 0.50, Size: 100}}},
		"no-token":  {Asks: []LimitOrder{{Price: 0.49, Size: 100}}},
	}

	_, err := m.FavoredSidePNLFromOrderBooks(now, books)
	if err == nil {
		t.Fatalf("expected error when no best ask > 0.5")
	}
}
