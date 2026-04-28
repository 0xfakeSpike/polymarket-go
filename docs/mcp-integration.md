# MCP integration with `polymarket-go`

This repository now separates **SDK** and **entrypoints**, and includes a runnable bridge command in `cmd/polymarket-mcp`.

## Recommended architecture

- Keep all Polymarket business logic in SDK packages.
  - Public import path: `github.com/0xfakeSpike/polymarket-go`
  - Compatibility implementation path: `github.com/0xfakeSpike/polymarket-go/polymarket`
- MCP tool handlers should be thin adapters:
  - parse tool arguments
  - call SDK (`polymarket.Client`)
  - map response/errors into MCP result format

## Typical tool mapping

- `search_events(query, limit)` -> `Client.SearchEventsWithQuery`
- `get_order_book(token_id)` -> `Client.GetOrderBook`
- `get_market(slug)` -> `Client.GetMarketBySlug`
- `get_event(event_id)` -> `Client.GetEvent`

## Minimal handler pattern (pseudo)

```go
func handleSearchEvents(args map[string]any) (any, error) {
  query := asString(args["query"])
  limit := asIntDefault(args["limit"], 10)

  c, err := polymarket.NewPublicClient()
  if err != nil {
    return nil, err
  }
  events, err := c.SearchEventsWithQuery(query)
  if err != nil {
    return nil, err
  }
  if len(events) > limit {
    events = events[:limit]
  }
  return events, nil
}
```

## Runtime and security notes

- For read-only MCP tools, use `NewPublicClient` and avoid private keys.
- For trade/write tools, inject credentials via environment variables and never return secrets in tool output.
- Keep request timeouts enabled (SDK default timeout is 30s).

## Bridge command usage

Current `cmd/polymarket-mcp` is a lightweight stdio JSON bridge that can sit behind a full MCP server adapter.

Input example (one JSON per line):

```json
{"tool":"search_events","params":{"query":"election","limit":5}}
```

Supported tool names:

- `search_events`
- `get_orderbook`
