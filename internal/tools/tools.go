// Package tools defines the small CLOB-v2 tool surface shared by CLI and MCP adapters.
package tools

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/0xfakeSpike/polymarket-go"
	"github.com/0xfakeSpike/polymarket-go/internal/tools/invoke"
)

// Tool is a JSON-parameterized adapter around a polymarket.Client method.
type Tool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ReadOnly    bool   `json:"read_only"`
	Run         func(*polymarket.Client, json.RawMessage) (any, error)
}

// Info is the serializable view of a Tool.
type Info struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ReadOnly    bool   `json:"read_only"`
}

type clientCallParams struct {
	Method string          `json:"method"`
	Args   json.RawMessage `json:"args"`
}

type methodsParams struct {
	Long bool `json:"long"`
}

// getMarketsByAnnualizedReturnToolParams are JSON object params for the named MCP/CLI tool.
type getMarketsByAnnualizedReturnToolParams struct {
	Limit        int     `json:"limit"`
	MaxPages     int     `json:"max_pages"`
	MinBestAsk   float64 `json:"min_best_ask"`
	NowRFC3339   string  `json:"now_rfc3339,omitempty"`
}

var registry = map[string]Tool{
	"client_call": {
		Name:        "client_call",
		Description: "Invoke any exported polymarket.Client method with JSON array arguments.",
		ReadOnly:    false,
		Run: func(c *polymarket.Client, raw json.RawMessage) (any, error) {
			var p clientCallParams
			if err := decodeParams(raw, &p); err != nil {
				return nil, err
			}
			if p.Method == "" {
				return nil, fmt.Errorf("missing method")
			}
			if len(p.Args) == 0 {
				p.Args = []byte("[]")
			}
			return invoke.Invoke(c, p.Method, p.Args)
		},
	},
	"get_markets_by_annualized_return": {
		Name:        "get_markets_by_annualized_return",
		Description: "Rank CLOB markets by annualized return from now until settlement (read-only; uses /markets + order books).",
		ReadOnly:    true,
		Run: func(c *polymarket.Client, raw json.RawMessage) (any, error) {
			var p getMarketsByAnnualizedReturnToolParams
			if err := decodeParams(raw, &p); err != nil {
				return nil, err
			}
			ap := &polymarket.AnnualizedReturnMarketsParams{
				Limit:      p.Limit,
				MaxPages:   p.MaxPages,
				MinBestAsk: p.MinBestAsk,
			}
			if s := p.NowRFC3339; s != "" {
				t, err := time.Parse(time.RFC3339Nano, s)
				if err != nil {
					t, err = time.Parse(time.RFC3339, s)
					if err != nil {
						return nil, fmt.Errorf("invalid now_rfc3339: %w", err)
					}
				}
				ap.Now = t
			}
			return c.GetMarketsByAnnualizedReturn(ap)
		},
	},
	"methods": {
		Name:        "methods",
		Description: "List exported polymarket.Client methods available to client_call.",
		ReadOnly:    true,
		Run: func(_ *polymarket.Client, raw json.RawMessage) (any, error) {
			var p methodsParams
			if err := decodeParams(raw, &p); err != nil {
				return nil, err
			}
			names := invoke.ListClientMethods()
			if !p.Long {
				return names, nil
			}
			type row struct {
				Name string `json:"name"`
				Sig  string `json:"sig,omitempty"`
			}
			out := make([]row, 0, len(names))
			for _, name := range names {
				sig, _ := invoke.MethodHelp(name)
				out = append(out, row{Name: name, Sig: sig})
			}
			return out, nil
		},
	},
}

// List returns tool metadata sorted by name.
func List() []Info {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	sort.Strings(names)

	out := make([]Info, 0, len(names))
	for _, name := range names {
		t := registry[name]
		out = append(out, Info{
			Name:        t.Name,
			Description: t.Description,
			ReadOnly:    t.ReadOnly,
		})
	}
	return out
}

// Call executes a registered tool with JSON object params.
func Call(client *polymarket.Client, name string, params json.RawMessage) (any, error) {
	tool, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown tool %q", name)
	}
	return tool.Run(client, params)
}

func decodeParams(raw json.RawMessage, dst any) error {
	if len(raw) == 0 || string(raw) == "null" {
		raw = []byte("{}")
	}
	if err := json.Unmarshal(raw, dst); err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}
	return nil
}
