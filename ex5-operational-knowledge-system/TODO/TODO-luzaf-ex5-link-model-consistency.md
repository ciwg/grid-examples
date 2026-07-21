# TODO luzaf - ex5 link model consistency

## Decision Intent Log

ID: DI-luzaf
Date: 2026-07-21 12:05:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Split the review findings so typed-link validation and responsibility-link consistency are tracked as their own ex5 fix TODO.
Intent: Make the graph trustworthy enough that browser, CLI, Neovim, and docs all describe the same link model.
Constraints: Focus on write-time endpoint validation and responsibility-link consistency; keep docs and tests in the same implementation pass.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-vurab-ex5-review-followups.md`, `ex5-operational-knowledge-system/TODO/TODO-luzaf-ex5-link-model-consistency.md`

## Goal

Fix the typed-link model so link writes are structurally valid and the
responsibility record behaves consistently with the documented graph model.

## Tasks

- [x] luzaf.1 Validate typed-link endpoints and types on write instead of accepting dangling or mistyped graph edges.
- [x] luzaf.2 Make responsibility link behavior consistent with the documented typed-link model across service, browser, Neovim, and docs.

## Status

- done
- derived from the 2026-07-21 extensive ex5 review
