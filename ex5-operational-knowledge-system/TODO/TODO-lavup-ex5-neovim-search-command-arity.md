# TODO lavup - ex5 Neovim search command arity

## Decision Intent Log

ID: DI-lavup
Date: 2026-07-22 00:20:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track the `:OksSearch` command-surface mismatch as its own `016` child TODO.
Intent: Keep the next terminal fix focused on the real user-facing gap where the documented multi-token search syntax may not match the actual Ex command registration.
Constraints: Reuse the existing `:OksSearch` implementation and shared `/api/search` contract; do not invent a second Neovim search command.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-lavup-ex5-neovim-search-command-arity.md`, `ex5-operational-knowledge-system/TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md`

## Goal

Make the shipped `:OksSearch` Ex command match the documented structured
search syntax.

## Links

- Parent follow-on: `TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md` (`016`)

## Tasks

- [x] lavup.1 Fix `:OksSearch` command registration so documented multi-token syntax works.
- [x] lavup.2 Add regression coverage for the actual Ex command surface.
- [x] lavup.3 Update docs if the command shape changes.
