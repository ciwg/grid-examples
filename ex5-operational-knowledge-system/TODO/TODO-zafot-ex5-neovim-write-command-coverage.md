# TODO zafot - ex5 Neovim write command coverage

## Decision Intent Log

ID: DI-zafot
Date: 2026-07-22 01:11:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track direct Ex-command coverage of the Neovim write-side actions as its own `016` child TODO.
Intent: Make sure `:OksApproveItem`, `:OksApproveRun`, and `:OksSupersedeItem` are covered through the actual shipped command surface, not only through the underlying Lua functions.
Constraints: Keep this slice focused on command-level coverage, not on broader Neovim feature expansion.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-zafot-ex5-neovim-write-command-coverage.md`, `ex5-operational-knowledge-system/TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md`, `ex5-operational-knowledge-system/nvim/approve_item_test.go`, `ex5-operational-knowledge-system/nvim/approve_run_test.go`, `ex5-operational-knowledge-system/nvim/supersede_item_test.go`

## Goal

Exercise the actual write-side `:Oks...` commands directly in headless Neovim
tests.

## Links

- Parent follow-on: `TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md` (`016`)

## Tasks

- [x] zafot.1 Convert item approval coverage to use `:OksApproveItem`.
- [x] zafot.2 Convert run approval coverage to use `:OksApproveRun`.
- [x] zafot.3 Convert item supersede coverage to use `:OksSupersedeItem`.
- [x] zafot.4 Keep marker and behavior tests aligned with the shipped commands.
