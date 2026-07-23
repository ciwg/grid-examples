# TODO sulok - ex5 example vs runtime product boundary

## Decision Intent Log

ID: DI-sulok
Date: 2026-07-22 21:12:08 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track a future decision on how `ex5` should relate to a broader reusable PromiseGrid runtime or product boundary, instead of remaining only a strong example system.
Intent: Make the eventual split between “example application” and “general runtime/product” explicit so PromiseGrid claims can advance without overloading `ex5` with every future substrate responsibility.
Constraints: Avoid speculative framework work before there is a clear boundary and a justified first reusable runtime surface.
Affects: `ex5-operational-knowledge-system/*`, `docs/thought-experiments/*`, `ex5-operational-knowledge-system/TODO/TODO.md`

ID: DI-rasok
Date: 2026-07-22 21:15:46 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Make the ex5 boundary explicit in the docs: `promisegrid/*` is the reusable substrate boundary, while `service/*` plus the shipped embodiments remain the ex5 application/runtime layer rather than the final generalized PromiseGrid product.
Intent: Keep PromiseGrid claims honest now that the repo has real reusable substrate slices without overstating ex5-specific workflows as a generalized runtime.
Constraints: Document the split directly in the high-visibility technical docs; do not rename the module or claim a broader product boundary than the code proves.
Affects: `docs/thought-experiments/TE-nuzek-ex5-example-vs-substrate-boundary.md`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/architecture.md`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Decide whether and how `ex5` should remain an example app while a broader
PromiseGrid runtime/product boundary becomes explicit elsewhere.

## Tasks

- [x] sulok.1 Decide whether the next generalized runtime step should stay inside `ex5` or move into a separate reusable module boundary. See `../../docs/thought-experiments/TE-nuzek-ex5-example-vs-substrate-boundary.md`.
- [x] sulok.2 Define what still properly belongs to the example system after a broader runtime/product line is introduced.
- [x] sulok.3 Align repo docs so `ex5` claims stay honest about example scope versus generalized runtime scope.

## Status

- completed
- the docs now state `promisegrid/*` as reusable substrate and keep `service/*` plus embodiments as ex5 application/runtime ownership
