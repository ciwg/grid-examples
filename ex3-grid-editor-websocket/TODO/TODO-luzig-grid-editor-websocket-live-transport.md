# TODO luzig - grid-editor websocket live transport

## Decision Intent Log

ID: DI-vubih
Date: 2026-07-14 16:55:02 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Make `ex3-grid-editor-websocket` use websocket transport for the browser live `sync` and `awareness` flows while keeping metadata and publish on their existing HTTP endpoints and keeping the local loopback-only mutation rule.
Intent: Make the copied `ex3` example honestly behave like the websocket-oriented variant without collapsing the existing protocol split or widening the security surface.
Constraints: `live-document` and `live-awareness` remain distinct flows; websocket frames are transport carriers, not new protocol meaning; metadata and publish remain HTTP; the existing Neovim sidecar path can keep using HTTP polling for now; remote unauthenticated mutation stays out of scope.
Affects: `ex3-grid-editor-websocket/service`, `ex3-grid-editor-websocket/web`, `ex3-grid-editor-websocket/README.md`, `ex3-grid-editor-websocket/docs/practical-implementation.md`, `ex3-grid-editor-websocket/service/testdata/browser-harness.mjs`, `ex3-grid-editor-websocket/service/*tests*`

ID: DI-bitus
Date: 2026-07-14 16:55:02 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Extend ex3's websocket live transport to the Neovim sidecar so both browser and sidecar use websocket for live `sync` and `awareness`, while metadata and publish remain HTTP.
Intent: Finish the websocket transition for ex3's live editing path instead of leaving the browser and Neovim embodiments on different live transports.
Constraints: Preserve the sidecar's stdin/stdout control contract with Neovim; keep `live-document` and `live-awareness` as separate transport channels; keep polling only as a compatibility fallback when websocket is unavailable; do not widen metadata or publish into websocket in this patch.
Affects: `ex3-grid-editor-websocket/cmd/grid-nvim-sidecar/**`, `ex3-grid-editor-websocket/service/*tests*`, `ex3-grid-editor-websocket/README.md`, `ex3-grid-editor-websocket/docs/practical-implementation.md`
Supersedes: DI-vubih

ID: DI-vipat
Date: 2026-07-14 20:25:55 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Document in the ex3 README that the websocket transport split is aligned with current PromiseGrid dev-guide direction, while remaining a repo-local provisional implementation rather than a frozen upstream app API.
Intent: Keep ex3 readers from confusing "PromiseGrid-aligned" with "already standardized," while making the transport and capability-token direction explicit for future multi-machine work.
Constraints: Do not claim a frozen upstream websocket or editor-auth spec; keep the note tied to the existing ex3 split between websocket live transport and HTTP metadata/publish surfaces.
Affects: `ex3-grid-editor-websocket/README.md`, `ex3-grid-editor-websocket/TODO/TODO-luzig-grid-editor-websocket-live-transport.md`

## Goal

Enable the browser live editing path in `ex3-grid-editor-websocket` to use
websocket transport for document sync and awareness updates while preserving
the current protocol boundaries and existing HTTP-only surfaces that are not
part of the live session.

## Tasks

- [x] luzig.1 Add websocket upgrade and frame handling to the ex3 relay for live browser traffic.
- [x] luzig.2 Add websocket live endpoints for `sync` and `awareness` while preserving the loopback-only mutation rule.
- [x] luzig.3 Switch the browser live clients to prefer websocket transport and keep a polling fallback only when websocket is unavailable.
- [x] luzig.4 Rebuild the bundled browser asset so the served ex3 app matches the source changes.
- [x] luzig.5 Verify browser-to-relay live transport works over websocket and does not break the existing sidecar interoperability path.

## Evidence

- `service/websocket.go` now upgrades HTTP requests to websocket and reads/writes JSON websocket frames for ex3's live browser path.
- `service/live_socket.go` exposes separate `sync-socket` and `awareness-socket` handlers that preserve the live-document/live-awareness split.
- `web/src/automerge-relay.js` and `web/src/relay-awareness.js` now prefer websocket transport for browser live traffic and fall back to polling only when websocket is unavailable.
- `cmd/grid-nvim-sidecar/src/helper.mjs` now prefers websocket transport for sidecar live traffic and falls back to polling only when websocket is unavailable.
- `service/testdata/browser-harness.mjs` reports the selected transport mode, and `service/interoperability_test.go` now asserts websocket is in use for both the browser and Neovim sides of the browser/Neovim interoperability test.
- `web/app.js` was rebuilt from the updated browser source with `npm run build`.
- Verification passed with `go test ./...`, `errcheck ./...`, and `npm test`.
