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

ID: DI-tuzok
Date: 2026-07-22 12:24:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Close TODO `016` as a completed historical umbrella instead of leaving it open once all concrete child slices are shipped.
Intent: Keep the ex5 backlog honest by removing a phantom “still open” embodiment project when the tracked Neovim follow-ons listed here are already complete.
Constraints: Future Neovim work should reopen as new concrete TODOs rather than by keeping this umbrella artificially open.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md`

## Goal

Track future Neovim embodiment work for `ex5` beyond the implemented live-draft
phase 1.

## Tasks

- [x] nuvok.1 Define the scope of a Neovim operational embodiment for `ex5`.
- [x] nuvok.2 Decide that the embodiment is staged, and implement the first read/write live-draft phase under `TODO/TODO-fudok-ex5-neovim-live-draft-phase1.md`.
- [x] nuvok.3 Define the first richer post-phase-1 workflow surface as read-only inspector navigation for item metadata, revisions, approvals, and related runs under `TODO/TODO-lonuk-ex5-neovim-item-inspector-phase.md`.
- [x] nuvok.4 Define the next richer Neovim workflow surface as direct read-only run inspection under `TODO/TODO-ravok-ex5-neovim-run-inspector-phase.md`.
- [x] nuvok.5 Define the next richer Neovim workflow surface as read-only typed-link browsing under `TODO/TODO-zalor-ex5-neovim-typed-link-browsing-phase.md`.
- [x] nuvok.6 Define the next richer Neovim workflow surface after typed-link browsing as read-only search and browse under `TODO/TODO-givot-ex5-neovim-search-browse-phase.md`.
- [x] nuvok.7 Define the next richer Neovim workflow surface after search/browse as a read-only pending-review view under `TODO/TODO-lorav-ex5-neovim-pending-review-phase.md`.
- [x] nuvok.8 Define the next richer Neovim workflow surface after pending review as a small item approval action under `TODO/TODO-vamor-ex5-neovim-item-approval-phase.md`.
- [x] nuvok.9 Define the next richer Neovim workflow surface after item approval as a small run approval action under `TODO/TODO-bafor-ex5-neovim-run-approval-phase.md`.
- [x] nuvok.10 Define the next richer Neovim workflow surface after run approval as a small item supersede action under `TODO/TODO-pudor-ex5-neovim-item-supersede-phase.md`.
- [x] nuvok.11 Define the next richer terminal-first follow-on after item supersede as CLI evidence upload, typed-link creation, structured/problem search, grouped problem review, and responsibility detail parity under `TODO/TODO-zanub-ex5-cli-evidence-upload.md`, `TODO/TODO-vuteg-ex5-cli-typed-link-creation.md`, `TODO/TODO-mifot-ex5-cli-structured-search-filters.md`, `TODO/TODO-nuvaz-ex5-cli-problem-review.md`, and `TODO/TODO-jubav-ex5-cli-responsibility-detail.md`.
- [x] nuvok.12 Define the next richer terminal-first follow-on after that as CLI pending-review aggregation under `TODO/TODO-vabok-ex5-cli-pending-review.md`.
- [x] nuvok.13 Define the next richer terminal-first follow-on after CLI pending review under `TODO/TODO-zovam-ex5-terminal-followon-slice-definition.md`, with the adjacent doc cleanup in `TODO/TODO-fudab-ex5-terminal-doc-current-state-cleanup.md` and the next concrete slice recorded as `TODO/TODO-ravum-ex5-cli-review-queue-rendering.md`.
- [x] nuvok.14 Define the next richer Neovim workflow surface after related-run handoffs as a grouped problem-review view under `TODO/TODO-sivok-ex5-neovim-problem-review-phase.md`.
- [x] nuvok.15 Define the next richer Neovim workflow surface after grouped problem review as structured search filters under `TODO/TODO-fanub-ex5-neovim-structured-search-filters.md`.
- [x] nuvok.16 Define the next richer Neovim workflow follow-on after structured search filters as a real `:OksSearch` command-surface repair under `TODO/TODO-lavup-ex5-neovim-search-command-arity.md`.
- [x] nuvok.17 Define the adjacent terminal test-hardening slice as direct Ex-command coverage under `TODO/TODO-rozaf-ex5-neovim-ex-command-coverage.md`.
- [x] nuvok.18 Define the adjacent high-level doc cleanup slice under `TODO/TODO-dorun-ex5-terminal-doc-surface-drift.md`.
- [x] nuvok.19 Define the next richer Neovim authoring follow-on as durable revision snapshot creation under `TODO/TODO-jabup-ex5-neovim-revision-snapshot-phase.md`.
- [x] nuvok.20 Define the adjacent write-side command coverage slice under `TODO/TODO-zafot-ex5-neovim-write-command-coverage.md`.
- [x] nuvok.21 Define the adjacent terminal authoring-handoff doc slice under `TODO/TODO-vogar-ex5-terminal-authoring-handoff-docs.md`.
- [x] nuvok.22 Define the next shell-only authoring follow-on as CLI durable revision snapshot parity under `TODO/TODO-muvok-ex5-cli-revision-snapshot-parity.md`.
- [x] nuvok.23 Define the adjacent older inspect command coverage slice under `TODO/TODO-taruv-ex5-neovim-inspect-command-coverage.md`.
- [x] nuvok.24 Define the adjacent terminal capability matrix drift slice under `TODO/TODO-razim-ex5-terminal-capability-matrix-drift.md`.
- [x] nuvok.25 Define the adjacent architecture terminal-surface drift slice under `TODO/TODO-favun-ex5-architecture-terminal-surface-drift-2.md`.

## Status

- deferred
- desired for real team and customer workflows
- intentionally separate from inventory TODO `007`
- phase 1 now exists as a thin live-draft embodiment over the local HTTP runtime
- item inspection now exists as the first richer follow-on over projected item detail
- direct run inspection now exists for evidence and approval review over projected run detail
- typed-link browsing now exists over item, run, place, resource, and responsibility detail
- terminal review queues now reject omitted run `approvals` fields instead of silently inventing fake unreviewed work under `TODO/TODO-davur-ex5-review-queue-approvals-contract.md`
- terminal run review now hands off into related item, place, resource, and responsibility context, and the older inspect commands now have headless behavior coverage under `TODO/TODO-vunep-ex5-run-context-handoffs.md` and `TODO/TODO-zorik-ex5-neovim-inspect-behavior-tests.md`
- related-run sections inside the older Neovim inspectors now emit direct `:OksInspectRun` handoffs under `TODO/TODO-josav-ex5-neovim-related-run-handoffs.md`
- grouped problem review now exists as a dedicated Neovim hotspot view under `TODO/TODO-sivok-ex5-neovim-problem-review-phase.md`
- Neovim search now supports the same shared structured-filter vocabulary the CLI already uses under `TODO/TODO-fanub-ex5-neovim-structured-search-filters.md`
- the next Neovim follow-ons are now split into explicit child TODOs for `:OksSearch` command-surface repair, direct Ex-command coverage, and high-level terminal-doc drift under `TODO/TODO-lavup-ex5-neovim-search-command-arity.md`, `TODO/TODO-rozaf-ex5-neovim-ex-command-coverage.md`, and `TODO/TODO-dorun-ex5-terminal-doc-surface-drift.md`
- the next Neovim follow-ons are now split into explicit child TODOs for durable revision snapshots, write-side Ex-command coverage, and terminal authoring-handoff docs under `TODO/TODO-jabup-ex5-neovim-revision-snapshot-phase.md`, `TODO/TODO-zafot-ex5-neovim-write-command-coverage.md`, and `TODO/TODO-vogar-ex5-terminal-authoring-handoff-docs.md`
- durable revision snapshots now exist inside the Neovim live-draft embodiment, and the direct write-side `:Oks...` commands now have headless command-path coverage under `TODO/TODO-jabup-ex5-neovim-revision-snapshot-phase.md` and `TODO/TODO-zafot-ex5-neovim-write-command-coverage.md`
- the next terminal follow-ons are now split into explicit child TODOs for CLI durable revision snapshot parity, older inspect command coverage, terminal capability matrix drift, and architecture-doc drift under `TODO/TODO-muvok-ex5-cli-revision-snapshot-parity.md`, `TODO/TODO-taruv-ex5-neovim-inspect-command-coverage.md`, `TODO/TODO-razim-ex5-terminal-capability-matrix-drift.md`, and `TODO/TODO-favun-ex5-architecture-terminal-surface-drift-2.md`
- CLI durable revision snapshots now exist, the older `:OksInspect...` commands now have direct headless command-path coverage, and the remaining terminal matrix and architecture drift are closed under `TODO/TODO-muvok-ex5-cli-revision-snapshot-parity.md`, `TODO/TODO-taruv-ex5-neovim-inspect-command-coverage.md`, `TODO/TODO-razim-ex5-terminal-capability-matrix-drift.md`, and `TODO/TODO-favun-ex5-architecture-terminal-surface-drift-2.md`
- read-only search and browse now exists over the shared `/api/search` projection
- read-only pending review now exists over draft-item and run-review slices from the shared search projections
- a small item approval action now exists over the existing item approval API
- CLI run, item, responsibility, place, and resource drilldown summaries now exist over the shared detail routes
- CLI pending-review and problem-review terminal summaries now exist over the shared review projections
- remaining terminal follow-ons are now split into explicit tracked gaps for review-queue contract strictness, run-context handoffs, architecture-doc drift, and Neovim inspect behavior coverage
- a small run approval action now exists over the existing run approval API
- a small item supersede action now exists over the existing item supersede API
- CLI evidence upload now exists over the shared run evidence route
- CLI typed-link creation now exists over the shared link route
- CLI structured/problem search now exists over the shared search route
- CLI grouped problem review now exists over the shared problem-review route
- CLI responsibility detail now exists over the shared responsibility route
- CLI pending-review aggregation now exists over the same shared search projections as `:OksPending`
- no concrete unfinished child slices remain under this umbrella; later Neovim follow-on work should file new focused TODOs instead of keeping TODO `016` open. Source: `DI-tuzok`.

## Result

TODO `016` now serves as the historical umbrella for the Neovim embodiment
wave that shipped across its child TODOs. It is complete and should no longer
appear as active backlog. Source: `DI-tuzok`.
