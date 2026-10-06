# TODO maduh - Grid TUI Bubbles usability components

## Decision Intent Log

ID: DI-digis
Date: 2026-10-06 16:14:14 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track practical Charm Bubbles additions for Ex3 Grid TUI, with
priority on features that materially improve usability and discoverability.
Intent: Make the terminal collaborator feel like a usable product without
changing Ex3's shared document, awareness, relay, Automerge, WebSocket,
browser, or Neovim behavior.
Constraints: Keep additions local to Grid TUI; preserve the existing custom
Bubble Tea menu layer; add deterministic tests for each delivered component;
do not add a component merely to increase the Charm library count.
Affects: `tui/`, `docs/testing.md`, `README.md`, this TODO, and both TODO
indexes.

ID: DI-vujub
Date: 2026-10-06 17:02:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Open focused Bubbles panels from the existing Grid TUI menu bar for
help, relay result selection, file importing, metadata inspection, and
activity review; show a spinner only while a real local command is active.
Intent: Improve discoverability and terminal ergonomics without replacing the
custom menu vocabulary or changing the shared Ex3 collaboration contract.
Constraints: Local transient UI state only; no fake progress indicator; retain
the textarea as the shared-document editor; verify with deterministic tests
and VHS recordings.
Affects: `tui/workspace.go`, `tui/panels.go`, `tui/workspace_test.go`,
`README.md`, `docs/testing.md`, `demos/`, and `scripts/`.

## Highest-impact additions

- [x] maduh.1 Add Bubbles Help and key bindings for a discoverable shortcut
  overlay. Replace the footer-only shortcut reminder while preserving the
  existing Alt-key menu behavior.
- [x] maduh.2 Add a Bubbles List for choosing relay-known documents, published
  versions, and catalog results instead of requiring users to type identifiers.
- [x] maduh.3 Add Bubbles File picker to the import flow so users do not have
  to enter an absolute source path manually.

## Supporting usability additions

- [x] maduh.4 Add a Bubbles Spinner for connection, import, export, publish,
  and relay-refresh work that is genuinely in progress.
- [x] maduh.5 Add a Bubbles Table for catalog and published-version metadata
  when those records have enough fields to benefit from columns.
- [x] maduh.6 Add a Bubbles Viewport for the activity/history surface while
  do not replace the editor textarea's scrolling behavior without separate
  design and interaction validation.
- [ ] maduh.7 Add a Bubbles Progress bar only after an import, export, or
  publish operation exposes meaningful incremental progress.

## Deferred because they do not currently improve collaboration

- [ ] maduh.8 Keep Bubbles Paginator deferred unless a long, bounded result
  set proves clearer as pages than as a searchable list.
- [ ] maduh.9 Keep Bubbles Timer and Stopwatch deferred; neither is currently
  a collaboration or editor usability need.

## Design boundary

The top application menu remains Grid TUI code built on Bubble Tea's
Model/Update/View loop and styled with Lip Gloss. Bubbles should supply focused
components inside that experience: textarea and text input already do so;
Help, List, File picker, Spinner, Table, Viewport, and Progress may follow
when their user-visible job is clear.
