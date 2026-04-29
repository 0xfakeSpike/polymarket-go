# Contributing

## Prerequisites

- Go **1.24+** (see `go.mod`)
- Git

## Setup

```bash
git clone https://github.com/0xfakeSpike/polymarket-go.git
cd polymarket-go
go mod download
```

## Checks before opening a PR

```bash
make fmt
make vet
make test
```

## Layout (where changes usually go)

| Area | Path |
|------|------|
| SDK implementation | `polymarket/` |
| Public API surface (re-exports) | Repository root `*.go` |
| CLI | `internal/cli/pmctl/`, `cmd/pmctl/` |
| MCP stdio server | `internal/mcp/stdio/`, `cmd/polymarket-mcp/` |
| Shared tools (`pmctl` / MCP) | `internal/tools/`, `internal/tools/invoke/` |
| Examples | `examples/` |
| User docs | `docs/`, `README.md` |

## Pull requests

- Keep changes focused and easy to review.
- Add or update tests when behavior or contracts change.
- Update `README.md` or `docs/` when user-visible behavior changes.
- Record notable user-facing changes under **`[Unreleased]`** in `CHANGELOG.md` (Keep a Changelog style).

## Commits

Use clear subject lines; explain non-obvious decisions in the body when helpful.
