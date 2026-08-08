# Corrective-action review

TE ID: TE-nufip

## Status

decided

## Decision under test

Whether a rejected quarantine case should start a reusable corrective-action
domain package and review workflow, reuse maintenance findings, or remain only
a terminal quarantine event.

## Assumptions

- A quarantine rejection is a durable terminal event with explicit evidence.
- The runtime retains and dispatches workflow artifacts but does not own domain
  policy.
- Alice, Bob, and Carol are cooperative operators; Mallory may provide
  untrusted artifacts or evidence but cannot append local records without the
  required activation and route evidence.
- Corrective work may concern suppliers, procedures, documentation, training,
  or equipment; it is not limited to a physical resource repair.

## Alternatives

1. Add a reusable corrective-action package plus a review workflow linked to a
   rejected quarantine case.
2. Reuse maintenance findings and the maintenance-round workflow.
3. Keep rejection as the final record with no corrective-action domain record.

## Scenario analysis

### Normal rejected receipt

Alice rejects an inbound pallet because its identity documentation conflicts
with the delivery. Bob needs to assign a follow-up review and later retain its
outcome.

Alternative 1 records a separate corrective-action case linked to the rejected
quarantine case and its evidence. It can apply equally to supplier correction,
procedure correction, or training correction.

Alternative 2 incorrectly frames a documentation discrepancy as maintenance of
a resource. Alternative 3 retains the rejection but loses accountable follow-up
and completion evidence.

### Interrupted review and concurrency

Carol opens a corrective-action review while Bob attempts another review for
the same rejected case. The runtime stops before one run completes.

Alternative 1 can retain separate append-only action events and allow a future
policy to decide whether duplicate actions are accepted, merged, or rejected.
The workflow run distinguishes an incomplete attempt from a durable action.

Alternative 2 inherits unrelated maintenance semantics. Alternative 3 offers
no review run or action identity to inspect.

### Mixed versions and long horizon

A later policy adds stronger evidence for closure, while older action records
remain necessary for audit. Some actions outlive their original supplier or
procedure revision.

Alternative 1 versions a focused pCID-defined action contract while retaining
older artifacts and records. Alternative 2 couples corrective-action evolution
to maintenance. Alternative 3 has no contract to evolve.

### Trust boundary and scale

Mallory proposes an artifact that closes actions automatically. At many sites,
local reviewers accept different evidence sources and authorities.

Alternative 1 keeps every action and review locally activated and evidenced;
package declarations do not become automatic authority. It keeps future
closure promises separate from the initial review. Alternative 2 broadens the
maintenance trust surface. Alternative 3 leaves policy only in operator habit.

## Conclusions

Alternative 3 is rejected because a rejected case can require durable,
accountable follow-up beyond its terminal quarantine event.

Alternative 2 is rejected because maintenance findings are about resource
service and cannot represent supplier, procedure, or training correction.

Alternative 1 survives and is recommended: a reusable corrective-action
package should own its records, while the first corrective-action-review
workflow opens a linked action from a rejected quarantine case. Closure remains
a separate later workflow so review does not imply completion.

## Output to decision framing

Surviving design: a corrective-action package plus a
`corrective-action-review` workflow that opens, but does not close, an action.

Remaining decisions:

1. Action record model: append-only events with a rebuildable view, or a
   mutable action record plus transitions.
2. Linkage: require a rejected quarantine case ID, or allow the first workflow
   to open actions from any evidence subject.
3. Initial action fields and workflow input/output contract.
4. Names, exact source paths, adapter, schemas, test paths, and temporary-path
   pattern.

## Decision status

Locked by DI-hiboj.
