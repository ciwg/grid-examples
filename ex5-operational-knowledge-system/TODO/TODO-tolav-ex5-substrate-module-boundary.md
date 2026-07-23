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

ID: DI-zufek
Date: 2026-07-22 21:51:02 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Lock `146A` and keep `promisegrid/*` nested inside the `ex5-operational-knowledge-system` module boundary for now, while treating the current `promisegrid/records`, `promisegrid/transport`, and `promisegrid/store` directories as the real semantic substrate boundary.
Intent: Keep PromiseGrid packaging evidence-first. The semantic substrate line is now honest and useful, but there is still one consumer set, one release line, and an intentionally partial substrate after `145A`, so a separate module boundary would be more symbolic than justified.
Constraints: Close `146` as an intentional stay-nested decision; do not imply that a standalone module boundary is already warranted before second-consumer, versioning, or release evidence appears.
Affects: `../../docs/thought-experiments/TE-muvok-ex5-substrate-module-boundary.md`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `ex5-operational-knowledge-system/TODO/TODO.md`
Supersedes: DI-tolav

## Goal

Decide whether the long-term reusable PromiseGrid substrate should keep living
inside `ex5-operational-knowledge-system` or graduate to a clearer standalone
boundary.

## Tasks

- [x] tolav.1 Define the criteria for when nested `promisegrid/*` is no longer the right packaging boundary. See `../../docs/thought-experiments/TE-muvok-ex5-substrate-module-boundary.md`.
- [x] tolav.2 Compare staying inside ex5 versus a separate package/module boundary once more substrate slices are proven. Locked to `146A`: stay nested for now.
- [x] tolav.3 Update docs and repo claims if the substrate packaging boundary changes. The packaging boundary is intentionally unchanged; docs now state that explicitly.

## Status

- completed
- created from the remaining “serious reference implementation, but not yet the final generalized substrate/product boundary” PromiseGrid alignment gap
- TE complete: `TE-muvok` recommends staying inside the `ex5` module boundary for now because the semantic substrate line is real, but independent packaging is not yet justified by consumer or release evidence.
- locked to `146A`: keep `promisegrid/*` nested inside the `ex5` module until second-consumer, versioning, or release evidence justifies a wider package boundary
