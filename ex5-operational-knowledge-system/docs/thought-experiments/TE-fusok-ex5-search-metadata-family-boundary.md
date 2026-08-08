# Ex5 Search-Metadata Family Boundary

TE ID: `TE-fusok`
## Status
decided

## Decision under test

Whether `knowledge-search-metadata` should become its own frozen
PromiseGrid-native durable family in `ex5`, or whether search metadata should
instead remain a derived projection over the already-frozen durable families.

Related TODO:

- `099` - `ex5-operational-knowledge-system/TODO/TODO-lomuk-ex5-sixth-frozen-protocol-family.md`

## Assumptions

- `ex5` now has five frozen PromiseGrid families:
  `knowledge-item`, `knowledge-approval`, `knowledge-evidence`,
  `knowledge-link`, and `knowledge-responsibility`.
- The current search behavior is implemented as a projection over current
  runtime state, not as a separate append-only durable stream.
- Search already indexes fields that come from multiple families and multiple
  record types:
  - item titles, summaries, bodies, statuses, and responsibility IDs
  - run outcomes, notes, evidence facts, and approval notes
  - place/resource names and summaries
  - responsibility titles, summaries, and role keys
- Browser, CLI, and Neovim should remain on the current local HTTP adapter
  during this decision.

## Alternatives

### Alternative A

Freeze `knowledge-search-metadata` as a standalone durable family that stores
search-oriented latest-state metadata as its own signed append-only stream.

This would make search metadata explicitly durable and separately signed.

### Alternative B

Freeze a narrower `knowledge-search-metadata` family that stores only selected
search hints or indexing facts, while the rest of search remains derived from
the main families.

This would try to split the difference between full standalone search metadata
and fully derived search.

### Alternative C

Do **not** freeze `knowledge-search-metadata` as a separate durable family.
Instead, treat search metadata as a derived projection over the already-frozen
families and update `ex5` docs/backlog to reflect that decision.

## Scope and systems affected

- `protocols/knowledge-search-metadata.md`
- `docs/promisegrid-implementation-claims.md`
- `README.md`
- `docs/architecture.md`
- `docs/practical-implementation.md`
- `TODO/TODO-lomuk-ex5-sixth-frozen-protocol-family.md`
- `TODO/TODO.md`
- search implementation in `service/app.go`
- possible new runtime storage if a standalone family is chosen

## Scenario analysis

### Scenario 1: normal operator search

Alice searches for:
- a procedure title
- a receiving problem by evidence fact
- a responsibility by role
- a run by approval note

Alternative A:

- creates a separate durable search stream that must be kept aligned with all
  source families
- duplicates metadata already present in items, runs, responsibilities,
  places, and resources
- makes search read performance potentially more direct
- creates a strong obligation to prove that the search stream stays coherent
  across all upstream mutations

Alternative B:

- still creates a search-specific durable stream
- but now must define exactly which search hints are authoritative and which are
  still derived
- risks a hard-to-explain hybrid contract where some search behavior is
  PromiseGrid-native and some is merely projected

Alternative C:

- keeps search behavior derived from the already-frozen durable families
- avoids duplicate durable authority for titles, summaries, notes, facts, and
  context labels
- means search remains a projection concern rather than an artifact family

Result:

- C best matches the current search behavior and avoids duplicated durable
  authority.

### Scenario 2: replay verification and integrity

Bob restarts after many knowledge-item revisions, run records, evidence
additions, links, and responsibilities.

Alternative A:

- requires a second-order verification rule: not only must the source families
  replay correctly, the separately stored search metadata stream must match the
  derived search state
- creates two durable truths for the same operator-visible labels

Alternative B:

- still requires verifying some durable search-hint stream
- but adds ambiguity because only selected fields would be authoritative

Alternative C:

- keeps the integrity story simpler: verify the source families, then derive
  search metadata from them
- preserves one authoritative durable source for each fact

Result:

- C creates the fewest new verification obligations.

### Scenario 3: mixed-version nodes and migration

Carol runs an older ex5 node while Dave runs a newer one.

Alternative A:

- forces all nodes to agree on the exact durable search schema immediately
- creates version-skew risk around what counts as search metadata and when it is
  emitted

Alternative B:

- creates even more version-skew risk because old and new nodes may disagree on
  which search hints are authoritative

Alternative C:

- keeps mixed-version compatibility centered on the already-frozen families
- allows search improvements to remain projection-level behavior as long as they
  are derived honestly from the durable families

Result:

- C is safer for staged migration.

### Scenario 4: long-horizon evolution

Ellen wants `ex5` fully on-grid, but without freezing the wrong abstraction.

Alternative A:

- treats search as a first-class durable artifact family
- could make sense if search metadata is exchanged between peers or independently
  stored in CAS later
- but freezes a duplicated abstraction before that need is proven

Alternative B:

- risks an awkward partial family that later has to be replaced or expanded
- is the least conceptually clean option

Alternative C:

- says clearly that search metadata is derived operational projection state,
  not a separate durable family
- leaves room for later peer-visible search/index exchange if a future TE proves
  that need
- keeps the on-grid story centered on durable operational facts rather than
  derived retrieval views

Result:

- C preserves the clearest long-term architecture.

### Scenario 5: scale and storage cost

Frank expects bigger histories and heavier search usage.

Alternative A:

- adds another append-only log containing data that is largely duplicated from
  the source families
- increases write amplification because many upstream events would need search
  metadata refresh artifacts

Alternative B:

- adds another log plus ongoing ambiguity about what belongs there
- still duplicates some source-family data

Alternative C:

- keeps storage growth concentrated in the actual durable operational families
- allows indexes/caches to evolve without pretending they are separate durable
  families

Result:

- C is the lighter and cleaner operational choice.

## Conclusions

Rejected:

- Alternative A: a standalone durable search-metadata family would duplicate
  already-authoritative facts and create a second-order verification burden.
- Alternative B: a partial search-hints family would be even less coherent,
  because it would split search semantics between durable and derived layers.

Surviving:

- Alternative C: do not freeze `knowledge-search-metadata` as a separate durable
  family; keep search metadata derived from the already-frozen families and
  update docs/backlog accordingly.

## Implications for TODOs and pending DIs

- `099` may no longer be an implementation TODO if Alternative C is locked.
- If Alternative C is locked, the right follow-on is:
  - rewrite the `knowledge-search-metadata` protocol stub and claims/docs to say
    it is derived projection state, not a standalone durable family
  - close `099`
- If Alternative C is rejected, a new DI must define exactly what durable
  search-metadata artifacts are emitted and why they are not merely projections.

## Decision status

Alternative C locked by `DI-fusok`: do not freeze `knowledge-search-metadata`
as a separate durable family; keep search metadata derived from the
already-frozen families and update docs/backlog accordingly.
