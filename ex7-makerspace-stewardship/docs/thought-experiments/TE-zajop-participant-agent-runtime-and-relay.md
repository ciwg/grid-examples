# Participant-Agent Runtime, Signed Ingress, and Relay Carriage

TE ID: TE-zajop

## Status

decided

## Decision under test

How an Ex7 participant agent receives signed records, applies its own projection, and exchanges exact bytes without a central runtime.

## Alternatives

### A. Per-participant local agent plus dumb relay

Each participant runs a local Ex7 agent with its own framed record/blob store and local projection. Signed-record ingress accepts only exact canonical bytes, verifies envelope/family/author evidence, retains unknown or locally untrusted bytes, and projects only records matching that agent's policy. The participant's logged-in UI uses account bootstrap only to request a connection to its own agent. A relay exposes append/pull of exact frames/records and neither signs, validates semantics, computes state, nor decides trust.

### B. One shared application runtime

All users submit records to one service that stores history and computes shared current state.

### C. Semantic relay

Agents sign records, but the relay validates families and publishes the authoritative makerspace state.

## Scenario analysis

With A, Alice's agent can retain Alice's promise while offline and later push/pull through any relay. Carol's agent receives the same bytes and independently decides whether Carol/Alice keys and local makerspace policy make them projectable. A relay outage delays carriage, not authorship or local evidence. Unknown pCIDs remain bytes, not rejected semantics.

B reintroduces a central durable-history/state authority. C makes the relay a hidden authority and gives it semantic power over a community it does not own. Both contradict the no-central-server boundary.

For partitions and forks, A allows nodes to continue locally and reconcile records later. Conflicting evidence stays visible for local assessment rather than being overwritten. A can use multiple relays or direct exchange without changing record meaning.

## Conclusions

B and C are rejected. A is recommended: local participant agents with signed exact-byte ingress, account-only UI bootstrap, and non-authoritative relay carriage.

## Output to decision framing

**A is recommended.** Runtime code must split the current HTTP process into participant-agent APIs and a relay API, with separate storage/projection per agent and no relay semantic state.

## Decision status

locked: Alternative A by DI-sinov.
