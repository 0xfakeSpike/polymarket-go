// Package stdio implements the lightweight JSON-line stdio bridge used by cmd/polymarket-mcp.
package stdio

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/0xfakeSpike/polymarket-go"
	"github.com/0xfakeSpike/polymarket-go/internal/tools"
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

type Server struct {
	client *polymarket.Client
	in     io.Reader
	out    io.Writer
}

func New(in io.Reader, out io.Writer) (*Server, error) {
	var client *polymarket.Client
	var err error
	if pk := os.Getenv("POLYMARKET_MCP_PRIVATE_KEY"); pk != "" {
		client, err = polymarket.NewClient(pk)
	} else {
		client, err = polymarket.NewPublicClient()
	}
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
	data, err := tools.Call(s.client, req.Tool, req.Params)
	if err != nil {
		s.writeResp(response{OK: false, Error: err.Error()})
		return
	}
	s.writeResp(response{OK: true, Data: data})
}

func (s *Server) writeResp(v response) {
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Fprintf(s.out, "{\"ok\":false,\"error\":\"marshal response: %v\"}\n", err)
		return
	}
	fmt.Fprintln(s.out, string(b))
}
