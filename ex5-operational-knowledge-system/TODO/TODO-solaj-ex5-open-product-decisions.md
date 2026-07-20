# TODO solaj - ex5 open product decisions

## Decision Intent Log

ID: DI-solaj
Date: 2026-07-20 11:56:32 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Record the remaining ex5 product and architecture decisions as an explicit pending TODO instead of leaving them only in README/doc open-items sections.
Intent: Make the unresolved ex5 scope choices visible in the same local queue as the implementation work so future feature slices do not drift away from the still-open decisions.
Constraints: This TODO records open questions only; it does not lock answers to them; implementation that depends on these answers should stay behind explicit follow-up decisions.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-solaj-ex5-open-product-decisions.md`

## Goal

Track the ex5 product and architecture decisions that are still unresolved and
that should be answered before the bigger collaboration/editor direction is
treated as closed.

## Open decisions

- whether to fully port the `ex3` websocket collaboration model
- whether collaborative editing is truly core or optional
- whether `ex5` should eventually include another editor embodiment like Neovim

## Tasks

- [ ] solaj.1 Decide whether `ex5` keeps the current local HTTP live-draft model or ports the full `ex3` websocket collaboration transport.
- [ ] solaj.2 Decide whether collaborative editing is required for the product or optional beside a more traditional revision/workflow model.
- [ ] solaj.3 Decide whether a future non-browser editor embodiment belongs in `ex5`, and if so whether Neovim is the first target.
