# TODO tizaf - grid-editor phase 1

## Decision Intent Log

ID: DI-favok
Date: 2026-07-13 22:40:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Implement approved Phase 1 as one milestone, with browser-first UX and Neovim core collaboration/status parity rather than strict feature-for-feature parity.
Intent: Deliver the highest-value collaborative-editor polish in one pass while keeping Neovim viable for demos without blocking on identical UI richness across embodiments.
Constraints: Browser should receive the fuller settings, onboarding, and menu UX; Neovim must still receive core peer visibility, remote selections, status clarity, and shortcut/help discoverability.
Affects: `ex2-grid-editor/web`, `ex2-grid-editor/nvim`, `ex2-grid-editor/README.md`

ID: DI-vasul
Date: 2026-07-13 22:40:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Phase 1 preferences will persist locally through a PromiseGrid-shaped abstraction, ship with demo-friendly defaults, and expose a lightweight onboarding/help surface.
Intent: Make Grid Editor feel good in demos now while leaving a clean seam for future PromiseGrid-synced preferences instead of baking local-storage behavior directly into every UI interaction.
Constraints: Preferences remain local in this phase; supported settings include display name, color, theme, line numbers, font size, dyslexia-friendly spacing/font mode, and shortcut overrides where supported.
Affects: `ex2-grid-editor/web`, `ex2-grid-editor/nvim`, `ex2-grid-editor/docs`, `ex2-grid-editor/README.md`

Goal: Finish the first collaborative-editor UX milestone so browser and Neovim both present clear collaboration, status, and editor controls on top of the existing CRDT relay.

- [x] tizaf.1 Add Phase 1 browser settings, menu, onboarding, and status surfaces.
- [x] tizaf.2 Add presence aging, peer lifecycle feedback, and clearer collaboration visibility.
- [x] tizaf.3 Add Neovim core parity improvements for selections, status, and help.
- [x] tizaf.4 Add tests and docs for the new Phase 1 behavior.
- [ ] tizaf.5 Run a manual demo pass and capture any remaining Phase 1 polish gaps before closing TODO 006.
  Current gaps from the manual pass:
  - Chrome still needs confirmation that remote cursor line color renders correctly instead of a tiny black line.
  - Toolbar formatting needs a clear product decision on whether raw `<u>...</u>` text in the editor is acceptable or whether inline rendered underline is required.
