# TODO pavur - ex5 browser direct contract integration coverage

## Decision Intent Log

ID: DI-pavur
Date: 2026-07-22 19:58:07 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track a follow-on pass to strengthen browser direct-contract integration coverage beyond the current mock bridge smoke tests.
Intent: Verify more of the shipped Chrome/Chromium extension and native-host boundary directly so browser tests do not overstate how much of the real embodiment path is covered.
Constraints: Keep tests deterministic and explicit about which layers are real versus mocked.
Affects: `ex5-operational-knowledge-system/web/*`, `ex5-operational-knowledge-system/chrome-extension/*`, `ex5-operational-knowledge-system/service/*`, `ex5-operational-knowledge-system/docs/*`, `ex5-operational-knowledge-system/TODO/TODO.md`

ID: DI-vasem
Date: 2026-07-22 20:19:56 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add deterministic script-level contract tests for the shipped `background.js` and `content.js` paths, covering startup/readiness, one-shot RPC forwarding, and live-port forwarding/disconnect behavior together, while keeping the browser smoke suite explicitly documented as a page-level UI mock.
Intent: Strengthen PromiseGrid alignment by testing the real browser embodiment boundary directly without escalating into a fragile full Chrome/native-host registration harness.
Constraints: Keep the tests deterministic and local, use the real shipped extension scripts instead of reimplemented test doubles for the covered boundary, and state clearly which browser smoke coverage remains mocked.
Affects: `ex5-operational-knowledge-system/chrome-extension/assets_test.go`, `ex5-operational-knowledge-system/web/browser_smoke_test.go`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `docs/thought-experiments/TE-lavem-ex5-browser-direct-contract-integration-coverage.md`, `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Add stronger coverage for the real browser direct-contract boundary instead of
relying mainly on a synthetic in-page mock bridge.

## Tasks

- [x] pavur.1 Identify which shipped browser layers are only covered through `withMockBrowserBridge(...)`. See `../../docs/thought-experiments/TE-lavem-ex5-browser-direct-contract-integration-coverage.md`.
- [x] pavur.2 Add the next strongest deterministic coverage for extension/native-host behavior and fail-closed startup paths.
- [x] pavur.3 Update docs or test notes so the remaining mocked versus real coverage boundary is explicit.

## Status

- completed
- deterministic extension/native-host contract tests now cover readiness, one-shot RPC forwarding, and live-port forwarding/disconnect behavior, while the browser smoke suite states its remaining page-level mock boundary explicitly
