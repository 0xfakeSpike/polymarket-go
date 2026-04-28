package main

import (
	"fmt"
	"os"

	"github.com/0xfakeSpike/polymarket-go/internal/mcpbridge"
)

// This command provides a tiny stdio bridge intended for MCP server adapters.
// Input: one JSON request per line. Output: one JSON response per line.
func main() {
	srv, err := mcpbridge.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	if err := srv.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
