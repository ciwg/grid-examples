# Workflow Dispatch Route Evidence

TE ID: TE-niliv

## Status

decided

## Decision under test

Determine which explicit local route-promise evidence must exist before an
active workflow may dispatch its Docker-confined worker, and whether a denied
dispatch creates a durable workflow-run event.

## Assumptions and trust model

- `TE-dovek` remains locked: Docker is the only production worker backend and
  does not receive ambient host authority.
- `TE-ravuk` remains locked: installed package claims are bootstrap metadata,
  not live promises; enabled binding, receive, and delivery records are local
  evidence.
- The first dispatch slice is local, deterministic, non-networked, and does
  not claim that a router delivered a message across a peer boundary.
- Alice may publish all required local evidence for a workflow agent. Mallory
  may install or activate a package without publishing a receive promise, or
  may leave stale/conflicting route evidence in local CAS.
- Existing workflow-run events are append-only CAS lifecycle evidence after a
  run is accepted for dispatch.

## Alternatives

### A. Gate Docker dispatch on workflow activity and binding plus receive promise

The runtime requires the workflow artifact to be active, an enabled binding for
the worker package, and an enabled receive promise for the workflow input pCID.
Delivery promises remain planner-only evidence.

### B. Gate Docker dispatch on the complete planned-route evidence

The runtime requires active workflow state plus enabled binding, receive, and
delivery promises for the workflow input pCID before it invokes Docker.

### C. Keep dispatch independent of route-promise evidence

The runtime continues treating active artifact state and adapter declaration as
sufficient to invoke Docker; route promises remain introspection-only.

## Scenario analysis

### Normal local workflow start

Alice activates a workflow, binds its app label to the adapter package, publishes
a receive promise for its input pCID, and has a local routing role publish a
delivery promise. A and B permit dispatch; C also permits it. B makes the
execution gate match the plan the runtime presents to operators. A treats direct
local invocation as special even though the runtime has no separate direct-path
promise record.

### Missing or withdrawn promise

Mallory activates an installed workflow but no app agent currently promises to
receive the input pCID. A and B refuse before Docker starts. C still starts the
worker and therefore turns package activation into an implicit promise. If the
receive promise is disabled after a previous run, A and B refuse the next start.
If the delivery promise is disabled, only B also refuses, matching the planner's
non-executable result.

### Incomplete CAS evidence and restart

An interrupted write, a malformed route-promise record, or competing heads must
not make a route executable after replay. A and B rely on the existing registry
projection and fail closed. C ignores that evidence entirely. A failed
pre-dispatch eligibility check should return a local non-commitment error before
creating a workflow-run event, because no worker dispatch was accepted or
attempted; a Docker invocation that later fails remains a durable failed run.

### Concurrent updates and mixed versions

Alice may disable a receive promise while another actor starts a run. Both A and
B need one deterministic eligibility check immediately before worker invocation;
the slice does not claim distributed locking or peer agreement. Older nodes can
continue their historical behavior, but new nodes must not infer eligibility
from package claims. B adds one more required record but avoids a disagreement
between route-plan output and execution eligibility.

### Long-horizon evolution and trust boundaries

C preserves a dangerous split where the runtime says a route is not executable
but still executes the worker. A leaves delivery evidence as an operator-only
signal. B makes the local runtime's invocation a consequence of the same
voluntary evidence it uses for route selection, while still keeping sandboxing,
signatures, trust, and network delivery separate. Future direct-delivery or
signed-routing pCIDs can add different evidence without weakening this slice.

### Scale and operations

A performs two keyed local lookups; B performs three. Both are negligible
against Docker startup and preserve a rebuildable in-memory projection. C is
cheapest but creates operational surprises and weakens the evidence model.

## Conclusions

Alternative C is rejected: it contradicts the locked rule that package claims do
not create live promise authority.

Alternative A survives only if direct local dispatch is intentionally defined as
outside the delivery-promise model. That exception is not currently specified.

Alternative B is recommended: require active workflow state plus enabled
binding, receive, and delivery evidence immediately before Docker invocation.
On missing or invalid evidence, return a pre-dispatch local refusal without a
workflow-run event; after the runtime accepts dispatch, retain existing
append-only lifecycle events for success or failure.

## Decisions still requiring DF

1. Choose A or B for the dispatch eligibility gate. Recommendation: B.
2. Choose whether denied pre-dispatch eligibility returns without a run event
   (recommended) or creates a durable failed run event.
3. Approve the exact implementation and test paths after the selected gate is
   known.

## Implications for open TODOs and pending DIs

- `puvok` can close its remaining Docker-dispatch task only after the selected
  gate is recorded in a DI and tested across enabled, disabled, and replayed
  evidence.
- The slice must not claim network delivery, durable identity, or execution
  authority beyond the existing Docker confinement decision.

## Decision status

locked by DI-bidam.

## Refinements

None.
