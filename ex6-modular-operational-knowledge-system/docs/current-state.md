# EX6 Current State

This document is intentionally blunt.

`ex6` is **not done**.

It currently contains the runtime foundation, one small built-in example
package, and the first six real first-party built-in packages. It does not
yet contain the full operational knowledge product as a set of ex6 packages.
Source: `DI-moksu`; `DI-lupok`; `DI-lorup`; `DI-vakod`; `DI-pamuk`;
`DI-figar`; `DI-tusav`; `DI-sivuk`.

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
- one real first-party built-in package: `context`
- one real first-party built-in package: `knowledge`
- one real first-party built-in package: `runs`
- one real first-party built-in package: `links`
- one real first-party built-in package: `procedures`
- one real first-party built-in package: `training`
- first-party package source directories for the current 6-package spine
- a starter package template for outside authors
- basket-mediated durable writes for installed executable packages
- runnable built-in and installed-package examples

Source: `DI-moksu`; `DI-lupok`; `DI-rovum`; `DI-nupad`.

## Not Implemented

- browser embodiment
- Neovim embodiment
- ex5 review/search/operate parity
- full signing/verification discipline
- richer PromiseGrid relay/peer behavior
- scaffold command for new packages
- later domain packages: `maintenance`, `receiving`, `inventory`

Source: `DI-moksu`; `DI-lupok`.

## Why The Gap Exists

The current repo state reflects a deliberate kernel-first slice:

- first lock the basket contract
- then build real eggs against that contract

That first step happened.
The second step has now started with the built-in `context`, `knowledge`,
`runs`, `links`, and `procedures` eggs, but the rest is still mostly open.
Source: `DI-moksu`; `DI-lupok`; `DI-lorup`; `DI-vakod`; `DI-pamuk`;
`DI-figar`; `DI-tusav`.

## Next Product Work

Best next product work:

1. make package creation easier than manual template copying
2. add the next domain packages: `maintenance`, `receiving`, `inventory`
3. harden relay/peer behavior beyond the current batch shell

Source: `DI-moksu`; `DI-rovum`; `DI-nupad`; `DI-sivuk`.
