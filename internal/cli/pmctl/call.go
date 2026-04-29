package pmctl

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/0xfakeSpike/polymarket-go/internal/tools"
)

func (a App) runMethods(args []string) error {
	fs := flag.NewFlagSet("methods", flag.ContinueOnError)
	fs.SetOutput(a.Stderr)
	long := fs.Bool("long", false, "include one-line reflect signatures")
	if err := fs.Parse(args); err != nil {
		return err
	}
	result, err := tools.Call(nil, "methods", mustMarshalJSON(map[string]any{
		"long": *long,
	}))
	if err != nil {
		return err
	}
	return a.printJSON(result)
}

func (a App) runCall(args []string) error {
	fs := flag.NewFlagSet("call", flag.ContinueOnError)
	fs.SetOutput(a.Stderr)
	public := fs.Bool("public", true, "use public read-only client (no private key)")
	pk := fs.String("private-key", "", "hex private key for authenticated client (or env PMCTL_PRIVATE_KEY)")
	argsJSON := fs.String("args", "[]", "JSON array of arguments in parameter order (context.Context is injected and omitted)")
	argsFile := fs.String("args-file", "", "read JSON args from file (overrides -args)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	raw := []byte(*argsJSON)
	if *argsFile != "" {
		b, err := os.ReadFile(*argsFile)
		if err != nil {
			return err
		}
		raw = b
	}
	raw = []byte(strings.TrimSpace(string(raw)))
	if len(raw) == 0 {
		raw = []byte("[]")
	}

	c, err := newClientFromFlags(*public, *pk)
	if err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("usage: pmctl call [flags] <MethodName>\nexample: pmctl call GetOK\nexample: pmctl call -args '[\"TOKEN_ID\"]' GetOrderBook")
	}
	method := fs.Args()[0]

	params, err := json.Marshal(map[string]any{
		"method": method,
		"args":   json.RawMessage(raw),
	})
	if err != nil {
		return err
	}
	result, err := tools.Call(c, "client_call", params)
	if err != nil {
		return err
	}
	return a.printJSON(result)
}
