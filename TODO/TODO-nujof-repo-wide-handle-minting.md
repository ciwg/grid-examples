# TODO nujof — Make handle minting repo-wide

## Decision Intent Log

### DI-nujof

- ID: DI-nujof
- Date: 2026-08-07 20:52:12 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Move the proven ex6 `mint-handle` utility and its append-only reservation ledger to the `grid-examples` repository root, making them the shared allocator for all current and future exercises.
- Intent: Future exercises must not depend on ex6-local coordination state; one repository-owned tool and ledger preserve collision-safe allocations across the exercise corpus.
- Constraints: Preserve every existing reservation; retain file history with `git mv`; do not alter the dev-guide tool; make this an isolated commit and push; update the allocator to scan the whole repository.
- Affects: `tools/mint-handle`, `tools/README.md`, `TODO/handle-namespace.tsv`, `TODO/TODO.md`, and this TODO record.

## Tasks

- [x] nujof.1 Move the allocator, documentation, and ledger to repository-root ownership.
- [x] nujof.2 Verify repository-root scanning and non-mutating command validation without creating an extra reservation.
- [ ] nujof.3 Commit and push the isolated migration.
