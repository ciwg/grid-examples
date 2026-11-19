# TODO budad - Grid Glow help reader

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

Provide a Glow-powered terminal reader for the Ex3 README and protocol/help
documents so operators can inspect the implementation boundary without leaving
the terminal workflow.

**Order:** implement after TODO gahin so Glow is opened from the Grid TUI Help
menu instead of becoming another disconnected command.

## Tasks

- [x] budad.1 Add a Glow launcher with an explicit document selector.
- [x] budad.2 Add Make targets and README instructions.
- [x] budad.3 Verify missing Glow is reported with a usable install command.
