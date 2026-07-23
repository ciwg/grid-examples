# TODO tolav - ex5 substrate module boundary

## Decision Intent Log

ID: DI-tolav
Date: 2026-07-22 21:24:27 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track a future decision on whether the growing `promisegrid/*` substrate should remain nested inside the `ex5` module or move toward a clearer standalone package or module boundary.
Intent: Keep the current code honest about reusable substrate without assuming that the final packaging boundary must stay permanently nested under the example app.
Constraints: Do not split packaging just for aesthetics; only widen the boundary when the extracted substrate slices are stable enough to justify their own ownership line.
Affects: `ex5-operational-knowledge-system/promisegrid/*`, `ex5-operational-knowledge-system/*`, `docs/thought-experiments/*`, `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Decide whether the long-term reusable PromiseGrid substrate should keep living
inside `ex5-operational-knowledge-system` or graduate to a clearer standalone
boundary.

## Tasks

- [ ] tolav.1 Define the criteria for when nested `promisegrid/*` is no longer the right packaging boundary.
- [ ] tolav.2 Compare staying inside ex5 versus a separate package/module boundary once more substrate slices are proven.
- [ ] tolav.3 Update docs and repo claims if the substrate packaging boundary changes.

## Status

- open
- created from the remaining “serious reference implementation, but not yet the final generalized substrate/product boundary” PromiseGrid alignment gap
