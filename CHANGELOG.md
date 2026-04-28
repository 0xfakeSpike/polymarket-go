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

## [v1.0.1]

### Fixed
- Removed private module dependency (`github.com/0xfakespike/everything`) that broke CI/release on GitHub runners.
- Replaced external bool helper usage in search APIs with an internal helper.
- Regenerated module metadata (`go.mod`/`go.sum`) to keep dependencies fully public.

## [v1.0.2]

### Fixed
- Split Homebrew publishing into two formulas (`polymarket-go` and `polymarket-mcp`) so GoReleaser can match build artifacts correctly during formula generation.

## [v1.0.3]

### Fixed
- Split GoReleaser archives per binary and aligned formula IDs to archive IDs, fixing Homebrew artifact matching in release workflow.
