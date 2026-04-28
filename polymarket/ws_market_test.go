package polymarket

import "testing"

func TestParseMarketChannelMessage_book(t *testing.T) {
	data := []byte(`{"event_type":"book","asset_id":"a","market":"m","bids":[{"price":"0.5","size":"1"}],"asks":[],"timestamp":"1","hash":"h"}`)
	msg, err := parseMarketChannelMessage(data)
	if err != nil {
		t.Fatal(err)
	}
	if msg.Book == nil || msg.Book.AssetID != "a" || len(msg.Book.Bids) != 1 {
		t.Fatalf("%+v", msg)
	}
}

func TestParseMarketChannelMessage_priceChange(t *testing.T) {
	data := []byte(`{"event_type":"price_change","market":"m","price_changes":[{"asset_id":"a","price":"0.5","size":"1","side":"BUY","hash":"x"}],"timestamp":"1"}`)
	msg, err := parseMarketChannelMessage(data)
	if err != nil || msg.PriceChange == nil || len(msg.PriceChange.PriceChanges) != 1 {
		t.Fatalf("err=%v %+v", err, msg)
	}
}

func TestParseMarketChannelMessage_lastTrade(t *testing.T) {
	data := []byte(`{"event_type":"last_trade_price","asset_id":"a","market":"m","price":"0.5","size":"1","side":"SELL","timestamp":"1"}`)
	msg, err := parseMarketChannelMessage(data)
	if err != nil || msg.LastTradePrice == nil || msg.LastTradePrice.Price != "0.5" {
		t.Fatalf("%+v", msg)
	}
}

func TestParseMarketChannelMessage_tickSize(t *testing.T) {
	data := []byte(`{"event_type":"tick_size_change","asset_id":"a","market":"m","old_tick_size":"0.01","new_tick_size":"0.001","timestamp":"1"}`)
	msg, err := parseMarketChannelMessage(data)
	if err != nil || msg.TickSizeChange == nil {
		t.Fatalf("%+v", msg)
	}
}

func TestParseMarketChannelMessage_bestBidAsk(t *testing.T) {
	data := []byte(`{"event_type":"best_bid_ask","asset_id":"a","market":"m","best_bid":"0.1","best_ask":"0.2","spread":"0.1","timestamp":"1"}`)
	msg, err := parseMarketChannelMessage(data)
	if err != nil || msg.BestBidAsk == nil {
		t.Fatalf("%+v", msg)
	}
}

func TestParseMarketChannelMessage_newMarket(t *testing.T) {
	data := []byte(`{"event_type":"new_market","id":"1","question":"q?","market":"m","slug":"s","assets_ids":["x"],"outcomes":["Y","N"],"timestamp":"1"}`)
	msg, err := parseMarketChannelMessage(data)
	if err != nil || msg.NewMarket == nil || msg.NewMarket.Slug != "s" {
		t.Fatalf("%+v", msg)
	}
}

func TestParseMarketChannelMessage_resolved(t *testing.T) {
	data := []byte(`{"event_type":"market_resolved","id":"1","market":"m","assets_ids":["x"],"winning_asset_id":"x","winning_outcome":"Yes","timestamp":"1"}`)
	msg, err := parseMarketChannelMessage(data)
	if err != nil || msg.Resolved == nil {
		t.Fatalf("%+v", msg)
	}
}

func TestParseMarketChannelMessage_skipUnknown(t *testing.T) {
	msg, err := parseMarketChannelMessage([]byte(`{"event_type":"future_type"}`))
	if err != nil || !msg.empty() {
		t.Fatalf("msg=%+v err=%v", msg, err)
	}
	msg, err = parseMarketChannelMessage([]byte(`not json`))
	if err != nil || !msg.empty() {
		t.Fatalf("msg=%+v err=%v", msg, err)
	}
}

func TestParseMarketChannelMessage_bookBadJSON(t *testing.T) {
	_, err := parseMarketChannelMessage([]byte(`{"event_type":"book","bids":"not-array"}`))
	if err == nil {
		t.Fatal("expected error")
	}
}
