# TODO turem - ex5 sample data and newcomer guide

## Decision Intent Log

ID: DI-turem
Date: 2026-07-22 22:16:57 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track a newcomer-facing ex5 improvement wave for checked-in sample data under `sample-data/` and one canonical operator guide that is enough for a new user to understand and use the system.
Intent: Make ex5 easier to approach as a serious reference implementation without blurring operator guidance into implementation/reference docs or turning the sample corpus into generator-shaped ambiguity.
Constraints: Keep the sample corpus checked in and deterministic; keep one primary user guide; preserve separate deeper browser, terminal, API, and architecture/reference docs.
Affects: `docs/thought-experiments/*`, `ex5-operational-knowledge-system/sample-data/*`, `ex5-operational-knowledge-system/docs/user-guide.md`, `ex5-operational-knowledge-system/docs/*`, `ex5-operational-knowledge-system/TODO/TODO.md`

ID: DI-rubav
Date: 2026-07-22 22:31:52 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Implement TODO 148 as one checked-in shared newcomer runtime under `sample-data/newcomer-runtime/`, one fail-closed loader script at `scripts/load-sample-data.sh`, and one operator-complete `docs/user-guide.md` that walks newcomers through the sample world before handing them off to deeper browser and terminal guides.
Intent: Make the first newcomer path concrete and reproducible by anchoring the guide to one rich shared corpus with four storylines, one persisted live draft, and one real attachment instead of asking operators to imagine their own data model from empty forms.
Constraints: Keep the sample corpus deterministic; keep all four storylines in one runtime; include at least one receiving problem thread, one inventory discrepancy, one training history thread, one maintenance draft, one real evidence attachment, and one active draft; fail closed if the loader target already contains data.
Affects: `ex5-operational-knowledge-system/sample-data/newcomer-runtime/**`, `ex5-operational-knowledge-system/scripts/load-sample-data.sh`, `ex5-operational-knowledge-system/docs/user-guide.md`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/product-overview.md`, `ex5-operational-knowledge-system/docs/browser-ui-guide.md`, `ex5-operational-knowledge-system/docs/terminal-capability-matrix.md`, `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Make ex5 understandable and usable for newcomers through one rich checked-in
sample corpus and one comprehensive operator guide.

## Tasks

- [x] turem.1 Define the right sample-corpus and newcomer-guide boundary. See `../../docs/thought-experiments/TE-nobav-ex5-sample-data-and-newcomer-guide-boundary.md`.
- [x] turem.2 Lock the exact sample corpus scope and the exact user-guide coverage boundary.
- [x] turem.3 Implement the checked-in corpus and expand the user guide to match.

## Status

- closed
- created from the newcomer-need for realistic sample data and one canonical operator guide
- TE complete: `TE-nobav` recommends one rich checked-in sample corpus plus one canonical newcomer-ready user guide
- implemented as one shared newcomer runtime with four storylines, one active maintenance draft, one real receiving attachment, one loader script, and one operator-complete guide
