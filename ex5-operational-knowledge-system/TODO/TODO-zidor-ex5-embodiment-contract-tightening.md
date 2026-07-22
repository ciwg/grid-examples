# TODO zidor - decide when ex5 embodiments should move beyond the local HTTP adapter contract

## Decision Intent Log

ID: DI-vabek
Date: 2026-07-22 10:23:43 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Keep browser, CLI, and Neovim described through the current local HTTP adapter until the first relay-visible exchange layer and additive CAS layer are actually implemented, then tighten embodiment/runtime language in a later slice.
Intent: Preserve honest shipped docs today while still naming the concrete runtime milestone that justifies later embodiment-contract tightening.
Constraints: Do not restate current embodiments as if they already speak a direct PromiseGrid peer/runtime contract; do not combine this timing decision with the first peer-exchange or CAS implementation slices.
Affects: `ex5-operational-knowledge-system/TODO/TODO-zidor-ex5-embodiment-contract-tightening.md`, `docs/thought-experiments/TE-lavok-ex5-embodiment-contract-tightening-timing.md`, `ex5-operational-knowledge-system/docs/promisegrid-embodiment-staging.md`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `ex5-operational-knowledge-system/docs/architecture.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Define the point at which browser, CLI, and Neovim should stop being described
primarily through the local HTTP adapter and instead bind more directly to the
shipped PromiseGrid runtime contract.

## Why this exists

The current docs correctly describe the local HTTP API as the embodiment
adapter. The remaining question was when that description should tighten, not
whether it should be rewritten early.

## Tasks

- [x] zidor.1 Run the required TE for embodiment-contract tightening timing and
  scope.
- [x] zidor.2 Lock the staged boundary between local HTTP adapter behavior and
  direct runtime contract behavior.
- [x] zidor.3 Define the first embodiment-facing migration slice, if any.
- [x] zidor.4 Update the external and repo docs once that boundary is settled.

## Status

- done
- embodiment tightening waits for implemented peer-exchange and additive CAS
  layers
- current browser, CLI, and Neovim docs remain honestly adapter-first until
  then
