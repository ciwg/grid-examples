# Persistence Evidence Repair

## Decision Intent Log

Temporary identifier notice: `tools/mint-handle` is unavailable in this checkout. This file and its DI identifiers must be renamed to minted proquint handles when the helper is restored.

### DI-pending-mint-ex7-001

- ID: DI-pending-mint-ex7-001
- Date: 2026-08-06 09:25:55
- Author: jj@thesalleys.com (JJ)
- Status: active
- Decision: Persist an append-only event before applying its associated in-memory state mutation, and require `fsync` plus close before the mutating request succeeds.
- Intent: A reported successful observation, hold clearance, loan, or return must survive restart; a failed append must not create phantom live state.
- Constraints: Keep the single-process local-demo scope. Do not add multi-writer coordination, authentication, or a new top-level PromiseGrid action kind.
- Affects: `service/store.go`, `service/app.go`, `service/app_test.go`, `service/server_test.go`.

### DI-pending-mint-ex7-002

- ID: DI-pending-mint-ex7-002
- Date: 2026-08-06 09:25:55
- Author: jj@thesalleys.com (JJ)
- Status: active
- Decision: Store and replay the complete area-policy version and text accepted for every off-site loan.
- Intent: Historical loan evidence must remain faithful when an area's current policy later changes, including for non-Woodworking tools.
- Constraints: Use the existing `PolicyVersion` and `Policy` vocabulary. The repair does not add policy authoring or a separate per-tool policy model.
- Affects: `service/types.go`, `service/store.go`, `service/app.go`, `service/app_test.go`.

### DI-pending-mint-ex7-003

- ID: DI-pending-mint-ex7-003
- Date: 2026-08-06 09:25:55
- Author: jj@thesalleys.com (JJ)
- Status: active
- Decision: Fail startup closed on malformed, incomplete, or oversized local evidence, preserving the evidence file untouched.
- Intent: The service must not conceal or continue from an unresolved gap in its local evidence history.
- Constraints: The reader's explicit maximum event size must accommodate every valid photo event accepted by HTTP validation while bounding allocation.
- Affects: `service/store.go`, `service/app_test.go`.

### DI-pending-mint-ex7-004

- ID: DI-pending-mint-ex7-004
- Date: 2026-08-07 20:04:58
- Author: jj@thesalleys.com (JJ)
- Status: active
- Decision: Replay a legacy loan event without an accepted-policy snapshot as explicitly incomplete evidence, retaining only its known borrower, deadline, and creation facts.
- Intent: Keep a currently loaned tool visible to the trusted volunteer group without falsely claiming that a later or inferred policy was accepted.
- Constraints: Never rewrite historical event bytes or substitute current area policy. New loan events continue to require complete snapshots.
- Affects: `service/types.go`, `service/store.go`, `service/app.go`, `service/app_test.go`, `web/app.js`.

### DI-pending-mint-ex7-005

- ID: DI-pending-mint-ex7-005
- Date: 2026-08-07 20:07:55
- Author: jj@thesalleys.com (JJ)
- Status: active
- Decision: Represent whether a loan has a recorded accepted-policy snapshot with `Loan.TermsComplete`.
- Intent: Make the known complete-versus-legacy-incomplete distinction explicit without introducing a broader status vocabulary.
- Constraints: `true` is required for new loans; `false` is reserved for replayed legacy loans whose event lacks a snapshot.
- Affects: `service/types.go`, `service/store.go`, `service/app.go`, `service/app_test.go`, `web/app.js`.

## Scope

Implements the decisions recorded in `../docs/thought-experiments/TE-pending-mint-ex7-persistence-evidence-repair.md`.
