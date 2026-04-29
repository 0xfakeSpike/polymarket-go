# Homebrew release guide (first time)

This project uses GoReleaser to publish GitHub Releases and update a Homebrew tap formula automatically.

## 0) If you see `404 Not Found` for `homebrew-tap`

GoReleaser calls `GET https://api.github.com/repos/0xfakeSpike/homebrew-tap`. A **404** almost always means the repository **does not exist yet** (it is not the same as a private repo returning 403).

Create it first (any visibility is fine for a tap):

- Open [github.com/new](https://github.com/new)
- Owner: `0xfakeSpike`
- Repository name: **`homebrew-tap`** (must match `.goreleaser.yaml`)
- Initialize with a README, then create the repository.

Until the tap exists, you can still publish **GitHub Release binaries only** by setting a repository (or organization) variable **`SKIP_HOMEBREW_TAP`** to **`true`** in `polymarket-go` settings. The workflow passes it to GoReleaser so Homebrew commits are skipped.

## 1) Create a Homebrew tap repository

Create this repository on GitHub first:

- `0xfakeSpike/homebrew-tap`

It can be empty, but should be initialized with a README.

## 2) Create a fine-grained token for tap updates

Create a GitHub token that can write to `0xfakeSpike/homebrew-tap`, then add it to this repo secrets:

- Secret name: `HOMEBREW_TAP_GITHUB_TOKEN`
- Repository: `0xfakeSpike/polymarket-go`

`GITHUB_TOKEN` is provided by GitHub Actions automatically for release uploads in this repo.

## 3) Verify release config locally

```bash
go install github.com/goreleaser/goreleaser/v2@latest
goreleaser check
```

## 4) Create and push a release tag

```bash
git tag v1.0.0
git push origin v1.0.0
```

This triggers `.github/workflows/release.yml`:

- builds `pmctl` and `polymarket-mcp`
- creates GitHub release artifacts
- updates Homebrew formula in `0xfakeSpike/homebrew-tap`

## 5) Install via brew

After workflow succeeds:

```bash
brew tap 0xfakeSpike/tap
brew install polymarket-go
brew install polymarket-mcp
```

## Troubleshooting

- If you see `404 Not Found` for `.../repos/0xfakeSpike/homebrew-tap`, create that repository (see section 0). A mistyped name also produces 404.
- If formula update fails, check `HOMEBREW_TAP_GITHUB_TOKEN` permissions.
- If tag workflow does not start, ensure tag matches `v*`.
- If GitHub upload fails with `422 ... already_exists`, the release already has those asset names (often from a partial run). This repo sets `release.replace_existing_artifacts: true` in `.goreleaser.yaml` so GoReleaser removes conflicting assets and retries. Alternatively delete the broken release on GitHub and re-run the workflow or push a new patch tag.
- If Homebrew package name conflicts, rename `brews.name` in `.goreleaser.yaml`.
