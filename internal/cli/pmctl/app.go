package pmctl

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/0xfakeSpike/polymarket-go/internal/tools"
)

type App struct {
	Stdout io.Writer
	Stderr io.Writer
}

func (a App) Run(args []string) int {
	if len(args) < 1 {
		a.usage()
		return 2
	}

	cmds := map[string]func([]string) error{
		"tools":   a.runTools,
		"tool":    a.runTool,
		"methods": a.runMethods,
		"call":    a.runCall,
	}
	run, ok := cmds[args[0]]
	if !ok {
		a.usage()
		return 2
	}
	if err := run(args[1:]); err != nil {
		a.fail(err)
		return 1
	}
	return 0
}

func (a App) runTools(args []string) error {
	fs := flag.NewFlagSet("tools", flag.ContinueOnError)
	fs.SetOutput(a.Stderr)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		return fmt.Errorf("usage: pmctl tools")
	}
	return a.printJSON(tools.List())
}

func (a App) runTool(args []string) error {
	fs := flag.NewFlagSet("tool", flag.ContinueOnError)
	fs.SetOutput(a.Stderr)
	public := fs.Bool("public", true, "use public read-only client (no private key)")
	pk := fs.String("private-key", "", "hex private key for authenticated client (or env PMCTL_PRIVATE_KEY)")
	paramsJSON := fs.String("params", "{}", "JSON object tool params")
	paramsFile := fs.String("params-file", "", "read JSON params from file (overrides -params)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("usage: pmctl tool [flags] <tool_name>")
	}

	raw := []byte(*paramsJSON)
	if *paramsFile != "" {
		b, err := os.ReadFile(*paramsFile)
		if err != nil {
			return err
		}
		raw = b
	}
	raw = []byte(strings.TrimSpace(string(raw)))
	if len(raw) == 0 {
		raw = []byte("{}")
	}

	c, err := newClientFromFlags(*public, *pk)
	if err != nil {
		return err
	}
	result, err := tools.Call(c, fs.Arg(0), raw)
	if err != nil {
		return err
	}
	return a.printJSON(result)
}

func (a App) printJSON(v any) error {
	enc := json.NewEncoder(a.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func (a App) usage() {
	fmt.Fprintf(a.Stderr, `pmctl - Polymarket CLI

Usage:
  pmctl tools
  pmctl tool [flags] <tool_name>
  pmctl methods [-long]
  pmctl call [flags] <ClientMethod>   # JSON array args; see "pmctl methods -long"

Examples:
  pmctl tool -params '{"long":true}' methods
  pmctl call GetOK
  pmctl call -args '["<token_id>"]' GetOrderBook
  pmctl call -args '["<condition_id>"]' GetClobMarketInfo
  pmctl call -args '[{"limit":10,"max_pages":3,"min_best_ask":0.5}]' GetMarketsByAnnualizedReturn
  pmctl call -public=false -private-key "$PMCTL_PRIVATE_KEY" -args '[...]' CreateOrder
`)
}

func (a App) fail(err error) {
	fmt.Fprintln(a.Stderr, "error:", err)
}

func mustMarshalJSON(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
