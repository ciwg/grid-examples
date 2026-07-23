# ex5 PromiseGrid implementation claims

This document states what the shipped `ex5-operational-knowledge-system`
implementation currently promises in PromiseGrid terms, and what it does not
yet promise at the runtime/wire layer. Source: `DI-sobek`.

It exists so the current ex5 runtime can be described honestly without
blurring together:

- the current shipped implementation
- the PromiseGrid framing that ships with this example
- the future-scope work that ex5 still has not implemented yet, such as a
  browser-side direct non-HTTP embodiment contract and broader ERP-style
  planning behavior. Source: `DI-murev`.

## Current status

`ex5` ships as part of the PromiseGrid example set, and it now has eight
PromiseGrid-native runtime families, but not yet across the whole operational
model. Source: `DI-sobek`; `DI-mibor`; `DI-vosul`; `DI-kavup`; `DI-ribof`;
`DI-votek`; `DI-sarib`; `DI-vamok`; `DI-faruv`.

Within the current shipped scope, that runtime slice is now complete: the
remaining gaps are explicit future-scope choices, not hidden migration debt.
Source: `DI-lavek`.

Today it ships:

- one local Go runtime
- append-only operational event history
- local durable draft manifests plus CAS-backed draft bodies, and durable
  attachment storage
- websocket-preferred shared live-draft carriage for the browser, with the
  existing HTTP live routes preserved as compatibility paths
- direct local Unix-socket embodiment contracts for CLI and Neovim, with HTTP
  kept as fallback and compatibility
- projected read/query views over that history
- browser over the local HTTP adapter, plus CLI and Neovim over a direct local
  Unix-socket contract
- one frozen `knowledge-item` profile selected from the exact shipped protocol
  bytes
- one frozen `knowledge-approval` profile selected from the exact shipped
  protocol bytes
- one frozen `knowledge-evidence` profile selected from the exact shipped
  protocol bytes
- one frozen `knowledge-link` profile selected from the exact shipped protocol
  bytes
- one frozen `knowledge-responsibility` profile selected from the exact shipped
  protocol bytes
- one frozen `operational-run` profile selected from the exact shipped
  protocol bytes
- one frozen `operational-place` profile selected from the exact shipped
  protocol bytes
- one frozen `operational-resource` profile selected from the exact shipped
  protocol bytes
- one local signed-envelope runtime slice for durable knowledge-item
  create/revision/lifecycle events
- one local signed-envelope runtime slice for durable knowledge-item and run
  approval artifacts
- one local signed-envelope runtime slice for durable evidence metadata plus
  attachment references
- one local signed-envelope runtime slice for durable performed run records
- one local signed-envelope runtime slice for durable operational place records
- one local signed-envelope runtime slice for durable operational resource
  records
- one local signed-envelope runtime slice for durable typed links
- one local signed-envelope runtime slice for first-class responsibility
  records
- additive CAS-backed sidecar storage for signed family envelopes by envelope
  CID
- additive CAS-backed sidecar storage for copied evidence blobs by blob CID
- authoritative CAS-backed replay/export envelope bytes for the eight frozen
  families, with one-time manifest backfill for older runtimes
- runtime capability metadata exposing the shipped peer-exchange format,
  relay-feed format, blob-transfer support, CAS support, and per-embodiment
  transport semantics through `Meta`
- startup verification of the signed knowledge-item envelope log against the
  replayed item event history
- startup verification of the signed knowledge-approval envelope log against
  the replayed approval event history
- startup verification of the signed knowledge-evidence envelope log against
  the replayed evidence event history
- startup verification of the signed operational-run envelope log against the
  replayed run event history
- startup verification of the signed operational-place envelope log against the
  replayed place event history
- startup verification of the signed operational-resource envelope log against
  the replayed resource event history
- startup verification of the signed knowledge-link envelope log against the
  replayed link event history
- startup verification of the signed knowledge-responsibility envelope log
  against the replayed responsibility event history

Today it does **not** yet ship:

- a browser-side direct non-HTTP embodiment contract

## What the shipped implementation does promise

### 1. One shared local runtime contract

