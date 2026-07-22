# TODO luzom - ex5 CLI place and resource drilldown polish

## Decision Intent Log

ID: DI-luzom
Date: 2026-07-21 16:05:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Keep `oks-cli show-place` and `oks-cli show-resource` on the current detail routes, but render those projections in a more operator-useful terminal layout.
Intent: Make shell-first place and resource drilldowns useful for real review work by surfacing hierarchy, links, and related runs directly, without inventing a terminal-only backend surface.
Constraints: Reuse only `GET /api/places/{id}` and `GET /api/resources/{id}`, stay read-only, preserve honest terminal behavior, and add tests plus docs in the same slice.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-luzom-ex5-cli-place-resource-drilldown-polish.md`, `ex5-operational-knowledge-system/cmd/oks-cli/main.go`, `ex5-operational-knowledge-system/cmd/oks-cli/main_test.go`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`

## Goal

Make CLI place/resource inspection expose contextual drilldown information
clearly enough for terminal-first operator review.

## Tasks

- [x] luzom.1 Keep the CLI on the existing place/resource detail routes.
- [x] luzom.2 Render hierarchy, related runs, and links in a more useful terminal layout.
- [x] luzom.3 Add CLI regression coverage for place/resource detail rendering.
- [x] luzom.4 Update ex5 docs to describe the improved CLI drilldown behavior.
