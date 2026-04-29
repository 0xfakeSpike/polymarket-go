# MCP bridge (`polymarket-mcp`)

`polymarket-mcp` reads **one JSON request per line** from stdin and writes **one JSON response per line** to stdout. Host it behind any MCP server that can spawn a process and map tool calls to this protocol.

## Request and response

**Request** (single line):

```json
{
  "tool": "<tool_name>",
  "params": { }
}
```

`params` is always a JSON **object** (use `{}` when empty).

**Response** (single line):

```json
{
  "ok": true,
  "data": { }
}
```

or on failure:

```json
{
  "ok": false,
  "error": "message"
}
```

## Environment

| Variable | When to set |
|----------|-------------|
| `POLYMARKET_MCP_PRIVATE_KEY` | Hex private key (`0x` optional). Enables `NewClient` (L2 bootstrap, trading). Omit for read-only `NewPublicClient`. |

## Tools

All tools share the same definitions as `pmctl tool` (see `internal/tools`).

### `search_events`

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `query` | string | yes | Search text. |
| `limit` | int | no | Max events (default `10`). |

```json
{"tool":"search_events","params":{"query":"election","limit":5}}
```

### `get_orderbook`

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `token_id` | string | yes | CLOB token id. |

```json
{"tool":"get_orderbook","params":{"token_id":"<CLOB_TOKEN_ID>"}}
```

### `methods`

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `long` | bool | no | If `true`, include per-method reflect signatures. |

```json
{"tool":"methods","params":{"long":true}}
```

### `client_call`

Invoke any **exported** `polymarket.Client` method by name. Arguments are a JSON **array** in Go parameter order; `context.Context` parameters are injected and must not appear in `args`.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `method` | string | yes | Exported method name, e.g. `GetOrderBook`. |
| `args` | array | no | Defaults to `[]`. |

```json
{"tool":"client_call","params":{"method":"GetOrderBook","args":["<CLOB_TOKEN_ID>"]}}
```

Methods whose parameters include **functions** or **non-empty interfaces** (e.g. WebSocket handlers) cannot be called this way.

## Security

- Treat `POLYMARKET_MCP_PRIVATE_KEY` like production credentials; scope host filesystem and logs.
- Do not echo secrets in `data` payloads from custom wrappers; built-in tools return API-shaped data only.

## Related

- CLI equivalents: [README.md](../README.md#cli-pmctl)
- Homebrew install: [homebrew-release.md](homebrew-release.md)
