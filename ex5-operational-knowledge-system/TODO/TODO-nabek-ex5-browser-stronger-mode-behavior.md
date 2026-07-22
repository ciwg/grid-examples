# TODO nabek - ex5 browser stronger mode behavior

## Decision Intent Log

ID: DI-nabek
Date: 2026-07-22 22:14:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track stronger browser mode behavior as a follow-on ex5 UI TODO.
Intent: Reduce the remaining “all surfaces are still present at once” feeling by making inactive modes compress or recede more aggressively while keeping every browser area reachable.
Constraints: Preserve the current single-page browser and do not remove any workflow; the goal is stronger behavioral separation, not a rewrite into a multi-page app.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-nabek-ex5-browser-stronger-mode-behavior.md`, `ex5-operational-knowledge-system/web/**`, `ex5-operational-knowledge-system/docs/**`, `ex5-operational-knowledge-system/README.md`

## Goal

Make Review, Author, Operate, Create, and Browse feel more distinct in
behavior, not just in tinting and section labels.

## Tasks

- [x] nabek.1 Review where inactive browser modes still compete too strongly for attention.
- [x] nabek.2 Define the smallest stronger mode behaviors that keep all features reachable while lowering simultaneous cognitive load.
- [x] nabek.3 Implement those stronger mode transitions or disclosures in the browser shell.
