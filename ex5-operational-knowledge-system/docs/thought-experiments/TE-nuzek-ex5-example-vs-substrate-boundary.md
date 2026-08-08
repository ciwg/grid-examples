# ex5 example versus PromiseGrid substrate boundary

TE ID: TE-nuzek
## Status
decided

## Decision under test

How the repo should describe the boundary between the reusable PromiseGrid
substrate now living under `promisegrid/*` and the remaining ex5-specific
application/runtime layer under `service/*`.

This TE corresponds to TODO `sulok.1` / `sulok.2` / `sulok.3`.

## Assumptions

- `ex5` remains a real example application, not the final generalized
  PromiseGrid product.
- The repo now has reusable substrate evidence beyond just records.
- The docs should state the boundary directly instead of implying either “only
  an example” or “already the final generalized runtime.”

## Alternatives

### A. Keep the current soft wording

Leave the boundary mostly implicit and let the code layout speak for itself.

### B. State the split explicitly

Document `promisegrid/*` as reusable substrate and `service/*` plus the shipped
embodiments as the ex5 application/runtime layer built on top of it.

### C. Rebrand ex5 as the generalized runtime now

Describe `ex5` itself as the reusable runtime/product boundary.

## Scenario analysis

### Scenario 1: a new reader opens the README

- A leaves the boundary interpretive and softer than the code now justifies.
- B tells the truth quickly: reusable substrate exists, but ex5 remains the
  operational-knowledge application.
- C overstates the generality of ex5-specific projections and workflows.

### Scenario 2: later substrate extraction continues

- A makes each further extraction feel ad hoc.
- B gives a stable narrative that future substrate slices can extend.
- C pressures the repo into framework claims the code still does not prove.

## Conclusions

Rejected:

- Alternative A: too implicit now that the substrate boundary is real.
- Alternative C: overstates ex5 into a generalized runtime/product.

Surviving:

- Alternative B: state the example-vs-substrate split explicitly.

## Implications for TODOs and pending DIs

- TODO `142` is locked to Alternative `B`.
- The README and PromiseGrid claims/docs should name `promisegrid/*` as the
  reusable substrate boundary and keep `service/*` / embodiments as ex5 app
  ownership.
