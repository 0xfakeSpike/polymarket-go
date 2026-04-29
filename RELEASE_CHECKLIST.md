# Release Checklist

Use this checklist before tagging a release (for example `v1.0.0`).

## Quality gates

- [ ] `go mod tidy` produces no unexpected changes.
- [ ] `make fmt` passes.
- [ ] `make vet` passes.
- [ ] `make test` passes.
- [ ] CLI smoke test:
  - [ ] `go run ./cmd/pmctl search-events -q "election" -limit 1`
  - [ ] `go run ./cmd/pmctl orderbook -token-id "<token_id>"`
- [ ] MCP bridge smoke test:
  - [ ] `echo '{"tool":"search_events","params":{"query":"election","limit":1}}' | go run ./cmd/polymarket-mcp`

## Docs and compatibility

- [ ] `README.md` reflects current commands and paths.
- [ ] `CHANGELOG.md` has release notes under a new version heading.
- [ ] `VERSIONING.md` still matches release policy.
- [ ] Any deprecations are explicitly documented.

## Tag and publish

- [ ] Create release commit.
- [ ] Create git tag (example: `git tag v1.0.0`).
- [ ] Push commits and tags.
- [ ] Publish GitHub release notes from `CHANGELOG.md`.

## Homebrew tap (方案 A: required for `brew install`)

- [ ] Tap repository exists: `https://github.com/0xfakeSpike/homebrew-tap` (see `docs/homebrew-release.md`).
- [ ] `polymarket-go` has Actions secret **`HOMEBREW_TAP_GITHUB_TOKEN`** with **Contents: write** on that tap repo.
- [ ] Repository variable **`SKIP_HOMEBREW_TAP`** is **unset** or **`false`** (otherwise formulas are not pushed).
- [ ] After pushing a `v*` tag, confirm the **release** workflow is green and the tap repo received formula commits.
