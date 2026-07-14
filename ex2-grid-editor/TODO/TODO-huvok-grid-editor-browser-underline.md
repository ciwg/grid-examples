# TODO huvok - grid-editor browser inline underline rendering

## Decision Intent Log

ID: DI-naruv
Date: 2026-07-14 09:20:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Fix the browser underline feature by rendering `<u>...</u>` as visible inline underline inside the CodeMirror editor while preserving the raw source bytes for saving, export, preview, and protocol traffic.
Intent: Make underline visibly work in the live browser editor before the demo without changing the persisted document format or inventing a new wire-level formatting rule.
Constraints: Raw document text must still store literal `<u>` tags; markdown export must stay raw; preview and HTML export behavior must stay intact; the fix is browser-editor rendering only.
Affects: `ex2-grid-editor/web/src/editor.js`, `ex2-grid-editor/web/src/underline.js`, `ex2-grid-editor/web/style.css`, `ex2-grid-editor/web/src/*.test.mjs`

ID: DI-timab
Date: 2026-07-14 09:55:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Reopen the underline bug because the current browser behavior is still not accepted as visibly underlined enough for the demo.
Intent: Keep the bug open until the user confirms the browser experience really shows underline, not just raw tags, parser correctness, or partial rendering behavior.
Constraints: Passing tests are not enough; acceptance depends on visible browser output.
Affects: `ex2-grid-editor/TODO/TODO.md`, `ex2-grid-editor/TODO/TODO-huvok-grid-editor-browser-underline.md`, pending browser rendering follow-up work

Goal: Make underline visibly work in the browser editor so the demo no longer shows raw underline tags as plain text with no visual underline.

- [x] huvok.1 Add a pure underline-range parser for browser tests.
- [x] huvok.2 Render `<u>...</u>` as inline underline inside CodeMirror without changing saved bytes.
- [x] huvok.3 Add regression tests and run browser build/tests.
- [ ] huvok.4 Make browser-visible underline clearly visible enough for demo acceptance.

Current status:
- The current implementation is not accepted as visibly underlined.
- Parser/build behavior alone is not considered sufficient.
- TODO 011 remains open.
