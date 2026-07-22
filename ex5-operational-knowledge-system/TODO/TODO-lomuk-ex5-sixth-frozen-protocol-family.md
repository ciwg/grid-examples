# TODO lomuk - freeze and claim the sixth ex5 PromiseGrid protocol family

## Decision Intent Log

ID: DI-lomuk
Date: 2026-07-22 10:07:42 -0700
Author: jj@thesalleys.com (JJ)
Status: superseded
Decision: Leave `knowledge-search-metadata` open after the grouped `knowledge-link` and `knowledge-responsibility` batch, and require a dedicated TE before freezing any search-specific durable boundary.
Intent: Avoid freezing a synthetic search family before ex5 has a clean durable search-metadata event boundary to justify it.
Constraints: Do not force `knowledge-search-metadata` into the same grouped batch as links and responsibilities; keep the current search behavior projection-driven until a later TE resolves its durable scope.
Affects: `ex5-operational-knowledge-system/TODO/TODO-lomuk-ex5-sixth-frozen-protocol-family.md`, `docs/thought-experiments/TE-vusab-ex5-link-responsibility-search-family-order.md`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `ex5-operational-knowledge-system/README.md`

ID: DI-fusok
Date: 2026-07-22 10:18:21 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Do not freeze `knowledge-search-metadata` as a separate durable PromiseGrid family; treat search metadata as derived projection state over the already-frozen families and close TODO 099 on that basis.
Intent: Avoid duplicating durable authority for titles, summaries, tags, notes, facts, and context labels that are already carried by the operational families ex5 has frozen.
Constraints: No new signed search-metadata log is added in this slice; search remains projection-driven over current runtime state; docs and claims must stop describing search metadata as a pending sixth durable family.
Affects: `ex5-operational-knowledge-system/TODO/TODO-lomuk-ex5-sixth-frozen-protocol-family.md`, `docs/thought-experiments/TE-fusok-ex5-search-metadata-family-boundary.md`, `ex5-operational-knowledge-system/protocols/knowledge-search-metadata.md`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `ex5-operational-knowledge-system/docs/architecture.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/TODO/TODO.md`
Supersedes: DI-lomuk

## Goal

Resolve the PromiseGrid boundary for search metadata and make the shipped docs
honest about it.

## Why this exists

The earlier backlog still described search metadata as if it were the next
durable family. The TE for `099` resolved that the current ex5 search behavior
is derived projection state, not a separate durable family.

## Tasks

- [x] lomuk.1 Run the required TE for the search-metadata family boundary.
- [x] lomuk.2 Lock the family scope and implementation claim.
- [x] lomuk.3 Rewrite the protocol stub and claims/docs to describe search
  metadata as derived projection state.
- [x] lomuk.4 Close the backlog item without adding a sixth signed-envelope
  runtime slice.

## Status

- done
- search metadata remains derived projection state; no sixth signed family was
  added
