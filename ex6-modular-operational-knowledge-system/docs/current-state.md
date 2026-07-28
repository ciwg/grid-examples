# EX6 Current State

This document is intentionally blunt.

`ex6` is **not done**.

It currently contains the runtime foundation, one small built-in example
package, and the first nine real first-party built-in packages. It does not
yet contain the full operational knowledge product as a set of ex6 packages.
Source: `DI-moksu`; `DI-lupok`; `DI-lorup`; `DI-vakod`; `DI-pamuk`;
`DI-figar`; `DI-tusav`; `DI-sivuk`; `DI-ramek`; `DI-zibek`; `DI-lavom`.

## Implemented

- `moks package list`
- `moks package inspect <package-id>`
- `moks package install <dir>`
- `moks relay export <path>`
- `moks relay import <path>`
- `moks relay serve <addr>`
- `moks relay pull <url>`
- `moks relay push <url>`
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
- one real first-party built-in package: `maintenance`
- one real first-party built-in package: `receiving`
- one real first-party built-in package: `inventory`
- first-party package source directories for the current 9-package spine
- a starter package template for outside authors
- basket-mediated durable writes for installed executable packages
- runnable built-in and installed-package examples
- exact-byte relay import dedupe and relay batch metadata validation
- live HTTP peer exchange on top of the current relay batch shape

Source: `DI-moksu`; `DI-lupok`; `DI-rovum`; `DI-nupad`; `DI-sibok`; `DI-nasek`.

## Not Implemented

- browser embodiment
- Neovim embodiment
- ex5 review/search/operate parity
- full signing/verification discipline
- peer discovery, trust, and richer multi-peer relay behavior
- scaffold command for new packages
- later domain packages beyond the current first-pass ex5 set

Source: `DI-moksu`; `DI-lupok`.

## Why The Gap Exists

The current repo state reflects a deliberate kernel-first slice:

- first lock the basket contract
- then build real eggs against that contract

That first step happened.
The second step has now started with the built-in `context`, `knowledge`,
`runs`, `links`, `procedures`, `training`, `maintenance`, and `receiving`
eggs plus `inventory`, but the rest is still mostly open.
Source: `DI-moksu`; `DI-lupok`; `DI-lorup`; `DI-vakod`; `DI-pamuk`;
`DI-figar`; `DI-tusav`; `DI-sivuk`; `DI-ramek`; `DI-zibek`; `DI-lavom`.

## Next Product Work

Best next product work:

1. make package creation easier than manual template copying
2. harden and deepen the current package set
3. harden relay/peer behavior beyond the current batch shell

Source: `DI-moksu`; `DI-rovum`; `DI-nupad`; `DI-sibok`; `DI-sivuk`; `DI-ramek`; `DI-zibek`; `DI-lavom`.