Browser, CLI, and Neovim all read and write one shared ex5 runtime model. The
embodiments are not separate durable systems. The browser still projects
through the local HTTP adapter, while CLI and Neovim now prefer a direct local
Unix-socket contract over that same runtime. Browser live drafting still
prefers websocket carriage under the HTTP adapter, while Neovim live drafting
prefers the local socket and keeps HTTP as fallback. Source: `DI-fudok`;
`DI-ravum`; `DI-sobek`; `DI-bavuk`; `DI-noruv`; `DI-favel`.

### 1a. One dedicated remote relay surface

`ex5` now also ships a separate `operational-relay` binary whose only job is
to persist and serve origin-aware relay-feed history plus CID-addressed blobs
under `/relay/v1`. It does not replace the local embodiment adapter, and it
does not become the main application runtime. Source: `DI-rovik`;
`DI-tasov`; `DI-nulav`.

### 2. Append-only durable operational history

The runtime keeps durable operational history as append-only events and then
projects that history into current-state views for operator workflows. Source:
`DI-radok`; `DI-zuvob`; `DI-sobek`.

### 3. Durable history is separate from live drafting and other transient UI state

The current shared live draft is not itself the durable historical record.
Durable revisions, runs, approvals, evidence, and typed links are distinct
from transient working state such as draft presence, current filters, or local
UI focus. Source: `DI-lusov`; `DI-zoruk`; `DI-dazim`; `DI-sobek`.

### 4. The local HTTP API is still the browser adapter and compatibility surface

`GET /api/*` and `POST /api/*` routes are the shipped browser adapter surface
and the compatibility fallback for CLI and Neovim. The durable frozen families
underneath them are now also reachable through the direct local Unix-socket
contract used by the two terminal embodiments. Source: `DI-sobek`;
`DI-mibor`; `DI-vosul`; `DI-kavup`; `DI-votek`; `DI-sarib`; `DI-vamok`;
`DI-faruv`; `DI-favel`.

### 5. `knowledge-item`, `knowledge-approval`, `knowledge-evidence`, `knowledge-link`, `knowledge-responsibility`, `operational-run`, `operational-place`, and `operational-resource` are the current frozen families

`knowledge-item`, `knowledge-approval`, `knowledge-evidence`,
`knowledge-link`, `knowledge-responsibility`, `operational-run`,
`operational-place`, and `operational-resource` now select runtime behavior
through their computed `pCID`s, and the runtime signs and verifies durable
artifacts for all eight families. The other named ex5 families remain
documented framing and staged migration targets for now. Search metadata
remains derived projection state instead of a separate durable family. Source:
`DI-mibor`; `DI-vosul`; `DI-kavup`; `DI-ribof`; `DI-votek`; `DI-sarib`;
`DI-vamok`; `DI-fusok`; `DI-pivul`.

## What the shipped implementation does not yet promise

The current shipped ex5 runtime does not yet promise:

- that all ex5 durable families are already frozen and PromiseGrid-native at
  runtime
- that the local HTTP route names are the PromiseGrid peer contract
- that the browser already bypasses the local HTTP adapter
- that ephemeral presence or derived projections are durable PromiseGrid
  families

## Done now vs. remaining

Done now:

- `knowledge-item` is frozen as the first ex5 family
- `knowledge-approval` is frozen as the second ex5 family
- `knowledge-evidence` is frozen as the third ex5 family
- `knowledge-link` is frozen as the fourth ex5 family
- `knowledge-responsibility` is frozen as the fifth ex5 family
- `operational-run` is frozen as the sixth ex5 family
- `operational-place` is frozen as the seventh ex5 family
- `operational-resource` is frozen as the eighth ex5 family
- the runtime exports and bootstrap-imports whole-family signed
  `knowledge-item`, `knowledge-approval`, `knowledge-evidence`,
  `knowledge-link`, `knowledge-responsibility`, `operational-run`,
  `operational-place`, and `operational-resource` records plus their
  compatibility events over the local HTTP adapter
- the runtime now also imports origin-aware unseen peer history for those
  families into non-empty runtimes and dedupes it by
  `(origin_peer_id, origin_sequence)`
