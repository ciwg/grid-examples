# TODO bafor - ex5 Neovim run approval phase

## Decision Intent Log

ID: DI-bafor
Date: 2026-07-21 10:45:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add a small Neovim run approval action that posts directly to the existing run approval API and refreshes the relevant terminal view afterward.
Intent: Let terminal-first reviewers act on run review work from Neovim after finding it in `:OksPending` or `:OksInspectRun`, without broadening the editor into a full workflow surface all at once.
Constraints: Stay on the current local HTTP runtime, keep the action run-only in this slice, reuse the existing run approval API, and refresh pending or run-inspector state after approval so the terminal view does not stay stale.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md`, `ex5-operational-knowledge-system/TODO/TODO-bafor-ex5-neovim-run-approval-phase.md`, `ex5-operational-knowledge-system/nvim/lua/oks/init.lua`, `ex5-operational-knowledge-system/nvim/assets_test.go`, `ex5-operational-knowledge-system/nvim/approve_run_test.go`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/architecture.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`

## Goal

Add one safe write-side Neovim review action for run approvals, using the
existing ex5 run approval API.

## Tasks

- [x] bafor.1 Define the next Neovim follow-on after item approval as a small run approval action.
- [x] bafor.2 Add a Neovim command that posts a run approval through the existing HTTP API.
- [x] bafor.3 Refresh the relevant Neovim inspect or pending context after approval so the terminal view does not stay stale.
- [x] bafor.4 Add Neovim regression coverage for the run approval action and command markers.
- [x] bafor.5 Update the ex5 docs to describe the new terminal-side run approval action honestly.
