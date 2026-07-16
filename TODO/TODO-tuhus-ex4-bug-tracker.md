# TODO tuhus - ex4 bug tracker

## Decision Intent Log

ID: DI-tuhus
Date: 2026-07-16 16:00:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add `ex4-bug-tracker/` as a new sibling example that stands on its own as a usable bug tracker and as a concrete example program in this repo.
Intent: Expand the example set beyond the grid editor while keeping the new app understandable, durable-first, and consistent with the repo's example-driven structure.
Constraints: Leave `ex1-order-flow/`, `ex2-grid-editor/`, and `ex3-grid-editor-websocket/` unchanged except for shared root TODO bookkeeping; keep the first slice small enough to implement and verify in one pass.
Affects: `TODO/TODO.md`, `TODO/TODO-tuhus-ex4-bug-tracker.md`, `ex4-bug-tracker/**`

## Goal

Add a browser-first bug tracker with a simple CLI, append-only issue history,
real uploaded attachments, a single-team v1 workflow, and a built-in seam for
future multi-team support.

## Tasks

- [x] tuhus.1 Create the `ex4-bug-tracker/` module, runtime layout, and local TODO/docs corpus.
- [x] tuhus.2 Implement durable issue storage, browser HTTP routes, and attachment handling.
- [x] tuhus.3 Implement the engineer-focused CLI commands.
- [x] tuhus.4 Add browser UI for queue, detail, create, comment, assignment, status, and attachments.
- [x] tuhus.5 Verify the example with Go tests and targeted runtime checks.

## Evidence

- `ex4-bug-tracker/` exists as a sibling nested module under the repo root.
- The example stores runtime state under `ex4-bug-tracker/.bug-tracker/`.
- The browser UI and CLI both operate on the same issue history model.
- Verification passes with `go test ./...` from `ex4-bug-tracker/`.
