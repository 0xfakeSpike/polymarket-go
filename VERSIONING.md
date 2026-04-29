# Versioning

This project follows [Semantic Versioning 2.0](https://semver.org/).

## API stability

| Range | Guarantee |
|-------|-----------|
| **v1.x.y** | No breaking changes to the public API of **`github.com/0xfakeSpike/polymarket-go`**. Additive changes and fixes are allowed. |
| **v2+** | Used when a breaking change to the public API is required. |

## Import paths

- **Preferred:** `github.com/0xfakeSpike/polymarket-go` (root package).
- **Also published:** `github.com/0xfakeSpike/polymarket-go/polymarket` — same types; kept for existing importers. Prefer the root path for new code.

## Deprecations

Deprecated symbols are marked in GoDoc. Deprecated behavior remains available for at least one minor **v1** line; removal happens only in a **major** bump.
