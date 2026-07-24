# TODO vuram - ex5 external browser-check repo

## Decision Intent Log

ID: DI-jutek
Date: 2026-07-23 09:42:58 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track a follow-on wave to define an external browser-check repo for `grid-examples`, starting with the `ex5` demo path and keeping all such browser automation outside the main product repo.
Intent: Add durable browser-demo confidence without polluting `grid-examples` with a second testing surface that is really about presenter-grade browser interaction checks.
Constraints: Keep the browser-check repo external to `grid-examples`, align it explicitly to the `grid-examples` examples, and choose a harness that is strong on visible interaction proofs rather than only lower-level correctness.
Affects: `docs/thought-experiments/TE-rasem-ex5-external-browser-check-repo.md`, `ex5-operational-knowledge-system/TODO/TODO.md`

ID: DI-vuram
Date: 2026-07-23 09:42:58 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Use `~/lab/cswg/grid-examples-browser-checks/` as the external repo root, keep one subdirectory per example, and start `ex5` with a lean Playwright harness that checks the visible browser demo path without always-on recordings.
Intent: Give `grid-examples` one explicit external home for browser confidence work while keeping the first `ex5` slice practical, disk-light, and focused on the exact presenter interactions that were failing manually.
Constraints: The browser-check repo stays outside `grid-examples`; the first `ex5` slice uses Playwright; the initial harness stays lean with no always-on video or traces; the first tests cover visible `Current Record` and hotspot drilldown behavior.
Affects: `/home/jj/lab/cswg/grid-examples-browser-checks/README.md`, `/home/jj/lab/cswg/grid-examples-browser-checks/.gitignore`, `/home/jj/lab/cswg/grid-examples-browser-checks/ex5/*`, `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Decide the external repo shape and first harness for browser-demo validation
that goes beyond what should live in the main `grid-examples` repo.

## Tasks

- [x] vuram.1 Lock the external repo root and example subdirectory shape.
- [x] vuram.2 Lock the first harness choice for `ex5`.
- [x] vuram.3 Lock the first `ex5` browser checks to implement there.
- [x] vuram.4 Record how the external repo stays aligned to `grid-examples` over time.

## Status

- completed
- external browser checks now live under `~/lab/cswg/grid-examples-browser-checks/`, with a lean Playwright `ex5` slice that stays outside `grid-examples`
