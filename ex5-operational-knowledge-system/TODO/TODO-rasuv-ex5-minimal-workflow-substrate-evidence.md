# TODO rasuv - ex5 minimal workflow substrate evidence

## Decision Intent Log

ID: DI-rasuv
Date: 2026-07-22 21:24:27 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track a future PromiseGrid alignment pass to determine whether any app-agnostic workflow/operation layer is actually justified beyond the current ex5-specific projections and flows.
Intent: Avoid either freezing all workflow semantics inside ex5 forever or extracting a speculative generic workflow substrate before the evidence is there.
Constraints: Prefer “no extraction” over speculative framework work; any candidate shared workflow surface must be narrower than ex5 review/search/approval composition unless reuse is explicitly proven.
Affects: `ex5-operational-knowledge-system/service/*`, `ex5-operational-knowledge-system/promisegrid/*`, `docs/thought-experiments/*`, `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Decide whether there is a real minimal PromiseGrid workflow substrate beyond
records, transport, and persistence, or whether those higher layers should stay
ex5-specific.

## Tasks

- [ ] rasuv.1 Audit the current typed operation families and review which ones are fundamentally ex5-specific versus plausibly reusable.
- [ ] rasuv.2 Reject or define the smallest justified shared workflow layer, with explicit proof requirements.
- [ ] rasuv.3 Align implementation claims so ex5 does not overstate workflow generality before a real substrate exists.

## Status

- open
- created from the remaining “app-agnostic workflow layers are still intentionally unextracted” PromiseGrid alignment gap
