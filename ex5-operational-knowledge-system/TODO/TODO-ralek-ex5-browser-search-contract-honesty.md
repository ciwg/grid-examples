# TODO ralek - ex5 browser search contract honesty

## Decision Intent Log

ID: DI-ralek
Date: 2026-07-21 23:40:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track the remaining browser/search contract mismatch as its own ex5 follow-on TODO.
Intent: Align the shipped browser search wording and behavior so "known record" lookup, free-text matching, and `problem=true` review behave the way the UI and docs currently imply.
Constraints: Preserve the existing shared `/api/search` route; keep browser, CLI, Neovim, and docs aligned on the same search contract instead of inventing browser-only semantics.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-ralek-ex5-browser-search-contract-honesty.md`, `ex5-operational-knowledge-system/service/**`, `ex5-operational-knowledge-system/web/**`, `ex5-operational-knowledge-system/docs/**`, `ex5-operational-knowledge-system/README.md`

## Goal

Make the browser's known-record search honest and reliable by aligning the UI
copy, free-text matching rules, and `problem=true` result behavior.

## Tasks

- [x] ralek.1 Decide whether shared free-text search should index record IDs directly for places, resources, responsibilities, items, and runs.
- [x] ralek.2 Decide whether `problem=true` should suppress non-run groups from browser review results or merely filter runs while leaving context groups visible.
- [x] ralek.3 Implement the chosen shared search behavior and add server/browser coverage for the exact contract.
- [x] ralek.4 Update browser-facing docs so they describe the shipped search behavior exactly.
