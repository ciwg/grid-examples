# TODO mepuk - operational knowledge system

## Decision Intent Log

ID: DI-kesuv
Date: 2026-07-20 10:12:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add `ex5-operational-knowledge-system/` as a new independent sibling example that demonstrates collaborative knowledge documents plus durable operational history in one PromiseGrid-shaped application.
Intent: Make `ex5` the flagship operational reference app for this repo while keeping its runtime code self-contained instead of linking back to earlier examples.
Constraints: Keep `ex5` independent of `ex3` and `ex4` at runtime; ground the design in the actual grid technology already proven in the repo; provide equal browser and CLI embodiments.
Affects: `TODO/TODO.md`, `TODO/TODO-mepuk-operational-knowledge-system.md`, `ex5-operational-knowledge-system/**`

## Goal

Build `ex5-operational-knowledge-system/` as a standalone example app that
combines collaboratively editable operational documents with durable
event-sourced workflow history for responsibilities, approvals, evidence, and
performed procedure runs.

## Tasks

- [x] mepuk.1 Create the `ex5-operational-knowledge-system/` module, local TODO/docs corpus, and protocol docs.
- [x] mepuk.2 Implement durable event storage and projections for responsibilities, knowledge items, runs, evidence, approvals, and links.
- [x] mepuk.3 Implement equal CLI and HTTP/browser surfaces over the shared service model.
- [x] mepuk.4 Document the architecture, storage model, and user-facing behavior around the implemented example.
- [x] mepuk.5 Verify the example with Go tests and targeted runtime checks.

## Evidence

- `ex5-operational-knowledge-system/` now exists as a standalone nested module.
- The example stores durable runtime data under `.operational-knowledge-system/`.
- Browser and CLI both target the same local Go runtime and HTTP API.
- Verification passes with `go test ./...` and `errcheck ./...` from `ex5-operational-knowledge-system/`.
- Live smoke passed with `go run ./cmd/operational-knowledge`, `curl` on `/api/meta` and `/api/dashboard`, and `go run ./cmd/oks-cli` creating responsibilities, items, runs, and search results against the running server.
