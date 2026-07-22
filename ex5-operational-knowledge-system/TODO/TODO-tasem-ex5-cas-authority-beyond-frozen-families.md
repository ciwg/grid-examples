# TODO tasem - extend CAS authority beyond the eight frozen families in ex5

## Decision Intent Log

ID: DI-zunep
Date: 2026-07-22 14:10:17 -0700
Status: active
Decision: Implement `111A`: shared live drafts become the next
CAS-authoritative local runtime state through a manifest-plus-CAS storage
model, while remaining outside the frozen peer-visible family set.
Intent: Move the last meaningful persisted local-only state onto stronger
content-addressed storage without falsely promoting drafts into a new
PromiseGrid family before their protocol meaning is frozen.
Constraints: Keep presence ephemeral, keep search metadata derived, preserve
the current shared-draft product semantics, and support backfill from older
`drafts/*.json` manifests during migration.
Affects: service/persistence.go; service/app.go; service/app_test.go;
service/types.go; docs/promisegrid-implementation-claims.md;
docs/promisegrid-cas-staging.md; README.md; docs/practical-implementation.md

## Goal

Move `ex5` closer to a uniform grid-native storage model by deciding how CAS
should become authoritative for the still-unfrozen runtime state instead of
remaining limited to the eight frozen family envelopes.

## Why this exists

The current runtime now replays and exports the eight frozen families
authoritatively from CAS, but other runtime state still depends on
compatibility event replay and local projections.

## Tasks

- [x] tasem.1 Run the required TE for authoritative CAS adoption beyond the
  frozen families.
- [x] tasem.2 Lock which remaining runtime state should become CAS-backed
  first, and how compatibility replay coexists during migration.
- [x] tasem.3 Implement the chosen broader CAS authority step.
- [x] tasem.4 Update docs, claims, and tests to reflect the new replay/read
  source-of-truth model.

## Status

- closed
- created from the post-109 PromiseGrid review
