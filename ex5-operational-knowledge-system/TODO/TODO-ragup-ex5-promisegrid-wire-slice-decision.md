# TODO ragup - decide ex5 PromiseGrid wire slice after doc alignment

## Decision Intent Log

ID: DI-ragup
Date: 2026-07-21 00:00:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track the later ex5 decision about whether to open a real PromiseGrid implementation slice for frozen `pCID` handling and signed envelopes.
Intent: Keep the doc-honesty pass separate from the harder runtime decision so ex5 does not overclaim PromiseGrid completeness before the user explicitly chooses that next step.
Constraints: Do not implement transport or envelope changes in this TODO; this is a follow-on decision point after the docs and implementation-claims pass.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-ragup-ex5-promisegrid-wire-slice-decision.md`

## Goal

Decide whether ex5 should open a real PromiseGrid implementation slice around
frozen `pCID` handling, signed envelopes, and related wire/runtime boundaries.

## Tasks

- [ ] ragup.1 Re-read the PromiseGrid dev guide after the ex5 docs are aligned.
- [ ] ragup.2 Decide whether the next ex5 slice should stay local-runtime-only or begin real wire-level PromiseGrid work.
- [ ] ragup.3 If the answer is yes, define the first narrow implementation slice and its runtime/storage impact.

## Status

- open
- depends on completed doc-honesty and implementation-claims work
