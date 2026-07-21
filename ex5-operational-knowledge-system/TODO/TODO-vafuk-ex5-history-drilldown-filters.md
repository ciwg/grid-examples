# TODO vafuk - ex5 history drilldown filters

## Decision Intent Log

ID: DI-vafuk
Date: 2026-07-20 21:08:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add outcome-aware history filtering and one-click context drilldowns for receiving and inventory review so operators can move from a place/resource/responsibility into filtered runs without manually rebuilding the search.
Intent: Make the new context review panels actionable by turning them into direct history navigation, especially for “show me receiving problems here” and “show me counts for this bin” workflows.
Constraints: Stay within the current local HTTP search model; do not add ERP/MRP quantity logic; keep the work focused on search/filtering and browser drilldown; update tests and docs in the same slice.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-vafuk-ex5-history-drilldown-filters.md`, `ex5-operational-knowledge-system/service/types.go`, `ex5-operational-knowledge-system/service/app.go`, `ex5-operational-knowledge-system/service/server.go`, `ex5-operational-knowledge-system/service/app_test.go`, `ex5-operational-knowledge-system/service/server_test.go`, `ex5-operational-knowledge-system/web/index.html`, `ex5-operational-knowledge-system/web/app.js`, `ex5-operational-knowledge-system/web/assets_test.go`, `ex5-operational-knowledge-system/web/browser_smoke_test.go`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`, `ex5-operational-knowledge-system/docs/http-api-guide.md`

## Goal

Add outcome-based search filters and context drilldown actions so operators can jump directly into receiving and inventory run history by place/resource/responsibility.

## Intended repo paths

- `ex5-operational-knowledge-system/TODO/TODO.md`
- `ex5-operational-knowledge-system/TODO/TODO-vafuk-ex5-history-drilldown-filters.md`
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
- `ex5-operational-knowledge-system/docs/practical-implementation.md`
- `ex5-operational-knowledge-system/docs/http-api-guide.md`

## Intended runtime path patterns

- `t.TempDir()/**`
  - class: `test`
  - actions: `read/write`
  - purpose: server/browser smoke temporary roots
  - lifecycle: test-only; auto-cleaned by the Go test harness

- `127.0.0.1:<ephemeral-test-port>`
  - class: `test`
  - actions: `listen`
  - purpose: headless browser smoke server during `go test ./...`
  - lifecycle: test-only; closed by the test process

## Tasks

- [x] vafuk.1 Add `outcome` to structured search filters for run history.
- [x] vafuk.2 Add browser drilldown actions from place/resource/responsibility context into filtered receiving/count history.
- [x] vafuk.3 Add tests and docs for the new search/drilldown behavior.
