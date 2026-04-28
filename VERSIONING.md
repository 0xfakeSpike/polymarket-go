# Versioning Policy

This repository follows [Semantic Versioning](https://semver.org/).

## Stability contract

- `v1.x.y`
  - No breaking API changes for public root import path: `github.com/0xfakeSpike/polymarket-go`
  - Backward-compatible additions and fixes are allowed.
- Breaking changes require a major bump (`v2+`).

## Compatibility path

The subpackage path `github.com/0xfakeSpike/polymarket-go/polymarket` remains available for compatibility, but new integrations should prefer the root package.

## Deprecation process

1. Mark symbols or paths as deprecated in docs/comments.
2. Keep behavior working for at least one minor release line.
3. Remove only in the next major release.
