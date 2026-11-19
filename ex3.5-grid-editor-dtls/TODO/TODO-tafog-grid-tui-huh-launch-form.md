# TODO tafog - Grid TUI Huh launch form

## Decision Intent Log

ID: DI-sidob
Date: 2026-10-03 00:00:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add a Huh-powered `LaunchForm` to Grid TUI for interactive startup.
Intent: Let a person select their display name, peer color, relay URL, and document ID before joining Ex3 collaboration, while retaining explicit flags for automation.
Constraints: The form must complete before the sidecar connects; it must not change Automerge, WebSocket, relay, or browser behavior; colors come from a readable palette rather than requiring a hex code.
Affects: `tui/workspace.go`, `cmd/grid-tui/main.go`, `tui/workspace_test.go`, `go.mod`, `docs/charm-land-features.md`.

## Goal

Make `grid-tui` approachable for a first-time interactive user without
compromising reproducible script launches.

## Tasks

- [x] tafog.1 Add a Huh `LaunchForm` with display-name, color, relay, and document fields.
- [x] tafog.2 Show the form only when one or more required launch flags are omitted.
- [x] tafog.3 Keep fully specified flags as the non-interactive launch path.
- [x] tafog.4 Add deterministic form/default/bypass tests.
- [x] tafog.5 Run focused and full Ex3 verification.
