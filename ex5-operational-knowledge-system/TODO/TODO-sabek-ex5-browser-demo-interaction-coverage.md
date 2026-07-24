# TODO sabek - ex5 browser demo interaction coverage

## Decision Intent Log

ID: DI-sabek
Date: 2026-07-23 09:26:59 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track a follow-on pass for more extensive browser interaction coverage around the actual demo path, especially visible lane switching, hotspot drilldowns, and detail-pane handoffs.
Intent: Catch the class of browser issues that only show up when a human follows the live demo recipe, instead of relying mainly on lower-level success-path checks.
Constraints: Keep the tests deterministic, prefer exact on-screen strings and visible-state assertions, and make the covered demo path explicit instead of hiding it behind generic smoke coverage.
Affects: `ex5-operational-knowledge-system/web/*`, `ex5-operational-knowledge-system/docs/user-guide.md`, `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Strengthen browser tests so the live demo path is covered as an obvious user
interaction flow, not just as underlying data correctness.

## Tasks

- [ ] sabek.1 Add coverage for visible review-lane switching between `Draft queue`, `Problem hotspots`, and `Known record search`.
- [ ] sabek.2 Add coverage for hotspot drilldowns that proves the user can see the landing state without guessing or scrolling blindly.
- [ ] sabek.3 Add coverage for detail-pane updates that proves `Current Record` changes visibly when the user clicks `Inspect` on a search result or hotspot-related record.
- [ ] sabek.4 Re-sweep the browser demo sheet against the tested on-screen strings after the interaction coverage lands.

## Status

- open
- created from live demo prep failures where the browser path worked internally but was not obvious enough to trust as a presenter flow
