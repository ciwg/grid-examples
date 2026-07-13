# TODO kubiv - grid-editor foundation

## Decision Intent Log

ID: DI-lodug
Date: 2026-07-12 17:39:29 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Build `ex2-grid-editor/` as the combined `grid-editor` app repo with a nested Go module, one local Go service, a browser embodiment, and a Neovim embodiment.
Intent: Keep the example app self-contained while making the PromiseGrid-facing contract explicit and shared across both embodiments.
Constraints: Keep all implementation files under `ex2-grid-editor/`; do not use `internal/`; treat browser and Neovim as embodiments of one app contract rather than separate apps.
Affects: `ex2-grid-editor/go.mod`, `ex2-grid-editor/cmd/grid-editor`, `ex2-grid-editor/service`, `ex2-grid-editor/web`, `ex2-grid-editor/nvim`

ID: DI-tofug
Date: 2026-07-12 17:39:29 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Define repo-local draft PromiseGrid-facing specs in `protocols/live-document.md` and `protocols/live-awareness.md`, and keep them as separate protocol families.
Intent: Give the example app an explicit local source of truth for live collaboration while upstream PromiseGrid live-editor specs remain incomplete.
Constraints: The repo-local specs are authoritative for this repo until superseded; document sync and awareness stay separate because they have different cadence and durability pressure.
Affects: `ex2-grid-editor/protocols`, `ex2-grid-editor/protocol`, `ex2-grid-editor/document`, `ex2-grid-editor/awareness`

ID: DI-jilin
Date: 2026-07-12 17:39:29 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Use `Ed25519` as the v1 durable identity baseline, keep the shared runtime/core in Go, use a thin JS browser UI, and keep local HTTP adapter plumbing internal-only.
Intent: Make the first slice actually runnable across browser and Neovim while preserving a PromiseGrid-shaped signed message core.
Constraints: The local service persists its signing seed and append-only message log under `.grid-editor/`; the browser UI does not define protocol truth; the Neovim plugin uses the local service rather than speaking peer protocol directly.
Affects: `ex2-grid-editor/identity`, `ex2-grid-editor/store`, `ex2-grid-editor/service`, `ex2-grid-editor/web`, `ex2-grid-editor/nvim`, `.grid-editor/...`

Goal: Build the first runnable `grid-editor` slice with repo-local live protocol drafts, a signed Go runtime, a browser client, and a Neovim client.

- [ ] kubiv.1 Write the local TODO, TE, and architecture docs that lock the first slice.
- [ ] kubiv.2 Create the nested Go module and shared protocol/core packages.
- [ ] kubiv.3 Implement the local Go service with signed document and awareness message ingestion.
- [ ] kubiv.4 Implement the browser embodiment with internal local-service adapters.
- [ ] kubiv.5 Implement the Neovim embodiment with internal local-service adapters.
- [ ] kubiv.6 Add tests for envelope parsing, document convergence ordering, and awareness ordering.

