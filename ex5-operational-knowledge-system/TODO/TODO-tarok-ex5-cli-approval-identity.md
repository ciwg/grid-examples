# TODO tarok - ex5 CLI approval identity

## Decision Intent Log

ID: DI-tarok
Date: 2026-07-21 12:05:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Split the review findings so CLI approval identity handling is tracked as its own ex5 fix TODO.
Intent: Preserve trustworthy approval history by making CLI-created approval records use the real actor identity instead of a hardcoded placeholder.
Constraints: Focus on CLI payload shape, docs, and regression coverage in the same implementation pass.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-vurab-ex5-review-followups.md`, `ex5-operational-knowledge-system/TODO/TODO-tarok-ex5-cli-approval-identity.md`

## Goal

Fix CLI approval commands so durable approval records preserve the actual
approver identity.

## Tasks

- [ ] tarok.1 Make CLI approval actor identity explicit instead of hardcoding `boss`, and cover the corrected payload shape in tests.

## Status

- open
- derived from the 2026-07-21 extensive ex5 review
