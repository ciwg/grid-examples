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

ID: DI-mibor
Date: 2026-07-21 00:00:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Lock TODO 093 by choosing the staged PromiseGrid migration path, with `knowledge-item` as the first frozen family and the wire/runtime slice deferred until after that freeze and implementation claim exist.
Intent: Keep ex5 fully on the PromiseGrid path while still following the dev guide's spec-first rule. The first real runtime/wire work should now be narrow and tied to one frozen family rather than a broad rewrite.
Constraints: Keep websocket transport deferred; do not start a broad relay-visible rewrite in this TODO; move the active work to TODO 094.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-ragup-ex5-promisegrid-wire-slice-decision.md`, `ex5-operational-knowledge-system/TODO/TODO-mibor-ex5-first-frozen-protocol-family.md`, `docs/thought-experiments/TE-lafiz-ex5-promisegrid-wire-slice-decision.md`

## Goal

Decide whether ex5 should open a real PromiseGrid implementation slice around
frozen `pCID` handling, signed envelopes, and related wire/runtime boundaries.

## Tasks

- [x] ragup.1 Re-read the PromiseGrid dev guide after the ex5 docs are aligned.
- [x] ragup.2 Complete a thought experiment on whether the next ex5 slice should stay local-runtime-only or begin real wire-level PromiseGrid work.
- [x] ragup.3 Lock the final decision and then either close this TODO or promote the proposed follow-on TODO.

## Decision

Do not open the real wire/runtime PromiseGrid slice yet.

Instead, the next PromiseGrid-aligned ex5 step is:

1. choose one narrow protocol family
2. freeze that protocol contract
3. publish an implementation promise claim against it
4. only then open the first signed-envelope runtime slice

Thought experiment:

- `TE-lafiz` - [docs/thought-experiments/TE-lafiz-ex5-promisegrid-wire-slice-decision.md](/home/jj/lab/cswg/grid-examples/docs/thought-experiments/TE-lafiz-ex5-promisegrid-wire-slice-decision.md)

Follow-on:

- `094` - `TODO/TODO-mibor-ex5-first-frozen-protocol-family.md`

## Status

- done
- decision locked and promoted to TODO `094`
