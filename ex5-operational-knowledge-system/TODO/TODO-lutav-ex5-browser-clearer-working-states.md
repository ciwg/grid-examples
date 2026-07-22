# TODO lutav - ex5 browser clearer working states

## Decision Intent Log

ID: DI-lutav
Date: 2026-07-22 22:46:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track stronger browser working-state separation as its own ex5 UI TODO.
Intent: Make Review, Author, Operate, Create, and Browse feel more like distinct working states instead of one page with compressed sections.
Constraints: Preserve the single-page browser and keep every current surface reachable; focus on behavioral clarity, not a route rewrite.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-lutav-ex5-browser-clearer-working-states.md`, `ex5-operational-knowledge-system/web/**`, `ex5-operational-knowledge-system/docs/**`, `ex5-operational-knowledge-system/README.md`

## Goal

Make each browser mode feel more like its own working state while preserving
the current single-page structure.

## Tasks

- [x] lutav.1 Review where the current browser modes still behave too similarly.
- [x] lutav.2 Define the smallest stronger state transitions or collapses that make each mode more distinct.
- [x] lutav.3 Implement the clearer working-state behavior without hiding functionality permanently.
