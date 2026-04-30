# Homebrew tap and releases

Releases are built with **GoReleaser** on **GitHub Actions** when you push a tag matching `v*`. Each release publishes:

- GitHub Release assets (`pmctl`, `polymarket-mcp`, archives, `checksums.txt`)
- Formula updates in the tap repository **`0xfakeSpike/homebrew-tap`** (`polymarket-cli`, `polymarket-mcp`)

## Prerequisites

### 1. Tap repository

Create **`0xfakeSpike/homebrew-tap`** on GitHub (name must match `.goreleaser.yaml`):

1. [Create repository](https://github.com/new) under owner `0xfakeSpike`.
2. Repository name: **`homebrew-tap`**.
3. Add a README when prompted.

### 2. Actions secret on `polymarket-go`

Fine-grained PAT (or classic PAT) with **`Contents: Read and write`** on **`0xfakeSpike/homebrew-tap`**:

| Secret name | Value |
|-------------|--------|
| `HOMEBREW_TAP_GITHUB_TOKEN` | Token with push access to the tap repo |

GitHub Actions provides `GITHUB_TOKEN` for uploads to **`0xfakeSpike/polymarket-go`**; the extra secret is only for cross-repo formula commits.

### 3. Optional: binaries without tap

If formulas must not be pushed yet, set a repository **variable** on `polymarket-go`:

| Variable | Value |
|----------|--------|
| `SKIP_HOMEBREW_TAP` | `true` |

Remove it or set to `false` when the tap is ready.

## Publish a version

```bash
git tag v1.0.7
git push origin v1.0.7
```

Wait for the **release** workflow on `polymarket-go`, then verify formulas on `homebrew-tap`.

## Install for end users

```bash
brew tap 0xfakeSpike/tap
brew install polymarket-cli polymarket-mcp
```

### Renamed from `polymarket-go` (tap formula)

The Homebrew formula that installs **`pmctl`** is now **`polymarket-cli`** (the GitHub repo and Go module remain `polymarket-go`). If you still have the old formula:

```bash
brew uninstall polymarket-go
brew install polymarket-cli
```

After the next release, remove or replace the obsolete `polymarket-go.rb` in **`homebrew-tap`** so users are not offered two formulas for the same binary.

## Local GoReleaser check

```bash
go install github.com/goreleaser/goreleaser/v2@latest
goreleaser check
```

## Reference

| Symptom | What to verify |
|---------|----------------|
| Tap push fails | `HOMEBREW_TAP_GITHUB_TOKEN` scopes and expiry; tap repo exists and name is `homebrew-tap`. |
| Release asset upload errors | `.goreleaser.yaml` sets `release.replace_existing_artifacts` for safe retries; or delete the draft/broken release and tag again. |
| Workflow not triggered | Tag must match `v*`. |

Configuration: `.goreleaser.yaml`, workflow: `.github/workflows/release.yml`.
