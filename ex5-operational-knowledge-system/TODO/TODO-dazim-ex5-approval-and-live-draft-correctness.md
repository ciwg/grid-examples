# TODO dazim - ex5 approval and live-draft correctness

## Decision Intent Log

ID: DI-dazim
Date: 2026-07-21 12:05:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Split the review findings so approval correctness and live-draft correctness are tracked as one focused ex5 fix TODO.
Intent: Keep revision-truth and collaborative-draft truth from drifting apart in subtle ways.
Constraints: Focus on stale-revision approval behavior and empty-body draft updates; keep regression coverage in the same implementation pass.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-vurab-ex5-review-followups.md`, `ex5-operational-knowledge-system/TODO/TODO-dazim-ex5-approval-and-live-draft-correctness.md`

## Goal

Fix the workflow-correctness paths where the current system can misstate item
approval state or refuse a legitimate collaborative edit.

## Tasks

- [x] dazim.1 Make knowledge-item approval status changes revision-aware so approving an old revision cannot silently mark a newer draft as approved.
- [x] dazim.2 Allow live drafts to be intentionally cleared to an empty body and add regression coverage for browser and Neovim flows.

## Status

- done
- derived from the 2026-07-21 extensive ex5 review
