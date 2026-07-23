# TODO sovek - ex5 generalized runtime substrate

## Decision Intent Log

ID: DI-sovek
Date: 2026-07-22 18:12:55 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track a future-scope pass to separate `ex5` as an example application from any broader reusable PromiseGrid runtime substrate.
Intent: Leave room for a more generalized runtime/product layer later instead of overloading `ex5` indefinitely as both example app and full substrate.
Constraints: This is the broadest and least safe scope expansion; keep it explicitly downstream of the smaller embodiment and contract refinements.
Affects: repo architecture, `ex5-operational-knowledge-system/*`, future shared runtime modules, `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Decide whether the repo should extract a more generalized PromiseGrid runtime
layer beyond `ex5`, and what boundary should remain example-specific.

## Tasks

- [ ] sovek.1 Define the candidate split between `ex5` app logic and reusable runtime substrate.
- [ ] sovek.2 Evaluate migration cost, abstraction risk, and repository impact.
- [ ] sovek.3 Decide whether to open a real extraction wave.

## Status

- open
- future-scope PromiseGrid refinement
