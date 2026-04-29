package main

import (
	"fmt"
	"os"

	"github.com/0xfakeSpike/polymarket-go/internal/mcp/stdio"
)

// This command serves a standards-compatible MCP server over stdio.
func main() {
	srv, err := stdio.New(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	if err := srv.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
