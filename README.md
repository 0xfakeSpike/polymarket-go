# polymarket-go

`polymarket-go` is a Go SDK for Polymarket Gamma + CLOB APIs, with a repository layout ready for public use:

- reusable SDK public entrypoint in repository root package
- implementation package in `polymarket/` (kept for compatibility)
- executable CLI in `cmd/pmctl`
- MCP bridge command in `cmd/polymarket-mcp`
- MCP integration guide in `docs/mcp-integration.md`

## Project layout

```text
.
├── cmd/
│   └── pmctl/                # CLI entrypoint
│   └── polymarket-mcp/       # MCP stdio bridge entrypoint
├── internal/
│   ├── cli/pmctl/            # CLI app layer (business logic)
│   └── mcpbridge/            # MCP bridge runtime layer
├── examples/
│   ├── public-search/
│   └── orderbook/
├── docs/
│   └── mcp-integration.md    # MCP integration approach
├── polymarket/               # compatibility + implementation package
├── doc.go                    # root SDK package docs
├── sdk_bridge.go             # root SDK bridge exports
├── VERSIONING.md             # semver and deprecation rules
├── go.mod
└── README.md
```

## Install SDK

```bash
go get github.com/0xfakeSpike/polymarket-go
```

```go
package main

import (
  "fmt"

  "github.com/0xfakeSpike/polymarket-go"
)

func main() {
  c, err := polymarket.NewPublicClient()
  if err != nil {
    panic(err)
  }
  events, err := c.SearchEventsWithQuery("election")
  if err != nil {
    panic(err)
  }
  fmt.Println("events:", len(events))
}
```

## Use CLI

```bash
go run ./cmd/pmctl search-events -q "trump" -limit 5
go run ./cmd/pmctl orderbook -token-id "<CLOB_TOKEN_ID>"
```

Expose **any exported** `Client` method without writing a wrapper per command:

```bash
go run ./cmd/pmctl methods
go run ./cmd/pmctl methods -long
go run ./cmd/pmctl call GetOK
go run ./cmd/pmctl call GetOrderBook -args '["<CLOB_TOKEN_ID>"]'
```

Arguments are a **JSON array** in parameter order; `context.Context` is injected automatically. Methods that take **functions or handler interfaces** (for example WebSocket runners) cannot be used through `call` and should use the Go SDK instead.

## MCP integration

See `docs/mcp-integration.md` for how to expose this SDK as MCP tools for Cursor/Claude/Desktop clients.

### MCP bridge quickstart

```bash
echo '{"tool":"search_events","params":{"query":"election","limit":3}}' | \
  go run ./cmd/polymarket-mcp
```

## Examples

```bash
go run ./examples/public-search
go run ./examples/orderbook "<CLOB_TOKEN_ID>"
```

## Versioning

See `VERSIONING.md`.

## Contributing

See `CONTRIBUTING.md`.

## Changelog

See `CHANGELOG.md`.

## Release

See `RELEASE_CHECKLIST.md`.

## Homebrew

```bash
brew tap 0xfakeSpike/tap
brew install polymarket-go
brew install polymarket-mcp
```

Homebrew publishing uses a tap repo **`0xfakeSpike/homebrew-tap`**; create it first, then follow `docs/homebrew-release.md` (方案 A).
