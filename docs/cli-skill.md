# CLI Skill (`pmctl`)

This guide is the practical playbook for using `pmctl` in daily workflows.

## Quick Rules

- Run `pmctl methods -long` first to inspect callable SDK methods and signatures.
- Use `pmctl call -args '<json array>' <ClientMethod>` as the default invocation path.
- `pmctl tool` runs registry tools: `methods`, `client_call`, and the read-only helper `get_markets_by_annualized_return`.
- `tool` params are JSON **object**; `call` args are JSON **array**.
- Prefer read-only mode (`-public=true`, default). Use private key only when required.

## Command Patterns

```bash
pmctl tools
pmctl tool -params '{"limit":5,"max_pages":1}' get_markets_by_annualized_return
pmctl methods -long
pmctl call GetOK
pmctl call -args '["<CLOB_TOKEN_ID>"]' GetOrderBook
pmctl call -args '["<CONDITION_ID>"]' GetClobMarketInfo
pmctl call -args '[{"limit":10,"max_pages":3,"min_best_ask":0.5}]' GetMarketsByAnnualizedReturn
```

## Auth and Safety

- Public-only queries: keep default `-public=true`.
- Authenticated operations: `-public=false` with `-private-key` or `PMCTL_PRIVATE_KEY`.
- Never print full private keys in logs or examples.

## Troubleshooting

- `invalid params`: JSON is malformed or wrong shape.
  - `tool` requires object: `{"k":"v"}`
  - `call` requires array: `["arg1", 2]`
- `unknown tool`: run `pmctl tools` for valid names (`methods`, `client_call`, `get_markets_by_annualized_return`).
- reflection call errors: verify `pmctl methods -long`, method name, and argument order/types.
