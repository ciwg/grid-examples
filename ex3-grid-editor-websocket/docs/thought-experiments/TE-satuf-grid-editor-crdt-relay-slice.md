# grid-editor CRDT relay slice

TE ID: TE-satuf
## Status
decided

## Decision under test

How should `grid-editor` move from full-document replacement to a real collaborative editor while staying PromiseGrid-shaped and Go-centered?

## Assumptions

- `promisegrid-dev-guide` and `wire-lab` remain the orientation sources of truth.
- The current repo already has a signed relay skeleton, a browser embodiment, and a Neovim embodiment.
- The old collaboration behavior that actually worked came from Automerge plus CodeMirror plus `collab-awareness`.
- A real Go Automerge replica is not already present in the workspace.

## Alternatives

1. Keep full-document replacement and only tune polling.
2. Make the Go service the canonical CRDT host.
3. Use Automerge replicas in embodiments and keep the Go service as a non-canonical relay/signing agent.
4. Stop until a full upstream PromiseGrid live-CRDT stack exists.

## Scenario analysis

### Normal browser to browser editing

Alternative 1 keeps dropped text and races because whole-document replacement is still the unit of exchange.

Alternative 2 fixes merge correctness, but the relay quietly becomes the authority and the browser becomes a thin client.

Alternative 3 preserves local CRDT truth in each browser, lets the relay sign and persist exact message bytes, and keeps the working Automerge behavior closest to the old system.

Alternative 4 avoids short-term design debt but leaves the current broken browser experience in place.

### Concurrent edits and stale local views

Alternative 1 fails directly because overlapping full-document writes race.

Alternative 2 can converge, but only by centralizing the replica into the relay.

Alternative 3 converges if the relay carries signed Automerge sync messages and clients retain per-peer sync state.

Alternative 4 does not solve the existing bug.

### Persistence and replay

Alternative 1 persists snapshots, not the actual collaborative intent stream.

Alternative 2 can persist one canonical CRDT, but this couples durability to relay authority.

Alternative 3 persists signed envelopes to CAS and optionally mirrors them into a debug log, so replay and verification share one signed-object path without claiming canonical document ownership.

Alternative 4 defers the problem.

### Browser and Neovim together

Alternative 1 keeps both embodiments equally weak.

Alternative 2 makes it easier to keep browser and Neovim aligned, but only by collapsing them into relay clients.

Alternative 3 is the best long-term shape, but the current workspace only has a real Automerge implementation for the browser and the old Node-based Neovim helper.

Alternative 4 keeps the design pure at the cost of shipping nothing useful.

### Long-horizon evolution

Alternative 1 would have to be discarded entirely.

Alternative 2 makes later in-browser grid agents harder because the relay already owns the document.

Alternative 3 leaves room for later Go/WASM browser bridges and later COSE/CWT signing without invalidating the core relay/browser CRDT shape.

Alternative 4 depends on upstream timing.

## Conclusions

- Reject alternative 1 because polling a broken unit of exchange does not fix correctness.
- Reject alternative 2 because it violates the non-canonical relay constraint.
- Reject alternative 4 because the current app already needs a working CRDT path.
- Keep alternative 3.

## Surviving implementation consequences

- `live-document` must carry signed Automerge sync-message bytes, not whole-document replacement payloads.
- The relay must persist signed canonical envelope bytes to CAS and use those bytes as the replay truth path.
- The browser should move to CodeMirror plus Automerge immediately.
- Neovim remains transitional in this slice because the workspace does not already contain a Go Automerge replica.

## Implications for TODOs and DIs

- Lock the CRDT relay slice in `TODO/TODO-sanup-grid-editor-crdt-relay.md`.
- Supersede the old browser `textarea` path with CodeMirror and Automerge.
- Keep the current Neovim path explicit as compatibility code until a real sidecar implementation exists.

## Decision status

locked via DI-ramuv, DI-zegov, and DI-lumek
