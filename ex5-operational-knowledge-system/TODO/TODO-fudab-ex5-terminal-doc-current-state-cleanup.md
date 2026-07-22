# TODO fudab - ex5 terminal doc current-state cleanup

## Decision Intent Log

ID: DI-fudab
Date: 2026-07-21 16:35:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Rewrite the remaining terminal-behavior docs so shipped CLI/Neovim capabilities read primarily as current state, while preserving any useful historical context.
Intent: Make `ex5` terminal docs easier to use as an operator guide now that many formerly “next phase” terminal slices have already landed.
Constraints: Keep the docs honest, do not overstate parity, preserve relevant historical context where it still helps, and link this cleanup back to deferred TODO `016`.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-fudab-ex5-terminal-doc-current-state-cleanup.md`, `ex5-operational-knowledge-system/TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`, `ex5-operational-knowledge-system/docs/architecture.md`

## Goal

Make the terminal sections read like a clean current-state guide instead of a
mixed rollout log.

## Links

- Parent follow-on: `TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md` (`016`, `nuvok.13`)

## Tasks

- [x] fudab.1 Identify the terminal sections that still read as rollout history.
- [x] fudab.2 Rewrite them as current-state guidance without overstating parity.
- [x] fudab.3 Keep any remaining historical notes clearly secondary.
