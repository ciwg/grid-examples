# TODO jufip - Stop Grid TUI idle redraws

## Decision Intent Log

ID: DI-tubol
Date: 2026-10-04 17:30:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Replace Grid TUI's continuous typing animation with a static typing marker.
Intent: Keep remote typing awareness visible without making the terminal workspace
continuously repaint and throb while idle.
Constraints: Preserve the existing peer name, color, cursor, and typing-state
semantics; do not change relay, Automerge, WebSocket, or awareness protocols.
Affects: `tui/workspace.go`, `tui/workspace_test.go`, `docs/testing.md`,
`docs/charm-land-features.md`, this TODO, and the root TODO index.

## Goal

Make the Grid TUI usable in ordinary terminals by eliminating its unconditional
15-FPS redraw loop.

## Tasks

- [x] jufip.1 Replace the animated marker with a stable typing marker.
- [x] jufip.2 Cover the stable remote-typing marker with a model test.
- [x] jufip.3 Correct the Charm feature and test documentation.
