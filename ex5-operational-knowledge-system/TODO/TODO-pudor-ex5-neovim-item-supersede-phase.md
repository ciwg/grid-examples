# TODO pudor - ex5 Neovim item supersede phase

## Decision Intent Log

ID: DI-pudor
Date: 2026-07-21 11:00:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add a small Neovim item supersede action that posts directly to the existing item supersede API and refreshes the relevant terminal view afterward.
Intent: Let terminal-first reviewers complete the next obvious item lifecycle step from Neovim after approval and pending-review work, without turning the editor into a broad mutation surface all at once.
Constraints: Stay on the current local HTTP runtime, keep the action item-only in this slice, reuse the existing item supersede API, and refresh live, inspect, or pending-review state after superseding so the terminal view does not stay stale.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md`, `ex5-operational-knowledge-system/TODO/TODO-pudor-ex5-neovim-item-supersede-phase.md`, `ex5-operational-knowledge-system/nvim/lua/oks/init.lua`, `ex5-operational-knowledge-system/nvim/assets_test.go`, `ex5-operational-knowledge-system/nvim/supersede_item_test.go`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/architecture.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`

## Goal

Add one safe write-side Neovim lifecycle action for superseding a knowledge
item, using the existing ex5 item supersede API.

## Tasks

- [x] pudor.1 Define the next Neovim follow-on after run approval as a small item supersede action.
- [x] pudor.2 Add a Neovim command that posts an item supersede through the existing HTTP API.
- [x] pudor.3 Refresh the relevant Neovim live, inspect, or pending context after superseding so the terminal view does not stay stale.
- [x] pudor.4 Add Neovim regression coverage for the item supersede action and command markers.
- [x] pudor.5 Update the ex5 docs to describe the new terminal-side item supersede action honestly.
