# CLI Skill (`pmctl`)

This guide is the practical playbook for using `pmctl` in daily workflows.

## Quick Rules

- Run `pmctl methods -long` first to inspect callable SDK methods and signatures.
- Use `pmctl call -args '<json array>' <ClientMethod>` as the default invocation path.
- `pmctl tool` is reserved for registry-level tools (`methods` and `client_call`).
- `tool` params are JSON **object**; `call` args are JSON **array**.
- Prefer read-only mode (`-public=true`, default). Use private key only when required.

## Command Patterns

```bash
pmctl tools
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
- `unknown tool`: only `methods` and `client_call` are valid tool names.
- reflection call errors: verify `pmctl methods -long`, method name, and argument order/types.
