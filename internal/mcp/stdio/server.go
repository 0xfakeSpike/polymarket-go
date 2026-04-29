// Package stdio implements an MCP stdio server for cmd/polymarket-mcp.
package stdio

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/0xfakeSpike/polymarket-go"
	"github.com/0xfakeSpike/polymarket-go/internal/tools"
)

const mcpProtocolVersion = "2024-11-05"

type rpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type rpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Result  any             `json:"result,omitempty"`
	Error   *rpcErrorObject `json:"error,omitempty"`
}

type rpcErrorObject struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type toolsListResult struct {
	Tools []toolDescriptor `json:"tools"`
}

type toolDescriptor struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	InputSchema map[string]any `json:"inputSchema"`
}

type toolsCallParams struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

type initializeResult struct {
	ProtocolVersion string         `json:"protocolVersion"`
	Capabilities    map[string]any `json:"capabilities"`
	ServerInfo      serverInfo     `json:"serverInfo"`
}

type serverInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type toolsCallResult struct {
	Content []contentItem `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

type contentItem struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Server struct {
	client *polymarket.Client
	in     *bufio.Reader
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
	return &Server{client: client, in: bufio.NewReader(in), out: out}, nil
}

func (s *Server) Run() error {
	for {
		raw, err := s.readFrame()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		var req rpcRequest
		if err := json.Unmarshal(raw, &req); err != nil {
			s.writeResponse(rpcResponse{
				JSONRPC: "2.0",
				Error:   &rpcErrorObject{Code: -32700, Message: fmt.Sprintf("parse error: %v", err)},
			})
			continue
		}
		s.handle(req)
	}
}

func (s *Server) handle(req rpcRequest) {
	if req.JSONRPC != "2.0" {
		s.writeError(req.ID, -32600, "invalid request: jsonrpc must be 2.0")
		return
	}

	switch req.Method {
	case "initialize":
		if len(req.ID) == 0 {
			return
		}
		s.writeResponse(rpcResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: initializeResult{
				ProtocolVersion: mcpProtocolVersion,
				Capabilities: map[string]any{
					"tools": map[string]any{},
				},
				ServerInfo: serverInfo{
					Name:    "polymarket-mcp",
					Version: "1.0.0",
				},
			},
		})
	case "notifications/initialized":
		// Notification from client after initialize; no response.
		return
	case "ping":
		if len(req.ID) == 0 {
			return
		}
		s.writeResponse(rpcResponse{JSONRPC: "2.0", ID: req.ID, Result: map[string]any{}})
	case "tools/list":
		if len(req.ID) == 0 {
			return
		}
		s.writeResponse(rpcResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  toolsListResult{Tools: listToolDescriptors()},
		})
	case "tools/call":
		if len(req.ID) == 0 {
			return
		}
		var p toolsCallParams
		if err := json.Unmarshal(req.Params, &p); err != nil {
			s.writeError(req.ID, -32602, fmt.Sprintf("invalid params: %v", err))
			return
		}
		if p.Name == "" {
			s.writeError(req.ID, -32602, "invalid params: missing tool name")
			return
		}
		if len(p.Arguments) == 0 || string(p.Arguments) == "null" {
			p.Arguments = []byte("{}")
		}
		data, err := tools.Call(s.client, p.Name, p.Arguments)
		if err != nil {
			b, _ := json.Marshal(map[string]any{"error": err.Error()})
			s.writeResponse(rpcResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result: toolsCallResult{
					Content: []contentItem{{Type: "text", Text: string(b)}},
					IsError: true,
				},
			})
			return
		}
		out, err := json.Marshal(data)
		if err != nil {
			s.writeError(req.ID, -32603, fmt.Sprintf("marshal tool result: %v", err))
			return
		}
		s.writeResponse(rpcResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: toolsCallResult{
				Content: []contentItem{{Type: "text", Text: string(out)}},
			},
		})
	default:
		if len(req.ID) == 0 {
			return
		}
		s.writeError(req.ID, -32601, fmt.Sprintf("method not found: %s", req.Method))
	}
}

func listToolDescriptors() []toolDescriptor {
	infos := tools.List()
	out := make([]toolDescriptor, 0, len(infos))
	for _, info := range infos {
		out = append(out, toolDescriptor{
			Name:        info.Name,
			Description: info.Description,
			InputSchema: toolSchema(info.Name),
		})
	}
	return out
}

func toolSchema(name string) map[string]any {
	switch name {
	case "search_events":
		return objectSchema(
			map[string]any{"query": map[string]any{"type": "string"}, "limit": map[string]any{"type": "integer", "minimum": 1}},
			[]string{"query"},
		)
	case "get_orderbook":
		return objectSchema(
			map[string]any{"token_id": map[string]any{"type": "string"}},
			[]string{"token_id"},
		)
	case "rank_markets_by_annualized_return":
		return objectSchema(
			map[string]any{
				"tag_slug":              map[string]any{"type": "string"},
				"keyword":               map[string]any{"type": "string"},
				"events_limit":          map[string]any{"type": "integer", "minimum": 1},
				"limit":                 map[string]any{"type": "integer", "minimum": 1},
				"min_annualized_return": map[string]any{"type": "number"},
			},
			nil,
		)
	case "methods":
		return objectSchema(
			map[string]any{"long": map[string]any{"type": "boolean"}},
			nil,
		)
	case "client_call":
		return objectSchema(
			map[string]any{
				"method": map[string]any{"type": "string"},
				"args":   map[string]any{"type": "array"},
			},
			[]string{"method"},
		)
	default:
		return objectSchema(nil, nil)
	}
}

func objectSchema(properties map[string]any, required []string) map[string]any {
	s := map[string]any{
		"type":                 "object",
		"additionalProperties": false,
	}
	if properties == nil {
		properties = map[string]any{}
	}
	s["properties"] = properties
	if len(required) > 0 {
		s["required"] = required
	}
	return s
}

func (s *Server) writeError(id json.RawMessage, code int, message string) {
	s.writeResponse(rpcResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error:   &rpcErrorObject{Code: code, Message: message},
	})
}

func (s *Server) readFrame() ([]byte, error) {
	contentLength := -1
	for {
		line, err := s.in.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(strings.ToLower(parts[0]))
		value := strings.TrimSpace(parts[1])
		if key == "content-length" {
			n, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("invalid content-length: %w", err)
			}
			contentLength = n
		}
	}
	if contentLength < 0 {
		return nil, fmt.Errorf("missing content-length header")
	}
	body := make([]byte, contentLength)
	if _, err := io.ReadFull(s.in, body); err != nil {
		return nil, err
	}
	return body, nil
}

func (s *Server) writeResponse(v rpcResponse) {
	b, err := json.Marshal(v)
	if err != nil {
		return
	}
	if _, err := fmt.Fprintf(s.out, "Content-Length: %d\r\n\r\n%s", len(b), b); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "mcp write response: %v\n", err)
		return
	}
}
