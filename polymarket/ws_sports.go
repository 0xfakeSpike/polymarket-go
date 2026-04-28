package polymarket

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/coder/websocket"
)

const sportsPongPayload = "pong"

type sportsIncoming struct {
	typ  websocket.MessageType
	data []byte
	err  error
}

type sportsConn struct {
	conn *websocket.Conn

	readCtx    context.Context
	readCancel context.CancelFunc
	feed       chan sportsIncoming

	writeMu   sync.Mutex
	closeOnce sync.Once
}

func newSportsConn(ctx context.Context, dial *websocket.DialOptions) (*sportsConn, error) {
	if dial == nil {
		dial = &websocket.DialOptions{}
	}
	conn, _, err := websocket.Dial(ctx, SportsWebSocketURL, dial)
	if err != nil {
		return nil, fmt.Errorf("sports websocket dial: %w", err)
	}
	conn.SetReadLimit(clobWSMaxReadBytes)
	readCtx, readCancel := context.WithCancel(context.Background())
	s := &sportsConn{
		conn:       conn,
		readCtx:    readCtx,
		readCancel: readCancel,
		feed:       make(chan sportsIncoming, 256),
	}
	go s.readLoop()
	return s, nil
}

func (s *sportsConn) writeText(ctx context.Context, b []byte) error {
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	return s.conn.Write(ctx, websocket.MessageText, b)
}

func (s *sportsConn) readLoop() {
	defer close(s.feed)
	for {
		typ, data, err := s.conn.Read(s.readCtx)
		if err != nil {
			select {
			case s.feed <- sportsIncoming{err: err}:
			case <-s.readCtx.Done():
			}
			return
		}
		if typ == websocket.MessageText && bytes.Equal(bytes.TrimSpace(data), []byte("ping")) {
			_ = s.writeText(s.readCtx, []byte(sportsPongPayload))
			continue
		}
		select {
		case s.feed <- sportsIncoming{typ: typ, data: data}:
		case <-s.readCtx.Done():
			return
		}
	}
}

func (s *sportsConn) recv(ctx context.Context) (websocket.MessageType, []byte, error) {
	select {
	case <-ctx.Done():
		return 0, nil, ctx.Err()
	case m, ok := <-s.feed:
		if !ok {
			return 0, nil, io.EOF
		}
		if m.err != nil {
			return 0, nil, m.err
		}
		return m.typ, m.data, nil
	}
}

func (s *sportsConn) Close() error {
	if s == nil {
		return nil
	}
	var out error
	s.closeOnce.Do(func() {
		if s.readCancel != nil {
			s.readCancel()
		}
		if s.conn != nil {
			out = s.conn.Close(websocket.StatusNormalClosure, "")
		}
	})
	return out
}

// SportsChannelUpdate matches the sports WebSocket JSON payload (slug identifies the match).
type SportsChannelUpdate struct {
	Slug              string `json:"slug"`
	Live              bool   `json:"live,omitempty"`
	Ended             bool   `json:"ended,omitempty"`
	Score             string `json:"score,omitempty"`
	Period            string `json:"period,omitempty"`
	Elapsed           string `json:"elapsed,omitempty"`
	LastUpdate        string `json:"last_update,omitempty"`
	FinishedTimestamp string `json:"finished_timestamp,omitempty"`
	Turn              string `json:"turn,omitempty"`
}

// SportsMessage is one decoded sports update (non-empty slug).
type SportsMessage struct {
	Update SportsChannelUpdate
}

// SportsHandler receives each sports result message.
type SportsHandler func(m SportsMessage) error

// RunSportsWebSocket connects to the production sports channel, answers server pings automatically,
// and calls h for each decoded update until ctx is done or an error occurs.
func (c *Client) RunSportsWebSocket(ctx context.Context, dial *websocket.DialOptions, h SportsHandler) error {
	if c == nil {
		return fmt.Errorf("polymarket: nil Client")
	}
	if h == nil {
		return fmt.Errorf("polymarket: nil SportsHandler")
	}
	s, err := newSportsConn(ctx, dial)
	if err != nil {
		return err
	}
	defer s.Close()

	for {
		typ, data, err := s.recv(ctx)
		if err != nil {
			return err
		}
		if typ != websocket.MessageText {
			continue
		}
		var u SportsChannelUpdate
		if err := json.Unmarshal(data, &u); err != nil {
			return fmt.Errorf("sports websocket decode: %w", err)
		}
		if u.Slug == "" {
			continue
		}
		if err := h(SportsMessage{Update: u}); err != nil {
			return err
		}
	}
}
