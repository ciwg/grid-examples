# EX6 Current State

This document is intentionally blunt.

`ex6` is **not done**.

It currently contains the runtime foundation and one small built-in example
package. It does not yet contain the full operational knowledge product as a
set of ex6 packages. Source: `DI-moksu`; `DI-lupok`.

## Implemented

- `moks package list`
- `moks package inspect <package-id>`
- `moks package install <dir>`
- `moks relay export <path>`
- `moks relay import <path>`
- package manifest validation
- installed-package self-check
- append-only durable history
- CAS object storage
- one built-in example package: `ops-note`
- first-party package source directories for the initial 5-package split
- a starter package template for outside authors

Source: `DI-moksu`; `DI-lupok`.

## Not Implemented

- first-party `context` package
- first-party `knowledge` package
- first-party `runs` package
- first-party `links` package
- first-party `procedures` package
- browser embodiment
- Neovim embodiment
- ex5 review/search/operate parity
- full signing/verification discipline
- richer PromiseGrid relay/peer behavior

Source: `DI-moksu`; `DI-lupok`.

## Why The Gap Exists

The current repo state reflects a deliberate kernel-first slice:

- first lock the basket contract
- then build real eggs against that contract

That first step happened.
The second step is still mostly open. Source: `DI-moksu`; `DI-lupok`.

## Next Product Work

Best next package set:

1. `context`
2. `knowledge`
3. `runs`
4. `links`
5. `procedures`

Recommended source layout:

- `packages/context/`
- `packages/knowledge/`
- `packages/runs/`
- `packages/links/`
- `packages/procedures/`

Source: `DI-moksu`.
