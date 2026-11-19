# TODO pukur - grid-editor conference demo fixes

## Decision Intent Log

ID: DI-holoz
Date: 2026-07-19 15:05:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Tighten ex3's demo-time sync path by pushing websocket subscribers immediately on relay ingest, reducing cross-relay peer polling latency, and proving the live browser path's websocket transport in the visible demo surface.
Intent: Remove the multi-second delay and ambiguity that made the live demo look unreliable, while keeping the signed relay message path intact instead of inventing a direct client shortcut.
Constraints: Keep relay-to-relay replication on the signed `/api/peer/messages` feed; keep browser and Neovim live editing on websocket; do not change the document/auth protocol split or bypass the relay with peer-to-peer browser traffic.
Affects: `ex3-grid-editor-websocket/service`, `ex3-grid-editor-websocket/cmd/grid-relay`, `ex3-grid-editor-websocket/cmd/grid-editor`, `ex3-grid-editor-websocket/web`, `ex3-grid-editor-websocket/README.md`

ID: DI-dogub
Date: 2026-07-19 15:05:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add a visible PromiseGrid demo surface in ex3 with a live relay-backed message trace, clickable decoded-message inspection, and an on-screen browser-to-relay architecture diagram.
Intent: Make the conference demo explainable while editing, so viewers can see the actual signed envelopes and relay flow instead of treating the collaboration as an opaque black box.
Constraints: The trace must come from relay-observed envelopes, not synthetic client-local events; the decoded view must preserve access to raw envelope/payload base64; the diagram should stay simple and presentation-friendly.
Affects: `ex3-grid-editor-websocket/service`, `ex3-grid-editor-websocket/web`, `ex3-grid-editor-websocket/README.md`

## Goal

Stabilize the ex3 conference demo and make the messaging path visible enough
to explain live while editing across browser and Neovim clients.

## Tasks

- [x] pukur.1 Fix relay-side synchronization lag so same-relay and cross-relay edits appear immediately during the demo.
- [x] pukur.2 Fix cursor/file-position propagation lag by pushing awareness changes immediately through the relay websocket path.
- [x] pukur.3 Keep the live browser transport on websocket and show that transport mode explicitly in the demo UI.
- [x] pukur.4 Add a live PromiseGrid message panel that shows relay-observed envelopes for the current document during editing.
- [x] pukur.5 Allow clicking a message to inspect its decoded payload and raw envelope/payload base64.
- [x] pukur.6 Add a simple on-screen browser-to-relay data-flow diagram for presentation use.

## Evidence

- `service/server.go`, `service/app.go`, and `service/live_socket.go` now wake websocket subscribers immediately when sync or awareness envelopes land.
- `cmd/grid-relay/main.go` and `cmd/grid-editor/main.go` now poll peer relays at demo-friendly cadence instead of a multi-second delay.
- `service/trace.go` exposes real relay-observed document traffic, including decoded payload metadata and raw base64 envelope/payload bytes.
- `web/index.html`, `web/style.css`, and `web/src/main.js` now render the transport badge, the on-screen architecture diagram, the live trace list, and the clickable inspector view.
