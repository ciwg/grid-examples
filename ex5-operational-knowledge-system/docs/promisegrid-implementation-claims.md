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

`ex5` ships as part of the PromiseGrid example set, but it is not yet
PromiseGrid-complete at the runtime/wire layer. Source: `DI-sobek`.

Today it ships:

- one local Go runtime
- append-only operational event history
- local durable draft and attachment storage
- projected read/query views over that history
- browser, CLI, and Neovim embodiments over the same local HTTP adapter

Today it does **not** yet ship:

- wire-visible `grid([42(pCID), payload, proof])` envelopes
- frozen runtime behavior selected by a shipped `pCID`
- relay-visible peer exchange
- signing/verification behavior as part of the ex5 operational workflow

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

### 4. The local HTTP API is an embodiment adapter

`GET /api/*` and `POST /api/*` routes are the shipped local adapter surface for
browser, CLI, and Neovim. They are the current implementation contract for
those embodiments, but they are not claimed as the final PromiseGrid peer
contract. Source: `DI-sobek`.

### 5. Protocol-family language is shipped PromiseGrid framing, not yet a frozen wire promise

References to protocol families, protocol docs, and `pCID`-selected meaning
are part of the shipped PromiseGrid framing for ex5. They do not yet mean that
the current runtime emits or consumes frozen PromiseGrid wire artifacts.
Source: `DI-sobek`.

## What the shipped implementation does not yet promise

The current shipped ex5 runtime does not yet promise:

- that a specific frozen protocol document is the runtime-selected source of
  public wire meaning
- that envelopes are signed or verified as part of the ex5 workflow
- that the local HTTP route names are a stable PromiseGrid peer contract
- that relay transport, peer exchange, or CAS-backed envelope storage are
  already implemented

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

## Follow-on decision

The next real PromiseGrid implementation question is tracked separately in
[TODO ragup](../TODO/TODO-ragup-ex5-promisegrid-wire-slice-decision.md): if
and when ex5 should begin a real wire-level slice around frozen `pCID`
handling, signed envelopes, and related runtime/storage behavior. Source:
`DI-ragup`; `DI-sobek`.
