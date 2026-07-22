# TODO sivok - ex5 Neovim problem-review phase

## Decision Intent Log

ID: DI-sivok
Date: 2026-07-21 23:55:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add a read-only Neovim grouped problem-review buffer over the existing `/api/problem-review` projection.
Intent: Close the biggest remaining terminal-side review gap by giving Neovim users the same grouped hotspot review shape that browser and CLI users already have.
Constraints: Reuse the existing grouped problem-review route, keep the change read-only, render direct handoffs into the existing run and entity inspectors, and link the slice back to deferred TODO `016`.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-sivok-ex5-neovim-problem-review-phase.md`, `ex5-operational-knowledge-system/TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md`, `ex5-operational-knowledge-system/nvim/lua/oks/init.lua`, `ex5-operational-knowledge-system/nvim/problem_review_test.go`, `ex5-operational-knowledge-system/nvim/assets_test.go`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`, `ex5-operational-knowledge-system/docs/terminal-capability-matrix.md`, `ex5-operational-knowledge-system/docs/user-guide.md`

## Goal

Give Neovim users a grouped hotspot review view over the existing problem
review projection.

## Links

- Parent follow-on: `TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md` (`016`)

## Tasks

- [x] sivok.1 Define the grouped hotspot review shape for Neovim.
- [x] sivok.2 Add a read-only `:OksProblemReview` buffer over `/api/problem-review`.
- [x] sivok.3 Add headless behavior coverage and marker coverage.
- [x] sivok.4 Update terminal-first docs and matrix entries.
