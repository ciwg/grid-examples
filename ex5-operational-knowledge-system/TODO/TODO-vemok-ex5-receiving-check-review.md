# TODO vemok - ex5 receiving check review

## Decision Intent Log

ID: DI-vemok
Date: 2026-07-20 17:22:00
Status: active
Decision: Add a broader `receiving_check` knowledge-item and run workflow to `ex5`, with browser review panels and context history that treat receiving work as an operational-memory flow rather than an inventory-only special case.
Intent: Cover inbound inspection and intake work for parts, tools, kits, deliveries, and similar operational assets without turning `ex5` into a full inventory or ERP subsystem.
Constraints: Use the locked name `receiving_check`; stay within the current local HTTP live-draft model per `DI-tabiv`; keep the feature standalone inside `ex5`; update docs and tests in the same slice.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-vemok-ex5-receiving-check-review.md`, `ex5-operational-knowledge-system/service/types.go`, `ex5-operational-knowledge-system/service/app.go`, `ex5-operational-knowledge-system/service/app_test.go`, `ex5-operational-knowledge-system/service/server_test.go`, `ex5-operational-knowledge-system/web/index.html`, `ex5-operational-knowledge-system/web/app.js`, `ex5-operational-knowledge-system/web/assets_test.go`, `ex5-operational-knowledge-system/web/browser_smoke_test.go`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`

ID: DI-ravip
Date: 2026-07-20 17:46:00
Status: active
Decision: Strengthen the `receiving_check` browser proof by driving the real UI into run detail during headless smoke testing, so `Receiving review` is asserted directly instead of only inferring coverage from item history.
Intent: Close the remaining evidence gap for the new receiving review panel without changing shipped product behavior.
Constraints: Keep this as a test-harness improvement only; do not introduce new browser product controls or routing just to satisfy the smoke test.
Affects: `ex5-operational-knowledge-system/TODO/TODO-vemok-ex5-receiving-check-review.md`, `ex5-operational-knowledge-system/web/browser_smoke_test.go`

## Goal

Make receiving and inbound inspection work a first-class `ex5` workflow with its
own kind, clearer browser review panels, and contextual drilldown from places,
resources, and related operational records.

## Intended repo paths before coding

- `ex5-operational-knowledge-system/service/types.go`
- `ex5-operational-knowledge-system/service/app.go`
- `ex5-operational-knowledge-system/service/app_test.go`
- `ex5-operational-knowledge-system/service/server_test.go`
- `ex5-operational-knowledge-system/web/index.html`
- `ex5-operational-knowledge-system/web/app.js`
- `ex5-operational-knowledge-system/web/assets_test.go`
- `ex5-operational-knowledge-system/web/browser_smoke_test.go`
- `ex5-operational-knowledge-system/README.md`
- `ex5-operational-knowledge-system/docs/features-guide.md`
- `ex5-operational-knowledge-system/docs/practical-implementation.md`
- `ex5-operational-knowledge-system/TODO/TODO.md`
- `ex5-operational-knowledge-system/TODO/TODO-vemok-ex5-receiving-check-review.md`

## Intended runtime path patterns before coding

- `t.TempDir()/**`
  - Class: `temp`
  - Actions: `read/write`
  - Purpose: deterministic Go service and HTTP tests
  - Cleanup plan: Go test harness owns cleanup
- `127.0.0.1:<ephemeral-test-port>`
  - Class: `temp`
  - Actions: `listen`
  - Purpose: headless browser smoke test against a local test server
  - Cleanup plan: test server closes during test teardown

## Tasks

- [x] vemok.1 Add the `receiving_check` knowledge-item and run kind to the service model, meta surface, validation, and browser forms.
- [x] vemok.2 Add receiving-focused browser review panels and contextual run-history drilldown.
- [x] vemok.3 Add tests and docs for the receiving-check workflow.
- [x] vemok.4 Prove `Receiving review` directly in browser smoke coverage.
