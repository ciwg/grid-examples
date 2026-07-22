# TODO nobek - decide and stage CAS-backed envelope storage for ex5 PromiseGrid families

## Decision Intent Log

ID: DI-ribek
Date: 2026-07-22 10:23:43 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Introduce CAS-backed storage as an additive sidecar for signed family envelopes and copied evidence blobs, while keeping the current family logs and copied attachment tree as compatibility/manifests during the first migration pass.
Intent: Give ex5 a portable, content-addressable storage target for both envelopes and evidence blobs without risking a one-shot storage rewrite or forcing embodiment changes at the same time.
Constraints: Do not replace the current family logs or attachment tree in the first CAS pass; do not tighten embodiment contracts in the same slice; keep the first relay-visible exchange scope from TODO 100 separate from the storage migration decision.
Affects: `ex5-operational-knowledge-system/TODO/TODO-nobek-ex5-cas-envelope-storage.md`, `docs/thought-experiments/TE-nadok-ex5-cas-storage-migration-order.md`, `ex5-operational-knowledge-system/docs/promisegrid-cas-staging.md`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `ex5-operational-knowledge-system/docs/architecture.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Define when and how ex5 should move from local append-only family logs toward
CAS-backed envelope storage as part of the shipped PromiseGrid runtime.

## Why this exists

The current implementation persists signed family logs locally and stores
evidence bytes under the runtime root, but CAS-backed envelope/blob storage is
still outside the shipped ex5 operational workflow.

## Tasks

- [x] nobek.1 Run the required TE for CAS-backed storage scope and migration
  order.
- [x] nobek.2 Lock what stays in compatibility logs and what moves to CAS.
- [x] nobek.3 Define the first staged implementation slice.
- [x] nobek.4 Add storage-boundary docs and migration notes.

## Status

- done
- first CAS step is additive sidecar storage for signed envelopes and copied
  evidence blobs
- current family logs and copied attachment paths remain during migration
