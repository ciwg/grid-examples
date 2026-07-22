# TODO zorik - ex5 Neovim inspect behavior tests

## Decision Intent Log

ID: DI-zorik
Date: 2026-07-21 19:05:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add real behavior tests for the older Neovim inspect commands instead of relying mainly on marker presence.
Intent: Bring the older Neovim inspect surfaces up to the same test depth as the newer search, pending, approval, and supersede terminal features.
Constraints: Prefer headless Neovim behavior tests, keep the checks deterministic, and link the work back to deferred TODO `016`.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-zorik-ex5-neovim-inspect-behavior-tests.md`, `ex5-operational-knowledge-system/TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md`, `ex5-operational-knowledge-system/nvim/assets_test.go`, `ex5-operational-knowledge-system/nvim/*.go`, `ex5-operational-knowledge-system/nvim/lua/oks/init.lua`

## Goal

Add real behavior coverage for Neovim inspect flows.

## Links

- Parent follow-on: `TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md` (`016`)

## Tasks

- [x] zorik.1 Identify which inspect surfaces still rely only on marker tests.
- [x] zorik.2 Add headless behavior tests for the missing inspect flows.
- [x] zorik.3 Keep the lighter asset test only as a shipped-command backstop.
