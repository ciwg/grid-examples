# TODO pogul - ex5 grouped problem review

## Decision Intent Log

ID: DI-pogul
Date: 2026-07-20 21:20:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add a grouped problem-review view for receiving and inventory work so operators can see repeated problem runs summarized by place and resource instead of reconstructing patterns manually from individual records.
Intent: Make it easy to answer where recurring receiving and count problems are happening, while staying inside the current operational-memory model and avoiding ERP/MRP quantity-ledger logic.
Constraints: Reuse the existing run/evidence model; keep the summary bounded to receiving-check and inventory-audit problem runs; expose it through the local HTTP/browser surface; update docs and tests in the same slice.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-pogul-ex5-grouped-problem-review.md`, `ex5-operational-knowledge-system/service/types.go`, `ex5-operational-knowledge-system/service/app.go`, `ex5-operational-knowledge-system/service/server.go`, `ex5-operational-knowledge-system/service/app_test.go`, `ex5-operational-knowledge-system/service/server_test.go`, `ex5-operational-knowledge-system/web/index.html`, `ex5-operational-knowledge-system/web/app.js`, `ex5-operational-knowledge-system/web/assets_test.go`, `ex5-operational-knowledge-system/web/browser_smoke_test.go`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`

## Goal

Add a browser-visible problem-review summary that groups repeated receiving and
inventory problems by place and resource, using the existing run/evidence model
so operators can find hotspots without turning `ex5` into an ERP/MRP system.

## Intended repo paths

- `ex5-operational-knowledge-system/TODO/TODO.md`
- `ex5-operational-knowledge-system/TODO/TODO-pogul-ex5-grouped-problem-review.md`
- `ex5-operational-knowledge-system/service/types.go`
- `ex5-operational-knowledge-system/service/app.go`
- `ex5-operational-knowledge-system/service/server.go`
- `ex5-operational-knowledge-system/service/app_test.go`
- `ex5-operational-knowledge-system/service/server_test.go`
- `ex5-operational-knowledge-system/web/index.html`
- `ex5-operational-knowledge-system/web/app.js`
- `ex5-operational-knowledge-system/web/assets_test.go`
- `ex5-operational-knowledge-system/web/browser_smoke_test.go`
- `ex5-operational-knowledge-system/README.md`
- `ex5-operational-knowledge-system/docs/features-guide.md`
- `ex5-operational-knowledge-system/docs/http-api-guide.md`
- `ex5-operational-knowledge-system/docs/practical-implementation.md`

## Intended runtime path patterns

- `t.TempDir()/**`
  - class: `test`
  - actions: `read/write`
  - purpose: service/server/browser smoke temporary roots
  - lifecycle: test-only; auto-cleaned by the Go test harness

- `127.0.0.1:<ephemeral-test-port>`
  - class: `test`
  - actions: `listen`
  - purpose: headless browser smoke server during `go test ./...`
  - lifecycle: test-only; closed by the test process

## Tasks

- [x] pogul.1 Add an app/server summary of problematic receiving and inventory runs grouped by place and resource.
- [x] pogul.2 Add a browser review panel that renders grouped problem hotspots and links back into the existing inspector flow.
- [x] pogul.3 Add tests and docs for the grouped problem-review path.
