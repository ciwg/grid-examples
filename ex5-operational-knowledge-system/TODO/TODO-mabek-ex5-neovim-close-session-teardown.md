# TODO mabek - ex5 Neovim close session teardown

## Decision Intent Log

ID: DI-mabek
Date: 2026-07-21 16:10:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track the Neovim `:OksClose` teardown bug as its own ex5 review TODO.
Intent: Make the Neovim embodiment close cleanly instead of leaving behind a detached `acwrite` buffer that no longer participates in the live-draft session.
Constraints: Stay inside `ex5`; keep the current Neovim phase thin; fix behavior and add real regression coverage instead of only marker checks.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-mabek-ex5-neovim-close-session-teardown.md`, `ex5-operational-knowledge-system/nvim/lua/oks/init.lua`, `ex5-operational-knowledge-system/nvim/assets_test.go`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`

## Goal

Make `:OksClose` behave like a real session-close action for the Neovim
embodiment instead of only clearing timers and internal state.

## Review finding

`M.close()` currently stops polling and clears state, but it does not wipe the
live draft buffer or read-only inspector buffer. That leaves the user looking
at a dead `acwrite` buffer whose session hooks are gone, which is misleading
for a command explicitly named `:OksClose`.

## Tasks

- [x] mabek.1 Decide and implement the exact close behavior for the live-draft buffer and any open inspector buffer.
- [x] mabek.2 Add Neovim-side regression coverage for the close path instead of relying only on marker tests.
- [x] mabek.3 Update Neovim docs so `:OksClose` is described honestly.

## Status

- done
- derived from the 2026-07-21 deep ex5 review
