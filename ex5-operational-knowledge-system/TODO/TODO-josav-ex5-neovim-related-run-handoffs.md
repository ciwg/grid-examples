# TODO josav - ex5 Neovim related-run handoffs

## Decision Intent Log

ID: DI-josav
Date: 2026-07-21 23:10:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add direct `:OksInspectRun` handoff hints anywhere Neovim inspectors list related runs.
Intent: Keep editor-side review navigation consistent so users can move from item, place, resource, or responsibility context into specific runs without mentally reconstructing the next command.
Constraints: Reuse the existing inspectors, keep the change read-only, and link the slice back to deferred TODO `016`.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-josav-ex5-neovim-related-run-handoffs.md`, `ex5-operational-knowledge-system/TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md`, `ex5-operational-knowledge-system/nvim/lua/oks/init.lua`, `ex5-operational-knowledge-system/nvim/inspect_test.go`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`

## Goal

Make Neovim related-run sections consistently actionable.

## Links

- Parent follow-on: `TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md` (`016`)

## Tasks

- [x] josav.1 Add `:OksInspectRun` hints to item related-run sections.
- [x] josav.2 Add `:OksInspectRun` hints to place/resource/responsibility related-run sections.
- [x] josav.3 Extend headless Neovim behavior coverage.
- [x] josav.4 Update terminal-first docs.
