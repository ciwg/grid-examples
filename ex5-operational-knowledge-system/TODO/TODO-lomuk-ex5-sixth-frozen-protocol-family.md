# TODO lomuk - freeze and claim the sixth ex5 PromiseGrid protocol family

## Decision Intent Log

ID: DI-lomuk
Date: 2026-07-22 10:07:42 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Leave `knowledge-search-metadata` open after the grouped `knowledge-link` and `knowledge-responsibility` batch, and require a dedicated TE before freezing any search-specific durable boundary.
Intent: Avoid freezing a synthetic search family before ex5 has a clean durable search-metadata event boundary to justify it.
Constraints: Do not force `knowledge-search-metadata` into the same grouped batch as links and responsibilities; keep the current search behavior projection-driven until a later TE resolves its durable scope.
Affects: `ex5-operational-knowledge-system/TODO/TODO-lomuk-ex5-sixth-frozen-protocol-family.md`, `docs/thought-experiments/TE-vusab-ex5-link-responsibility-search-family-order.md`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `ex5-operational-knowledge-system/README.md`

## Goal

Freeze the remaining durable search-metadata layer as a PromiseGrid-native
family after the more operationally central link and responsibility families
are locked.

## Why this exists

The current PromiseGrid claims explicitly call out search-metadata as still
unfrozen. Search remains a projection/query capability over local runtime state
rather than a frozen family.

## Tasks

- [ ] lomuk.1 Run the required TE for the search-metadata family boundary.
- [ ] lomuk.2 Lock the family scope and implementation claim.
- [ ] lomuk.3 Freeze the protocol doc and add the signed-envelope runtime
  slice.
- [ ] lomuk.4 Extend replay verification, tests, and docs.

## Status

- open
- deferred until a dedicated search-boundary TE follows the grouped `knowledge-link` / `knowledge-responsibility` slice
