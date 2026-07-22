# TODO zumor - ex5 browser staged operate flows

## Decision Intent Log

ID: DI-zumor
Date: 2026-07-22 22:48:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track the remaining browser operate heaviness as its own ex5 UI TODO.
Intent: Make the operate area feel less like a large transaction console by staging action choices ahead of the full forms.
Constraints: Preserve the generic run, evidence, and review forms and their manual overrides; improve the main path through staging, not capability removal.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-zumor-ex5-browser-staged-operate-flows.md`, `ex5-operational-knowledge-system/web/**`, `ex5-operational-knowledge-system/docs/**`, `ex5-operational-knowledge-system/README.md`

## Goal

Make browser operation work start from a small set of action choices before the
full transaction forms take over.

## Tasks

- [x] zumor.1 Review the operate workspace for where the full forms still appear too early.
- [x] zumor.2 Define the smallest staged action entry that reduces operate heaviness while preserving the forms.
- [x] zumor.3 Implement the staged operate flow without removing any current browser action path.
