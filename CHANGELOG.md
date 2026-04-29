# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project follows [Semantic Versioning](https://semver.org/).

## [Unreleased]

## [v1.0.9]

### Changed
- Move CLI skill/playbook from `.cursor/skills` into public docs as `docs/cli-skill.md`.
- README documentation index now links the CLI skill guide for all users.

## [v1.0.8]

### Added
- New shared tool `rank_markets_by_annualized_return` for CLI/MCP to rank open markets by favored-side annualized return from live order books.
- Optional filter `min_annualized_return` to keep only markets above a target annualized threshold.

### Changed
- README and MCP integration docs now include usage and parameter examples for the annualized ranking tool.

## [v1.0.7]

### Added
- `pmctl methods`, `pmctl call`, and MCP `client_call` for invoking any exported `*Client` method with JSON arguments.
- `SKIP_HOMEBREW_TAP` repository variable to skip Homebrew formula commits when the tap is not configured.
- Root import `github.com/0xfakeSpike/polymarket-go`; `pmctl` and `polymarket-mcp` binaries; shared tool registry; examples under `examples/`.
- Docs: MCP integration, Homebrew release, versioning.

### Changed
- Public layout oriented around the root SDK module; `polymarket/` remains a supported import path.
- CLI uses `pmctl tool -params '<json>' <name>` for named tools instead of separate top-level subcommands per workflow.
- Homebrew documentation assumes a dedicated tap repository and token (see `docs/homebrew-release.md`).

## [v1.0.6]

### Changed
- Homebrew release documentation and release checklist aligned with tap-first publishing.

### Notes
- Tag when `0xfakeSpike/homebrew-tap` exists and `HOMEBREW_TAP_GITHUB_TOKEN` is set so GoReleaser can push `polymarket-go` and `polymarket-mcp` formulas.

## [v1.0.1]

### Fixed
- Dropped private module dependency that broke CI and public `go get`.
- Search helpers inlined where needed; public dependency graph is fully fetchable.

## [v1.0.2]

### Fixed
- Homebrew: two formulas (`polymarket-go`, `polymarket-mcp`) so artifact names match GoReleaser builds.

## [v1.0.3]

### Fixed
- Homebrew: per-binary archives and aligned formula IDs for correct artifact resolution.

## [v1.0.4]

### Fixed
- GoReleaser `release.replace_existing_artifacts` so release asset re-uploads do not fail with HTTP 422.
