# TODO sarib - freeze and claim the fifth ex5 PromiseGrid protocol family

## Goal

Freeze `knowledge-responsibility` as the next ex5 PromiseGrid-native durable
family after `knowledge-link`, keeping the migration family-by-family and
adapter-preserving.

## Why this exists

Responsibilities are still part of the local runtime and projections, but they
are not yet frozen as a PromiseGrid-native family in the shipped runtime.

## Tasks

- [ ] sarib.1 Run the required TE for the `knowledge-responsibility` family
  boundary.
- [ ] sarib.2 Lock the family scope and implementation claim.
- [ ] sarib.3 Freeze the protocol doc and add the signed-envelope runtime
  slice.
- [ ] sarib.4 Extend replay verification, tests, and docs.

## Status

- open