- the runtime now also exports and imports incremental relay-feed batches for
  those families by origin-aware cursor, and evidence blobs now move through
  separate CID-addressed relay blob routes instead of being inlined into every
  ongoing feed batch
- canonical durable IDs for those peer-visible entities now come from the
  create-event envelope CID, and the old short IDs are preserved only as
  aliases for display, replay compatibility, and embodiment transition
- search metadata remains derived projection state over those families, not a
  sixth signed family
- the runtime computes all eight `pCID`s from the exact shipped spec bytes
- the runtime signs and verifies durable knowledge-item create/revision/status
  artifacts
- the runtime signs and verifies durable knowledge-item and run approval
  artifacts
- the runtime signs and verifies durable evidence metadata plus attachment
  references
- the runtime signs and verifies durable performed run artifacts
- the runtime signs and verifies durable place artifacts
- the runtime signs and verifies durable resource artifacts
- the runtime signs and verifies durable typed-link artifacts
- the runtime signs and verifies durable responsibility-creation artifacts
- the runtime now rehydrates the eight frozen family envelopes from CAS
  authoritatively during replay/export, while keeping compatibility event
  replay for still-unfrozen state
- the runtime now reloads shared live draft bodies authoritatively from CAS
  through per-item local draft manifests, including one-time backfill of older
  manifest files that only carried inline draft text
- bootstrap peer exchange still carries inline CID-keyed evidence blobs and
  re-materializes them into a local compatibility attachment path on import
- ongoing relay-feed exchange now requires evidence blobs to be staged into
  local CAS by CID before evidence-bearing feed batches import successfully
- the browser still projects through the current local HTTP adapter, while CLI
  and Neovim now project through the local Unix-socket contract on top of
  those signed families

Remaining:

- keep embodiment/runtime language honest while browser-side non-HTTP
  embodiment work and other deferred product follow-on work remain outside the
  current PromiseGrid slice

## Current implementation claim

The current ex5 implementation claims live in
[CHANGELOG.md](../CHANGELOG.md). Source: `DI-mibor`; `DI-vosul`; `DI-kavup`;
`DI-votek`; `DI-sarib`.

## How to read the other ex5 docs

- [README](../README.md): current product/runtime summary with an honest
  PromiseGrid boundary
- [Architecture](./architecture.md): current topology and the local/runtime vs.
  future PromiseGrid boundary
- [HTTP API Guide](./http-api-guide.md): current local embodiment adapter
- [Practical Implementation](./practical-implementation.md): current storage,
  projection, and embodiment details
- [operational-knowledge-system-spec-v0.1.md](../operational-knowledge-system-spec-v0.1.md):
  shipped PromiseGrid reference spec prose for ex5, not the fully implemented
  wire contract

## Follow-on work

The next staged PromiseGrid work is no longer a sixth durable
`knowledge-search-metadata` family. Search metadata is settled as derived
projection state, so the next work is the peer/storage layer backlog that
follows the frozen operational families. The relay-visible slice now ships as
whole-family bootstrap export/import plus incremental relay-feed export/import
for items, approvals, evidence, runs, places, resources, links, and
responsibilities. Bootstrap exchange still inlines evidence blobs by CID,
while incremental relay-feed exchange keeps blobs on separate CID-addressed
relay routes.
CAS now ships as an additive sidecar for signed envelopes and copied evidence
blobs plus authoritative reload for shared live draft bodies through local
manifests, and the embodiment-tightening step now ships through capability
metadata plus adapter-over-runtime doc updates while keeping HTTP as the sole
embodiment adapter. The current peer layer is no longer bootstrap-only; it now
uses origin-aware dedupe and local sequence projection for ongoing import, and
its peer-visible entities now use create-envelope CIDs as the durable IDs while
preserving the old short IDs as aliases. Source: `DI-fusok`; `DI-guzab`;
`DI-voruk`; `DI-ribek`; `DI-lavuz`; `DI-vabek`; `DI-rovuz`; `DI-tivor`;
`DI-vamok`; `DI-faruv`; `DI-ruzok`; `DI-rumek`; `DI-loruk`; `DI-pivul`;
`DI-zunep`; `DI-bavuk`; `DI-pazek`.
