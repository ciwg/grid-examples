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

## Highest-impact additions

- [ ] maduh.1 Add Bubbles Help and key bindings for a discoverable shortcut
  overlay. Replace the footer-only shortcut reminder while preserving the
  existing Alt-key menu behavior.
- [ ] maduh.2 Add a Bubbles List for choosing relay-known documents, published
  versions, and catalog results instead of requiring users to type identifiers.
- [ ] maduh.3 Add Bubbles File picker to the import flow so users do not have
  to enter an absolute source path manually.

## Supporting usability additions

- [ ] maduh.4 Add a Bubbles Spinner for connection, import, export, publish,
  and relay-refresh work that is genuinely in progress.
- [ ] maduh.5 Add a Bubbles Table for catalog and published-version metadata
  when those records have enough fields to benefit from columns.
- [ ] maduh.6 Evaluate a Bubbles Viewport for the activity/history surface;
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
