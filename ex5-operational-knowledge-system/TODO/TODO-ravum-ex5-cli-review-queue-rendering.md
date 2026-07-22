# TODO ravum - ex5 CLI review queue rendering

## Decision Intent Log

ID: DI-ravum
Date: 2026-07-21 17:40:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Make `oks-cli pending-review` and `oks-cli problem-review` render terminal-first summaries over the existing shared review routes instead of falling back to raw JSON output.
Intent: Finish the next obvious shell-review gap after CLI drilldown polish so terminal users can triage review queues as easily as they inspect individual records.
Constraints: Reuse the existing `/api/search` and `/api/problem-review` projections, stay read-only, keep queue semantics aligned with Neovim and browser review logic, and link this slice back to deferred TODO `016`.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-ravum-ex5-cli-review-queue-rendering.md`, `ex5-operational-knowledge-system/TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md`, `ex5-operational-knowledge-system/cmd/oks-cli/main.go`, `ex5-operational-knowledge-system/cmd/oks-cli/main_test.go`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`

## Goal

Render the CLI review queues as shell-usable summaries instead of raw JSON.

## Links

- Parent follow-on: `TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md` (`016`, `nuvok.13`)

## Tasks

- [ ] ravum.1 Define the terminal summary shape for `pending-review` and `problem-review`.
- [ ] ravum.2 Add CLI renderers over the existing review projections.
- [ ] ravum.3 Add regression coverage for the new queue rendering.
- [ ] ravum.4 Update the ex5 terminal docs.
