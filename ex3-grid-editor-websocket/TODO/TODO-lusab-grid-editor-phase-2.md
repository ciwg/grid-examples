# TODO lusab - grid-editor phase 2

## Decision Intent Log

ID: DI-dovoz
Date: 2026-07-13 15:10:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Implement approved Phase 2 as a browser-heavy document workflow slice on top of the existing CRDT relay, while keeping Neovim compatibility intact and avoiding new PromiseGrid-native backend semantics that belong to later phases.
Intent: Land the practical document workflow and export value now without blocking on permissions, shared document registries, or synced metadata models that belong to the later PromiseGrid-native feature phase.
Constraints: Browser gets the richer workflow surfaces; Neovim remains compatible with the same document IDs and relay flow; document metadata, recent docs, templates, and export UX may be local/browser-managed in this phase.
Affects: `ex2-grid-editor/web`, `ex2-grid-editor/docs`, `ex2-grid-editor/README.md`

ID: DI-nuvif
Date: 2026-07-13 15:10:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Phase 2 document workflow features will use a local registry/preferences layer for titles, timestamps, recent docs, templates, snapshots, and export history, while the relay remains focused on CRDT message exchange and awareness.
Intent: Keep workflow metadata coherent and testable now, and preserve a clean seam for later migration of selected metadata into PromiseGrid-native document exchange or publishing surfaces.
Constraints: The relay API stays simple in this phase; browser workflow features must still operate correctly when opening documents solely by shared URL or doc ID; new local data must not interfere with CRDT text sync.
Affects: `ex2-grid-editor/web`, `ex2-grid-editor/service`, `ex2-grid-editor/docs`

ID: DI-vafuk
Date: 2026-07-13 19:40:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Rename the browser export surface to `Export / Exchange` while keeping publish/import inside that panel until a later layout pass changes the overall toolbar.
Intent: Make publish/import easier to discover during manual review without reshuffling the whole Phase 2 toolbar and overlay structure on the eve of the server demo.
Constraints: This is a wording/discoverability change only; the existing export and exchange actions stay in the same panel; broader search/new-document UX polish remains open in `lusab.5`.
Affects: `ex2-grid-editor/web/index.html`, `ex2-grid-editor/README.md`, `ex2-grid-editor/docs/grid-editor-ui-example.md`

ID: DI-zosuf
Date: 2026-07-14 09:45:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Make `Preview` switch the browser editor into preview-only mode, while `Split View` switches into side-by-side editor-plus-preview mode.
Intent: Prevent the toolbar from feeling broken when `Split View` is already active by giving the two buttons distinct visible outcomes instead of sharing a vague preview toggle.
Constraints: This is a browser Phase 2 workflow behavior fix only; document content, markdown rendering, and export behavior must remain unchanged; the fix should be covered by pure helper tests.
Affects: `ex2-grid-editor/web/src/main.js`, `ex2-grid-editor/web/src/panes.js`, `ex2-grid-editor/web/src/panes.test.mjs`, `ex2-grid-editor/docs/grid-editor-ui-example.md`

ID: DI-rusok
Date: 2026-07-20 20:17:05 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Close the remaining Phase 2 browser workflow polish gaps by clarifying the `Find / Replace` label, adding explicit new-document feedback, and keeping the existing `Preview` versus `Split View` split while documenting it clearly.
Intent: Resolve the remaining manual confusion around search, new-doc creation, and preview/split behavior without changing the underlying relay or workflow model.
Constraints: Keep the feature scope browser-local; preserve the existing preview and split behavior semantics; update docs and tests in the same pass.
Affects: `ex3-grid-editor-websocket/web`, `ex3-grid-editor-websocket/docs`, `ex3-grid-editor-websocket/README.md`, `ex3-grid-editor-websocket/TODO`

Goal: Make documents easier to create, preview, navigate, export, and demo without changing the core CRDT relay contract.

- [x] lusab.1 Add local document registry, recent docs, title/metadata, and sample/template flows.
- [x] lusab.2 Add markdown preview, split views, search/replace, go-to-line, and document navigation tools.
- [x] lusab.3 Add export, share, snapshot, bookmark, and audit-report surfaces.
- [x] lusab.4 Add tests and docs for the Phase 2 workflow behavior.
- [x] lusab.5 Run a manual browser workflow pass for confusing labels and flow polish before closing TODO 007.
  Resolved in the follow-up pass:
  - `Find / Replace` now names the search tool directly and the overlay explains the actions more clearly.
  - `New Shared Doc` now gives immediate success feedback after creating and opening the new document id.
  - `Export / Exchange` remains explicitly named and documented as the publish/import surface.
  - `Preview` remains preview-only while `Split View` remains side-by-side, and the docs now say that plainly.
