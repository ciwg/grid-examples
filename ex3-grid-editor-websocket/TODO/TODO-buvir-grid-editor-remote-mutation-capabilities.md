# TODO buvir - grid-editor remote mutation capabilities

## Decision Intent Log

ID: DI-povip
Date: 2026-07-14 20:25:55 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add a repo-local remote-client admission model for `ex3` where an operator-configured bootstrap secret mints short-lived relay-signed document mutation capabilities, and require those capabilities for non-loopback HTTP and websocket mutation.
Intent: Let `ex3` demonstrate real multi-machine browser and Neovim collaboration over published ports without turning websocket into protocol meaning or leaving a long-lived bearer secret on every live mutation path.
Constraints: Keep `live-document`, `live-awareness`, `document-metadata`, and `publish-document` as separate protocol surfaces; keep websocket as carriage only; treat the bootstrap secret as local provisional admission, not a frozen PromiseGrid app API; preserve loopback mutation as a no-secret local fast path.
Affects: `ex3-grid-editor-websocket/service/**`, `ex3-grid-editor-websocket/web/src/**`, `ex3-grid-editor-websocket/cmd/grid-nvim-sidecar/**`, `ex3-grid-editor-websocket/cmd/grid-relay/main.go`, `ex3-grid-editor-websocket/cmd/grid-editor/main.go`, `ex3-grid-editor-websocket/compose.yaml`, `ex3-grid-editor-websocket/README.md`, `ex3-grid-editor-websocket/docs/docker-simulation.md`

ID: DI-talih
Date: 2026-07-14 20:25:55 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Update ex3 docs to reflect the newer upstream PromiseGrid snapshot, explicitly noting that app-facing auth/API guidance remains provisional and that POC20 semantic-model work is now distinct from the POC21 DevOps track.
Intent: Keep ex3's multi-machine documentation aligned with the July 13/14 upstream guide refresh without overstating what PromiseGrid has frozen.
Constraints: Do not claim a settled universal app auth API; keep the note scoped to ex3's current remote mutation bootstrap and websocket carriage story.
Affects: `ex3-grid-editor-websocket/README.md`, `ex3-grid-editor-websocket/docs/docker-simulation.md`, `ex3-grid-editor-websocket/TODO/TODO-buvir-grid-editor-remote-mutation-capabilities.md`

## Goal

Make `ex3-grid-editor-websocket` work as a real multi-machine demo by replacing
loopback-only remote mutation assumptions with short-lived document-scoped
capabilities and normal published-port Docker networking.

## Tasks

- [x] buvir.1 Add relay-side bootstrap-secret and capability-token support for remote mutation.
- [x] buvir.2 Require remote capabilities on HTTP mutation and websocket live paths while preserving loopback local mutation.
- [x] buvir.3 Teach the browser and Neovim sidecar to bootstrap a remote session and present capabilities automatically.
- [x] buvir.4 Move the Docker demo off host networking and document the new cross-machine flow.
- [x] buvir.5 Verify browser and Neovim interoperability locally and over non-loopback HTTP addresses.

## Evidence

- `service/bootstrap.go` and `service/mutation_capability.go` now mint short-lived relay-signed document mutation capabilities from an operator-configured bootstrap token.
- `service/server.go` and `service/live_socket.go` now preserve loopback local mutation while requiring bearer/websocket capabilities for non-loopback mutation.
- `web/src/live-transport.js`, `web/src/automerge-relay.js`, `web/src/relay-awareness.js`, and `web/src/main.js` now bootstrap remote browser sessions and carry the returned capabilities on HTTP and websocket mutation paths.
- `cmd/grid-nvim-sidecar/src/helper.mjs`, `nvim/lua/grid_editor/init.lua`, and `scripts/grid-editor-nvim` now support remote bootstrap tokens for the Neovim embodiment.
- `compose.yaml`, `README.md`, `docs/docker-simulation.md`, and `docs/practical-implementation.md` now describe the published-port Docker demo and the provisional PromiseGrid-aligned bootstrap/capability story.
- Verification passed with `go test ./...`, `errcheck ./...`, `npm test`, `npm run build` in `web/`, and `npm run build` in `cmd/grid-nvim-sidecar/`.
