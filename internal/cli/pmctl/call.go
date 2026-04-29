package pmctl

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/0xfakeSpike/polymarket-go/internal/clientcall"
)

func (a App) runMethods(args []string) error {
	fs := flag.NewFlagSet("methods", flag.ContinueOnError)
	fs.SetOutput(a.Stderr)
	long := fs.Bool("long", false, "include one-line reflect signatures")
	if err := fs.Parse(args); err != nil {
		return err
	}
	names := clientcall.ListClientMethods()
	if !*long {
		return a.printJSON(names)
	}
	type row struct {
		Name string `json:"name"`
		Sig  string `json:"sig,omitempty"`
	}
	out := make([]row, 0, len(names))
	for _, n := range names {
		sig, err := clientcall.MethodHelp(n)
		if err != nil {
			sig = ""
		}
		out = append(out, row{Name: n, Sig: sig})
	}
	return a.printJSON(out)
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
	if len(fs.Args()) < 1 {
		return fmt.Errorf("usage: pmctl call [flags] <MethodName>\nexample: pmctl call GetOK\nexample: pmctl call GetOrderBook -args '[\"TOKEN_ID\"]'")
	}
	method := fs.Args()[0]

	result, err := clientcall.Invoke(c, method, raw)
	if err != nil {
		return err
	}
	return a.printJSON(result)
}
