# TODO jabup - ex5 Neovim revision snapshot phase

## Decision Intent Log

ID: DI-jabup
Date: 2026-07-22 01:10:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track durable revision creation from Neovim as its own `016` child TODO.
Intent: Close the biggest remaining terminal authoring gap by letting an editor-first user turn a live draft into a durable knowledge-item revision without leaving the Neovim embodiment.
Constraints: Reuse the existing revision-creation API and staged Neovim model; do not reopen the transport decision or broaden this into a full browser-parity project.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-jabup-ex5-neovim-revision-snapshot-phase.md`, `ex5-operational-knowledge-system/TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md`, `ex5-operational-knowledge-system/docs/terminal-capability-matrix.md`, `ex5-operational-knowledge-system/docs/user-guide.md`

## Goal

Add durable revision creation to the Neovim embodiment.

## Links

- Parent follow-on: `TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md` (`016`)

## Tasks

- [x] jabup.1 Define the smallest Neovim command shape for durable revision creation.
- [x] jabup.2 Reuse the existing revision snapshot route from the editor.
- [x] jabup.3 Refresh the relevant Neovim context after snapshot creation.
- [x] jabup.4 Add headless coverage and docs.
