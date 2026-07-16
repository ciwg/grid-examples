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

ID: DI-kotip
Date: 2026-07-16 00:00:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Accept the current ex2 browser underline behavior after forcing the CodeMirror underline decoration to inherit normal text color and cursor behavior so underline no longer presents like a link.
Intent: Close TODO 011 now that ex2 matches the accepted ex3 underline presentation while preserving the raw `<u>...</u>` storage model and unchanged preview/export/protocol behavior.
Constraints: Stored document bytes remain raw `<u>...</u>` markup; preview/export/protocol behavior stays unchanged; the fix is browser-editor presentation only.
Affects: `ex2-grid-editor/web/style.css`, `ex2-grid-editor/TODO/TODO.md`, `ex2-grid-editor/TODO/TODO-huvok-grid-editor-browser-underline.md`, `ex2-grid-editor/TODO/TODO-tizaf-grid-editor-phase-1.md`, `ex2-grid-editor/README.md`, `ex2-grid-editor/docs/practical-implementation.md`, `ex2-grid-editor/docs/architecture.md`, `ex2-grid-editor/docs/grid-editor-ui-example.md`

Goal: Make underline visibly work in the browser editor so the demo no longer shows raw underline tags as plain text with no visual underline.

- [x] huvok.1 Add a pure underline-range parser for browser tests.
- [x] huvok.2 Render `<u>...</u>` as inline underline inside CodeMirror without changing saved bytes.
- [x] huvok.3 Add regression tests and run browser build/tests.
- [x] huvok.4 Make browser-visible underline clearly visible enough for demo acceptance.

Current status:
- The current implementation stores raw underline markup correctly and now
  forces inline underline to inherit normal text styling in the browser editor.
- The user has accepted the current browser-visible underline behavior.
- TODO 011 is closed.
