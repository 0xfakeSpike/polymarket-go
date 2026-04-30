# polymarket-go

Go client for the [Polymarket](https://polymarket.com) **CLOB** API, plus **`pmctl`** (CLI) and **`polymarket-mcp`** (standards-compatible MCP stdio server).

## Features

- **SDK** — import `github.com/0xfakeSpike/polymarket-go`; `Client` mirrors the CLOB API plus read-only helpers for the Polymarket **Gamma** profile and **Data** API (current/closed positions, user activity). No API key required for those public HTTP routes.
- **CLI (`pmctl`)** — unified reflective method calls via **`call`** plus **`methods`** to list signatures.
- **MCP server** — same tool registry as the CLI, exposed as native MCP over stdio; optional authenticated client via env.
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
	book, err := c.GetOrderBook("<CLOB_TOKEN_ID>")
	if err != nil {
		panic(err)
	}
	fmt.Println("bids:", len(book.Bids), "asks:", len(book.Asks))
}
```

Prefer the **root import** above. The implementation also lives under `github.com/0xfakeSpike/polymarket-go/polymarket` for users who want the internal package path explicitly.

## Install — binaries

### Homebrew

```bash
brew tap 0xfakeSpike/tap
brew install polymarket-cli polymarket-mcp
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
pmctl tool -params '{"limit":5,"max_pages":1,"tag_slug":"crypto"}' get_markets_by_annualized_return
pmctl methods -long | head -20
pmctl call GetOK
pmctl call -args '["<CLOB_TOKEN_ID>"]' GetOrderBook
pmctl call -args '["<CONDITION_ID>"]' GetClobMarketInfo
pmctl call -args '[{"limit":10,"max_pages":3,"events_page_limit":100,"tag_slug":"crypto","min_best_ask":0.5}]' GetMarketsByAnnualizedReturn
```

`call` injects `context.Context` where needed; methods that take **functions** or **handler interfaces** (e.g. WebSocket runners) are not supported through reflection — use the SDK in Go.

## MCP (`polymarket-mcp`)

Native MCP stdio server (JSON-RPC + `Content-Length` framing). Specification and tool schemas: **[docs/mcp-integration.md](docs/mcp-integration.md)**.

**Environment**

| Variable | Effect |
|----------|--------|
| `POLYMARKET_MCP_PRIVATE_KEY` | If set, `NewClient` (trading / L2 bootstrap). If unset, public client only. |

## Examples

```bash
go run ./examples/orderbook "<CLOB_TOKEN_ID>"
```

## Repository layout

```text
cmd/pmctl              CLI entrypoint
cmd/polymarket-mcp     MCP stdio entrypoint
internal/cli/pmctl   CLI wiring (flags, stdout/stderr)
internal/mcp/stdio   MCP stdio server (native MCP protocol)
internal/tools       Shared tool registry (methods, client_call, get_markets_by_annualized_return)
internal/tools/invoke Reflection helpers for client_call
polymarket/            CLOB client implementation
examples/              Runnable examples
docs/                  User and operator guides
```

## Documentation

| Document | Content |
|----------|---------|
| [docs/cli-skill.md](docs/cli-skill.md) | Practical `pmctl` command playbook for daily usage. |
| [docs/mcp-integration.md](docs/mcp-integration.md) | MCP wire format, tools, parameters, security. |
| [docs/homebrew-release.md](docs/homebrew-release.md) | Tap repository, tokens, tags, optional skip of formula push. |
| [CHANGELOG.md](CHANGELOG.md) | Release history. |

## Security

- Never commit private keys. Use env vars or secret managers.
- MCP and CLI can perform trading when a private key is supplied; restrict access to the process and logs.

## License

See [LICENSE](LICENSE).
