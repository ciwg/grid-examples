# TODO fanok - ex5 browser success-path coverage

## Decision Intent Log

ID: DI-fanok
Date: 2026-07-21 23:42:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track the missing browser success-path integration coverage as its own ex5 quality TODO.
Intent: Close the gap between strong shared server tests and thinner browser success-path proof so browser-only regressions can be caught earlier.
Constraints: Keep coverage deterministic and headless; prefer real browser wiring over purely marker-based assertions when exercising shipped create/run/evidence/approval/snapshot/supersede flows.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-fanok-ex5-browser-success-path-coverage.md`, `ex5-operational-knowledge-system/web/**`, `ex5-operational-knowledge-system/docs/**`

## Goal

Add real browser success-path coverage for the highest-value shipped actions,
not just marker rendering and failure handling.

## Tasks

- [x] fanok.1 Identify the minimum set of browser success flows that most need end-to-end proof.
- [x] fanok.2 Add headless browser coverage for those flows using the shipped UI and stubbed shared routes.
- [x] fanok.3 Document the resulting browser coverage honestly anywhere the current docs overstate or understate it.
