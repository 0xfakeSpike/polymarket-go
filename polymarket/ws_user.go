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

const clientPingInterval = 10 * time.Second

type userSubscribeWire struct {
	Auth    APIKeyCredentials `json:"auth"`
	Type    string            `json:"type"`
	Markets []string          `json:"markets,omitempty"`
}

type userIncoming struct {
	typ  websocket.MessageType
	data []byte
	err  error
}

type userConn struct {
	conn *websocket.Conn

	readCtx    context.Context
	readCancel context.CancelFunc
	feed       chan userIncoming

	writeMu   sync.Mutex
	closeOnce sync.Once
	startOnce sync.Once
}

func newUserConn(ctx context.Context, dial *websocket.DialOptions) (*userConn, error) {
	if dial == nil {
		dial = &websocket.DialOptions{}
	}
	conn, _, err := websocket.Dial(ctx, UserWebSocketURL, dial)
	if err != nil {
		return nil, fmt.Errorf("user websocket dial: %w", err)
	}
	conn.SetReadLimit(clobWSMaxReadBytes)
	readCtx, readCancel := context.WithCancel(context.Background())
	u := &userConn{
		conn:       conn,
		readCtx:    readCtx,
		readCancel: readCancel,
		feed:       make(chan userIncoming, 256),
	}
	return u, nil
}

func (u *userConn) startReadAndPing() {
	u.startOnce.Do(func() {
		go u.readLoop()
		go u.clientPingLoop()
	})
}

func (u *userConn) writeText(ctx context.Context, b []byte) error {
	u.writeMu.Lock()
	defer u.writeMu.Unlock()
	return u.conn.Write(ctx, websocket.MessageText, b)
}

