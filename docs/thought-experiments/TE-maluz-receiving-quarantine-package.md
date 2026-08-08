# Receiving exception and quarantine package

TE ID: TE-maluz

## Status

decided

## Decision under test

Whether Ex6 should represent an inbound receiving exception merely as a
receiving disposition, as a reusable quarantine domain package composed by a
receiving-exception workflow, or as a runtime-level rule.

## Assumptions

- The Ex6 runtime captures, activates, dispatches, and retains workflow
  artifacts, but does not own receiving-domain policy.
- A workflow artifact composes package capabilities for one operational
  outcome and remains separately versioned and removable.
- The current `receiving-check` workflow records a receipt and can record a
  `quarantine` receiving disposition, but that disposition has no independent
  lifecycle, owner, or evidence contract.
- Packages declare durable record families and protocol claims; the runtime
  mediates their persistence. Package declarations are bootstrap and routing
  hints, not automatic authority to perform work.
- Alice, Bob, and Carol are cooperative operators. Mallory may supply an
  untrusted imported artifact or misleading evidence, but cannot bypass the
  runtime's append-only record history.

## Alternatives

1. **Receiving-disposition workflow.** Add a `receiving-exception` workflow
   that records a failed receipt and a `quarantine` value in the existing
   receiving disposition family.
2. **Reusable quarantine package and workflow.** Add a quarantine package with its
   own durable case and lifecycle records. Add a receiving-exception workflow that
   composes `context`, `receiving`, `quarantine`, and `runs` to open the case
   after an exception.
3. **Runtime-enforced quarantine.** Teach the generic runtime to infer or block
   receiving and inventory work whenever an input is quarantined.

## Scenario analysis

### Normal receiving exception

Alice receives an inbound pallet, finds a seal mismatch, records inspection
evidence, and needs to prevent ordinary inventory handling until a review
resolves the problem.

Alternative 1 records that a receiving decision was `quarantine`, but leaves
the pallet's hold, accountable owner, evidence, and eventual release or
rejection implicit in a notes field or a later unrelated disposition.

Alternative 2 creates a specific quarantine case linked to the receiving item
and its evidence. The receiving-exception workflow can promise only that it opened
the case; a later release or rejection workflow can make its own explicit
promise. This preserves the distinction between observing a defect, placing a
hold, and deciding the final disposition.

Alternative 3 makes the runtime understand the operational meaning of a seal
mismatch. That couples the runtime to one domain and does not generalize cleanly
to other domain holds.

### Incomplete write or interrupted workflow

Bob's receiving-exception run records its input and begins, then the machine
stops after receipt evidence is retained but before the quarantine operation
completes.

Alternative 1 makes it difficult to distinguish "a failed receipt exists" from
"a quarantine hold was actually established." Operators may infer the latter
from an incomplete run.

Alternative 2 permits the runtime's existing durable run state and the package
records to show the exact boundary: evidence retained, case absent, case open,
or later terminal resolution. Retrying can use the retained input without
inventing a completed hold.

Alternative 3 has the same persistence problem while adding hidden runtime
state that must be rebuilt consistently with receiving records.

### Concurrent actors and mixed versions

Alice opens a quarantine case while Carol is reviewing the same receiving item.
Later, Bob receives a newer workflow version but retains the older
one for audit.

Alternative 1 can accumulate contradictory free-form dispositions without a
clear record relationship or lifecycle boundary. A newer workflow has no
stable quarantine protocol to inspect.

Alternative 2 lets each case and lifecycle event carry a package-defined
protocol and stable identifiers. Both workflow versions can remain retained;
each records exactly the pCID-defined contract it used. A future policy can
resolve conflicts explicitly instead of silently overwriting a disposition.

Alternative 3 centralizes conflict behavior in the runtime, forcing every
workflow version to inherit one runtime policy even when its domain rules
differ.

### Trust boundary

Mallory sends Alice a workflow artifact that claims it can release quarantined
material. Alice may retain and inspect it, but has not activated it or made any
acceptance/delivery promise for it.

Alternative 2 keeps the trust decision local and explicit: a case record is
only appended through an active, authorized route and workflow; the imported
artifact is merely a retained bootstrap hint until Alice activates it and the
needed route promises exist. The package does not turn a claim into authority.

Alternative 1 has fewer components but cannot express what a release decision
must prove. Alternative 3 risks treating a generic runtime rule as authority
over a domain decision that belongs to the participating agents.

### Long horizon and scale

At scale, multiple receiving sites hold material for damage, temperature,
documentation, and identity exceptions. Some cases are later released,
rejected, or handed to a corrective-action workflow.

Alternative 1 remains simple for a one-time exception but grows a large,
unqueryable set of receiving disposition strings and ad-hoc cross-links.

Alternative 2 creates an additional package and schemas, but supplies a
reusable protocol boundary for queries, evidence, handoffs, and later eggs.
It avoids duplicating the same lifecycle in receiving, inventory, and
corrective-action workflows.

Alternative 3 grows the generic runtime with every domain-specific hold and
resolution rule, increasing coupling and migration cost.

## Conclusions

Alternative 3 is rejected: quarantine is domain semantics, not runtime
semantics.

Alternative 1 is sufficient only if quarantine is permanently a single
receiving-only label with no lifecycle, evidence requirements, or reuse. That
does not match the requested reusable operational capability.

Alternative 2 survives and is recommended. It is the PromiseGrid-aligned
boundary: an independent package defines durable quarantine meaning; a
receiving-exception workflow composes it; the runtime remains generic; and each
transition is represented by explicit local evidence and promises rather than
an inferred global state.

## Output to decision framing

The architecture choice is narrowed to alternative 2. The following decisions
remain before implementation:

1. The quarantine lifecycle record shape: an explicit case plus append-only
   transition records, or one append-only event family with a rebuildable
   case projection.
2. The initial terminal transitions: release and reject only, or release,
   reject, and transfer to a named corrective-action owner.
3. The initial workflow's boundary: only open a quarantine case from a
   receiving exception, or also support a terminal decision in the same
   workflow.
4. Names, exact source paths, schemas, adapter contract, test locations, and
   dynamic runtime-path patterns.

The user's decisions are recorded in DI-luval, DI-giriz, and DI-hogid in the
relevant TODO Decision Intent Log.

## Implications for open TODOs and pending DIs

- `TODO-zibek` remains complete for its original receiving-package scope; it
  should not be rewritten as if it had always included quarantine.
- A new tracked work item is required for the quarantine package and
  receiving-exception workflow.
- The resulting implementation will need new DIs for the package boundary,
  lifecycle contract, workflow contract, names, source paths, and runtime test
  paths.

## Decision status

Locked by DI-luval, DI-giriz, and DI-hogid.
