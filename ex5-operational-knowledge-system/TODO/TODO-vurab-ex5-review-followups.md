# TODO vurab - ex5 review followups

## Decision Intent Log

ID: DI-vurab
Date: 2026-07-21 11:55:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Capture the current ex5 review findings as one actionable follow-up TODO covering durability, workflow correctness, and embodiment-consistency fixes.
Intent: Turn the deep review into concrete engineering backlog instead of leaving the findings only in chat text.
Constraints: This TODO records the backlog only; it does not itself change production behavior. Follow-up implementation slices should stay small and continue to add docs and tests in the same pass.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-vurab-ex5-review-followups.md`

## Goal

Track the concrete fixes exposed by the deep ex5 review so the next engineering
passes address durability risks and model/UX mismatches, not just add new
features.

## Tasks

- [x] vurab.1 Track durability and replay fixes under `TODO/TODO-busor-ex5-durability-and-replay-safety.md`.
- [x] vurab.2 Track approval and live-draft correctness fixes under `TODO/TODO-dazim-ex5-approval-and-live-draft-correctness.md`.
- [x] vurab.3 Track typed-link consistency fixes under `TODO/TODO-luzaf-ex5-link-model-consistency.md`.
- [ ] vurab.4 Track browser problem-drilldown alignment under `TODO/TODO-vemur-ex5-problem-drilldown-alignment.md`.
- [ ] vurab.5 Track CLI approval identity fixes under `TODO/TODO-tarok-ex5-cli-approval-identity.md`.

## Status

- open
- derived from the 2026-07-21 extensive ex5 review
- focused on correctness and durability, not feature expansion
- acts as an umbrella over the more specific review-followup TODOs
