# TODO ravem - ex5 further PromiseGrid substrate extraction

## Decision Intent Log

ID: DI-ravem
Date: 2026-07-22 21:12:08 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track a future wave to extract more reusable PromiseGrid substrate beyond the first `promisegrid/records` slice now living inside `ex5`.
Intent: Close the remaining gap between a strong PromiseGrid-aligned example system and a more reusable substrate by moving proven non-example mechanics out of `ex5` app ownership when the boundary is clear.
Constraints: Keep extractions evidence-based; do not over-generalize projections or workflow-specific app logic just to increase substrate surface area.
Affects: `ex5-operational-knowledge-system/promisegrid/*`, `ex5-operational-knowledge-system/service/*`, `docs/thought-experiments/*`, `ex5-operational-knowledge-system/TODO/TODO.md`

ID: DI-vurem
Date: 2026-07-22 21:15:46 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Extract the next reusable PromiseGrid substrate as `promisegrid/transport/`, carrying peer-exchange and relay-feed wire structs plus origin-aware batch helpers on top of the already-reusable record types.
Intent: Move the next proven non-example mechanics out of `service/` without over-generalizing ex5-specific projections, workflows, or embodiment composition.
Constraints: Keep `service/` ownership of ex5 persistence and operator-facing workflows; do not turn this into a larger embodiment or framework rewrite.
Affects: `docs/thought-experiments/TE-sorav-ex5-next-promisegrid-transport-substrate.md`, `ex5-operational-knowledge-system/promisegrid/transport/*`, `ex5-operational-knowledge-system/service/*`, `ex5-operational-knowledge-system/docs/*`, `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Identify and extract the next proven reusable PromiseGrid substrate slices beyond
the current durable record core.

## Tasks

- [x] ravem.1 Decide which already-shipped ex5 mechanics now have enough reuse evidence to leave `service/` and move into `promisegrid/`. See `../../docs/thought-experiments/TE-sorav-ex5-next-promisegrid-transport-substrate.md`.
- [x] ravem.2 Preserve a clean boundary between reusable substrate and ex5-specific projections, review flows, and embodiment UX.
- [x] ravem.3 Align docs and implementation claims when the next substrate slices are extracted.

## Status

- completed
- peer-exchange and relay-feed wire types plus origin-aware transport helpers now live under `promisegrid/transport/`
