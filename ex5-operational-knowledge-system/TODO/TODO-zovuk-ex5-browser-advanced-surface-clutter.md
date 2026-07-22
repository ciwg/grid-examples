# TODO zovuk - ex5 browser advanced surface clutter

## Decision Intent Log

ID: DI-zovuk
Date: 2026-07-21 23:43:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track the remaining browser clutter in advanced search and operate override surfaces as its own ex5 UI TODO.
Intent: Keep all current override power while reducing how much schema-heavy detail the browser shows at once in the advanced search and operate areas.
Constraints: Do not remove any override fields; keep all current search and operate capabilities reachable; reduce visible complexity by staging, grouping, or relabeling rather than deleting controls.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-zovuk-ex5-browser-advanced-surface-clutter.md`, `ex5-operational-knowledge-system/web/**`, `ex5-operational-knowledge-system/docs/**`, `ex5-operational-knowledge-system/README.md`

## Goal

Reduce the remaining clutter concentrated in advanced search filters and
operate override controls without weakening current browser functionality.

## Tasks

- [x] zovuk.1 Review the current advanced search and operate override surfaces for the most distracting schema-heavy controls.
- [x] zovuk.2 Define the smallest staging/grouping pass that keeps all current escape hatches reachable.
- [x] zovuk.3 Implement the clutter-reduction pass and add browser coverage for the resulting interaction model.
- [x] zovuk.4 Update browser docs so the advanced controls are still understandable after the reduction pass.
