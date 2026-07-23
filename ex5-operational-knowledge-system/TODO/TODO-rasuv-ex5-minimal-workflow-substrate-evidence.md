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

ID: DI-luvem
Date: 2026-07-22 21:44:16 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Lock `145A` and keep workflow composition inside `service/` for now; do not extract a PromiseGrid workflow substrate until a future slice can be named without ex5-specific entity or review vocabulary and proven reusable across more than one application shape.
Intent: Keep PromiseGrid substrate evidence-first. The repo now has real reuse proof for records, transport, and persistence, but the shipped create/review/search/approval orchestration still reads as ex5 application logic rather than app-agnostic substrate.
Constraints: Close `145` as an intentional no-extraction decision; align implementation claims so they do not imply workflow substrate extraction is still the next justified PromiseGrid step.
Affects: `../../docs/thought-experiments/TE-vunek-ex5-minimal-workflow-substrate-evidence.md`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `ex5-operational-knowledge-system/TODO/TODO.md`
Supersedes: DI-rasuv

## Goal

Decide whether there is a real minimal PromiseGrid workflow substrate beyond
records, transport, and persistence, or whether those higher layers should stay
ex5-specific.

## Tasks

- [x] rasuv.1 Audit the current typed operation families and review which ones are fundamentally ex5-specific versus plausibly reusable. See `../../docs/thought-experiments/TE-vunek-ex5-minimal-workflow-substrate-evidence.md`.
- [x] rasuv.2 Reject or define the smallest justified shared workflow layer, with explicit proof requirements. Locked to `145A`: no workflow substrate extraction is justified yet.
- [x] rasuv.3 Align implementation claims so ex5 does not overstate workflow generality before a real substrate exists.

## Status

- completed
- created from the remaining “app-agnostic workflow layers are still intentionally unextracted” PromiseGrid alignment gap
- TE complete: `TE-vunek` recommends `145A` as the most PromiseGrid-aligned surviving choice because current workflow composition remains ex5-specific.
- locked to `145A`: workflow composition remains ex5 application logic until a future slice can be named without ex5-specific workflow vocabulary and proven reusable beyond one app
