# TODO rumav - ex5 browser direct contract above routes

## Decision Intent Log

ID: DI-rumav
Date: 2026-07-22 19:58:07 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track a follow-on pass to move more browser semantics off generic `method + path` forwarding and onto typed direct-contract runtime operations.
Intent: Improve PromiseGrid alignment by reducing the amount of browser traffic that still re-enters the runtime as adapter-shaped request forwarding over the native-messaging bridge.
Constraints: Preserve the current shipped browser embodiment and keep any migration explicit; do not reintroduce silent fallback to the older HTTP browser path.
Affects: `ex5-operational-knowledge-system/web/*`, `ex5-operational-knowledge-system/chrome-extension/*`, `ex5-operational-knowledge-system/service/*`, `ex5-operational-knowledge-system/docs/*`, `ex5-operational-knowledge-system/TODO/TODO.md`

ID: DI-lorim
Date: 2026-07-22 20:04:26 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Move the browser create/operate mutation slice onto one typed operation per workflow family, covering create place/resource/responsibility/item, record run, add evidence, record approvals, add revision, and supersede item.
Intent: Improve PromiseGrid alignment at the more important durable write boundary first, while keeping the broader browser direct-contract migration scoped and testable.
Constraints: Keep browser startup/catalog reads on their current paths for now, keep the browser fail-closed with no silent HTTP fallback, and route evidence attachments through the typed operation family instead of leaving them on raw page HTTP.
Affects: `ex5-operational-knowledge-system/web/app.js`, `ex5-operational-knowledge-system/service/local_socket.go`, `ex5-operational-knowledge-system/service/local_socket_types.go`, `ex5-operational-knowledge-system/service/browser_host_test.go`, `ex5-operational-knowledge-system/web/browser_smoke_test.go`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Raise more of the browser direct contract above route-shaped semantics so the
native-messaging embodiment carries runtime intents more directly.

## Tasks

- [x] rumav.1 Identify which browser reads and writes still tunnel through generic `request` forwarding. See `../../docs/thought-experiments/TE-zovek-ex5-browser-direct-contract-above-routes.md`.
- [x] rumav.2 Define the next typed browser/runtime operation slice beyond the current inspect/search/problem-review set. See `../../docs/thought-experiments/TE-zovek-ex5-browser-direct-contract-above-routes.md`.
- [x] rumav.3 Migrate the highest-value browser workflows and align the contract docs.

## Status

- completed
- browser create/operate mutations now use typed direct-contract operations instead of generic route-shaped request forwarding
