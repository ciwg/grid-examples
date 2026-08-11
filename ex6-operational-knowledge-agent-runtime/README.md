# Operational Knowledge Agent Runtime (OKAR)

`ex6-operational-knowledge-agent-runtime` is the OKAR project folder.
OKAR is the new runtime for the operational knowledge agent runtime effort.

It is not `ex5` version 2.

It is the start of a new standalone product where:

- the runtime owns shared grid-facing infrastructure
- packages own business behavior
- built-in and installed packages use the same runtime contract

Current status: this repo contains the runtime foundation, first-party package
families, installed-package support, relay carriage, and workflow artifacts.
Its 27 built-in package-record families have immutable specifications and fixed
CIDv1 pCIDs. It does **not** yet contain the full ex5 capability set as OKAR
packages. Source: `DI-moksu`; `DI-lupok`; `DI-lorup`; `DI-vakod`; `DI-pamuk`;
`DI-figar`; `DI-tusav`; `DI-sivuk`; `DI-ramek`; `DI-zibek`; `DI-lavom`;
`DI-okar`; `DI-jusij`; `DI-solan`.

## PromiseGrid scope

Ex6 records durable package evidence as canonical Grid CBOR bytes. A pCID names
one frozen shared record contract, not a package, workflow, executable, or
individual record. Workflows normally compose existing family pCIDs, so adding
or retiring a workflow does not require recompiling Ex6. A new interoperable
record meaning requires a new immutable specification and pCID; third-party
packages own those specifications outside the built-in registry.

The runtime preserves unknown-family records as exact bytes and does not infer
their semantics. Semantic author evidence, relay-carriage signatures, route
availability, approvals, and workflow execution remain separate local-policy
questions. Ex6 does not claim global identity, universal authority, consensus,
or automatic trust from package installation or a signature. Source:
`DI-sidoh`; `DI-jusij`; `DI-solan`.

Ex6's current local admission policy requires a valid semantic author signature
before a canonical package record enters append-only history through append or
import. This does not make a signature an authority grant or alter a family
pCID. See [Author-Signature Admission Policy](./docs/author-signature-admission-policy.md).
Source: `DI-damab`.

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
- implementation-local, pCID-selected CAS workflow lifecycle events with a disposable rebuilt cache
- peer-authenticated workflow artifact relay with separate received-lifecycle evidence
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
- 27 frozen built-in package-record families across `context`, `knowledge`,
  `runs`, `links`, `procedures`, `training`, `maintenance`, `receiving`,
  `inventory`, `quarantine`, corrective action, and operational notes
- runnable built-in and installed-package examples

Source: `DI-moksu`; `DI-lupok`; `DI-lorup`; `DI-vakod`; `DI-pamuk`;
`DI-figar`; `DI-tusav`; `DI-sivuk`; `DI-ramek`; `DI-zibek`; `DI-lavom`; `DI-rovum`; `DI-nupad`; `DI-zotem`; `DI-vemut`; `DI-kasud`; `DI-lutep`; `DI-zumep`; `DI-ravud`; `DI-luzef`; `DI-sovem`; `DI-fogem`; `DI-movek`; `DI-ravok`; `DI-rumek`; `DI-rutom`; `DI-novuk`.

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
- `moks workflow relay push <alias> <peer-id>` transfers an exact workflow artifact and the sender's signed lifecycle evidence to an allowed peer; the receiver retains the artifact plus local authenticated receipt metadata but does not import or activate it
- `moks workflow inbox list`, `inspect <artifact-cid>`, and `import <artifact-cid> <alias>` make retained receipt evidence visible and require a separate local import decision
- `moks workflow overview` is the read-only team briefing for readiness, inbox attention, current runs, and one safe next command
- claim attestation trust is local too: `moks relay peer classify ...` assigns class/weight, `moks relay peer federate ...` assigns federation, and `moks relay policy claim set-federated ...` decides what spread is enough before import

