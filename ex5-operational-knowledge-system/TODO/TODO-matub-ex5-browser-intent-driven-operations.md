# TODO matub - ex5 browser intent-driven operations

## Decision Intent Log

ID: DI-matub
Date: 2026-07-22 22:12:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track the remaining browser operation-form schema burden as its own ex5 UI TODO.
Intent: Reduce how much the browser asks operators to think in terms of item IDs, revisions, places, resources, responsibilities, and target types when they are trying to log work or capture review actions.
Constraints: Preserve all current browser operation capabilities and manual override paths; improve the main action path by making intent and context more prominent than the raw record model.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-matub-ex5-browser-intent-driven-operations.md`, `ex5-operational-knowledge-system/web/**`, `ex5-operational-knowledge-system/docs/**`, `ex5-operational-knowledge-system/README.md`

## Goal

Make browser operation flows read like the work the operator is trying to do
rather than the schema the system stores.

## Tasks

- [x] matub.1 Review the remaining operation-form fields that still expose too much record-model vocabulary.
- [x] matub.2 Define the smallest set of staged prompts, presets, or context-derived defaults that reduce operator thinking cost.
- [x] matub.3 Implement the intent-driven browser operation improvements without removing the current manual escape hatches.
