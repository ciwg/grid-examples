# TODO mivor - ex5 browser-first demo prep

## Decision Intent Log

ID: DI-ravot
Date: 2026-07-23 07:46:39 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track a browser-first ex5 demo-prep wave for one canonical live-demo path, one CLI proof slice, and the smallest honest recording-helper boundary that supports short repeatable videos.
Intent: Make the next ex5 demonstration reproducible and newcomer-aligned by anchoring it to the same checked-in sample world and operator guidance already shipped, instead of building a second demo-only story.
Constraints: Keep the browser as the primary demo embodiment; include a real CLI proof slice; keep recording helpers thin unless a broader helper layer is explicitly justified later; preserve alignment with `sample-data/newcomer-runtime/` and `docs/user-guide.md`.
Affects: `docs/thought-experiments/*`, `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-mivor-ex5-browser-first-demo-prep.md`, `ex5-operational-knowledge-system/docs/*`, `ex5-operational-knowledge-system/scripts/*`

ID: DI-luren
Date: 2026-07-23 08:05:17 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add the first live-demo presenter sheet as a compact subsection under the existing browser walkthrough in `docs/user-guide.md`, scoped to one browser-first intro flow plus one short CLI proof line.
Intent: Keep the live-demo sheet attached to the same newcomer sample tour it demonstrates, instead of creating a second top-level intro surface that could drift from the canonical guide.
Constraints: Scope is live-demo only; keep the browser as the main story; keep the CLI proof short; do not widen this slice into the broader recording-helper or boss-walkthrough work yet.
Affects: `ex5-operational-knowledge-system/TODO/TODO-mivor-ex5-browser-first-demo-prep.md`, `ex5-operational-knowledge-system/docs/user-guide.md`

## Goal

Prepare one canonical browser-first ex5 demo path, one short CLI proof slice,
and thin repeatability helpers for live use and short recordings.

## Tasks

- [x] mivor.1 Define the browser-first demo-prep boundary. See `../../docs/thought-experiments/TE-lurak-ex5-browser-first-demo-prep-boundary.md`.
- [ ] mivor.2 Lock the exact demo story arc, CLI proof slice, and helper boundary.
- [ ] mivor.3 Implement the demo TODO, script, and any approved recording helpers.

## Status

- open
- created from the need to prepare a browser-first ex5 demo today, with CLI proof and optional short-video support
- TE complete: `TE-lurak` recommends one browser-first demo pack anchored to the checked-in newcomer corpus, with a short CLI proof slice and only thin recording helpers
