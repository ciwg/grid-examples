# TODO lunav - ex5 browser direct contract completion

## Decision Intent Log

ID: DI-lunav
Date: 2026-07-22 20:34:17 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track a follow-on pass to move the remaining browser direct-contract traffic off generic `type:"request"` forwarding and onto typed runtime operations.
Intent: Finish the direct browser embodiment boundary so the native-messaging path carries runtime intents directly instead of re-entering the runtime as adapter-shaped HTTP requests where that mapping is still avoidable.
Constraints: Preserve the shipped Chrome/Chromium direct embodiment, keep fail-closed behavior intact, and continue favoring one typed operation per workflow family over generic mutation envelopes.
Affects: `ex5-operational-knowledge-system/web/*`, `ex5-operational-knowledge-system/service/*`, `ex5-operational-knowledge-system/docs/*`, `ex5-operational-knowledge-system/TODO/TODO.md`

ID: DI-ronav
Date: 2026-07-22 20:38:43 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Move the remaining browser read slice onto typed runtime operations by adding named operations for dashboard, collection lists, structured search, and live-state bootstrap, while keeping the local HTTP shell and `/api/meta` bootstrap unchanged.
Intent: Remove the browser's remaining day-to-day semantic dependence on generic `type:"request"` forwarding without turning the HTTP shell/bootstrap surface into an unnecessary rewrite target.
Constraints: Keep one named operation per runtime intent, preserve the shipped Chrome/Chromium direct embodiment and fail-closed posture, and leave `/api/meta` plus shell asset loading as the explicit bootstrap boundary.
Affects: `ex5-operational-knowledge-system/web/app.js`, `ex5-operational-knowledge-system/service/local_socket.go`, `ex5-operational-knowledge-system/service/browser_host_test.go`, `ex5-operational-knowledge-system/web/browser_smoke_test.go`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Complete the remaining browser direct-contract migration above route-shaped
forwarding.

## Tasks

- [x] lunav.1 Identify which browser reads and writes still fall back to `type:"request"` instead of typed operations. See `../../docs/thought-experiments/TE-novek-ex5-browser-direct-contract-completion.md`.
- [x] lunav.2 Define the next typed browser/runtime operation slice for those remaining paths.
- [x] lunav.3 Migrate the remaining browser direct-contract gaps and align the docs.

## Status

- completed
- dashboard, catalog refresh, structured search, and live-state bootstrap now use typed browser/runtime operations instead of generic request forwarding
