# TODO sarib - freeze and claim the fifth ex5 PromiseGrid protocol family

## Decision Intent Log

ID: DI-sarib
Date: 2026-07-22 10:07:42 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Freeze `knowledge-responsibility` as the fifth ex5 PromiseGrid-native durable family and implement it in the same grouped migration batch as `knowledge-link`, while keeping the current browser, CLI, and Neovim adapter contract in place.
Intent: Land the next clean first-class durable record family immediately after links so the PromiseGrid migration keeps following the existing runtime’s real event boundaries.
Constraints: Keep the migration additive; use the existing `responsibility_created` durable event boundary; do not force search-metadata into the same slice; ship replay verification, tamper detection, and claim/docs updates with the runtime work.
Affects: `ex5-operational-knowledge-system/TODO/TODO-sarib-ex5-fifth-frozen-protocol-family.md`, `docs/thought-experiments/TE-vusab-ex5-link-responsibility-search-family-order.md`, `ex5-operational-knowledge-system/protocols/knowledge-responsibility.md`, `ex5-operational-knowledge-system/protocols/**`, `ex5-operational-knowledge-system/service/**`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `ex5-operational-knowledge-system/docs/architecture.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/CHANGELOG.md`

## Goal

Freeze `knowledge-responsibility` as the next ex5 PromiseGrid-native durable
family after `knowledge-link`, keeping the migration family-by-family and
adapter-preserving.

## Why this exists

Responsibilities are still part of the local runtime and projections, but they
are not yet frozen as a PromiseGrid-native family in the shipped runtime.

## Tasks

- [x] sarib.1 Run the required TE for the `knowledge-responsibility` family
  boundary.
- [x] sarib.2 Lock the family scope and implementation claim.
- [x] sarib.3 Freeze the protocol doc and add the signed-envelope runtime
  slice.
- [x] sarib.4 Extend replay verification, tests, and docs.

## Status

- done
- fifth frozen family and signed responsibility runtime slice implemented in
  the grouped link/responsibility batch
