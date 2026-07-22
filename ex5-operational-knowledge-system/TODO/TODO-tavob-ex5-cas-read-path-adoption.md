# TODO tavob - adopt CAS-backed read and replay paths in ex5

## Decision Intent Log

ID: DI-rovud
Date: 2026-07-22 10:53:34 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Make CAS authoritative first for the five frozen family envelope bytes while keeping compatibility event replay active for still-unfrozen runtime state and allowing one-time manifest backfill for older runtimes.
Intent: Move ex5 toward a stricter PromiseGrid storage model now, without forcing an all-at-once cutover for places, resources, runs, and other still-unfrozen state.
Constraints: Family JSONL files may remain as manifests/indexes during migration; compatibility events remain active for unfrozen state; peer-visible evidence exchange remains separate follow-on work.
Affects: `ex5-operational-knowledge-system/TODO/TODO-tavob-ex5-cas-read-path-adoption.md`, `ex5-operational-knowledge-system/service/persistence.go`, `ex5-operational-knowledge-system/service/peer_exchange.go`, `ex5-operational-knowledge-system/service/app_test.go`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `ex5-operational-knowledge-system/docs/promisegrid-cas-staging.md`, `ex5-operational-knowledge-system/docs/architecture.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Promote CAS from additive sidecar storage to an authoritative read/replay
source for the signed PromiseGrid families and portable evidence blobs.

## Why this exists

ex5 currently dual-writes envelopes and evidence blobs into CAS, but startup
and replay still rely on the compatibility logs and local attachment paths as
the active source of truth.

## Tasks

- [x] tavob.1 Run the required TE for CAS read-path migration and recovery
  behavior.
- [x] tavob.2 Lock the authority order between CAS objects, family logs, and
  compatibility event logs.
- [x] tavob.3 Define how ex5 rebuilds runtime state when family logs are
  missing, partial, or inconsistent.
- [x] tavob.4 Implement at least one authoritative CAS-backed read/replay path.
- [x] tavob.5 Add corruption and recovery coverage for CAS-versus-log drift.

## Status

- done
- CAS now authoritatively rehydrates the five frozen family envelopes during
  replay and export
- older runtimes can backfill missing CAS envelope objects once from the
  manifest copy during migration
