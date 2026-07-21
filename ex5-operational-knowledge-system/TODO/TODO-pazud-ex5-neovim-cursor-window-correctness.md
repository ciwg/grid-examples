# TODO pazud - ex5 Neovim cursor-window correctness

## Decision Intent Log

ID: DI-pazud
Date: 2026-07-21 13:25:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track the Neovim live-draft cursor/head bug as its own ex5 embodiment fix TODO.
Intent: Keep Neovim presence and live-draft pushes tied to the actual draft window instead of whichever split is currently focused.
Constraints: Stay on the current HTTP live-draft embodiment; preserve the read-only inspector splits; add regression coverage and Neovim docs in the same pass.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-pazud-ex5-neovim-cursor-window-correctness.md`, `ex5-operational-knowledge-system/nvim/lua/oks/init.lua`, `ex5-operational-knowledge-system/nvim/assets_test.go`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`, `ex5-operational-knowledge-system/README.md`

## Goal

Ensure Neovim presence and body pushes compute cursor/head from the live-draft
window and buffer, not from the currently focused window after the user opens
an inspector split.

## Review finding

`current_cursor_offset()` reads lines from `M.state.bufnr` but gets the cursor
from `vim.api.nvim_win_get_cursor(0)`. After switching to an inspector or other
split, `0` can point at the wrong window, which makes the shared presence state
and pushed cursor offsets incorrect.

## Tasks

- [x] pazud.1 Track the actual live-draft window or otherwise compute cursor/head from the correct Neovim window-buffer pair.
- [x] pazud.2 Add regression coverage and document the corrected Neovim presence behavior.

## Status

- done
- derived from the 2026-07-21 deep ex5 review
