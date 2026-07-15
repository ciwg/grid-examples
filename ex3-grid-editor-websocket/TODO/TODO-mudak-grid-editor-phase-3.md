# TODO mudak - grid-editor phase 3

## Decision Intent Log

ID: DI-safor
Date: 2026-07-13 16:05:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Implement the first approved Phase 3 slice as a browser-heavy review, comments, outline, and history layer on top of the existing CRDT relay, while keeping the relay and Neovim flow unchanged.
Intent: Add the highest-value review and visibility tools now without blocking on PromiseGrid-native restore, permissions, or shared metadata models that belong to later phases.
Constraints: Review/history state may remain browser-local in this slice; the shared CRDT text path must stay intact; new UI should be optional and not crowd the main editor by default.
Affects: `ex2-grid-editor/web`, `ex2-grid-editor/docs`, `ex2-grid-editor/README.md`

ID: DI-lapek
Date: 2026-07-13 16:05:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Phase 3 browser review data will be stored through the existing local document registry seam, including comments, activity, recent participants, named versions, and lightweight diagnostics data.
Intent: Keep comments/history coherent and testable now, and preserve a clean seam for later migration of selected review features into PromiseGrid-native publish or audit flows.
Constraints: The relay API stays unchanged in this slice; saved versions are not restore actions; activity and comment surfaces must tolerate browser-only local persistence; @mentions resolve to stable participant IDs when known.
Affects: `ex2-grid-editor/web`, `ex2-grid-editor/docs`

Goal: Add review, comments, outline, and history tools that make the browser embodiment better for collaboration review and demos.

- [x] mudak.1 Extend the local document registry with comments, activity, participant history, named versions, and summary helpers.
- [x] mudak.2 Add browser review/history UI for inline comments, activity, outline, recent participants, and focus mode.
- [x] mudak.3 Add diagnostics, summary, and version naming surfaces without changing relay semantics.
- [x] mudak.4 Add tests and docs for the Phase 3 review/history behavior.
- [x] mudak.5 Run a manual browser pass for comment flow, outline clarity, and focus-mode polish before closing TODO 008.
