# TODO tasem - extend CAS authority beyond the six frozen families in ex5

## Goal

Move `ex5` closer to a uniform grid-native storage model by deciding how CAS
should become authoritative for the still-unfrozen runtime state instead of
remaining limited to the six frozen family envelopes.

## Why this exists

The current runtime now replays and exports the six frozen families
authoritatively from CAS, but other runtime state still depends on
compatibility event replay and local projections.

## Tasks

- [ ] tasem.1 Run the required TE for authoritative CAS adoption beyond the
  frozen families.
- [ ] tasem.2 Lock which remaining runtime state should become CAS-backed
  first, and how compatibility replay coexists during migration.
- [ ] tasem.3 Implement the chosen broader CAS authority step.
- [ ] tasem.4 Update docs, claims, and tests to reflect the new replay/read
  source-of-truth model.

## Status

- open
- created from the post-109 PromiseGrid review

