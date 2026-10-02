# TODO numop - Ex3 Neovim dashboard proof of concept

## Decision Intent Log

ID: DI-loril
Date: 2026-10-01 20:27:36
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Implement the Ex3 Neovim proof of concept as explicit native-Neovim
dashboard commands with a bounded in-memory local activity view.
Intent: Make the existing collaborative editing embodiment visibly coherent
for a demonstration while preserving its one-editor, one-sidecar boundary.
Constraints: Change only Ex3; do not change browser behavior, pCID-selected
protocol behavior, durable storage, trust policy, or the sidecar schema; do
not add a separate Bubble Tea application or automatic dashboard opening.
Affects: `ex3-grid-editor-websocket/nvim/lua/grid_editor/init.lua`,
`ex3-grid-editor-websocket/service/interoperability_test.go`,
`ex3-grid-editor-websocket/README.md`, `ex3-grid-editor-websocket/docs/testing.md`,
and `ex3-grid-editor-websocket/docs/thought-experiments/TE-gozug-nvim-poc-boundary.md`.

## Goal

Make Ex3's existing Neovim collaborative-editor embodiment more legible as a
native proof of concept without changing browser behavior, protocol behavior,
or the existing sidecar boundary.

## Tasks

- [x] numop.1 Lock the native-Neovim dashboard scope through Decision Framing (DI-loril).
- [x] numop.2 Implement the approved dashboard and bounded local activity view.
- [x] numop.3 Add deterministic coverage and update the Ex3 user documentation.

## Evidence

- TE-gozug records the alternatives and current recommendation.
- `go test ./service -run 'TestNeovimPluginRegistersPhaseOneCommands|TestNeovimPluginRendersRemoteDocumentAndPeerMarkers' -count=1` passes.
- `go test ./service -run TestNeovimLauncherEquivalentSessionsObserveEachOther -count=1` passes and proves two launcher-equivalent sessions observe each other.
- `errcheck ./...` passes from `ex3-grid-editor-websocket/`.
- The full `go test ./...` suite is currently blocked by the pre-existing
  `TestHeadlessBrowserRecoversFromBlankSnapshotState` expectation that the
  browser says `browser sync: websocket`; the unchanged browser renders
  `browser sync: polling · awareness: websocket` instead.
