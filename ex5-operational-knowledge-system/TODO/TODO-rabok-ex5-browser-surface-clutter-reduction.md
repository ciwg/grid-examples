# TODO rabok - ex5 browser surface clutter reduction

## Decision Intent Log

ID: DI-rabok
Date: 2026-07-22 22:45:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track the remaining browser clutter problem as its own ex5 UI TODO.
Intent: Reduce how many equally prominent surfaces the browser shows at once so the page stops reading like a broad control console.
Constraints: Preserve all current browser functionality and keep every surface reachable; reduce simultaneous exposure rather than deleting capabilities.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-rabok-ex5-browser-surface-clutter-reduction.md`, `ex5-operational-knowledge-system/web/**`, `ex5-operational-knowledge-system/docs/**`, `ex5-operational-knowledge-system/README.md`

## Goal

Make the browser feel less busy by reducing how many large, equal-weight
surfaces compete for attention at the same time.

## Tasks

- [x] rabok.1 Review the browser for panels, helper text, and actions that still compete too strongly at the same time.
- [x] rabok.2 Define the smallest reduction in simultaneous surface exposure that preserves all current capabilities.
- [x] rabok.3 Implement the clutter-reduction pass without removing any browser feature.
