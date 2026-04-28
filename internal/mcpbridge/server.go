package mcpbridge

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"

	"github.com/0xfakeSpike/polymarket-go"
)

type request struct {
	Tool   string          `json:"tool"`
	Params json.RawMessage `json:"params"`
}

type response struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
}

type searchEventsParams struct {
	Query string `json:"query"`
	Limit int    `json:"limit"`
}

type orderbookParams struct {
	TokenID string `json:"token_id"`
}

type Server struct {
	client *polymarket.Client
	in     io.Reader
	out    io.Writer
}

func New(in io.Reader, out io.Writer) (*Server, error) {
	client, err := polymarket.NewPublicClient()
	if err != nil {
		return nil, err
	}
	return &Server{client: client, in: in, out: out}, nil
}

func (s *Server) Run() error {
	scanner := bufio.NewScanner(s.in)
	for scanner.Scan() {
		var req request
		if err := json.Unmarshal(scanner.Bytes(), &req); err != nil {
			s.writeResp(response{OK: false, Error: fmt.Sprintf("invalid request: %v", err)})
			continue
		}
		s.handle(req)
	}
	return scanner.Err()
}

func (s *Server) handle(req request) {
	switch req.Tool {
	case "search_events":
		var p searchEventsParams
		if err := json.Unmarshal(req.Params, &p); err != nil {
			s.writeResp(response{OK: false, Error: err.Error()})
			return
		}
		if p.Limit <= 0 {
			p.Limit = 10
		}
		events, err := s.client.SearchEventsWithQuery(p.Query)
		if err != nil {
			s.writeResp(response{OK: false, Error: err.Error()})
			return
		}
		if len(events) > p.Limit {
			events = events[:p.Limit]
		}
		s.writeResp(response{OK: true, Data: events})
	case "get_orderbook":
		var p orderbookParams
		if err := json.Unmarshal(req.Params, &p); err != nil {
			s.writeResp(response{OK: false, Error: err.Error()})
			return
		}
		book, err := s.client.GetOrderBook(p.TokenID)
		if err != nil {
			s.writeResp(response{OK: false, Error: err.Error()})
			return
		}
		s.writeResp(response{OK: true, Data: book})
	default:
		s.writeResp(response{OK: false, Error: "unknown tool"})
	}
}

func (s *Server) writeResp(v response) {
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Fprintf(s.out, "{\"ok\":false,\"error\":\"marshal response: %v\"}\n", err)
		return
	}
	fmt.Fprintln(s.out, string(b))
}
