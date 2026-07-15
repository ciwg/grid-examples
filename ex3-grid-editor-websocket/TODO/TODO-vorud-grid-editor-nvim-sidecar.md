# TODO vorud - grid-editor Neovim sidecar

## Decision Intent Log

ID: DI-sulod
Date: 2026-07-13 10:15:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Implement the first real `grid-nvim-sidecar` slice as a Go launcher that runs a bundled Node-based Automerge helper over stdio, while keeping the public command path `cmd/grid-nvim-sidecar`.
Intent: Deliver a real Neovim CRDT replica now using the Automerge implementation already present in the local workspace instead of leaving Neovim on the old snapshot flow while waiting for a future native Go replica.
Constraints: The sidecar must own a local Automerge replica, speak the relay `/sync` and `/awareness` endpoints, and keep the command path and user-facing workflow stable; the Node helper is an implementation detail behind the Go command, not a new top-level product surface.
Affects: `ex2-grid-editor/cmd/grid-nvim-sidecar`, `ex2-grid-editor/nvim`, `ex2-grid-editor/README.md`

ID: DI-gafit
Date: 2026-07-13 10:15:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Migrate the Neovim plugin from direct `curl` polling to a stdio sidecar protocol carrying open, local-text, cursor, and awareness events.
Intent: Keep the plugin thin and move CRDT logic out of Lua while preserving visible peer cursors and cross-embodiment convergence with the browser path.
Constraints: The Lua plugin stays responsible only for buffer lifecycle, cursor rendering, and forwarding local text/cursor state; the sidecar protocol stays internal-only.
Affects: `ex2-grid-editor/nvim/lua/grid_editor/init.lua`, `ex2-grid-editor/cmd/grid-nvim-sidecar`

ID: DI-larok
Date: 2026-07-13 13:45:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Switch the relay document payload unit from interactive Automerge sync-session messages to durable Automerge change packets while keeping the existing signed relay envelope and `/sync` HTTP surface.
Intent: Make relay history replayable for late joiners and remove the peer-session coupling that caused empty-document regressions and stale-state loops in browser and Neovim replicas.
Constraints: The Go relay remains non-canonical and opaque to document text, `message_base64` remains the transport field name for now, and both browser and Neovim embodiments must treat relay document records as append-only Automerge changes.
Affects: `ex2-grid-editor/protocols/live-document.md`, `ex2-grid-editor/service`, `ex2-grid-editor/web`, `ex2-grid-editor/cmd/grid-nvim-sidecar`

ID: DI-ralov
Date: 2026-07-13 20:55:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: The browser embodiment must rebuild document state from relay history on each open instead of reusing a persisted local Automerge replica, and it must render remote cursor decorations with explicit per-peer colors from the relay awareness state.
Intent: Prevent stale browser-local replicas from generating unreplayable document changes after protocol/runtime revisions, and make peer cursor colors visibly trustworthy across browsers.
Constraints: Preference persistence for display name and color may remain local, but document CRDT state must come from the relay; remote cursor rendering must stay thin and client-local.
Affects: `ex2-grid-editor/web/src/automerge-relay.js`, `ex2-grid-editor/web/src/editor.js`, `ex2-grid-editor/web/style.css`

ID: DI-samuv
Date: 2026-07-13 21:10:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: The Neovim embodiment will provide a one-command launcher script and plugin defaults that auto-detect the repo root, default the sidecar to `go run`, and expose explicit status and peer-roster commands.
Intent: Remove the fragile manual setup flow and make Neovim users able to join a document, inspect peers, and understand session state without memorizing `setup()` details.
Constraints: The launcher stays repo-local, the plugin remains thin over the sidecar protocol, and peer/status visibility must derive from the same awareness feed used for cursor rendering.
Affects: `ex2-grid-editor/nvim`, `ex2-grid-editor/scripts/grid-editor-nvim`, `ex2-grid-editor/README.md`

Goal: Replace the Neovim snapshot client with a real sidecar-owned Automerge replica that interoperates with the browser relay path.

- [x] vorud.1 Write the sidecar TE and lock the hybrid helper decision.
- [x] vorud.2 Implement the bundled Node Automerge helper and Go launcher.
- [x] vorud.3 Rewire the Neovim plugin to speak the sidecar stdio protocol.
- [x] vorud.4 Add sidecar/browser interoperability safeguards and docs.
