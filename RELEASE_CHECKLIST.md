# Release checklist

Run through this list before tagging a new **`v*`** release.

## Code quality

- [ ] `go mod tidy` — no unexpected diff.
- [ ] `make fmt` — passes.
- [ ] `make vet` — passes.
- [ ] `make test` — passes.

## CLI and MCP smoke

- [ ] `go run ./cmd/pmctl tools`
- [ ] `go run ./cmd/pmctl tool -params '{"query":"election","limit":1}' search_events`
- [ ] `echo '{"tool":"search_events","params":{"query":"election","limit":1}}' | go run ./cmd/polymarket-mcp`

## Documentation

- [ ] `README.md` matches current commands and flags.
- [ ] `docs/` updated if wire format or tools changed.
- [ ] `CHANGELOG.md` — add a **`[vX.Y.Z]`** section with user-facing notes.

## Version and publish

- [ ] Choose semver (`VERSIONING.md`).
- [ ] Commit release notes.
- [ ] `git tag vX.Y.Z && git push origin vX.Y.Z`
- [ ] Confirm GitHub Actions **release** workflow succeeds.
- [ ] Publish or edit GitHub Release notes from `CHANGELOG.md` if you maintain them manually.

## Homebrew (if publishing formulas)

- [ ] Tap **`0xfakeSpike/homebrew-tap`** exists.
- [ ] `HOMEBREW_TAP_GITHUB_TOKEN` is set on **`polymarket-go`** with push access to the tap.
- [ ] `SKIP_HOMEBREW_TAP` unset or `false` unless intentionally skipping formula push.

See **[docs/homebrew-release.md](docs/homebrew-release.md)** for details.
