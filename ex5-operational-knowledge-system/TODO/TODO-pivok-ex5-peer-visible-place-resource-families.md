# TODO pivok - make place and resource context peer-visible in ex5

## Goal

Make the remaining operational context references in `ex5` peer-visible and
durable so exchanged runs and links do not stop at unresolved `place` and
`resource` references.

## Why this exists

`ex5` now exchanges six signed families and resolves peer-visible entity
identity correctly, but runs and links still preserve `place` and `resource`
references only as unresolved context outside the current peer-visible slice.

## Tasks

- [ ] pivok.1 Run the required TE for whether `place` and `resource` should be
  separate frozen families or one combined context family.
- [ ] pivok.2 Lock the family boundary and naming.
- [ ] pivok.3 Implement the peer-visible frozen family or families.
- [ ] pivok.4 Extend peer exchange, CAS authority, docs, and claims for the
  new context coverage.

## Status

- open
- created from the post-109 PromiseGrid review

