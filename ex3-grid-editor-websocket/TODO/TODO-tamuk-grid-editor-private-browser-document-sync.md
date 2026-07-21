# TODO tamuk - grid-editor private-browser document sync

## Decision Intent Log

ID: DI-ribaf
Date: 2026-07-20 19:59:35 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track the private/incognito browser document-sync failure as its own deferred ex3 follow-up instead of leaving it buried inside the broader Phase 1 manual-pass notes.
Intent: Presence and participant visibility currently work in private/incognito sessions while the shared document text can diverge or fail to appear, so the bug needs explicit follow-up ownership before later ex3 polish work can claim stable browser interoperability.
Constraints: Scope this TODO to browser-session-mode interoperability in ex3; do not assume private-mode storage isolation is the whole root cause; preserve the distinction between awareness/presence success and document-text sync failure during diagnosis.
Affects: `ex3-grid-editor-websocket/web`, `ex3-grid-editor-websocket/service`, `ex3-grid-editor-websocket/docs`, `ex3-grid-editor-websocket/TODO`

ID: DI-bonuv
Date: 2026-07-20 20:17:05 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Harden ex3 for fresh/private browser sessions by falling back to in-memory browser storage when local/session storage are blocked and by adding late-join browser regression coverage for shared document text over websocket.
Intent: Remove the most plausible browser-session-mode hazards that can make a fresh/private browser diverge from an already-open shared document, even before a full manual private-browser verification pass is rerun.
Constraints: This is a hardening and proof pass; keep the TODO open until the real private/incognito browser flow is manually rechecked.
Affects: `ex3-grid-editor-websocket/web`, `ex3-grid-editor-websocket/service/interoperability_test.go`, `ex3-grid-editor-websocket/docs`, `ex3-grid-editor-websocket/TODO`

ID: DI-sulor
Date: 2026-07-20 20:33:44 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Treat an empty browser document after websocket startup as a relay-catch-up failure when the relay already reports authoritative history, and recover by forcing one bounded HTTP sync replay from the authoritative starting offset.
Intent: Private/incognito sessions can still show peers and other relay metadata while opening with blank text, so the browser needs a second text-recovery path that does not depend on websocket replay succeeding perfectly on first load.
Constraints: Keep the normal websocket path primary; only invoke the HTTP catch-up when the relay already has history and the browser text is still empty; add regression coverage for the mismatch instead of closing the TODO on theory alone.
Affects: `ex3-grid-editor-websocket/web/src/automerge-relay.js`, `ex3-grid-editor-websocket/web/src/main.js`, `ex3-grid-editor-websocket/web/src/startup.js`, `ex3-grid-editor-websocket/web/src/*.test.mjs`, `ex3-grid-editor-websocket/docs`, `ex3-grid-editor-websocket/TODO`

ID: DI-vobek
Date: 2026-07-21 11:20:06 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: When ex3 startup still has blank text after priming from relay state, distrust the advertised snapshot offset and replay HTTP sync from offset `0` against a fresh empty Automerge base instead of continuing from the snapshot offset.
Intent: A stale or blank relay snapshot can otherwise cause private/incognito browsers to skip the real shared text forever while still showing presence and other relay metadata, because recovery starts after the history that actually contains the document body.
Constraints: Keep normal snapshot bootstrap intact for non-empty startup text; only reset to full-history replay when the browser is already in the blank-text recovery path; prove the fix with both client-level regression coverage and a headless browser startup test that injects a poisoned relay state response.
Affects: `ex3-grid-editor-websocket/web/src/automerge-relay.js`, `ex3-grid-editor-websocket/web/src/automerge-relay.test.mjs`, `ex3-grid-editor-websocket/service/browser_startup_test.go`, `ex3-grid-editor-websocket/README.md`, `ex3-grid-editor-websocket/docs/features-guide.md`, `ex3-grid-editor-websocket/TODO`

Goal: Diagnose and fix the ex3 bug where private/incognito browser sessions can show who is present and other collaboration state while failing to converge on the shared document text.

- [ ] tamuk.1 Reproduce the mismatch with at least one normal browser window and one private/incognito browser window against the same ex3 document and relay.
- [x] tamuk.2 Determine whether the failure is in bootstrap document loading, live sync replay, local draft seeding, storage/session isolation assumptions, or another browser-only path.
- [x] tamuk.3 Fix document convergence so private/incognito sessions receive the same shared text as normal browser sessions without regressing awareness/presence behavior.
- [x] tamuk.4 Add regression coverage for mixed normal/private browser sessions if the test harness can support it; otherwise document the remaining manual-proof requirement precisely.
- [x] tamuk.5 Update ex3 docs with the actual private/incognito support status and any remaining caveats.

Current status:
- fresh browser sessions now have an explicit late-join regression test for shared text over websocket
- browser-local registry and preference storage now fall back to in-memory state if local/session storage are blocked
- the real browser startup path now has direct helper coverage for participant-id and welcome-banner state when browser storage is restricted
- browser startup now forces one bounded HTTP sync catch-up when websocket startup leaves the editor blank even though the relay reports authoritative history
- blank startup recovery now ignores a stale empty snapshot offset and replays from full relay history instead of skipping the real shared text forever
- headless Chrome now probes the real built app and asserts that a fresh browser late-join renders shared document text in the DOM
- real manual private/incognito browser verification is still pending before TODO 016 can close fully
