package pmctl

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"github.com/0xfakeSpike/polymarket-go"
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

	switch args[0] {
	case "search-events":
		if err := a.runSearchEvents(args[1:]); err != nil {
			a.fail(err)
			return 1
		}
	case "orderbook":
		if err := a.runOrderbook(args[1:]); err != nil {
			a.fail(err)
			return 1
		}
	default:
		a.usage()
		return 2
	}

	return 0
}

func (a App) runSearchEvents(args []string) error {
	fs := flag.NewFlagSet("search-events", flag.ContinueOnError)
	fs.SetOutput(a.Stderr)
	query := fs.String("q", "", "search query")
	limit := fs.Int("limit", 10, "max events")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *query == "" {
		return fmt.Errorf("missing -q")
	}

	c, err := polymarket.NewPublicClient()
	if err != nil {
		return err
	}
	events, err := c.SearchEventsWithQuery(*query)
	if err != nil {
		return err
	}
	if *limit > 0 && len(events) > *limit {
		events = events[:*limit]
	}

	return a.printJSON(events)
}

func (a App) runOrderbook(args []string) error {
	fs := flag.NewFlagSet("orderbook", flag.ContinueOnError)
	fs.SetOutput(a.Stderr)
	tokenID := fs.String("token-id", "", "CLOB token id")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *tokenID == "" {
		return fmt.Errorf("missing -token-id")
	}

	c, err := polymarket.NewPublicClient()
	if err != nil {
		return err
	}
	book, err := c.GetOrderBook(*tokenID)
	if err != nil {
		return err
	}

	return a.printJSON(book)
}

func (a App) printJSON(v any) error {
	enc := json.NewEncoder(a.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func (a App) usage() {
	fmt.Fprintf(a.Stderr, `pmctl - Polymarket CLI

Usage:
  pmctl search-events -q "<query>" [-limit 10]
  pmctl orderbook -token-id "<token_id>"
`)
}

func (a App) fail(err error) {
	fmt.Fprintln(a.Stderr, "error:", err)
}
