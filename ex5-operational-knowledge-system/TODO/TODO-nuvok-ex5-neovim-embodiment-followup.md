# TODO nuvok - ex5 Neovim embodiment followup

## Decision Intent Log

ID: DI-nuvok
Date: 2026-07-20 21:35:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track the future `ex5` Neovim embodiment as its own deferred TODO instead of keeping it bundled under inventory follow-on work.
Intent: Keep the inventory backlog honest and keep a future Neovim embodiment visible as a separate embodiment project.
Constraints: This TODO is deferred; it does not imply that Neovim is implemented now, and it does not reopen the decision to port the full `ex3` websocket model into `ex5`.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-ranor-ex5-inventory-and-embodiment-followups.md`, `ex5-operational-knowledge-system/TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md`

ID: DI-nuvop
Date: 2026-07-20 22:30:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Keep this TODO open for Neovim follow-on work beyond the new live-draft phase 1 implementation.
Intent: Make the docs honest that `ex5` now has a real first Neovim phase while preserving a visible backlog for richer embodiment features.
Constraints: Follow-on scope remains separate from inventory TODO `007`; later Neovim work must stay aligned with the current local HTTP live-draft model unless a new decision changes that.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-fudok-ex5-neovim-live-draft-phase1.md`, `ex5-operational-knowledge-system/TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`

## Goal

Track future Neovim embodiment work for `ex5` beyond the implemented live-draft
phase 1.

## Tasks

- [x] nuvok.1 Define the scope of a Neovim operational embodiment for `ex5`.
- [x] nuvok.2 Decide that the embodiment is staged, and implement the first read/write live-draft phase under `TODO/TODO-fudok-ex5-neovim-live-draft-phase1.md`.
- [x] nuvok.3 Define the first richer post-phase-1 workflow surface as read-only inspector navigation for item metadata, revisions, approvals, and related runs under `TODO/TODO-lonuk-ex5-neovim-item-inspector-phase.md`.
- [ ] nuvok.4 Define the next richer Neovim workflow surface, such as approvals, run review beyond item-linked summaries, or typed-link browsing.

## Status

- deferred
- desired for real team and customer workflows
- intentionally separate from inventory TODO `007`
- phase 1 now exists as a thin live-draft embodiment over the local HTTP runtime
- item inspection now exists as the first richer follow-on over projected item detail
