# TODO vamor - ex5 Neovim item approval phase

## Decision Intent Log

ID: DI-vamor
Date: 2026-07-21 10:25:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add a small Neovim item approval action that resolves the current item revision through the existing item-detail API and then posts to the existing item-approval API.
Intent: Let terminal-first reviewers act on item review work from Neovim without inventing a new transport or jumping straight to broad workflow mutation in the editor.
Constraints: Stay on the current local HTTP runtime, keep the action item-only in this slice, preserve the revision-aware approval semantics already enforced by the server, and refresh the relevant Neovim view after approval instead of leaving stale review state behind.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md`, `ex5-operational-knowledge-system/TODO/TODO-vamor-ex5-neovim-item-approval-phase.md`, `ex5-operational-knowledge-system/nvim/lua/oks/init.lua`, `ex5-operational-knowledge-system/nvim/assets_test.go`, `ex5-operational-knowledge-system/nvim/approve_item_test.go`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/architecture.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`

## Goal

Add one safe write-side Neovim review action for item approvals, using the
existing ex5 item approval API and current revision lookup.

## Tasks

- [x] vamor.1 Define the next Neovim follow-on after pending review as a small item approval action.
- [x] vamor.2 Add a Neovim command that resolves the current revision and posts an item approval through the existing HTTP API.
- [x] vamor.3 Refresh the relevant Neovim live/inspect/pending context after approval so the terminal view does not stay stale.
- [x] vamor.4 Add Neovim regression coverage for the approval action and command markers.
- [x] vamor.5 Update the ex5 docs to describe the new terminal-side approval action honestly.
