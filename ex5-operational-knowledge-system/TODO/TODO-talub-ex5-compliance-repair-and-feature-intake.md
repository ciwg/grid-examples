# TODO talub - ex5 compliance repair and feature intake

## Decision Intent Log

ID: DI-talub
Date: 2026-07-20 17:08:00
Status: active
Decision: Retroactively document the recent `ex5` feature slices that were implemented before the full local decision/path intake was written down, and require every future `ex5` feature slice to start with a local TODO/DI entry that records intended repo paths and runtime path patterns before code edits begin.
Intent: Repair the traceability gap without pretending the earlier slices were truly decision-first, and prevent the same process failure from repeating on the next `ex5` feature.
Constraints: This repair may document past work, but it cannot rewrite history into a true decision-first sequence; future slices must stop and amend the local TODO before coding if they discover a new touched path or runtime pattern.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-talub-ex5-compliance-repair-and-feature-intake.md`, recent `ex5` slice TODO files `TODO-savuk-browser-navigation-and-docs.md`, `TODO-honus-ex5-workflow-search-and-automation.md`, `TODO-hozom-ex5-item-run-history.md`, `TODO-julos-ex5-context-run-history.md`, `TODO-pojul-ex5-inventory-discrepancy-review.md`

## Goal

Repair the recent `ex5` process record enough that the implemented slices are easy
to audit, and make the future feature-start rule explicit in the local queue.

## Repair scope

Recent slices that were implemented before the full local decision/path intake was
made explicit:

- `DI-vopuk` / `savuk` browser navigation and docs
- `DI-honus` workflow search and browser automation
- `DI-hozom` item run history
- `DI-julos` context run history
- `DI-pojul` inventory discrepancy review

## Retroactive path and runtime summary

Repo paths touched by those slices:

- `ex5-operational-knowledge-system/web/index.html`
- `ex5-operational-knowledge-system/web/app.js`
- `ex5-operational-knowledge-system/web/style.css`
- `ex5-operational-knowledge-system/web/assets_test.go`
- `ex5-operational-knowledge-system/web/browser_smoke_test.go`
- `ex5-operational-knowledge-system/service/app.go`
- `ex5-operational-knowledge-system/service/app_test.go`
- `ex5-operational-knowledge-system/service/server_test.go`
- `ex5-operational-knowledge-system/README.md`
- `ex5-operational-knowledge-system/docs/features-guide.md`
- `ex5-operational-knowledge-system/docs/practical-implementation.md`
- `ex5-operational-knowledge-system/TODO/TODO.md`
- `ex5-operational-knowledge-system/TODO/TODO-savuk-browser-navigation-and-docs.md`
- `ex5-operational-knowledge-system/TODO/TODO-honus-ex5-workflow-search-and-automation.md`
- `ex5-operational-knowledge-system/TODO/TODO-hozom-ex5-item-run-history.md`
- `ex5-operational-knowledge-system/TODO/TODO-julos-ex5-context-run-history.md`
- `ex5-operational-knowledge-system/TODO/TODO-pojul-ex5-inventory-discrepancy-review.md`

Runtime path patterns exercised by those slices:

- `t.TempDir()/**`
  - Class: `temp`
  - Actions: `read/write`
  - Purpose: deterministic Go service and HTTP tests
- `127.0.0.1:<ephemeral-test-port>`
  - Class: `temp`
  - Actions: `listen`
  - Purpose: headless browser smoke tests against a local test server
  - Example ports already used in this queue: `7046`, `7047`

## Future feature-start rule

Before the next `ex5` feature slice begins:

1. Create or update the local `TODO/TODO-<handle>-<slug>.md` file first.
2. Add the locking DI entry first.
3. Record intended repo paths before code edits begin.
4. Record intended runtime path patterns before code or tests that use them begin.
5. If implementation exposes a new path or runtime pattern, stop and amend the
   local TODO first, then continue.
6. Keep docs and tests in the same feature slice before commit.

## Tasks

- [x] talub.1 Record the recent `ex5` slices that need retroactive decision/path traceability repair.
- [x] talub.2 Summarize the repo paths and runtime path patterns those slices actually used.
- [x] talub.3 Make the future `ex5` feature-start rule explicit in the local TODO queue.
