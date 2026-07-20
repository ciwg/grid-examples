# TODO solaj - ex5 open product decisions

## Decision Intent Log

ID: DI-solaj
Date: 2026-07-20 11:56:32 -0700
Author: jj@thesalleys.com (JJ)
Status: superseded
Decision: Record the remaining ex5 product and architecture decisions as an explicit pending TODO instead of leaving them only in README/doc open-items sections.
Intent: Make the unresolved ex5 scope choices visible in the same local queue as the implementation work so future feature slices do not drift away from the still-open decisions.
Constraints: This TODO records open questions only; it does not lock answers to them; implementation that depends on these answers should stay behind explicit follow-up decisions.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-solaj-ex5-open-product-decisions.md`

ID: DI-tabiv
Date: 2026-07-20 16:23:38 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Keep `ex5` on its current local HTTP live-draft model instead of fully porting the `ex3` websocket collaboration transport; treat collaborative editing as optional rather than core; keep a future Neovim embodiment as a desirable follow-on because it matches real team and customer usage.
Intent: Narrow `ex5` to the features needed for the current operational knowledge demonstration while still preserving a credible future path for teams that operate heavily in Neovim.
Constraints: Do not port the full `ex3` websocket stack into `ex5` for the current phase; do not require collaborative editing for the product to be considered valid; do not imply that a Neovim embodiment exists yet.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-solaj-ex5-open-product-decisions.md`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`, `ex5-operational-knowledge-system/TODO/TODO-ranor-ex5-inventory-and-embodiment-followups.md`
Supersedes: DI-solaj

## Goal

Record and close the ex5 product and architecture decisions that define the
current collaboration/editor direction.

## Locked decisions

- do not fully port the `ex3` websocket collaboration model into `ex5`
- collaborative editing is optional, not core
- a future Neovim embodiment would be very valuable and is worth keeping as a follow-on

## Tasks

- [x] solaj.1 Decide whether `ex5` keeps the current local HTTP live-draft model or ports the full `ex3` websocket collaboration transport.
- [x] solaj.2 Decide whether collaborative editing is required for the product or optional beside a more traditional revision/workflow model.
- [x] solaj.3 Decide whether a future non-browser editor embodiment belongs in `ex5`, and if so whether Neovim is the first target.
