# TODO salup - ex5 CLI run, item, and responsibility drilldown polish

## Decision Intent Log

ID: DI-salup
Date: 2026-07-21 16:35:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Keep `oks-cli show-run`, `show-item`, and `show-responsibility` on the existing detail routes, but render them as operator-useful terminal drilldowns instead of raw JSON dumps.
Intent: Finish the next obvious terminal-consistency gap after place/resource drilldown polish so shell-first review flows do not fall back to raw JSON at the next handoff step.
Constraints: Reuse the existing detail routes, stay read-only, keep behavior aligned with the projected server shapes, and link this work back to deferred TODO `016`.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-salup-ex5-cli-run-item-responsibility-drilldown-polish.md`, `ex5-operational-knowledge-system/TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md`, `ex5-operational-knowledge-system/cmd/oks-cli/main.go`, `ex5-operational-knowledge-system/cmd/oks-cli/main_test.go`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`

## Goal

Bring the remaining CLI detail commands up to the same terminal drilldown
quality as place and resource detail.

## Links

- Parent follow-on: `TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md` (`016`, `nuvok.13`)

## Tasks

- [x] salup.1 Define the target terminal summary shape for run, item, and responsibility detail.
- [x] salup.2 Add CLI renderers over the existing detail routes.
- [x] salup.3 Add regression coverage for the new drilldown rendering.
- [x] salup.4 Update the ex5 terminal docs.
