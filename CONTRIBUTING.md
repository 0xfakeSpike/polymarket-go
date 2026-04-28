# Contributing

Thanks for contributing to `polymarket-go`.

## Prerequisites

- Go 1.22+
- Git

## Local development

```bash
go mod tidy
make test
```

## Project structure

- Root package `github.com/0xfakeSpike/polymarket-go`: public SDK entrypoint.
- `polymarket/`: compatibility and implementation package.
- `cmd/pmctl`: CLI command.
- `cmd/polymarket-mcp`: MCP bridge command.
- `internal/`: non-public runtime layers for executables.
- `examples/`: runnable examples.

## Pull request checklist

- Keep changes focused and reviewable.
- Add or update tests when behavior changes.
- Run:
  - `make fmt`
  - `make vet`
  - `make test`
- Update docs (`README.md`, `docs/`, `CHANGELOG.md`) when needed.

## Commit guidance

- Prefer clear commit messages that explain *why*.
- Avoid force-push on shared branches unless coordinated.
