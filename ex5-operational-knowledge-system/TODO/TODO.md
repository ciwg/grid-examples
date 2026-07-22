# TODO

## Open / Planned

- [x] 040 - Add ex5 Neovim run approval phase - `TODO/TODO-bafor-ex5-neovim-run-approval-phase.md`
- [x] 039 - Add ex5 Neovim item approval phase - `TODO/TODO-vamor-ex5-neovim-item-approval-phase.md`
- [x] 038 - Add ex5 Neovim pending review phase - `TODO/TODO-lorav-ex5-neovim-pending-review-phase.md`
- [x] 037 - Add ex5 Neovim search and browse phase - `TODO/TODO-givot-ex5-neovim-search-browse-phase.md`
- [x] 032 - Fix ex5 Neovim `:OksClose` session teardown behavior - `TODO/TODO-mabek-ex5-neovim-close-session-teardown.md`
- [x] 033 - Harden ex5 browser form and search error handling - `TODO/TODO-ruvot-ex5-browser-error-handling.md`
- [x] 034 - URL-encode ex5 CLI search queries correctly - `TODO/TODO-sifeg-ex5-cli-search-query-encoding.md`
- [x] 035 - Correct ex5 browser/CLI embodiment-equality overstatement in docs - `TODO/TODO-pobud-ex5-embodiment-doc-honesty.md`
- [x] 036 - Correct ex5 browser smoke coverage overstatement in docs/tests - `TODO/TODO-lanis-ex5-browser-smoke-doc-honesty.md`

## Completed

- [x] 031 - Clean up ex5 pre-existing errcheck debt - `TODO/TODO-zuvom-ex5-errcheck-cleanup.md`
- [x] 030 - Expand ex5 search so evidence and approval history is actually searchable - `TODO/TODO-farun-ex5-search-evidence-and-approval-history.md`
- [x] 029 - Harden ex5 browser participant identity startup in restricted storage environments - `TODO/TODO-mitob-ex5-browser-participant-identity-hardening.md`
- [x] 028 - Fix ex5 Neovim live-draft cursor reporting against the wrong window - `TODO/TODO-pazud-ex5-neovim-cursor-window-correctness.md`
- [x] 027 - Enforce ex5 evidence attachment size limits - `TODO/TODO-navos-ex5-evidence-attachment-size-enforcement.md`
- [x] 025 - Align ex5 browser problem drilldowns with review logic - `TODO/TODO-vemur-ex5-problem-drilldown-alignment.md`
- [x] 021 - Fix ex5 review findings across durability, workflow correctness, and embodiment consistency - `TODO/TODO-vurab-ex5-review-followups.md`
- [x] 026 - Fix ex5 CLI approval identity handling - `TODO/TODO-tarok-ex5-cli-approval-identity.md`
- [x] 024 - Fix ex5 typed-link model consistency - `TODO/TODO-luzaf-ex5-link-model-consistency.md`
- [x] 023 - Fix ex5 approval and live-draft correctness - `TODO/TODO-dazim-ex5-approval-and-live-draft-correctness.md`
- [x] 022 - Fix ex5 durability and replay safety - `TODO/TODO-busor-ex5-durability-and-replay-safety.md`
- [x] 020 - Add ex5 Neovim typed-link browsing phase - `TODO/TODO-zalor-ex5-neovim-typed-link-browsing-phase.md`
- [x] 019 - Add ex5 Neovim run inspector phase - `TODO/TODO-ravok-ex5-neovim-run-inspector-phase.md`
- [x] 018 - Add ex5 Neovim item inspector phase - `TODO/TODO-lonuk-ex5-neovim-item-inspector-phase.md`
- [x] 017 - Add ex5 Neovim live draft phase 1 - `TODO/TODO-fudok-ex5-neovim-live-draft-phase1.md`
- [x] 007 - Finish ex5 inventory follow-ons in the operational-memory lane - `TODO/TODO-ranor-ex5-inventory-and-embodiment-followups.md`
- [x] 015 - Add ex5 grouped problem review - `TODO/TODO-pogul-ex5-grouped-problem-review.md`
- [x] 014 - Add ex5 history drilldown filters - `TODO/TODO-vafuk-ex5-history-drilldown-filters.md`
- [x] 013 - Add ex5 context review facts - `TODO/TODO-zemok-ex5-context-review-facts.md`
- [x] 012 - Add ex5 receiving check workflow review - `TODO/TODO-vemok-ex5-receiving-check-review.md`
- [x] 011 - Repair ex5 decision-first and path-intake records - `TODO/TODO-talub-ex5-compliance-repair-and-feature-intake.md`
- [x] 010 - Add ex5 inventory discrepancy review - `TODO/TODO-pojul-ex5-inventory-discrepancy-review.md`
- [x] 008 - Add ex5 item run history drilldown - `TODO/TODO-hozom-ex5-item-run-history.md`
- [x] 006 - Expand ex5 workflow search and browser automation - `TODO/TODO-honus-ex5-workflow-search-and-automation.md`
- [x] 004 - Resolve ex5 open product decisions - `TODO/TODO-solaj-ex5-open-product-decisions.md`
- [x] 009 - Add ex5 context run history drilldown - `TODO/TODO-julos-ex5-context-run-history.md`
- [x] 003 - Polish ex5 browser navigation and docs - `TODO/TODO-savuk-browser-navigation-and-docs.md`
- [x] 002 - Build ex5 live operational workflow slice - `TODO/TODO-foluk-operational-knowledge-system-live-workflow.md`
- [x] 001 - Build ex5 operational knowledge system foundation - `TODO/TODO-radok-operational-knowledge-system-foundation.md`

