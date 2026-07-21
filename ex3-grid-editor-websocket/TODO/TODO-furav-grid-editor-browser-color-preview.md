# TODO furav - grid-editor browser color chooser visibility

## Decision Intent Log

ID: DI-pafob
Date: 2026-07-14 09:32:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Fix the Chrome color chooser visibility bug by adding a persistent in-app color preview swatch and visible hex value in the `You` card instead of relying only on the native `input[type=color]` preview.
Intent: Make the user’s chosen presence color obvious before the demo, even in Chrome where the native color control can collapse into a thin line with no meaningful preview.
Constraints: The underlying saved color value stays the same; awareness color propagation stays the same; this is a browser-UI visibility fix, not a protocol or relay change.
Affects: `ex2-grid-editor/web/index.html`, `ex2-grid-editor/web/style.css`, `ex2-grid-editor/web/src/main.js`, `ex2-grid-editor/web/src/color.js`, `ex2-grid-editor/web/src/color.test.mjs`

ID: DI-vonol
Date: 2026-07-14 09:46:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Keep the Chrome color chooser bug open because the added swatch/hex display is not accepted as the final fix for the native preview/visibility problem.
Intent: Track the bug honestly for the demo instead of treating a partial UI workaround as closure.
Constraints: The current swatch/hex change may stay in the worktree for evaluation, but TODO 012 remains open until the actual Chrome-side visibility problem is resolved to user satisfaction.
Affects: `ex2-grid-editor/TODO/TODO.md`, `ex2-grid-editor/TODO/TODO-furav-grid-editor-browser-color-preview.md`, pending browser UI follow-up work

ID: DI-zafuk
Date: 2026-07-14 10:20:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add a persistent `You` badge to the peer-badge row so the current user can always see their own chosen presence color without opening the color picker or reading the separate settings area.
Intent: Keep the demo-ready presence surface honest and immediately legible in Chrome even while the native color chooser bug remains open.
Constraints: The main peer count stays remote-only; this is a browser-UI visibility aid, not a protocol change; TODO 012 remains open until the Chrome-side visibility problem is fully satisfactory.
Affects: `ex2-grid-editor/web/src/main.js`, `ex2-grid-editor/web/style.css`, `ex2-grid-editor/TODO/TODO-furav-grid-editor-browser-color-preview.md`

ID: DI-fegom
Date: 2026-07-20 20:17:05 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Resolve the remaining Chrome color-picker visibility gap by making the visible swatch/name/value row the primary button for opening the native picker instead of relying on the collapsed native preview line.
Intent: Keep the real browser color chooser available while making the chosen color legible and clickable in Chrome even when the built-in preview strip is weak.
Constraints: Preserve the same saved hex value and awareness protocol behavior; keep the browser interaction accessible with or without `showPicker()`.
Affects: `ex3-grid-editor-websocket/web`, `ex3-grid-editor-websocket/docs`, `ex3-grid-editor-websocket/README.md`, `ex3-grid-editor-websocket/TODO`

Goal: Make the user’s current color obvious in the browser UI without opening the native color picker.

- [x] furav.1 Add a pure helper for visible color display state.
- [x] furav.2 Add a persistent swatch and hex value next to the browser color input.
- [x] furav.3 Add regression tests and rebuild the browser bundle.
- [x] furav.4 Resolve the actual Chrome native color-chooser visibility problem to user satisfaction.

Current status:
- The visible color row is now the primary picker button, so Chrome no longer depends on the tiny native preview line for discoverability.
- The same hex value and awareness color still flow through the existing browser and relay paths.
