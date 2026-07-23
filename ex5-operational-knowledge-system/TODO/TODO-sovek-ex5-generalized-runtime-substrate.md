# TODO sovek - ex5 generalized runtime substrate

## Decision Intent Log

ID: DI-sovek
Date: 2026-07-22 18:12:55 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track a future-scope pass to separate `ex5` as an example application from any broader reusable PromiseGrid runtime substrate.
Intent: Leave room for a more generalized runtime/product layer later instead of overloading `ex5` indefinitely as both example app and full substrate.
Constraints: This is the broadest and least safe scope expansion; keep it explicitly downstream of the smaller embodiment and contract refinements.
Affects: repo architecture, `ex5-operational-knowledge-system/*`, future shared runtime modules, `ex5-operational-knowledge-system/TODO/TODO.md`

ID: DI-ragiv
Date: 2026-07-22 19:49:37 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Open a thin reusable PromiseGrid substrate wave under `promisegrid/records/` for the frozen-family durable record core only, leaving app persistence, projections, and operational workflows in `ex5`.
Intent: Extract the reusable durable record truth first so identity, origin ordering, canonical durable ID derivation, and frozen-family signed-envelope replay/verification no longer live only as `ex5` service internals.
Constraints: Do not generalize search, review, approval workflow composition, or app projections in this wave; preserve current runtime behavior and adapter surfaces while the record core is extracted.
Affects: `docs/thought-experiments/TE-nivor-ex5-generalized-runtime-substrate-boundary.md`, `ex5-operational-knowledge-system/promisegrid/records/*`, `ex5-operational-knowledge-system/service/*.go`, `ex5-operational-knowledge-system/service/*_test.go`, `ex5-operational-knowledge-system/promisegrid/records/*_test.go`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/TODO/TODO.md`
Supersedes: DI-sovek

## Goal

Decide whether the repo should extract a more generalized PromiseGrid runtime
layer beyond `ex5`, and what boundary should remain example-specific.

## Tasks

- [x] sovek.1 Define the candidate split between `ex5` app logic and reusable runtime substrate. See `../../docs/thought-experiments/TE-nivor-ex5-generalized-runtime-substrate-boundary.md`.
- [x] sovek.2 Evaluate migration cost, abstraction risk, and repository impact.
- [x] sovek.3 Decide whether to open a real extraction wave.

## Status

- completed
- thin reusable record substrate extracted under `promisegrid/records/`
- TE `TE-nivor` completed on 2026-07-22 19:38:06 -0700
- locked through `DI-ragiv`
