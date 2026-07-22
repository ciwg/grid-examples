# TODO davur - ex5 review queue approvals contract

## Decision Intent Log

ID: DI-davur
Date: 2026-07-21 19:05:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Make terminal review queues reject missing `approvals` fields instead of treating them as real unreviewed work.
Intent: Keep CLI and Neovim review triage honest when shared `/api/search` payloads drift, rather than silently inventing fake review work from omitted fields.
Constraints: Reuse the existing shared `/api/search` projections, keep CLI and Neovim semantics aligned, and link this fix back to deferred TODO `016`.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-davur-ex5-review-queue-approvals-contract.md`, `ex5-operational-knowledge-system/TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md`, `ex5-operational-knowledge-system/cmd/oks-cli/main.go`, `ex5-operational-knowledge-system/cmd/oks-cli/main_test.go`, `ex5-operational-knowledge-system/nvim/lua/oks/init.lua`, `ex5-operational-knowledge-system/nvim/pending_test.go`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`

## Goal

Reject malformed shared review payloads consistently instead of silently
reclassifying them as unreviewed work.

## Links

- Parent follow-on: `TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md` (`016`)

## Tasks

- [x] davur.1 Define the strict queue behavior for missing versus empty approvals.
- [x] davur.2 Apply the same contract checks to CLI and Neovim pending review.
- [x] davur.3 Add regression coverage for omitted approvals fields.
- [x] davur.4 Update terminal review docs.
