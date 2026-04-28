package polymarket

import (
	"testing"
)

func TestParseUserChannelEvent_order(t *testing.T) {
	data := []byte(`{"event_type":"order","id":"1","owner":"o","market":"m","asset_id":"a","side":"BUY","original_size":"1","size_matched":"0","price":"0.5","type":"PLACEMENT","timestamp":"1"}`)
	msg, err := parseUserChannelEvent(data)
	if err != nil {
		t.Fatal(err)
	}
	if msg.Order == nil || msg.Order.ID != "1" || msg.Trade != nil {
		t.Fatalf("got %#v", msg)
	}
}

func TestParseUserChannelEvent_trade(t *testing.T) {
	data := []byte(`{"event_type":"trade","type":"TRADE","id":"t","taker_order_id":"x","market":"m","asset_id":"a","side":"BUY","size":"1","price":"0.5","status":"MATCHED","owner":"o","timestamp":"1"}`)
	msg, err := parseUserChannelEvent(data)
	if err != nil {
		t.Fatal(err)
	}
	if msg.Trade == nil || msg.Trade.ID != "t" || msg.Order != nil {
		t.Fatalf("got %#v", msg)
	}
}

func TestParseUserChannelEvent_skipUnknown(t *testing.T) {
	msg, err := parseUserChannelEvent([]byte(`{"event_type":"other"}`))
	if err != nil || msg.Order != nil || msg.Trade != nil {
		t.Fatalf("msg=%+v err=%v", msg, err)
	}
	msg, err = parseUserChannelEvent([]byte(`not json`))
	if err != nil || msg.Order != nil || msg.Trade != nil {
		t.Fatalf("msg=%+v err=%v", msg, err)
	}
}