func (u *userConn) readLoop() {
	defer close(u.feed)
	for {
		typ, data, err := u.conn.Read(u.readCtx)
		if err != nil {
			select {
			case u.feed <- userIncoming{err: err}:
			case <-u.readCtx.Done():
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
		case u.feed <- userIncoming{typ: typ, data: data}:
		case <-u.readCtx.Done():
			return
		}
	}
}

func (u *userConn) clientPingLoop() {
	_ = u.writeText(u.readCtx, []byte(CLOBWebSocketPing))
	t := time.NewTicker(clientPingInterval)
	defer t.Stop()
	for {
		select {
		case <-u.readCtx.Done():
			return
		case <-t.C:
			_ = u.writeText(u.readCtx, []byte(CLOBWebSocketPing))
		}
	}
}

func (u *userConn) sendSubscribe(ctx context.Context, auth APIKeyCredentials, conditionIDs []string) error {
	if auth.ApiKey == "" || auth.Secret == "" || auth.Passphrase == "" {
		return fmt.Errorf("user websocket: incomplete auth")
	}
	msg := userSubscribeWire{Auth: auth, Type: "user", Markets: conditionIDs}
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return u.writeText(ctx, b)
}

func (u *userConn) recv(ctx context.Context) (websocket.MessageType, []byte, error) {
	select {
	case <-ctx.Done():
		return 0, nil, ctx.Err()
	case m, ok := <-u.feed:
		if !ok {
			return 0, nil, io.EOF
		}
		if m.err != nil {
			return 0, nil, m.err
		}
		return m.typ, m.data, nil
	}
}

func (u *userConn) Close() error {
	if u == nil {
		return nil
	}
	u.closeOnce.Do(func() {
		if u.readCancel != nil {
			u.readCancel()
		}
		if u.conn != nil {
			_ = u.conn.Close(websocket.StatusNormalClosure, "")
		}
	})
	return nil
}

// UserChannelOrderEvent is a server "order" event on the user channel.
type UserChannelOrderEvent struct {
	EventType       string   `json:"event_type"`
	ID              string   `json:"id"`
	Owner           string   `json:"owner"`
	Market          string   `json:"market"`
	AssetID         string   `json:"asset_id"`
	Side            string   `json:"side"`
	OrderOwner      string   `json:"order_owner,omitempty"`
	OriginalSize    string   `json:"original_size"`
	SizeMatched     string   `json:"size_matched"`
	Price           string   `json:"price"`
	AssociateTrades []string `json:"associate_trades,omitempty"`
	Outcome         string   `json:"outcome,omitempty"`
	Type            string   `json:"type"`
	CreatedAt       string   `json:"created_at,omitempty"`
	Expiration      string   `json:"expiration,omitempty"`
	OrderType       string   `json:"order_type,omitempty"`
	Status          string   `json:"status,omitempty"`
	MakerAddress    string   `json:"maker_address,omitempty"`
	Timestamp       string   `json:"timestamp"`
}

// UserChannelTradeMakerOrder is an element of UserChannelTradeEvent.MakerOrders.
type UserChannelTradeMakerOrder struct {
	OrderID       string `json:"order_id"`
	Owner         string `json:"owner"`
	MakerAddress  string `json:"maker_address,omitempty"`
	MatchedAmount string `json:"matched_amount"`
	Price         string `json:"price"`
	FeeRateBps    string `json:"fee_rate_bps,omitempty"`
	AssetID       string `json:"asset_id"`
	Outcome       string `json:"outcome,omitempty"`
	Side          string `json:"side,omitempty"`
}

// UserChannelTradeEvent is a server "trade" event on the user channel.
type UserChannelTradeEvent struct {
	EventType       string                       `json:"event_type"`
	Type            string                       `json:"type"`
	ID              string                       `json:"id"`
	TakerOrderID    string                       `json:"taker_order_id"`
	Market          string                       `json:"market"`
	AssetID         string                       `json:"asset_id"`
	Side            string                       `json:"side"`
	Size            string                       `json:"size"`
	Price           string                       `json:"price"`
	FeeRateBps      string                       `json:"fee_rate_bps,omitempty"`
	Status          string                       `json:"status"`
	Matchtime       string                       `json:"matchtime,omitempty"`
	LastUpdate      string                       `json:"last_update,omitempty"`
	Outcome         string                       `json:"outcome,omitempty"`
	Owner           string                       `json:"owner"`
	TradeOwner      string                       `json:"trade_owner,omitempty"`
	MakerAddress    string                       `json:"maker_address,omitempty"`
	TransactionHash string                       `json:"transaction_hash,omitempty"`
	BucketIndex     int                          `json:"bucket_index,omitempty"`
	MakerOrders     []UserChannelTradeMakerOrder `json:"maker_orders,omitempty"`
	TraderSide      string                       `json:"trader_side,omitempty"`
	Timestamp       string                       `json:"timestamp"`
}

// UserChannelMessage is one parsed user-channel payload: exactly one of Order or Trade is non-nil.
type UserChannelMessage struct {
	Order *UserChannelOrderEvent
	Trade *UserChannelTradeEvent
}

// UserChannelHandler receives each order or trade from the authenticated user WebSocket.
type UserChannelHandler func(m UserChannelMessage) error

// RunUserWebSocket connects to the production user channel, subscribes with this client's L2 API key,
// maintains PING/PONG in the background, and calls h for each order or trade event until ctx is done or an error occurs.
// Unknown JSON frames are ignored. Cancel ctx to stop.
func (c *Client) RunUserWebSocket(ctx context.Context, dial *websocket.DialOptions, conditionIDs []string, h UserChannelHandler) error {
	if h == nil {
		return fmt.Errorf("polymarket: nil UserChannelHandler")
	}
	if err := c.requireL2(); err != nil {
		return err
	}
	u, err := newUserConn(ctx, dial)
	if err != nil {
		return err
	}
	defer u.Close()

	if err := u.sendSubscribe(ctx, *c.apiKeyCredentials, conditionIDs); err != nil {
		return err
	}
	u.startReadAndPing()

	for {
		typ, data, err := u.recv(ctx)
		if err != nil {
			return err
		}
		if typ != websocket.MessageText {
			continue
		}
		msg, perr := parseUserChannelEvent(data)
		if perr != nil {
			return perr
		}
		if msg.Order == nil && msg.Trade == nil {
			continue
		}
		if err := h(msg); err != nil {
			return err
		}
	}
}

func parseUserChannelEvent(data []byte) (UserChannelMessage, error) {
	data = unwrapJSONArrayWSMessage(data)
	var head struct {
		EventType string `json:"event_type"`
	}
	if err := json.Unmarshal(data, &head); err != nil {
		return UserChannelMessage{}, nil
	}
	switch head.EventType {
	case "order":
		var ev UserChannelOrderEvent
		if err := json.Unmarshal(data, &ev); err != nil {
			return UserChannelMessage{}, fmt.Errorf("user channel order event: %w", err)
		}
		return UserChannelMessage{Order: &ev}, nil
	case "trade":
		var ev UserChannelTradeEvent
		if err := json.Unmarshal(data, &ev); err != nil {
			return UserChannelMessage{}, fmt.Errorf("user channel trade event: %w", err)
		}
		return UserChannelMessage{Trade: &ev}, nil
	default:
		return UserChannelMessage{}, nil
	}
}
