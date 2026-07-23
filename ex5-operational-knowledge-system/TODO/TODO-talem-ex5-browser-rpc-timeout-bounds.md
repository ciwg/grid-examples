# TODO talem - ex5 browser rpc timeout bounds

## Decision Intent Log

ID: DI-talem
Date: 2026-07-22 20:34:17 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track a follow-on pass to bound browser direct-contract RPC waits and clean up pending requests when the extension/native-host path drops a reply.
Intent: Keep the browser fail-closed and honest under bridge failure without leaving UI actions hanging indefinitely on unresolved direct-contract promises.
Constraints: Preserve the current direct browser embodiment, keep error reporting explicit, and avoid reintroducing silent demotion to the older HTTP browser path.
Affects: `ex5-operational-knowledge-system/web/*`, `ex5-operational-knowledge-system/chrome-extension/*`, `ex5-operational-knowledge-system/docs/*`, `ex5-operational-knowledge-system/TODO/TODO.md`

ID: DI-zabem
Date: 2026-07-22 20:51:01 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Bound one-shot browser direct-contract RPCs with a 1000ms page-owned timeout that removes pending request state and rejects with an explicit direct-browser error when no reply returns.
Intent: Keep the Chrome/Chromium browser embodiment fail-closed and honest under dropped extension/native-host replies, without reviving the older HTTP browser path or leaving the UI waiting indefinitely.
Constraints: Timeout ownership must stay in the page layer that owns the user-visible promise; healthy replies and explicit extension errors must still resolve immediately; live-open acknowledgement remains separate under TODO `140`.
Affects: `ex5-operational-knowledge-system/web/app.js`, `ex5-operational-knowledge-system/web/browser_smoke_test.go`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Add bounded failure semantics for browser direct-contract one-shot RPCs.

## Tasks

- [x] talem.1 Define timeout and cleanup behavior for one-shot browser bridge RPCs. See `../../docs/thought-experiments/TE-suvik-ex5-browser-rpc-timeout-bounds.md`.
- [x] talem.2 Add regression coverage for lost reply, explicit error, and healthy reply cases.
- [x] talem.3 Update docs so browser direct-contract failure semantics are stated honestly.

## Status

- completed
- browser one-shot direct-contract RPCs now fail closed after 1000ms with local cleanup, while healthy replies and explicit bridge errors remain covered by the browser smoke and extension contract suites