Source: `DI-vemut`; `DI-kasud`; `DI-lutep`; `DI-movek`; `DI-ravok`; `DI-rumek`; `DI-novuk`; `DI-jifuk`; `DI-rufir`; `DI-sotad`.

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
- named scopes now exist for common views, including `scope direct-hops`, `scope deep-hops`, and `scope one-branch:<candidate-id>`
- operators can now define local scope aliases with `moks route scope list`, `moks route scope set <name> <kind> <target> ...`, and `moks route scope remove <name>`
- `moks route scope inspect <name>` now shows both the raw scope clauses and the fully expanded clause list that trace filtering will use
- scope inspection now also reports skipped branches such as alias cycles or unresolved scope references
- each expanded scope clause now carries provenance showing which alias or built-in scope chain produced it
- scope inspection now groups expanded clauses by provenance branch as well as showing the flat expanded list
- each grouped provenance branch now carries a short deterministic label and a human-readable summary
- skipped scope branches now also attach to the grouped provenance branch they came from
- `moks route scope inspect <name>` can now sort or filter grouped branches by depth, label, or summary
- `moks route scope inspect <name>` now also echoes the active branch query in the output whenever grouped branches are being sorted or filtered
- that same scope inspection output now includes a short query summary with matched-group count, hidden-group count, total-group count, and effective ordering
- scope inspection now also emits diagnostics when it defaulted to label ordering or when an active branch query matched zero grouped branches
- invalid branch-query values such as `depth abc` now appear in diagnostics as ignored filters instead of silently acting like no filter
- invalid sort values such as `sort weird` now also show up as ignored filters and explicitly explain the label-order fallback
- degenerate text filters such as whitespace-only `label` or `summary` values now also appear in diagnostics as ignored filters
- focused traces now report total, shown, and hidden step counts plus the active filter
- `moks route policy show [<protocol-pcid> [<role>]]`, `moks route policy set ...`, `moks route policy set-for <protocol-pcid> ...`, `moks route policy remove <protocol-pcid>`, `moks route policy set-for-role <protocol-pcid> <role> ...`, and `moks route policy remove-role <protocol-pcid> <role>` control planner preferences globally, per input protocol, and per input protocol plus route role
- a registered family now requires a matching `family-validator` route claim for its `protocol_pcid`
- relay export now carries those route registrations as batch metadata too
- parser and transform routes can declare `emits_protocols` to describe the next-hop protocols they produce

Source: `DI-rutom`; `DI-ruvot`; `DI-lafek`; `DI-fotav`; `DI-pabut`; `DI-matek`; `DI-posek`; `DI-rivuk`; `DI-lavik`; `DI-fobek`; `DI-povak`; `DI-rusom`; `DI-dovak`; `DI-buvok`; `DI-zafek`; `DI-rukav`; `DI-vatuk`; `DI-lupav`; `DI-sovak`; `DI-vobek`; `DI-zumok`; `DI-zamuk`; `DI-bemok`; `DI-rusek`; `DI-fusek`; `DI-zusek`; `DI-vusek`; `DI-busek`; `DI-yusek`; `DI-lusek`; `DI-musek`; `DI-nusek`; `DI-pusek`; `DI-tusek`; `DI-vusem`; `DI-zusev`.

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
- [PromiseGrid Implementation Claims](./docs/promisegrid-implementation-claims.md)
- [Frozen Package Family pCID Registry](./docs/protocols/package-family-pcid-registry.md)
- [Testing and Evidence Guide](./docs/testing.md)
- [Agent/Kernel Alignment](./docs/agent-kernel-alignment.md)
- [Current State](./docs/current-state.md)
- [Runnable Examples](./docs/runnable-examples.md)
- [Package Author Guide](./docs/package-author-guide.md)
- [EX5 Capability Map](./docs/ex5-capability-map.md)
- [Changelog](./CHANGELOG.md)
