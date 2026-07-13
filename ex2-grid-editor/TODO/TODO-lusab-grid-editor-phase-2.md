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

Goal: Make documents easier to create, preview, navigate, export, and demo without changing the core CRDT relay contract.

- [x] lusab.1 Add local document registry, recent docs, title/metadata, and sample/template flows.
- [x] lusab.2 Add markdown preview, split views, search/replace, go-to-line, and document navigation tools.
- [x] lusab.3 Add export, share, snapshot, bookmark, and audit-report surfaces.
- [x] lusab.4 Add tests and docs for the Phase 2 workflow behavior.
- [ ] lusab.5 Run a manual browser workflow pass for confusing labels and flow polish before closing TODO 007.
