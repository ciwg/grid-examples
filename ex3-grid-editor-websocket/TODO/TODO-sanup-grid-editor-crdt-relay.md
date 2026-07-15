# TODO sanup - grid-editor CRDT relay

## Decision Intent Log

ID: DI-ramuv
Date: 2026-07-13 09:10:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Replace full-document replacement with signed Automerge sync-message relay under `protocols/live-document.md`, persisted to CAS under `<data-root>/cas/**`, while keeping the relay non-canonical.
Intent: Fix correctness and latency problems by carrying real CRDT sync bytes over the PromiseGrid-shaped relay instead of racing whole-document snapshots.
Constraints: The relay signs and relays messages but does not own canonical document state; CAS stores canonical signed envelope bytes addressed by the signed object hash; `<data-root>/message-log.jsonl` remains only as a transitional append-only debug feed.
Affects: `ex2-grid-editor/protocols/live-document.md`, `ex2-grid-editor/cas`, `ex2-grid-editor/crdt`, `ex2-grid-editor/service`, `ex2-grid-editor/cmd/grid-relay`, `<data-root>/cas/**`, `<data-root>/message-log.jsonl`

ID: DI-zegov
Date: 2026-07-13 09:10:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Migrate the browser embodiment to CodeMirror 6 plus Automerge with a minimal `esbuild` pipeline under `web/`, reusing the rooted behavior of `collab-awareness` for real remote cursor rendering.
Intent: Make the browser path actually usable for concurrent editing while keeping the JS build narrow and limited to editor/CRDT glue.
Constraints: `web/src/**` is limited to editor, awareness, Automerge, and thin relay transport glue; built browser assets stay small and static; browser CRDT correctness must not depend on Go/WASM in this slice.
Affects: `ex2-grid-editor/web/index.html`, `ex2-grid-editor/web/style.css`, `ex2-grid-editor/web/app.js`, `ex2-grid-editor/web/src/**`, `ex2-grid-editor/web/package.json`, `ex2-grid-editor/web/package-lock.json`

ID: DI-lumek
Date: 2026-07-13 09:10:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Keep Ed25519 signing for the first CRDT slice, add explicit CRDT HTTP sync endpoints, and split the relay command from the existing app command while keeping Neovim in transitional compatibility mode until a real Go-side CRDT replica exists.
Intent: Land the working CRDT relay/browser path now without coupling it to the separate COSE/CWT migration or pretending the old Neovim polling client is already a real CRDT sidecar.
Constraints: The relay command lives under `cmd/grid-relay`; `cmd/grid-nvim-sidecar` may exist only as a transitional scaffold in this slice; existing Neovim direct-HTTP behavior remains compatibility code rather than the final CRDT design.
Affects: `ex2-grid-editor/cmd/grid-relay`, `ex2-grid-editor/cmd/grid-nvim-sidecar`, `ex2-grid-editor/service`, `ex2-grid-editor/nvim`

Goal: Build the first real CRDT relay/browser slice with signed Automerge sync messages, CAS persistence, and visible remote cursors while keeping the current Neovim path explicit about its transitional status.

- [x] sanup.1 Write the CRDT TE and update the TODO/decision logs.
- [x] sanup.2 Add CAS storage and CRDT message packages.
- [x] sanup.3 Reshape the service into a non-canonical signed relay with explicit sync endpoints.
- [x] sanup.4 Migrate the browser UI to CodeMirror plus Automerge relay sync.
- [x] sanup.5 Add the relay command and transitional Neovim sidecar scaffolding.
- [x] sanup.6 Add deterministic tests for CAS addressing and CRDT relay ingestion.
