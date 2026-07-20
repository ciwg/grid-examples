# TODO honus - ex5 workflow search and automation

## Decision Intent Log

ID: DI-honus
Date: 2026-07-20 11:56:32 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Group the safe ex5 feature backlog into one pending TODO covering richer workflow views, search/filtering, and stronger browser automation.
Intent: Keep the non-controversial ex5 improvement work visible and implementable without blocking on the larger collaboration/editor decisions.
Constraints: Stay within the existing ex5 product shape; do not expand into ERP/MRP planning logic here; each feature slice must keep docs and tests in sync.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-honus-ex5-workflow-search-and-automation.md`

## Goal

Improve the everyday operator experience in `ex5` with better record views,
search/filtering, workflow drilldown, and stronger browser automation.

## Tasks

- [x] honus.1 Add richer single-record views and timelines for places, resources, responsibilities, items, and runs.
- [x] honus.2 Expand search and filtering by kind, status, place, resource, and responsibility, with better grouped drilldown.
- [x] honus.3 Improve run, evidence, and approval review flows so operators can trace revision-to-run history more easily.
- [x] honus.4 Strengthen browser automation beyond embedded asset checks to cover real UI workflow behavior.
