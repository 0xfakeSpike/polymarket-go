package polymarket

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/coder/websocket"
)

type marketSubscribeWire struct {
	AssetsIDs            []string `json:"assets_ids"`
	Type                 string   `json:"type"`
	InitialDump          *bool    `json:"initial_dump,omitempty"`
	Level                *int     `json:"level,omitempty"`
	CustomFeatureEnabled *bool    `json:"custom_feature_enabled,omitempty"`
}

type marketIncoming struct {
	typ  websocket.MessageType
	data []byte
	err  error
}

type marketConn struct {
	conn *websocket.Conn

	readCtx    context.Context
	readCancel context.CancelFunc
	feed       chan marketIncoming

	writeMu   sync.Mutex
	closeOnce sync.Once
	startOnce sync.Once
}

func newMarketConn(ctx context.Context, dial *websocket.DialOptions) (*marketConn, error) {
	if dial == nil {
		dial = &websocket.DialOptions{}
	}
	conn, _, err := websocket.Dial(ctx, MarketWebSocketURL, dial)
	if err != nil {
		return nil, fmt.Errorf("market websocket dial: %w", err)
	}
	conn.SetReadLimit(clobWSMaxReadBytes)
	readCtx, readCancel := context.WithCancel(context.Background())
	m := &marketConn{
		conn:       conn,
		readCtx:    readCtx,
		readCancel: readCancel,
		feed:       make(chan marketIncoming, 256),
	}
	return m, nil
}

func (m *marketConn) startReadAndPing() {
	m.startOnce.Do(func() {
		go m.readLoop()
		go m.clientPingLoop()
	})
}

func (m *marketConn) writeText(ctx context.Context, b []byte) error {
	m.writeMu.Lock()
	defer m.writeMu.Unlock()
	return m.conn.Write(ctx, websocket.MessageText, b)
}

func (m *marketConn) readLoop() {
	defer close(m.feed)
	for {
		typ, data, err := m.conn.Read(m.readCtx)
		if err != nil {
			select {
			case m.feed <- marketIncoming{err: err}:
			case <-m.readCtx.Done():
			}
			return
		}
		if typ == websocket.MessageText {
			s := bytes.TrimSpace(data)
			if bytes.Equal(s, []byte("PONG")) || bytes.Equal(s, []byte("PING")) {
				continue
			}
		}
		select {
		case m.feed <- marketIncoming{typ: typ, data: data}:
		case <-m.readCtx.Done():
			return
		}
	}
}

func (m *marketConn) clientPingLoop() {
	_ = m.writeText(m.readCtx, []byte(CLOBWebSocketPing))
	ticker := time.NewTicker(clientPingInterval)
	defer ticker.Stop()
	for {
		select {
		case <-m.readCtx.Done():
			return
		case <-ticker.C:
			_ = m.writeText(m.readCtx, []byte(CLOBWebSocketPing))
		}
	}
}

func (m *marketConn) sendSubscribe(ctx context.Context, assetIDs []string, initialDump *bool, level *int, customFeature *bool) error {
	if len(assetIDs) == 0 {
		return fmt.Errorf("market websocket: assets_ids required")
	}
	msg := marketSubscribeWire{
		AssetsIDs:            assetIDs,
		Type:                 "market",
		InitialDump:          initialDump,
		Level:                level,
		CustomFeatureEnabled: customFeature,
	}
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return m.writeText(ctx, b)
}

func (m *marketConn) recv(ctx context.Context) (websocket.MessageType, []byte, error) {
	select {
	case <-ctx.Done():
		return 0, nil, ctx.Err()
	case msg, ok := <-m.feed:
		if !ok {
			return 0, nil, io.EOF
		}
		if msg.err != nil {
			return 0, nil, msg.err
		}
		return msg.typ, msg.data, nil
	}
}

func (m *marketConn) Close() error {
	if m == nil {
		return nil
	}
	m.closeOnce.Do(func() {
		if m.readCancel != nil {
			m.readCancel()
		}
		if m.conn != nil {
			_ = m.conn.Close(websocket.StatusNormalClosure, "")
		}
	})
	return nil
}

