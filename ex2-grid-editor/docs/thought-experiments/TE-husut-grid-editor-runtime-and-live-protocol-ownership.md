# grid-editor runtime and live protocol ownership

TE ID: TE-husut

## Status

decided

## Decision under test

How `grid-editor` should become a PromiseGrid-native collaborative editor example while upstream PromiseGrid lacks frozen live-collaboration specs.

## Assumptions

- `promisegrid-dev-guide/README.md` and `wire-lab/DEV-GUIDE-RESOURCES.md` are the orientation sources of truth for current PromiseGrid direction.
- This repo is allowed to write its own first-pass live-collaboration spec documents.
- The first meaningful slice must support both browser and Neovim.
- The implementation language for the shared runtime/core is Go.
- A local service is acceptable as embodiment-local plumbing so long as the peer-visible contract remains explicit and PromiseGrid-shaped.

## Alternatives

### Alt A: direct embodiment peers over inherited legacy protocols

Keep the browser and Neovim embodiments close to the legacy repos and let them keep using the old Automerge/Awareness wire contracts while new docs merely describe future PromiseGrid aspirations.

### Alt B: one local Go service with repo-local PromiseGrid-facing live protocols

Make one Go runtime the local embodiment-facing service, define repo-local `live-document` and `live-awareness` specs, and have browser and Neovim talk to the service through internal adapters while peer-visible collaboration uses the signed grid envelope shape.

### Alt C: wait for upstream live-document and live-awareness specs

Do not define local live specs here. Stop at design docs and defer implementation until upstream PromiseGrid publishes the necessary live collaboration specs.

### Alt D: browser UI and Neovim plugin each embed their own full peer protocol stack

Avoid the local service and make each embodiment directly own the signed peer protocol, identity handling, persistence, and peer sync logic.

## Scope and systems affected

- `ex2-grid-editor/` Go module and runtime
- browser UI under `web/`
- Neovim embodiment under `nvim/`
- repo-local protocol docs under `protocols/`
- local runtime data under `.grid-editor/`

## Scenario analysis

### S1 — normal browser to Neovim collaboration on one host

Alice opens the browser UI and Neovim against the same local `grid-editor` service. With Alt A, the two embodiments inherit legacy transport assumptions and the repo still has no honest PromiseGrid-facing contract. Alt B gives both embodiments a single local source of truth for document state, awareness state, identity, and protocol framing, while keeping the local adapter plumbing out of the public contract. Alt C does not produce a runnable example. Alt D duplicates identity, persistence, and sync logic across two embodiments immediately.

### S2 — two hosts exchanging signed live document traffic

Alice runs a local service on one host and Bob runs another. Their services need a peer-visible contract that can be logged, replayed, verified, and migrated later. Alt A still has no PromiseGrid-facing contract and preserves the wrong source of truth. Alt B yields signed repo-local draft protocols that can be exchanged, replayed, and superseded later by upstream frozen specs. Alt C blocks all real inter-host work. Alt D can also exchange signed traffic, but every embodiment has to reimplement durable identity, append-only logging, and peer ingestion separately.

### S3 — failure, restart, and log replay

Alice restarts her local service after a crash. Alt B makes the append-only message log and signing key live in one place, so the document and awareness projections can be rebuilt deterministically from that local evidence. Alt D must teach both browser and Neovim how to persist and rebuild the same message history, which creates duplication and drift pressure. Alt A inherits legacy stores that are not the intended PromiseGrid-facing truth. Alt C still has no runnable slice.

### S4 — mixed-version peers and future upstream migration

Alice upgrades to a newer `grid-editor` build while Bob still runs the previous one. Alt B keeps the peer contract explicit in repo-local specs and lets version drift be reasoned about at the protocol-doc level. It also gives this repo a clean place to record later supersedence toward upstream frozen specs. Alt A delays that clarity. Alt D multiplies migration surface across embodiments. Alt C avoids the problem only by not shipping anything.

### S5 — durability versus ephemeral pressure

Document content wants durable replay; awareness wants lower-cost ephemeral freshness. Alt B keeps `live-document` and `live-awareness` separate so the document stream can stay append-only and replayable while awareness can remain a latest-state projection with append-only local evidence if desired. Alt A preserves the split accidentally but without PromiseGrid-facing clarity. Alt D can also split them, but still duplicates logic. Alt C gives no implementation guidance.

### S6 — constrained host differences

The browser has DOM and fetch; Neovim has buffers and local shell/process access. Alt B isolates those differences inside embodiment-local adapters while the shared Go service owns signed protocol logic and persistence. Alt D pushes the same non-trivial cryptographic and protocol logic into each constrained embodiment separately. Alt A keeps the same split but around the wrong protocol center.

## Conclusions

Rejected:

- Alt A because it preserves the legacy protocols as de facto truth and fails the "PromiseGrid-native example" goal.
- Alt C because it blocks implementation entirely.
- Alt D because it duplicates the hardest state, identity, and protocol logic across embodiments too early.

Surviving alternative:

- Alt B.

## Implications for open TODOs and pending DIs

- `TODO-kubiv` should create a nested Go module in `ex2-grid-editor/`.
- The first protocol docs should live in `protocols/live-document.md` and `protocols/live-awareness.md`.
- The local runtime should persist identity and append-only message history under `.grid-editor/`.
- Browser and Neovim local adapter protocols are implementation detail, not the public contract surface.

## Decision status

locked by `DI-lodug`, `DI-tofug`, and `DI-jilin`

