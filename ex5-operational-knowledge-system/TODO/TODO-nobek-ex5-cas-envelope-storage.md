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

ID: DI-lavuz
Date: 2026-07-22 11:05:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Implement the first CAS step as dual-write sidecar storage for all signed family envelopes and copied evidence blobs, while preserving the current family logs, compatibility event log, and attachment paths as the active read model.
Intent: Make CAS real in the shipped runtime now so later peer-visible evidence exchange has a portable blob target, without forcing a read-path cutover or deleting compatibility storage.
Constraints: Keep startup replay on the current logs in this slice; preserve attachment paths in events and evidence views; expose the new blob identity explicitly so later peer-visible evidence work can bind to it.
Affects: `ex5-operational-knowledge-system/TODO/TODO-nobek-ex5-cas-envelope-storage.md`, `ex5-operational-knowledge-system/service/persistence.go`, `ex5-operational-knowledge-system/service/app.go`, `ex5-operational-knowledge-system/service/types.go`, `ex5-operational-knowledge-system/service/knowledge_evidence_envelopes.go`, `ex5-operational-knowledge-system/service/app_test.go`, `ex5-operational-knowledge-system/service/server_test.go`, `ex5-operational-knowledge-system/docs/promisegrid-cas-staging.md`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `ex5-operational-knowledge-system/docs/architecture.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`, `ex5-operational-knowledge-system/README.md`

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
- [x] nobek.5 Dual-write signed family envelopes into CAS objects by CID.
- [x] nobek.6 Dual-write copied evidence blobs into CAS objects by CID while
  keeping compatibility attachment paths.
- [x] nobek.7 Add CAS coverage and expose blob identity in the runtime view.

## Status

- done
- additive CAS sidecar storage now ships for signed envelopes and copied
  evidence blobs
- current family logs and copied attachment paths remain during migration
