# polymarket-go

Go client for [Polymarket](https://polymarket.com) **Gamma** (markets, search) and **CLOB** (order book, trading), plus **`pmctl`** (CLI) and **`polymarket-mcp`** (JSON-line stdio bridge for MCP hosts).

## Features

- **SDK** — import `github.com/0xfakeSpike/polymarket-go`; `Client` covers Gamma, Data API, CLOB, RFQ, and helpers aligned with common Polymarket client usage.
- **CLI (`pmctl`)** — named tools with JSON params, optional reflection **`call`** for any exported `Client` method, and **`methods`** to list signatures.
- **MCP bridge** — same tool registry as the CLI over stdin/stdout; optional authenticated client via env.
- **Examples** — under `examples/`.

## Requirements

- **Go** `1.24+` (see `go.mod`).

## Install — library

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

Prefer the **root import** above. The implementation also lives under `github.com/0xfakeSpike/polymarket-go/polymarket` for compatibility; new code should use the root path.

## Install — binaries

### Homebrew

```bash
brew tap 0xfakeSpike/tap
brew install polymarket-go polymarket-mcp
```

Tap setup and release automation: **[docs/homebrew-release.md](docs/homebrew-release.md)**.

### From source

```bash
go install github.com/0xfakeSpike/polymarket-go/cmd/pmctl@latest
go install github.com/0xfakeSpike/polymarket-go/cmd/polymarket-mcp@latest
```

## CLI (`pmctl`)

| Command | Purpose |
|--------|---------|
| `pmctl tools` | JSON list of registered tools (name, description, `read_only`). |
| `pmctl tool [flags] <name>` | Run one tool; `-params` is a JSON **object** (see [docs/mcp-integration.md](docs/mcp-integration.md)). |
| `pmctl methods [-long]` | List exported `Client` method names; `-long` adds reflect signatures (for `call`). |
| `pmctl call [flags] <Method>` | Call any exported `Client` method; `-args` is a JSON **array** in parameter order. |

**Client mode flags** (for `tool` and `call`): `-public` (default `true`), or `-public=false` with `-private-key` / **`PMCTL_PRIVATE_KEY`**.

Examples:

```bash
pmctl tools
pmctl tool -params '{"query":"election","limit":5}' search_events
pmctl tool -params '{"token_id":"<CLOB_TOKEN_ID>"}' get_orderbook
pmctl tool -params '{"tag_slug":"geopolitics","keyword":"iran","limit":10,"min_annualized_return":0.25}' rank_markets_by_annualized_return
pmctl methods -long | head -20
pmctl call GetOK
pmctl call -args '["<CLOB_TOKEN_ID>"]' GetOrderBook
```

`call` injects `context.Context` where needed; methods that take **functions** or **handler interfaces** (e.g. WebSocket runners) are not supported through reflection — use the SDK in Go.

## MCP (`polymarket-mcp`)

One JSON object per input line; one JSON response per line. Specification and tool schemas: **[docs/mcp-integration.md](docs/mcp-integration.md)**.

**Environment**

| Variable | Effect |
|----------|--------|
| `POLYMARKET_MCP_PRIVATE_KEY` | If set, `NewClient` (trading / L2 bootstrap). If unset, public client only. |

## Examples

```bash
go run ./examples/public-search
go run ./examples/orderbook "<CLOB_TOKEN_ID>"
```

## Repository layout

```text
cmd/pmctl              CLI entrypoint
cmd/polymarket-mcp     MCP stdio entrypoint
internal/cli/pmctl   CLI wiring (flags, stdout/stderr)
internal/mcp/stdio   MCP JSON-line server
internal/tools       Shared tool registry (search_events, get_orderbook, rank_markets_by_annualized_return, methods, client_call)
internal/tools/invoke Reflection helpers for client_call
polymarket/            Client implementation (same module, compatibility import path)
examples/              Runnable examples
docs/                  User and operator guides
```

## Documentation

| Document | Content |
|----------|---------|
| [docs/mcp-integration.md](docs/mcp-integration.md) | MCP wire format, tools, parameters, security. |
| [docs/homebrew-release.md](docs/homebrew-release.md) | Tap repository, tokens, tags, optional skip of formula push. |
| [CHANGELOG.md](CHANGELOG.md) | Release history. |

## Security

- Never commit private keys. Use env vars or secret managers.
- MCP and CLI can perform trading when a private key is supplied; restrict access to the process and logs.

## License

See [LICENSE](LICENSE).
