# TODO timav - ex5 staging abstraction cleanup

## Decision Intent Log

ID: DI-timav
Date: 2026-07-22 21:12:08 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track a future pass to replace remaining practical staging abstractions with more minimal long-term PromiseGrid substrate boundaries where the current staged shapes are no longer justified.
Intent: Close the remaining alignment gap where `ex5` still carries some expedient rollout abstractions that work correctly but are not yet the smallest long-term PromiseGrid shape.
Constraints: Keep behavior stable and evidence-based; do not remove a staging abstraction until the replacement boundary is clearer and no important operator honesty is lost.
Affects: `ex5-operational-knowledge-system/service/*`, `ex5-operational-knowledge-system/web/*`, `ex5-operational-knowledge-system/promisegrid/*`, `docs/thought-experiments/*`, `ex5-operational-knowledge-system/TODO/TODO.md`

ID: DI-nolav
Date: 2026-07-22 21:15:46 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Retire the remaining service-owned mirror structs for durable record rows, events, and peer/relay wire shapes by making `service/*` consume `promisegrid/records` and `promisegrid/transport` types directly where the shapes are already identical.
Intent: Remove rollout-era duplication now that the reusable substrate boundary is real, while keeping ex5 behavior stable and reducing future drift risk.
Constraints: Keep app behavior unchanged; do not mix this cleanup with a broader workflow/projection refactor.
Affects: `docs/thought-experiments/TE-fulok-ex5-staging-abstraction-cleanup.md`, `ex5-operational-knowledge-system/service/*`, `ex5-operational-knowledge-system/promisegrid/records/*`, `ex5-operational-knowledge-system/promisegrid/transport/*`, `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Review and retire the remaining rollout-era abstraction layers that are still
useful today but not obviously the final minimal PromiseGrid boundary.

## Tasks

- [x] timav.1 Identify which current abstractions are still justified rollout scaffolding versus true long-term substrate boundaries. See `../../docs/thought-experiments/TE-fulok-ex5-staging-abstraction-cleanup.md`.
- [x] timav.2 Replace the unjustified staging abstractions with smaller and more direct PromiseGrid-aligned contracts.
- [x] timav.3 Re-audit docs, tests, and implementation claims after each abstraction cleanup so the resulting boundary stays explicit.

## Status

- completed
- `service/*` now consumes substrate-owned record and transport shapes directly instead of keeping those mirror structs as rollout scaffolding
