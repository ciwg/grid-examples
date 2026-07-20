# TODO radok - operational knowledge system foundation

## Decision Intent Log

ID: DI-radok
Date: 2026-07-20 10:13:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Keep `ex5` as one local Go runtime with equal browser and CLI embodiments, explicit protocol docs, and self-contained runtime code copied or adapted from earlier example patterns rather than linked at runtime.
Intent: Preserve the repo's example independence while still reusing proven implementation shapes for collaborative documents and durable workflow state.
Constraints: Do not import `ex3` or `ex4` runtime packages; store durable operational truth under an `ex5`-local runtime root; keep collaborative document state and durable workflow history separate.
Affects: `ex5-operational-knowledge-system/go.mod`, `ex5-operational-knowledge-system/service/**`, `ex5-operational-knowledge-system/cmd/**`, `ex5-operational-knowledge-system/web/**`, `ex5-operational-knowledge-system/docs/**`

ID: DI-kovup
Date: 2026-07-20 10:14:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Model procedures, training, and maintenance content as hybrid-equal knowledge items with collaboratively editable document bodies plus structured metadata, while responsibilities remain their own first-class durable records.
Intent: Support operational and collaborative knowledge in one tool without collapsing everything into either plain documents or plain forms.
Constraints: The durable anchor for completed work is a performed procedure run linked to an exact revision plus evidence; responsibilities must stay separately addressable and linkable.
Affects: `ex5-operational-knowledge-system/service/types.go`, `ex5-operational-knowledge-system/service/app.go`, `ex5-operational-knowledge-system/protocols/**`, `ex5-operational-knowledge-system/docs/**`

ID: DI-zuvob
Date: 2026-07-20 10:15:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Use named workflow roles plus local team policy, with approvals, evidence, links, and runs recorded as append-only operational events projected into current query views.
Intent: Keep the public model flexible for multiple organizations while preserving a concrete, reviewable workflow story for the flagship example.
Constraints: Do not hardcode one company's authority layout; make HTTP a local adapter surface rather than the PromiseGrid-facing contract; treat live document editing as important but secondary to durable operational history.
Affects: `ex5-operational-knowledge-system/service/types.go`, `ex5-operational-knowledge-system/service/app.go`, `ex5-operational-knowledge-system/service/server.go`, `ex5-operational-knowledge-system/cmd/oks-cli/main.go`, `ex5-operational-knowledge-system/docs/**`

## Goal

Add the first complete runnable `ex5` foundation: protocol docs, durable
storage, shared service model, CLI/API/browser surfaces, and implementation
docs.

## Tasks

- [x] radok.1 Create the module, runtime layout, and protocol/docs corpus.
- [x] radok.2 Implement append-only operational event storage and projections.
- [x] radok.3 Implement CLI and HTTP/browser embodiments over the shared model.
- [x] radok.4 Verify the example and close the foundation task.

## Evidence

- The module includes `service/`, `cmd/`, `web/`, `docs/`, and `protocols/`.
- The service replays append-only operational events into responsibility, item, and run projections.
- The browser and CLI both operate over the same HTTP-backed service model.
- Verification passes with `go test ./...` and `errcheck ./...`.
- Live smoke passed through the running server with API probes plus CLI creation/list/search flows for responsibilities, items, and performed runs.
