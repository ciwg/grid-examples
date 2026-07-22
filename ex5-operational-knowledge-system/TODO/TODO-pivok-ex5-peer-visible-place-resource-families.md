# TODO pivok - make place and resource context peer-visible in ex5

## Decision Intent Log

ID: DI-pivul
Date: 2026-07-22 13:08:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Lock TODO `110` to Alternative A. `place` and `resource` become separate frozen PromiseGrid families named `operational-place` and `operational-resource`.
Intent: Finish the remaining peer-visible operational context as first-class signed families so runs and links stop carrying unresolved place/resource references outside the durable exchange set.
Constraints: Keep the current hierarchical place model and current resource-to-place reference model. Reuse the canonical create-envelope CID identity rule from `DI-loruk`. Preserve the current HTTP adapter and compatibility projections while extending CAS-authoritative replay and peer exchange.
Affects: `docs/thought-experiments/TE-puvok-ex5-place-resource-family-boundary.md`, `ex5-operational-knowledge-system/protocols/operational-place.md`, `ex5-operational-knowledge-system/protocols/operational-resource.md`, `ex5-operational-knowledge-system/protocols/profiles.go`, `ex5-operational-knowledge-system/service/app.go`, `ex5-operational-knowledge-system/service/persistence.go`, `ex5-operational-knowledge-system/service/types.go`, `ex5-operational-knowledge-system/service/peer_exchange.go`, `ex5-operational-knowledge-system/service/app_test.go`, `ex5-operational-knowledge-system/service/server_test.go`, and PromiseGrid claims/docs for TODO `110`

## Goal

Make the remaining operational context references in `ex5` peer-visible and
durable so exchanged runs and links do not stop at unresolved `place` and
`resource` references.

## Why this exists

Before this TODO landed, `ex5` exchanged six signed families and resolved
peer-visible entity identity correctly, but runs and links still preserved
`place` and `resource` references only as unresolved context outside the
peer-visible slice.

## Tasks

- [x] pivok.1 Run the required TE for whether `place` and `resource` should be
  separate frozen families or one combined context family.
- [x] pivok.2 Lock the family boundary and naming.
- [x] pivok.3 Implement the peer-visible frozen family or families.
- [x] pivok.4 Extend peer exchange, CAS authority, docs, and claims for the
  new context coverage.

## Status

- closed
- created from the post-109 PromiseGrid review
- `TE-puvok` completed; `DI-pivul` locks separate `operational-place` and
  `operational-resource` families
- the seventh and eighth frozen families now ship, and peer exchange carries
  place/resource context as first-class durable records
