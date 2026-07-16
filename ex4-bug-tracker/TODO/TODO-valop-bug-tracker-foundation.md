# TODO valop - bug tracker foundation

## Decision Intent Log

ID: DI-dajak
Date: 2026-07-16 16:00:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Implement `ex4-bug-tracker/` as a new nested Go module with a Go HTTP server, embedded browser assets, a small CLI, and runtime state rooted at `.bug-tracker/`.
Intent: Keep the new example easy to run and inspect while matching the repo's existing example structure.
Constraints: Avoid editor- and websocket-specific machinery; use a standalone module path `github.com/computerscienceiscool/grid-examples/ex4-bug-tracker`; keep runtime writes under `.bug-tracker/`.
Affects: `ex4-bug-tracker/go.mod`, `ex4-bug-tracker/cmd/**`, `ex4-bug-tracker/service/**`, `ex4-bug-tracker/web/**`, `ex4-bug-tracker/README.md`

ID: DI-nunit
Date: 2026-07-16 16:00:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Model bug tracking as append-only issue events with projected current state, and store uploaded attachment copies under the local runtime root.
Intent: Preserve a durable history for issue activity while keeping the current queue/detail views easy to compute.
Constraints: Timeline events must cover creation, comments, assignment changes, status changes, and attachments; uploaded files must be copied into app-managed storage instead of relying on external host paths.
Affects: `ex4-bug-tracker/service/app.go`, `ex4-bug-tracker/service/persistence.go`, `ex4-bug-tracker/service/server.go`, `ex4-bug-tracker/service/*_test.go`

ID: DI-ninuf
Date: 2026-07-16 16:00:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Use a browser-first queue/detail UI with built-in `reporter`, `triage`, and `engineer` identities, plus a narrow engineer CLI for assigned work.
Intent: Make the example usable from the browser while showing that the same workflow can be driven from a second embodiment.
Constraints: Fixed statuses are `New`, `Triaged`, `In Progress`, and `Resolved`; only one active assignee exists at a time; the CLI focuses on assigned work instead of full administration.
Affects: `ex4-bug-tracker/cmd/**`, `ex4-bug-tracker/service/**`, `ex4-bug-tracker/web/index.html`, `ex4-bug-tracker/web/app.js`, `ex4-bug-tracker/web/style.css`

ID: DI-gofub
Date: 2026-07-16 16:00:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Keep v1 single-team in behavior while storing a hidden built-in `team` field with default value `CORE`, and allow reopen by moving `Resolved` issues back to `Triaged` while clearing the assignee.
Intent: Leave a clean seam for future multi-team work without complicating the initial UI or queue flow.
Constraints: Do not expose team controls in the v1 browser or CLI; reopening must preserve prior history while resetting active ownership.
Affects: `ex4-bug-tracker/service/app.go`, `ex4-bug-tracker/service/server.go`, `ex4-bug-tracker/web/app.js`, `ex4-bug-tracker/README.md`

ID: DI-zogof
Date: 2026-07-16 17:05:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add a repeatable seeded demo path that starts the normal bug tracker server against a fresh temp data root and preloads a small believable issue set only when that runtime is empty.
Intent: Let people evaluate the app quickly without hand-entering starter issues, while keeping the demo on the same browser, CLI, and append-only workflow paths as normal usage.
Constraints: Do not invent a separate demo-only UI; keep the seeded data local to the runtime root; avoid duplicating seed data on repeated starts against the same existing root.
Affects: `ex4-bug-tracker/service/demo.go`, `ex4-bug-tracker/service/*_test.go`, `ex4-bug-tracker/cmd/bug-tracker/main.go`, `ex4-bug-tracker/scripts/run-demo.sh`, `ex4-bug-tracker/README.md`, `ex4-bug-tracker/docs/**`, `ex4-bug-tracker/INPROGRESS-FEATURES.md`

## Goal

Build the first complete `ex4-bug-tracker` slice with durable issue storage,
browser UI, CLI commands, and attachment handling.

## Tasks

- [x] valop.1 Add the new module, runtime root, and local docs/TODO structure.
- [x] valop.2 Implement append-only issue storage and projected current state.
- [x] valop.3 Implement browser HTTP routes and file upload/download behavior.
- [x] valop.4 Implement the browser queue/detail UI and engineer CLI.
- [x] valop.5 Verify the example with Go tests.
- [x] valop.6 Add a seeded demo launch path and operator docs.
