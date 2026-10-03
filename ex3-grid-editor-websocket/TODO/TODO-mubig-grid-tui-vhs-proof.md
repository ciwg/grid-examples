# TODO mubig - Grid TUI VHS collaboration proof

## Decision Intent Log

ID: DI-mutoh
Date: 2026-10-03 00:00:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add the appropriate Charm.land enhancements to Ex3: Harmonica for terminal presence motion, VHS for a reproducible collaboration recording, and Glow for Ex3 help and protocol reading.
Intent: Make the Ex3 Charm embodiment more legible and demonstrable without replacing the existing Automerge, WebSocket, or signed-relay collaboration path.
Constraints: Ex3 only; do not add Wish; keep each tool's behavior local to the terminal/demo surface.
Affects: `ex3-grid-editor-websocket/cmd/grid-tui`, `ex3-grid-editor-websocket/scripts`, `ex3-grid-editor-websocket/demos`, `ex3-grid-editor-websocket/docs`.

## Goal

Create a reproducible VHS recording that demonstrates the Ex3 Grid TUI against
a live relay and makes the collaboration UI easy to show without manual setup.

**Order:** implement after TODO gahin, TODO budad, and TODO molim so the
recording captures the stable final terminal interaction.

## Tasks

- [x] mubig.1 Add a checked-in VHS tape for the Grid TUI collaboration view.
- [x] mubig.2 Add a safe recording launcher that starts an isolated relay and cleans it up.
- [x] mubig.3 Document the generated recording artifact and verification command.
