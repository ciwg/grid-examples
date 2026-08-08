# TODO komud — Reconcile global handle reservations

## Decision Intent Log

### DI-komud

- ID: DI-komud
- Date: 2026-08-07 21:24:56 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Backfill the root handle ledger with every current repository-owned proquint handle, identified by a coordination artifact filename or a declared coordination ID field.
- Intent: The root allocator must have a complete historical reservation set before it becomes the shared source for future exercises.
- Constraints: Preserve existing artifact text and paths; append ledger entries only; do not reserve bare citations, templates, or absent artifacts such as `TE-dabol` and `TE-vudaf`; preserve previously reserved unused handles.
- Affects: `TODO/handle-namespace.tsv`, `TODO/TODO.md`, and this TODO record.

## Tasks

- [x] komud.1 Append each current artifact-owner handle missing from the root ledger.
- [x] komud.2 Verify every current artifact-owner handle is represented exactly once in the root ledger.
- [x] komud.3 Commit and push the reconciliation.
