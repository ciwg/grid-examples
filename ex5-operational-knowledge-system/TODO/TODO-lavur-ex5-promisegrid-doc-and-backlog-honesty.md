# TODO lavur - realign ex5 PromiseGrid docs and backlog to the shipped runtime

## Decision Intent Log

ID: DI-tivor
Date: 2026-07-22 10:47:03 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Realign the PromiseGrid-facing ex5 docs and backlog to the shipped runtime by removing stale pre-implementation staging language, explicitly distinguishing shipped bootstrap exchange and additive CAS from the still-open follow-on gaps, and reopening the TODO backlog around those remaining on-grid steps.
Intent: Keep the repo honest about what ex5 already ships versus what still blocks a stricter fully-on-grid claim, so the next migration wave starts from an accurate technical baseline.
Constraints: Do not claim non-bootstrap peer exchange, authoritative CAS read-path adoption, or peer-visible evidence portability before they exist; preserve the existing staged migration story where it is still true.
Affects: `ex5-operational-knowledge-system/TODO/TODO-lavur-ex5-promisegrid-doc-and-backlog-honesty.md`, `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `ex5-operational-knowledge-system/docs/promisegrid-peer-exchange-staging.md`, `ex5-operational-knowledge-system/docs/promisegrid-cas-staging.md`, `ex5-operational-knowledge-system/docs/promisegrid-embodiment-staging.md`, `ex5-operational-knowledge-system/docs/architecture.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`, `ex5-operational-knowledge-system/README.md`

## Goal

Make the PromiseGrid-facing docs and TODO backlog fully consistent with the
current shipped ex5 runtime and its remaining on-grid gaps.

## Why this exists

The current claims and staging docs still mix implemented state with older
staging language, and the TODO index no longer shows the remaining PromiseGrid
work as open follow-ons.

## Tasks

- [x] lavur.1 Fix doc contradictions around shipped peer exchange and CAS
  support.
- [x] lavur.2 Align the implementation-claims doc with the actual runtime
  behavior.
- [x] lavur.3 Keep the PromiseGrid staging docs consistent with the current
  migration state.
- [x] lavur.4 Keep the TODO backlog honest about what still remains to get ex5
  fully on the grid.

## Status

- done
- PromiseGrid-facing docs now distinguish shipped bootstrap exchange and
  additive CAS from the still-open fully-on-grid follow-on work
