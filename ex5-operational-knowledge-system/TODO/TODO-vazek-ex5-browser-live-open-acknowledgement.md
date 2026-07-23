# TODO vazek - ex5 browser live open acknowledgement

## Decision Intent Log

ID: DI-vazek
Date: 2026-07-22 20:34:17 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track a follow-on pass to make browser live transport report itself connected only after a real live-open acknowledgement instead of immediately on outbound send.
Intent: Tighten the browser live embodiment so connection status reflects confirmed direct-contract state rather than an optimistic local assumption.
Constraints: Preserve the current native-messaging live path, keep reconnect behavior explicit, and avoid hiding missing acknowledgements behind silent fallback.
Affects: `ex5-operational-knowledge-system/web/*`, `ex5-operational-knowledge-system/chrome-extension/*`, `ex5-operational-knowledge-system/docs/*`, `ex5-operational-knowledge-system/TODO/TODO.md`

ID: DI-talik
Date: 2026-07-22 21:00:23 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Mark the browser live lane connected only after the first successful runtime `live-state` or `live-conflict` reply, and clear connected state again on active live-lane errors before scheduling reconnect.
Intent: Keep browser live connection truth anchored to real runtime acknowledgement instead of optimistic local send, while preserving the existing native-messaging lane and explicit reconnect behavior.
Constraints: Do not invent a synthetic browser-only acknowledgement message; preserve the current direct live contract family; keep missing acknowledgement or live-lane failure fail-closed with no silent HTTP fallback.
Affects: `ex5-operational-knowledge-system/web/app.js`, `ex5-operational-knowledge-system/web/browser_smoke_test.go`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Make browser live-transport connected state depend on explicit acknowledgement.

## Tasks

- [x] vazek.1 Define what counts as a successful live-open acknowledgement for the browser embodiment. See `../../docs/thought-experiments/TE-morav-ex5-browser-live-open-acknowledgement.md`.
- [x] vazek.2 Align browser live state transitions and reconnect logic to that acknowledgement boundary.
- [x] vazek.3 Add coverage and doc updates for the new live-open truth semantics.

## Status

- completed
- browser live transport now stays in an opening state until the first `live-state` or `live-conflict` reply arrives, and active live-lane errors clear connected state before reconnect
