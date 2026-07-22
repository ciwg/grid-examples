# TODO zovam - ex5 terminal follow-on slice definition

## Decision Intent Log

ID: DI-zovam
Date: 2026-07-21 16:35:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Break TODO `016` into an immediately actionable next slice instead of leaving the remaining terminal follow-on as a vague umbrella.
Intent: Keep the next `ex5` terminal work small, concrete, and reviewable rather than letting `016` expand into an unbounded bucket again.
Constraints: Link this decision directly to deferred TODO `016`, choose the next slice from concrete existing gaps, and preserve the staged terminal-first strategy.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-zovam-ex5-terminal-followon-slice-definition.md`, `ex5-operational-knowledge-system/TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md`

## Goal

Choose and record the next concrete terminal-first slice under `016`.

## Links

- Parent follow-on: `TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md` (`016`, `nuvok.13`)

## Tasks

- [x] zovam.1 Compare the current remaining terminal gaps.
- [x] zovam.2 Choose the next smallest useful slice under `016`.
- [x] zovam.3 Update `016` so the next slice is explicit.

## Outcome

The next concrete terminal-first slice under `016` is:

- `052` - render `oks-cli pending-review` and `oks-cli problem-review` as
  terminal summaries instead of raw JSON

Why this next:

- it builds directly on the new CLI drilldown work from `049`
- it improves the highest-value shell review queues without inventing new API
  routes
- it stays smaller and safer than broader terminal parity work
