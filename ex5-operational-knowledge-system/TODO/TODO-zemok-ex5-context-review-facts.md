# TODO zemok - ex5 context review facts

## Decision Intent Log

ID: DI-zemok
Date: 2026-07-20 21:02:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Expand `ex5` place/resource/responsibility review panels to show receiving and inventory fact history from related runs, instead of only lightweight run summaries.
Intent: Make count/discrepancy history practical from context anchors like places and resources so operators can review what happened there without drilling into every run one by one.
Constraints: Stay in the current operational-memory lane; do not add ERP/MRP quantity logic; use the existing local HTTP live-draft model; update docs and tests in the same slice.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-zemok-ex5-context-review-facts.md`, `ex5-operational-knowledge-system/web/app.js`, `ex5-operational-knowledge-system/web/assets_test.go`, `ex5-operational-knowledge-system/web/browser_smoke_test.go`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`, `ex5-operational-knowledge-system/README.md`

## Goal

Show clearer receiving and inventory review facts from place/resource/responsibility context views, so operators can inspect count/discrepancy history without treating the browser inspector like raw JSON.

## Intended repo paths

- `ex5-operational-knowledge-system/TODO/TODO.md`
- `ex5-operational-knowledge-system/TODO/TODO-zemok-ex5-context-review-facts.md`
- `ex5-operational-knowledge-system/web/app.js`
- `ex5-operational-knowledge-system/web/assets_test.go`
- `ex5-operational-knowledge-system/web/browser_smoke_test.go`
- `ex5-operational-knowledge-system/docs/features-guide.md`
- `ex5-operational-knowledge-system/docs/practical-implementation.md`
- `ex5-operational-knowledge-system/README.md`

## Intended runtime path patterns

- `t.TempDir()/**`
  - class: `test`
  - actions: `read/write`
  - purpose: browser smoke temporary server/profile roots
  - lifecycle: test-only; auto-cleaned by the Go test harness

- `127.0.0.1:<ephemeral-test-port>`
  - class: `test`
  - actions: `listen`
  - purpose: headless browser smoke server during `go test ./...`
  - lifecycle: test-only; closed by the test process

## Tasks

- [x] zemok.1 Add richer receiving and inventory fact history sections for place/resource/responsibility inspector views.
- [x] zemok.2 Keep item/run review readable while reusing the same fact-formatting helpers.
- [x] zemok.3 Add tests and docs for the richer context review history.

## Status

- context detail views now show evidence-backed receiving and inventory fact history from related runs
- item review reuses the same richer receiving and inventory history entries
- docs and browser/service tests cover the richer context review path
