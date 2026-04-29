# CLI Skill (`pmctl`)

This guide is the practical playbook for using `pmctl` in daily workflows.

## Quick Rules

- Run `pmctl tools` first to discover supported named tools.
- Use `pmctl tool -params '<json object>' <tool_name>` for registry tools.
- Use `pmctl methods` / `pmctl call` when a named tool does not cover the task.
- `tool` params are JSON **object**; `call` args are JSON **array**.
- Prefer read-only mode (`-public=true`, default). Use private key only when required.

## Command Patterns

```bash
pmctl tools
pmctl tool -params '{"query":"election","limit":5}' search_events
pmctl tool -params '{"token_id":"<CLOB_TOKEN_ID>"}' get_orderbook
pmctl tool -params '{"tag_slug":"geopolitics","keyword":"iran","limit":10,"min_annualized_return":0.25}' rank_markets_by_annualized_return
pmctl methods -long
pmctl call GetOK
pmctl call -args '["<CLOB_TOKEN_ID>"]' GetOrderBook
```

## Auth and Safety

- Public-only queries: keep default `-public=true`.
- Authenticated operations: `-public=false` with `-private-key` or `PMCTL_PRIVATE_KEY`.
- Never print full private keys in logs or examples.

## Troubleshooting

- `invalid params`: JSON is malformed or wrong shape.
  - `tool` requires object: `{"k":"v"}`
  - `call` requires array: `["arg1", 2]`
- `unknown tool`: run `pmctl tools` and use exact name.
- reflection call errors: verify `pmctl methods -long`, method name, and argument order/types.
