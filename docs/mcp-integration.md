# MCP bridge (`polymarket-mcp`)

`polymarket-mcp` is a **standards-compatible MCP server** over stdio. It speaks MCP JSON-RPC with `Content-Length` framing and can be connected directly by MCP clients (including Codex).

## Transport and protocol

- Transport: stdio
- Framing: `Content-Length: <n>\r\n\r\n<json>`
- Protocol: MCP JSON-RPC (`initialize`, `tools/list`, `tools/call`, `ping`)
- Current server protocol version: `2024-11-05`

## Environment

| Variable | When to set |
|----------|-------------|
| `POLYMARKET_MCP_PRIVATE_KEY` | Hex private key (`0x` optional). Enables `NewClient` (L2 bootstrap, trading). Omit for read-only `NewPublicClient`. |

## Tools

Tool registry is intentionally minimal and unified: `methods` + `client_call` (see `internal/tools`). SDK calls should go through MCP `tools/call.arguments`.

### `methods`

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `long` | bool | no | If `true`, include per-method reflect signatures. |

MCP call arguments:

```json
{"long":true}
```

### `client_call` (unified method invocation)

Invoke any **exported** `polymarket.Client` method by name. Arguments are a JSON **array** in Go parameter order; `context.Context` parameters are injected and must not appear in `args`.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `method` | string | yes | Exported method name, e.g. `GetOrderBook`. |
| `args` | array | no | Defaults to `[]`. |

MCP call arguments:

```json
{"method":"GetOrderBook","args":["<CLOB_TOKEN_ID>"]}
```

Get annualized market ranking through the same endpoint:

```json
{"method":"GetMarketsByAnnualizedReturn","args":[{"limit":10,"max_pages":3,"min_best_ask":0.5}]}
```

Methods whose parameters include **functions** or **non-empty interfaces** (e.g. WebSocket handlers) cannot be called this way.

## Codex connection

Point your MCP client to launch `polymarket-mcp` directly via stdio (no adapter needed).

Example command:

```bash
polymarket-mcp
```

## Security

- Treat `POLYMARKET_MCP_PRIVATE_KEY` like production credentials; scope host filesystem and logs.
- Do not echo secrets in `data` payloads from custom wrappers; built-in tools return API-shaped data only.

## Related

- CLI equivalents: [README.md](../README.md#cli-pmctl)
- Homebrew install: [homebrew-release.md](homebrew-release.md)
