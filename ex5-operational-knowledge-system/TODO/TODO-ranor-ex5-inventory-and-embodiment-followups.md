# TODO ranor - ex5 inventory and embodiment followups

## Decision Intent Log

ID: DI-ranor
Date: 2026-07-20 11:56:32 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track inventory expansion and alternate embodiment work as explicit deferred ex5 follow-ons rather than leaving them implicit in scattered notes.
Intent: Preserve the current operational-memory scope while keeping the obvious future branches visible for later planning.
Constraints: This TODO is future-facing; it does not commit `ex5` to full ERP/MRP quantity logic or to a second editor embodiment in the current phase.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-ranor-ex5-inventory-and-embodiment-followups.md`

ID: DI-ranov
Date: 2026-07-20 21:35:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Close TODO `007` as the completed inventory follow-on bucket for the current operational-memory scope, and split the future Neovim embodiment into its own deferred TODO.
Intent: Keep the inventory track honest about what is implemented now, while preventing embodiment work from making inventory TODO `007` look perpetually unfinished.
Constraints: `ex5` still does not commit to ERP/MRP quantity logic; Neovim remains desirable but deferred as a separate embodiment follow-on.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-ranor-ex5-inventory-and-embodiment-followups.md`, `ex5-operational-knowledge-system/TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md`

## Goal

Track the inventory follow-on work around deeper operational-memory support and
record its completion for the current `ex5` scope.

## Tasks

- [x] ranor.1 Decide whether `ex5` should stay in the operational-memory lane for inventory or later grow quantity, reservation, or planning features.
- [x] ranor.2 Define the next inventory-shaped operator flows beyond audits, such as receiving checks, discrepancy review, and count-history drilldown.
- [x] ranor.3 Split the future Neovim embodiment into its own deferred TODO instead of leaving it inside inventory TODO `007`.

## Status

- `ex5` stays in the operational-memory lane for inventory in the current phase
- implemented follow-ons now include receiving review, discrepancy review,
  context fact history, history drilldown filters, and grouped problem review
- Neovim embodiment work is now tracked separately in
  `TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md`

## Completion note

TODO `007` is complete for the current intended scope. It does not imply full
ERP/MRP quantity, reservation, or planning logic, and it does not imply that a
Neovim embodiment has been implemented.
