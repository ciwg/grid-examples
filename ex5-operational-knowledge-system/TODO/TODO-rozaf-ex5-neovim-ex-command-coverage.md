# TODO rozaf - ex5 Neovim Ex command coverage

## Decision Intent Log

ID: DI-rozaf
Date: 2026-07-22 00:21:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track direct headless coverage of the actual `:Oks...` command surfaces as its own `016` child TODO.
Intent: Close the remaining test gap where Lua functions are covered but Ex command parsing and registration can still drift without detection.
Constraints: Keep this focused on command-surface coverage, not a broad rework of the Neovim test harness.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-rozaf-ex5-neovim-ex-command-coverage.md`, `ex5-operational-knowledge-system/TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md`

## Goal

Add headless tests that exercise the shipped `:Oks...` commands directly.

## Links

- Parent follow-on: `TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md` (`016`)

## Tasks

- [x] rozaf.1 Decide the minimum Ex-command set that needs direct headless coverage.
- [x] rozaf.2 Add command-level coverage for the chosen search/review/action commands.
- [x] rozaf.3 Keep marker tests and behavior tests aligned with the real command surface.
