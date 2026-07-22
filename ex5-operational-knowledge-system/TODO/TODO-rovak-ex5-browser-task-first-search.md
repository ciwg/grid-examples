# TODO rovak - ex5 browser task-first search

## Decision Intent Log

ID: DI-rovak
Date: 2026-07-22 22:13:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track the remaining browser search-discovery friction as its own ex5 UI TODO.
Intent: Make browser search feel more like task-first discovery and less like a query console built around backend facets.
Constraints: Preserve the existing search power and advanced filter coverage; improve the default discovery path without removing the current structured controls.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-rovak-ex5-browser-task-first-search.md`, `ex5-operational-knowledge-system/web/**`, `ex5-operational-knowledge-system/docs/**`, `ex5-operational-knowledge-system/README.md`

## Goal

Turn browser search into a clearer discovery surface for common operator tasks
while keeping advanced filtering available.

## Tasks

- [x] rovak.1 Review the current search labels, filters, and result entry points for task-language gaps.
- [x] rovak.2 Define a smaller task-first search layer such as presets, friendlier labels, or clearer advanced-filter separation.
- [x] rovak.3 Implement the smallest browser search changes that improve discovery without reducing power.
