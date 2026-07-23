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

## Goal

Add stronger coverage for the real browser direct-contract boundary instead of
relying mainly on a synthetic in-page mock bridge.

## Tasks

- [ ] pavur.1 Identify which shipped browser layers are only covered through `withMockBrowserBridge(...)`.
- [ ] pavur.2 Add the next strongest deterministic coverage for extension/native-host behavior and fail-closed startup paths.
- [ ] pavur.3 Update docs or test notes so the remaining mocked versus real coverage boundary is explicit.

## Status

- open
- created from the post-`133` review finding on browser direct-contract test realism