## Not Implemented / Deferred

- [ ] 005 - Port ex5 to websocket collaboration transport - `TODO/TODO-masad-ex5-websocket-collaboration-transport.md`
  Deferred by `DI-tabiv`: do not port the full `ex3` websocket model into `ex5` in the current phase.
- [ ] 016 - Plan ex5 Neovim embodiment follow-on beyond phase 1 - `TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md`
  Deferred by `DI-tabiv`, `DI-nuvok`, and `DI-fudok`: a thin Neovim live-draft phase now exists, but richer Neovim workflow work remains a separate follow-on project rather than part of inventory TODO `007`.

## Locked decisions affecting implementation

- `DI-tabiv`: keep the current local HTTP live-draft model; collaborative editing is optional; a future Neovim embodiment is desirable.
- `DI-fudok`: implement the first Neovim phase as a thin launcher/plugin over the existing live-draft HTTP API, not as a websocket sidecar.
- `DI-lonuk`: add a read-only Neovim inspector over existing item detail APIs before attempting richer in-editor workflow actions.
- `DI-ravok`: add a direct read-only Neovim run inspector over the existing run detail API before attempting write-side workflow actions in the editor.
- `DI-zalor`: add read-only typed-link browsing in Neovim over the existing entity detail APIs before attempting in-editor mutation of links or approvals.
- `DI-givot`: add read-only Neovim search and browse over the existing `/api/search` projection before attempting write-side review or approval actions in the editor.
- `DI-lorav`: add a read-only Neovim pending-review view over the existing search projections before attempting write-side approval actions in the editor.
- `DI-vamor`: add a small Neovim item approval action that resolves the current revision through the existing item API before posting to the existing approval endpoint.
- `DI-bafor`: add a small Neovim run approval action that reuses the existing run approval endpoint and refreshes the relevant terminal view afterward.
- `DI-vurab`: track the 2026-07-21 deep ex5 review findings as an explicit fix backlog covering attachment durability, event replay limits, revision-aware approvals, empty-body drafts, link validation, problem drilldown alignment, and CLI approval identity.
- `DI-busor`: handle durability and replay hazards before adding more workflow surface area.
- `DI-dazim`: fix revision and empty-draft correctness so approvals and live collaboration cannot silently misstate current state.
- `DI-luzaf`: make typed links structurally trustworthy and consistent across docs, browser, and Neovim.
- `DI-vemur`: make browser “problem” drilldowns actually match the grouped review logic they claim to represent.
- `DI-tarok`: preserve real approver identity in CLI-created approval records.
- `DI-talub`: recent `ex5` feature slices are retroactively documented for decision/path traceability, and every future `ex5` feature slice must start by recording its intended repo paths and runtime path patterns in a local TODO before code changes begin.
