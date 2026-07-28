# EX6 Modular Operational Knowledge System

`ex6-modular-operational-knowledge-system` is the new runtime for
the modular operational knowledge system effort.

It is not `ex5` version 2.

It is the start of a new standalone product where:

- the runtime owns shared grid-facing infrastructure
- packages own business behavior
- built-in and installed packages use the same runtime contract

Current status: this repo contains the runtime foundation, the first nine
first-party built-in packages (`context`, `knowledge`, `runs`, `links`,
`procedures`, `training`, `maintenance`, `receiving`, `inventory`), and one
small built-in example package. It does **not** yet contain the full ex5
capability set as ex6 packages. Source: `DI-moksu`; `DI-lupok`; `DI-lorup`;
`DI-vakod`; `DI-pamuk`; `DI-figar`; `DI-tusav`; `DI-sivuk`; `DI-ramek`;
`DI-zibek`; `DI-lavom`.

## What Exists Today

- standalone Go module
- `moks` CLI
- package manifest loading
- manifest plus self-check activation for installed packages
- built-in package registration
- command routing
- explicit protocol route registration derived from package claims
- exported route metadata in relay batches
- explicit route types: `direct`, `parser`, and `transform`
- family registration
- explicit per-package protocol implementation claims
- append-only durable history
- CAS storage
- relay batch export and import
- signed live relay batches
- per-record relay digest proofs
- per-record relay-carriage signatures by the exporting peer
- claim-level proofs for advertised implementation claims
- semantic author-level signatures on durable records
- third-party attestation support for implementation claims
- runtime-owned attestation policy and quorum for implementation claims
- weighted attestation trust and attester classes for implementation claims
- federated trust semantics with distinct federation spread for implementation claims
- explicit peer discovery and optional untrusted local seeding
- peer-policy promotion shortcuts over stored metadata
- unknown-family exact-byte retention
- the first nine first-party built-in packages:
  `context`, `knowledge`, `runs`, `links`, `procedures`, `training`, `maintenance`, `receiving`, and `inventory`
- runnable built-in and installed-package examples

Source: `DI-moksu`; `DI-lupok`; `DI-lorup`; `DI-vakod`; `DI-pamuk`;
`DI-figar`; `DI-tusav`; `DI-sivuk`; `DI-ramek`; `DI-zibek`; `DI-lavom`; `DI-rovum`; `DI-nupad`; `DI-zotem`; `DI-vemut`; `DI-kasud`; `DI-lutep`; `DI-zumep`; `DI-ravud`; `DI-luzef`; `DI-sovem`; `DI-fogem`; `DI-movek`; `DI-ravok`; `DI-rumek`; `DI-rutom`.

## What Does Not Exist Yet

- the full ex5 package set
- browser embodiment
- Neovim embodiment
- ex5 review/search/workflow parity
- full PromiseGrid proof/signature discipline
- richer author identity beyond the current local runtime signing root
- richer cross-runtime federation semantics beyond the current local labels and spread policy
- richer relay behavior beyond the current peer-card and explicit allow model

Source: `DI-moksu`; `DI-lupok`; `DI-zotem`; `DI-vemut`; `DI-kasud`; `DI-zumep`; `DI-ravud`; `DI-luzef`; `DI-sovem`; `DI-fogem`; `DI-movek`; `DI-ravok`; `DI-rumek`.

## Relay Trust Rule

- discovery is not trust
- `moks relay peer discover <peer-card-url>` only fetches and prints peer metadata plus next commands
- `moks relay peer discover <peer-card-url> seed` creates a local `no-pull` `no-push` peer entry
- `moks relay peer promote <peer-id> <pull|push|both>` upgrades local policy using that stored peer entry
- pull or push remain disabled until an explicit `moks relay peer allow ...` command changes policy
- claim attestation trust is local too: `moks relay peer classify ...` assigns class/weight, `moks relay peer federate ...` assigns federation, and `moks relay policy claim set-federated ...` decides what spread is enough before import

Source: `DI-vemut`; `DI-kasud`; `DI-lutep`; `DI-movek`; `DI-ravok`; `DI-rumek`.

## Route Inspection

- `moks route list` shows the current protocol routes that the runtime has derived from active package claims
- `moks route inspect <protocol-pcid>` returns the matching direct/parser/transform routes as JSON
- `moks route plan <protocol-pcid>` returns the kernel's preferred executable route plan plus ordered candidates, winner explanation, pairwise candidate comparisons, and explicit downstream-plan summaries for parser/transform hops
- `moks route plan <protocol-pcid> trace` adds the planner's exact step-by-step decision sequence
- `moks route plan <protocol-pcid> trace candidate <package-id:role:route-type>` or `... trace downstream <protocol-pcid>` narrows that trace to one path
- trace summaries now say which `protocol_pcid` they describe and whether they belong to the root plan or a downstream hop
- traced route-plan output now carries a top-level downstream trace summary list for nested protocol hops
- downstream trace summaries now carry stable `hop_path` labels so repeated same-`protocol_pcid` hops stay distinguishable
- downstream trace summaries now also carry short `hop_summary` strings for faster scanning
- trace summaries now carry `hop_depth` and `hop_index` metadata for sorting and filtering by distance from root
- `moks route plan <protocol-pcid> trace depth <n|n+>` narrows trace output by hop depth
- trace filters can now be combined, for example `trace depth 2 downstream <protocol-pcid>`
- focused traces now report total, shown, and hidden step counts plus the active filter
- `moks route policy show [<protocol-pcid> [<role>]]`, `moks route policy set ...`, `moks route policy set-for <protocol-pcid> ...`, `moks route policy remove <protocol-pcid>`, `moks route policy set-for-role <protocol-pcid> <role> ...`, and `moks route policy remove-role <protocol-pcid> <role>` control planner preferences globally, per input protocol, and per input protocol plus route role
- a registered family now requires a matching `family-validator` route claim for its `protocol_pcid`
- relay export now carries those route registrations as batch metadata too
- parser and transform routes can declare `emits_protocols` to describe the next-hop protocols they produce

Source: `DI-rutom`; `DI-ruvot`; `DI-lafek`; `DI-fotav`; `DI-pabut`; `DI-matek`; `DI-posek`; `DI-rivuk`; `DI-lavik`; `DI-fobek`; `DI-povak`; `DI-rusom`; `DI-dovak`; `DI-buvok`; `DI-zafek`; `DI-rukav`; `DI-vatuk`; `DI-lupav`; `DI-sovak`; `DI-vobek`; `DI-zumok`.

## Layout

- `cmd/moks/`
  runtime CLI entrypoint
- `kernel/`
  runtime package activation, routing service, family registry, relay export/import
- `packages/`
  package manifest and installed-package runner support
- `packages/<package-id>/`
  first-party package source trees and package planning docs
- `records/`
  ex6 durable envelope shape
- `store/`
  append-only history and CAS
- `grid/`
  relay batch types
- `builtin/`
  built-in example package(s)
- `templates/package/`
  starter template for outside package authors
- `docs/`
  ex6 reader and author documentation

Source: `DI-moksu`.

## Start Here

- [Architecture](./docs/architecture.md)
- [Agent/Kernel Alignment](./docs/agent-kernel-alignment.md)
- [Current State](./docs/current-state.md)
- [Runnable Examples](./docs/runnable-examples.md)
- [Package Author Guide](./docs/package-author-guide.md)
- [EX5 Capability Map](./docs/ex5-capability-map.md)
