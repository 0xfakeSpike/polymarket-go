# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project follows [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added
- Public root SDK import path via `github.com/0xfakeSpike/polymarket-go`.
- CLI entrypoint `cmd/pmctl` for search and orderbook workflows.
- MCP stdio bridge entrypoint `cmd/polymarket-mcp`.
- Internal runtime layers in `internal/cli/pmctl` and `internal/mcpbridge`.
- Usage examples in `examples/public-search` and `examples/orderbook`.
- Documentation for MCP integration and versioning policy.
- Initial Go module setup (`go.mod`/`go.sum`).

### Changed
- Repository layout reorganized to be suitable for a public SDK + tools project.
- `polymarket/` path kept as compatibility layer, while root package is now preferred.
