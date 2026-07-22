# TODO suvor - ex5 browser review author mode leakage

## Decision Intent Log

ID: DI-suvor
Date: 2026-07-21 23:41:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track the remaining Review-vs-Author leakage and duplicate live-draft loading as its own ex5 browser TODO.
Intent: Preserve the current item inspection and drafting capabilities while stopping Review-mode inspection from quietly dragging in extra authoring work and redundant live-draft requests.
Constraints: Keep item inspection fast, keep live drafting available from the same browser, and avoid removing any current authoring path.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-suvor-ex5-browser-review-author-mode-leakage.md`, `ex5-operational-knowledge-system/web/**`, `ex5-operational-knowledge-system/docs/**`, `ex5-operational-knowledge-system/README.md`

## Goal

Separate review-mode item inspection from authoring internals more cleanly so
Review stays calm and the browser stops doing redundant live-draft work.

## Tasks

- [x] suvor.1 Decide when item inspection should preload author state versus when drafting should remain an explicit transition.
- [x] suvor.2 Remove redundant `loadEditorItem()` / live-draft fetch paths during refresh and item drilldown.
- [x] suvor.3 Add browser coverage that proves the intended Review-to-Author transition and guards against duplicated loading.
- [x] suvor.4 Update browser docs if the Review/Author handoff becomes more explicit.
