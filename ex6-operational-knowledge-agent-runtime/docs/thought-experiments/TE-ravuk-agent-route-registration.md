# Agent Route Registration and Routing Promises

TE ID: TE-ravuk

## Status

needs DF

## Decision under test

Determine how OKAR should move from package-claim-derived route descriptions to the
agent/routing-role model: an app agent declares the pCIDs it promises to accept,
and a routing-role agent declares the conditions under which it promises to deliver
matching messages.

## Assumptions

- The current package manifest and implementation claims remain useful bootstrap
  metadata, but they are not the intended long-term representation of live agent
  availability.
- A pCID identifies a protocol specification, not one payload or one message type.
  Payload-defined message types remain application/protocol concerns.
- Top-level PromiseGrid semantics stay `promise`; registration, routing, parser
  forwarding, and conditions are pCID-defined payload semantics rather than new
  wire-level action kinds.
- A node can have more than one routing role, and agents may communicate directly
  without using a routing role.
- The first implementation is local and deterministic. It must leave a clear path
  for exported, signed, multi-node registration evidence without prematurely
  claiming network delivery exists.

## Alternatives

### A. Keep package claims as the live route registry

The runtime continues deriving routes only from installed package manifests and
implementation claims. Package activation is treated as equivalent to an agent
being available.

### B. Add one runtime-owned registration record per accepting agent

Each agent records its accepted pCIDs and conditions in a local runtime registry.
The planner treats registration as sufficient to select that agent; a separate
routing-role promise is implicit.

### C. Record paired acceptance and delivery promises in a runtime-owned registry

An app agent's acceptance promise and a routing-role agent's delivery promise are
separate, typed records. A planned routed path exists only where both records are
compatible. Package claims can seed bootstrap records, but live records become the
planner's source of truth.

## Scenario analysis

### Normal operation

Alice starts an inventory agent and promises to accept its declared inventory pCID.
Ron serves a routing role and promises to deliver that pCID to Alice, subject to
the same local trust and availability conditions. Under A, this relationship is
only inferred from Alice's package. Under B, Alice's availability is explicit but
Ron's responsibility remains invisible. Under C, the route plan can explain both
independent promises and their conditions.

For parser-first delivery, Carol's parser agent accepts an envelope pCID and
promises to emit a parsed downstream pCID; Alice accepts that downstream pCID.
Under C, Ron can promise delivery to Carol and Carol's forwarding promise is a
separate hop. This fits the existing direct/parser/transform route concepts while
making the responsible agents explicit.

### Failure, corruption, and incomplete writes

Alice can stop or fail after publishing acceptance. A must infer liveness from
package installation and therefore has no clear withdrawal or expiry model. B can
expire one registration but cannot distinguish an unavailable acceptor from a
routing service that no longer agrees to deliver. C allows each promise to carry
its own lifecycle and makes a partially recorded pair safely non-routable until
both records validate.

The registry must validate records independently, write atomically, and ignore an
incomplete or invalid record rather than producing a route from it.

### Concurrent actors and mixed-version nodes

Alice may renew acceptance while Ron updates delivery conditions. B risks one
record silently absorbing both concerns. C permits independent versioning and
deterministic compatibility checks. An older node can continue using package-seeded
routes while a newer node publishes explicit records, provided the bootstrap
translation is documented and the planner ranks explicit live evidence first.

### Long-horizon evolution and migration

A couples every future routing change to package-manifest evolution. B creates a
new local schema but later needs a second schema when routing roles become
distributed. C gives OKAR one durable conceptual model: promises made by agents,
stored locally first and exportable later. It creates obligations to specify record
identity, conditions, expiry, revocation, and conflict handling, but avoids making
package installation the permanent routing protocol.

### Trust-boundary changes

Mallory can advertise an acceptance promise for a valuable pCID. Under all
alternatives, trust policy must decide whether the statement is usable. C keeps
the acceptor and router identities distinct, so operators can trust Alice's
acceptance while refusing Ron's delivery promise, or trust one router but not
another. That separation aligns with the current trust/proof work and supports
multiple routing agents on one node.

### Scale and operational complexity

A has low immediate storage cost but forces agents to parse or receive package-wide
route declarations that do not express live ownership. B has a small local index
but hides router-side responsibility. C stores two small records per routed
relationship and requires compatibility indexing, but it reduces ambiguity,
permits targeted pCID lookup, and creates usable operational evidence. At scale,
expiry and compaction policies are required; they should be designed with the
registry rather than retrofitted after network exchange begins.

## Conclusions

Alternative A is rejected as the long-term design because it conflates installed
package capability with a running agent's promise and cannot represent routing-role
responsibility.

Alternative B improves liveness modelling but is rejected because it makes the
routing role implicit, which contradicts the architectural model that the router
is another agent with its own voluntary promise.

Alternative C survives and is recommended: use paired, explicitly identified
acceptance and delivery promises in a local runtime registry. The first slice
should be local, non-executing, and deterministic. Package claims should seed
documented bootstrap records only; they should not remain the permanent route
authority.

## Decisions still requiring DF

1. Should the first registry persist across runtime restarts, or begin in-memory
   with a persistence boundary designed but unimplemented?
2. Should package claims seed active bootstrap registrations automatically, or
   require an explicit operator command to import them?
3. Should the first record condition model support only identity, expiry, and
   enabled state, or also arbitrary pCID-defined condition payloads?
4. What names should the public types and CLI surface use for acceptance and
   delivery promises?

## Implications for open TODOs and pending DIs

- TODO puvok needs a new DI before implementation that locks the selected
  registration model, persistence boundary, bootstrap behavior, terminology, and
  touched paths.
- The existing claim-derived route table remains the compatibility input until the
  DF-selected bootstrap transition is implemented.
- Route execution, network exchange, and automatic trust remain out of scope for
  this first registration slice.

## Decision status

needs DF

## Refinements

None.
