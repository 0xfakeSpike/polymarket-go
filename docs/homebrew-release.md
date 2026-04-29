# Homebrew release guide (recommended: tap repository)

Releases use GoReleaser in GitHub Actions (push tag `v*`) to publish **GitHub Release** assets and commit **Homebrew formulas** into your tap repository.

## Recommended path: create the tap first (方案 A)

Do these **before** the first successful Homebrew publish.

### 1) Create the tap repository on GitHub

The tap must exist at exactly:

- `https://github.com/0xfakeSpike/homebrew-tap`

Steps:

1. Open [github.com/new](https://github.com/new).
2. **Owner:** `0xfakeSpike`
3. **Repository name:** `homebrew-tap` (must match `.goreleaser.yaml`; Homebrew convention is `homebrew-<something>`).
4. Visibility: **Public** is typical for taps (private works if your token can access it).
5. Check **Add a README**, then create the repository.

If this repo is missing, GoReleaser fails with **`404 Not Found`** on `GET .../repos/0xfakeSpike/homebrew-tap`.

### 2) Token for pushing formulas into the tap

Create a **fine-grained personal access token** that can write to **`0xfakeSpike/homebrew-tap`** (at minimum **Contents: Read and write**).

Add it to **`0xfakeSpike/polymarket-go`**:

- **Settings → Secrets and variables → Actions → New repository secret**
- Name: **`HOMEBREW_TAP_GITHUB_TOKEN`**
- Value: the token

The workflow already uses the default **`GITHUB_TOKEN`** for uploads to **`polymarket-go`** releases.

### 3) Tag to trigger release

```bash
git tag v1.0.6
git push origin v1.0.6
```

This runs `.github/workflows/release.yml`, which builds `pmctl` and `polymarket-mcp`, uploads archives to the GitHub Release, and pushes **`polymarket-go.rb`** and **`polymarket-mcp.rb`** into the tap.

### 4) Install via Homebrew

After the workflow succeeds:

```bash
brew tap 0xfakeSpike/tap
brew install polymarket-go
brew install polymarket-mcp
```

### 5) Optional: verify GoReleaser config locally

```bash
go install github.com/goreleaser/goreleaser/v2@latest
goreleaser check
```

---

## Optional: release binaries only (no tap)

If the tap is not ready yet, you can set a **repository variable** on `polymarket-go`:

- **Settings → Secrets and variables → Actions → Variables**
- Name: **`SKIP_HOMEBREW_TAP`**
- Value: **`true`**

The workflow passes this to GoReleaser; formula commits are skipped so the **GitHub Release** can still succeed. Remove the variable (or set to `false`) when you adopt **方案 A**.

---

## Troubleshooting

- **`404 ... homebrew-tap`:** Create the repo (section 1) or fix the name to match `.goreleaser.yaml`.
- **`401 Bad credentials` on tap:** Fix `HOMEBREW_TAP_GITHUB_TOKEN` scope or expiry.
- **`422 ... already_exists` on release assets:** Partial upload + retry. This repo sets `release.replace_existing_artifacts: true` in `.goreleaser.yaml`, or delete the broken release and tag again.
- **Tag workflow does not run:** Tag must match `v*`.
- **Formula name clash:** Rename `brews.name` in `.goreleaser.yaml`.
