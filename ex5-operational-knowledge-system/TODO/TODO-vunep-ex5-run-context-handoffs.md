# TODO vunep - ex5 run context handoffs

## Decision Intent Log

ID: DI-vunep
Date: 2026-07-21 19:05:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add better terminal handoffs from run review into related item, place, resource, and responsibility context instead of stopping at the run record itself.
Intent: Make queue-driven run review flows lead naturally into surrounding operational context for both CLI and Neovim users.
Constraints: Reuse the existing detail routes and inspector surfaces, keep the change terminal-first, and link it back to deferred TODO `016`.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-vunep-ex5-run-context-handoffs.md`, `ex5-operational-knowledge-system/TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md`, `ex5-operational-knowledge-system/cmd/oks-cli/main.go`, `ex5-operational-knowledge-system/cmd/oks-cli/main_test.go`, `ex5-operational-knowledge-system/nvim/lua/oks/init.lua`, `ex5-operational-knowledge-system/nvim/*.go`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`

## Goal

Let run review hand off cleanly into the related context records from terminal
surfaces.

## Links

- Parent follow-on: `TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md` (`016`)

## Tasks

- [x] vunep.1 Define the desired run-context handoff hints for CLI and Neovim.
- [x] vunep.2 Add the CLI run drilldown handoffs.
- [x] vunep.3 Add the Neovim run inspector handoffs.
- [x] vunep.4 Add regression coverage and docs.
