# TODO molim - Grid TUI Harmonica presence motion

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

Use Harmonica motion to make remote typing and presence changes easier to spot in
the Grid TUI peer legend without turning awareness into durable document state.

**Order:** implement after TODO gahin; the menu establishes the stable workspace
navigation around the animated peer legend.

## Tasks

- [x] molim.1 Add Harmonica as the Grid TUI presence-motion dependency.
- [x] molim.2 Animate a short, bounded typing/presence pulse in the peer legend.
- [x] molim.3 Add deterministic tests for presence-motion state transitions.
- [x] molim.4 Verify the animation leaves relay, Automerge, and awareness payloads unchanged.