// MarketOptions are optional fields for the initial market subscription (nil pointers omit JSON keys).
type MarketOptions struct {
	InitialDump          *bool
	Level                *int
	CustomFeatureEnabled *bool
}

// MarketHandler receives each decoded market-channel event (see [MarketChannelMessage]).
type MarketHandler func(m MarketChannelMessage) error

func parseMarketChannelMessage(data []byte) (MarketChannelMessage, error) {
	data = unwrapJSONArrayWSMessage(data)
	var head struct {
		EventType string `json:"event_type"`
	}
	if err := json.Unmarshal(data, &head); err != nil {
		return MarketChannelMessage{}, nil
	}
	switch head.EventType {
	case "book":
		var ev MarketBookEvent
		if err := json.Unmarshal(data, &ev); err != nil {
			return MarketChannelMessage{}, fmt.Errorf("market book: %w", err)
		}
		return MarketChannelMessage{Book: &ev}, nil
	case "price_change":
		var ev MarketPriceChangeEvent
		if err := json.Unmarshal(data, &ev); err != nil {
			return MarketChannelMessage{}, fmt.Errorf("market price_change: %w", err)
		}
		return MarketChannelMessage{PriceChange: &ev}, nil
	case "last_trade_price":
		var ev MarketLastTradePriceEvent
		if err := json.Unmarshal(data, &ev); err != nil {
			return MarketChannelMessage{}, fmt.Errorf("market last_trade_price: %w", err)
		}
		return MarketChannelMessage{LastTradePrice: &ev}, nil
	case "tick_size_change":
		var ev MarketTickSizeChangeEvent
		if err := json.Unmarshal(data, &ev); err != nil {
			return MarketChannelMessage{}, fmt.Errorf("market tick_size_change: %w", err)
		}
		return MarketChannelMessage{TickSizeChange: &ev}, nil
	case "best_bid_ask":
		var ev MarketBestBidAskEvent
		if err := json.Unmarshal(data, &ev); err != nil {
			return MarketChannelMessage{}, fmt.Errorf("market best_bid_ask: %w", err)
		}
		return MarketChannelMessage{BestBidAsk: &ev}, nil
	case "new_market":
		var ev MarketNewMarketEvent
		if err := json.Unmarshal(data, &ev); err != nil {
			return MarketChannelMessage{}, fmt.Errorf("market new_market: %w", err)
		}
		return MarketChannelMessage{NewMarket: &ev}, nil
	case "market_resolved":
		var ev MarketResolvedEvent
		if err := json.Unmarshal(data, &ev); err != nil {
			return MarketChannelMessage{}, fmt.Errorf("market market_resolved: %w", err)
		}
		return MarketChannelMessage{Resolved: &ev}, nil
	default:
		return MarketChannelMessage{}, nil
	}
}

// RunMarketWebSocket connects to the production market channel, subscribes to assetIDs,
// maintains PING/PONG in the background, and calls h for each decoded event until ctx is done or an error occurs.
// Frames that are not recognized market events are skipped.
func (c *Client) RunMarketWebSocket(ctx context.Context, dial *websocket.DialOptions, assetIDs []string, opts MarketOptions, h MarketHandler) error {
	if c == nil {
		return fmt.Errorf("polymarket: nil Client")
	}
	if h == nil {
		return fmt.Errorf("polymarket: nil MarketHandler")
	}
	m, err := newMarketConn(ctx, dial)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.sendSubscribe(ctx, assetIDs, opts.InitialDump, opts.Level, opts.CustomFeatureEnabled); err != nil {
		return err
	}
	m.startReadAndPing()
	for {
		typ, data, err := m.recv(ctx)
		if err != nil {
			return err
		}
		if typ != websocket.MessageText {
			continue
		}
		msg, perr := parseMarketChannelMessage(data)
		if perr != nil {
			return perr
		}
		if msg.empty() {
			continue
		}
		if err := h(msg); err != nil {
			return err
		}
	}
}
