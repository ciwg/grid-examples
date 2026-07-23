# ex5 staging abstraction cleanup

TE ID: TE-fulok
## Status
decided

## Decision under test

Which remaining rollout-era abstractions should be retired now that reusable
record and transport substrate types exist.

This TE corresponds to TODO `timav.1` / `timav.2` / `timav.3`.

## Assumptions

- Some `service/*` duplicate structs were useful while the substrate boundary
  was still being discovered.
- Once reusable substrate packages exist, keeping duplicate app-owned mirror
  structs becomes staging residue rather than useful architecture.
- Cleanup should reduce seams, not introduce another abstraction layer.

## Alternatives

### A. Keep the duplicate service-owned mirror types

Leave `service`-owned duplicates for events, signed record rows, and transport
wire types in place.

### B. Replace the mirror types with direct substrate aliases/use

Make `service/*` consume `promisegrid/records` and `promisegrid/transport`
types directly where the shapes are already identical.

### C. Rewrite broader runtime semantics at the same time

Use the cleanup as a trigger for a much broader refactor of workflows and
projections.

## Scenario analysis

### Scenario 1: maintaining record and transport code

- A keeps duplicate definitions that can drift.
- B removes duplicated shape ownership while keeping behavior stable.
- C adds more change surface than the cleanup requires.

### Scenario 2: future substrate extractions

- A leaves one more reason for `service/*` to feel like the substrate owner.
- B makes later extractions clearer because shared shapes already have one
  owner.
- C risks mixing substrate cleanup with unrelated architectural churn.

## Conclusions

Rejected:

- Alternative A: duplicate mirror types are no longer justified.
- Alternative C: too broad for a cleanup wave.

Surviving:

- Alternative B: replace the duplicate service-owned mirror types with direct
  substrate aliases/use.

## Implications for TODOs and pending DIs

- TODO `143` is locked to Alternative `B`.
- `service/*` should stop owning duplicate event/record/transport wire shapes
  where those shapes are already substrate-owned.
