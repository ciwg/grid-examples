# Legacy off-site loan evidence migration

**TE ID:** TE-pending-mint-ex7-legacy-loan-migration

## Status

decided

## Decision under test

How ex7 should handle pre-repair `loan` events that contain a borrower and deadline but do not contain the policy version and policy text accepted at checkout.

## Assumptions and scope

- Scope: ex7's local `events.jsonl` replay path and any user-visible indication of incomplete legacy loan evidence.
- Alice created a loan under the earlier ex7 implementation. Carol later changes the Woodworking policy. The old event cannot prove which terms Alice accepted.
- The original terms cannot be recovered from the event itself. The service must not describe any fallback terms as Alice's accepted policy.
- This is local-demo evidence migration only. It does not add authentication, remote replication, or a new PromiseGrid action kind.

## Alternatives

### A. Fail closed

Refuse startup when an old loan event lacks its policy snapshot.

### B. Replay with labelled incomplete evidence

Keep the loan active and record that its accepted policy terms are unavailable. Do not substitute the current policy; expose an explicit incomplete-evidence marker instead.

### C. Rewrite old events with the current area policy

Migrate old records by embedding today's policy as if it were the accepted snapshot.

## Scenario analysis

### Normal operation

Alice's legacy cordless-drill loan remains visible under B, so the makerspace can see that the tool is out. The record states that accepted terms are unavailable. A preserves strict evidence rules but hides operational state behind a startup failure. C makes the interface look complete while asserting a fact the event never recorded.

### Failure and incomplete evidence

If a legacy log is otherwise valid, B confines the evidence gap to the affected loan. A treats the gap as fatal to the entire service. C converts an acknowledged gap into false historical evidence. If the log is malformed beyond this known legacy schema, all alternatives still fail closed.

### Concurrent actors and mixed versions

A newer server can display a labelled legacy loan while an older server continues to read its old schema. B requires a stable marker in the replayed state rather than a rewritten event, so repeated restarts do not alter the historical bytes. A reduces mixed-version behavior to unavailability. C causes nodes with different current policies to invent different histories.

### Long-horizon evolution

B retains the known facts—borrower, due date, and creation time—while making the absent terms auditable forever. A requires manual repair before any long-lived evidence can be viewed. C permanently destroys the distinction between original and inferred terms.

### Trust boundary

Mallory cannot use B to claim that a later policy was accepted: the record must explicitly say it is incomplete. Neither A nor B authenticates legacy events; such cryptographic upgrades are outside this migration. C expands Mallory's opportunity to benefit from policy drift.

### Scale effects

B adds a small conditional projection at replay and display time. A has no migration code but can cause whole-service outages as legacy logs accumulate. C is mechanically simple but creates misleading history at every scale.

## Conclusions

Rejected: C falsifies accepted-term provenance. A remains viable for installations that prefer strict unavailability over incomplete operation.

Surviving alternatives: A and B. B is recommended because it preserves the known operational state without making an unsupported claim about accepted terms.

## Decision status

locked: replay legacy loan events with an explicit incomplete-evidence marker; do not substitute any current or inferred policy terms.

## Implications and future work

- New loan events continue to require complete policy snapshots.
- The app's model and browser must show incomplete legacy terms distinctly from a complete policy snapshot.
- Any migration that rewrites historical bytes requires a separate decision and is not authorized by this TE.
