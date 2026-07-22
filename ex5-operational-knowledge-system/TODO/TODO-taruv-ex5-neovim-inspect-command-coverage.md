# TODO taruv - ex5 Neovim inspect command coverage

## Decision Intent Log

ID: DI-taruv
Date: 2026-07-22 02:06:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track direct Ex-command coverage of the older Neovim inspect surfaces as its own `016` child TODO.
Intent: Make sure `:OksInspect`, `:OksInspectRun`, and `:OksInspectEntity` are covered through the real shipped command layer, not only through underlying Lua helpers.
Constraints: Keep this slice focused on command-surface coverage; do not broaden it into new inspect features or deeper workflow expansion.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-taruv-ex5-neovim-inspect-command-coverage.md`, `ex5-operational-knowledge-system/TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md`, `ex5-operational-knowledge-system/nvim/inspect_test.go`, `ex5-operational-knowledge-system/nvim/assets_test.go`

## Goal

Exercise the actual older `:OksInspect...` commands directly in headless Neovim tests.

## Links

- Parent follow-on: `TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md` (`016`)

## Tasks

- [x] taruv.1 Convert item inspection coverage to use `:OksInspect`.
- [x] taruv.2 Convert run inspection coverage to use `:OksInspectRun`.
- [x] taruv.3 Convert entity inspection coverage to use `:OksInspectEntity`.
- [x] taruv.4 Keep marker and behavior tests aligned with the shipped inspect commands.
