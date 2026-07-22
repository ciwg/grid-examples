# TODO votek - freeze and claim the fourth ex5 PromiseGrid protocol family

## Decision Intent Log

ID: DI-votek
Date: 2026-07-22 10:07:42 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Freeze `knowledge-link` as the fourth ex5 PromiseGrid-native durable family and implement it in the same grouped migration batch as `knowledge-responsibility`, while keeping browser, CLI, and Neovim on the current HTTP adapter.
Intent: Continue the PromiseGrid migration with the next already-first-class durable event family instead of stalling on broader transport or search-boundary work.
Constraints: Keep the migration additive; use the existing `link_added` durable event boundary; do not change embodiment contracts in this slice; land replay verification, tamper detection, and claim/docs updates with the runtime work.
Affects: `ex5-operational-knowledge-system/TODO/TODO-votek-ex5-fourth-frozen-protocol-family.md`, `docs/thought-experiments/TE-vusab-ex5-link-responsibility-search-family-order.md`, `ex5-operational-knowledge-system/protocols/knowledge-link.md`, `ex5-operational-knowledge-system/protocols/**`, `ex5-operational-knowledge-system/service/**`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `ex5-operational-knowledge-system/docs/architecture.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/CHANGELOG.md`

## Goal

Freeze `knowledge-link` as the next ex5 PromiseGrid-native durable family and
add the fourth local signed-envelope runtime slice without changing the current
browser, CLI, or Neovim adapter contract.

## Why this exists

The current external PromiseGrid summary and repo claims now show
`knowledge-item`, `knowledge-approval`, and `knowledge-evidence` as frozen.
Typed links are the next durable family still modeled only through the local
event log and read model.

## Tasks

- [x] votek.1 Run the required TE for the `knowledge-link` family boundary.
- [x] votek.2 Lock the family scope and implementation claim.
- [x] votek.3 Freeze the protocol doc and add the fourth signed-envelope
  runtime slice.
- [x] votek.4 Extend replay verification, tests, and docs.

## Status

- done
- fourth frozen family and signed link runtime slice implemented in the grouped
  link/responsibility batch
