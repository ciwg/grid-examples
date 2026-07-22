# ex5 PromiseGrid implementation claims

This document states what the shipped `ex5-operational-knowledge-system`
implementation currently promises in PromiseGrid terms, and what it does not
yet promise at the runtime/wire layer. Source: `DI-sobek`.

It exists so the current ex5 runtime can be described honestly without
blurring together:

- the current shipped implementation
- the PromiseGrid framing that ships with this example
- the wire-level work that ex5 has not implemented yet, such as frozen `pCID`
  handling and signed envelopes

## Current status

`ex5` ships as part of the PromiseGrid example set, and it now has five
PromiseGrid-native runtime families, but not yet across the whole operational
model. Source: `DI-sobek`; `DI-mibor`; `DI-vosul`; `DI-kavup`; `DI-ribof`;
`DI-votek`; `DI-sarib`.

Today it ships:

- one local Go runtime
- append-only operational event history
- local durable draft and attachment storage
- projected read/query views over that history
- browser, CLI, and Neovim embodiments over the same local HTTP adapter
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
- one local signed-envelope runtime slice for durable knowledge-item
  create/revision/lifecycle events
- one local signed-envelope runtime slice for durable knowledge-item and run
  approval artifacts
- one local signed-envelope runtime slice for durable evidence metadata plus
  attachment references
- one local signed-envelope runtime slice for durable typed links
- one local signed-envelope runtime slice for first-class responsibility
  records
- additive CAS-backed sidecar storage for signed family envelopes by envelope
  CID
- additive CAS-backed sidecar storage for copied evidence blobs by blob CID
- runtime capability metadata exposing the shipped peer-exchange format and CAS
  support through `Meta`
- startup verification of the signed knowledge-item envelope log against the
  replayed item event history
- startup verification of the signed knowledge-approval envelope log against
  the replayed approval event history
- startup verification of the signed knowledge-evidence envelope log against
  the replayed evidence event history
- startup verification of the signed knowledge-link envelope log against the
  replayed link event history
- startup verification of the signed knowledge-responsibility envelope log
  against the replayed responsibility event history

Today it does **not** yet ship:

- peer-visible `knowledge-evidence` exchange
- relay transport or CAS-backed envelope storage as part of the ex5
  operational workflow

## What the shipped implementation does promise

### 1. One shared local runtime contract

Browser, CLI, and Neovim all read and write one shared ex5 runtime model
through the same local HTTP surface. The embodiments are not separate durable
systems. Source: `DI-fudok`; `DI-ravum`; `DI-sobek`.

### 2. Append-only durable operational history

The runtime keeps durable operational history as append-only events and then
projects that history into current-state views for operator workflows. Source:
`DI-radok`; `DI-zuvob`; `DI-sobek`.

### 3. Durable history is separate from live drafting and other transient UI state

The current shared live draft is not itself the durable historical record.
Durable revisions, runs, approvals, evidence, and typed links are distinct
from transient working state such as draft presence, current filters, or local
UI focus. Source: `DI-lusov`; `DI-zoruk`; `DI-dazim`; `DI-sobek`.

### 4. The local HTTP API is still the embodiment adapter

`GET /api/*` and `POST /api/*` routes are the shipped local adapter surface for
browser, CLI, and Neovim. They are the current implementation contract for
those embodiments, but the durable `knowledge-item` family underneath them is
now also written as signed PromiseGrid-style envelopes in the local runtime.
Source: `DI-sobek`; `DI-mibor`.

### 5. `knowledge-item`, `knowledge-approval`, `knowledge-evidence`, `knowledge-link`, and `knowledge-responsibility` are the current frozen families

`knowledge-item`, `knowledge-approval`, `knowledge-evidence`,
`knowledge-link`, and `knowledge-responsibility` now select runtime behavior
through their computed `pCID`s, and the runtime signs and verifies durable
artifacts for all five families. The other named ex5 families remain
documented framing and staged migration targets for now. Search metadata
remains derived projection state instead of a separate durable family. Source:
`DI-mibor`; `DI-vosul`; `DI-kavup`; `DI-ribof`; `DI-votek`; `DI-sarib`;
`DI-fusok`.

## What the shipped implementation does not yet promise

The current shipped ex5 runtime does not yet promise:

- that all ex5 durable families are already frozen and PromiseGrid-native at
  runtime
- that the local HTTP route names are the PromiseGrid peer contract
- that relay transport, peer exchange, or CAS-backed envelope storage are
  already implemented for ex5

## Done now vs. remaining

Done now:

- `knowledge-item` is frozen as the first ex5 family
- `knowledge-approval` is frozen as the second ex5 family
- `knowledge-evidence` is frozen as the third ex5 family
- `knowledge-link` is frozen as the fourth ex5 family
- `knowledge-responsibility` is frozen as the fifth ex5 family
- the runtime exports and bootstrap-imports whole-family signed
  `knowledge-item`, `knowledge-approval`, `knowledge-link`, and
  `knowledge-responsibility` records plus their compatibility events over the
  local HTTP adapter
- bootstrap import preserves unresolved run approvals and links to
  place/resource/run endpoints explicitly instead of trimming the family logs
- search metadata remains derived projection state over those families, not a
  sixth signed family
- the runtime computes all five `pCID`s from the exact shipped spec bytes
- the runtime signs and verifies durable knowledge-item create/revision/status
  artifacts
- the runtime signs and verifies durable knowledge-item and run approval
  artifacts
- the runtime signs and verifies durable evidence metadata plus attachment
  references
- the runtime signs and verifies durable typed-link artifacts
- the runtime signs and verifies durable responsibility-creation artifacts
- the browser, CLI, and Neovim embodiments still project through the current
  local HTTP adapter on top of those signed families

Remaining:

- add peer-visible `knowledge-evidence` exchange once portable blob carriage
  exists
- decide whether later embodiments should ever bypass the local HTTP adapter
  instead of using it as the delivery surface over the richer runtime

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
follows the frozen operational families. The first relay-visible slice now
ships as whole-family bootstrap export/import for the four attachment-free
families, while peer-visible evidence still waits for the later
storage/carriage decision. CAS now ships as an additive sidecar for signed
envelopes and copied evidence blobs rather than a log replacement rewrite, and
the first embodiment-tightening step now ships through capability metadata plus
adapter-over-runtime doc updates. Source: `DI-fusok`; `DI-guzab`; `DI-voruk`;
`DI-ribek`; `DI-lavuz`; `DI-vabek`; `DI-rovuz`.
