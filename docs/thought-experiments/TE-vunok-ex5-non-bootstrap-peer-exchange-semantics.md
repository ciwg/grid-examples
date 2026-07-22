# Ex5 Non-Bootstrap Peer Exchange Semantics

TE ID: `TE-vunok`
## Status
needs DF

## Decision under test

What `ex5` should implement as the first non-bootstrap PromiseGrid
peer-exchange step now that bootstrap import/export already ships.

Related TODO:

- `103` - `ex5-operational-knowledge-system/TODO/TODO-rumek-ex5-peer-exchange-beyond-bootstrap.md`

## Assumptions

- `ex5` already ships bootstrap export/import for signed `knowledge-item`,
  `knowledge-approval`, `knowledge-link`, and `knowledge-responsibility`
  artifacts plus their compatibility events.
- Bootstrap import currently requires an empty runtime.
- Compatibility replay still depends on `OperationalEvent.Sequence`, and the
  signed family envelopes also bind to that sequence today.
- `ex5` still mints local entity IDs such as `ITEM-0001`, `RESP-0001`, and
  `RUN-0001` from runtime-local counters rather than from a peer-stable global
  namespace.
- The runtime currently keeps one local signing identity under the data root.
- The user wants progress toward a stricter fully-on-grid shape, not merely a
  cosmetic extension of the bootstrap path.

## Alternatives

### Alternative A

Implement non-bootstrap import only for the same runtime lineage.

The importer would accept bundles into a non-empty runtime only when:

- the imported artifacts come from the same signer identity or lineage marker
- imported event sequences strictly extend the current local tail
- entity IDs continue the same already-existing local namespace

This is effectively append-only continuation or replication, not general
multi-origin merge.

### Alternative B

Allow arbitrary whole-family import into already-populated runtimes now.

The importer would accept signed artifacts from any peer into a non-empty
runtime, attempt dedupe by envelope CID or matching event payloads, and
preserve unresolved references explicitly.

This tries to make `ex5` multi-peer immediately without first redesigning the
compatibility event identity model.

### Alternative C

Defer non-bootstrap import until `ex5` first introduces peer-stable event or
artifact identity plus origin-aware ordering semantics.

Under this choice, `103` would not implement non-bootstrap import on the
current compatibility-event model. Instead, the next work would define the
identity and ordering layer needed for honest multi-origin import, then add
ongoing exchange on top of that.

## Scope and systems affected

- `ex5-operational-knowledge-system/TODO/TODO-rumek-ex5-peer-exchange-beyond-bootstrap.md`
- `ex5-operational-knowledge-system/TODO/TODO.md`
- `ex5-operational-knowledge-system/service/peer_exchange.go`
- `ex5-operational-knowledge-system/service/types.go`
- `ex5-operational-knowledge-system/service/app.go`
- `ex5-operational-knowledge-system/service/persistence.go`
- `ex5-operational-knowledge-system/service/app_test.go`
- `ex5-operational-knowledge-system/service/server_test.go`
- PromiseGrid implementation claims and peer-exchange staging docs

## Scenario analysis

### Scenario 1: same-lineage continuation after temporary separation

Alice runs one `ex5` node, exports a bundle, continues local work, then later
imports a bundle from the same node lineage into a restored or lagging peer.

Alternative A:

- fits this case directly
- can use the current signer identity plus strict sequence continuation as a
  narrow acceptance rule
- adds a useful replication-like path without solving general multi-peer merge

Alternative B:

- also accepts this case
- but pays for much broader ambiguity handling than the case actually needs

Alternative C:

- defers even this narrower continuation case
- keeps the repo from adding a limited same-lineage mode before the broader
  identity problem is solved

Result:

- A handles the narrow same-lineage case cleanly.

### Scenario 2: two independent peers create local work before exchange

Bob and Carol each start from empty runtimes. Bob creates `ITEM-0001`; Carol
also creates `ITEM-0001`. Both produce local event sequence `1`.

Alternative A:

- rejects import because the bundles do not share one lineage
- avoids pretending the current model can merge them

Alternative B:

- immediately collides on both event sequence semantics and local entity ID
  namespaces
- would need extra conflict rules that do not exist in the current model
- risks silently rewriting or misbinding durable history

Alternative C:

- defers implementation until the repo has a peer-stable identity and ordering
  layer
- keeps the runtime from claiming arbitrary multi-origin import too early

Result:

- B fails this scenario honestly on the current model; A and C avoid lying.

### Scenario 3: duplicate delivery and replay

Dave imports the same bundle twice or receives overlapping history segments
from the same peer.

Alternative A:

- can define replay in a narrow way: only accept strict same-lineage tail
  continuation and reject already-seen or non-contiguous sequences
- keeps duplicate handling simple

Alternative B:

- needs durable peer-stable duplicate detection that is stronger than the
  current local sequence model
- envelope CID alone is not enough when compatibility replay and entity IDs can
  still collide across peers

Alternative C:

- first creates the identity layer that duplicate detection actually needs
- postpones implementation until replay semantics can be defended

Result:

- A can support simple continuation replay; B still overreaches.

### Scenario 4: trust-boundary expansion

Ellen imports artifacts from Frank, whose runtime is validly signed but not
lineage-compatible with hers.

Alternative A:

- treats signer and lineage mismatch as a hard rejection
- keeps the first non-bootstrap step bounded to one logical runtime lineage

Alternative B:

- accepts cross-lineage artifacts before `ex5` can explain how local IDs,
  event ordering, and projections survive that import honestly

Alternative C:

- delays cross-lineage import until those semantics are actually designed

Result:

- A is bounded but narrow; C is broader and cleaner in the long run.

### Scenario 5: long-horizon migration toward fully on-grid behavior

Steve wants `ex5` eventually to behave as a real multi-peer PromiseGrid app,
not just a local app with nicer export/import.

Alternative A:

- gives immediate practical value
- but risks becoming a stepping stone whose semantics must later be
  superseded or tightly caveated

Alternative B:

- aims at the target shape too early
- but the current runtime does not yet have the identity model required to
  make that target honest

Alternative C:

- forces the next work to solve the real blocker first: peer-stable identity
  and ordering for imported durable history
- is slower in the short term, but cleaner for the intended end state

Result:

- C is the strongest long-horizon choice if the priority is fully on-grid
  semantics rather than incremental convenience.

## Conclusions

Rejected alternatives:

- Alternative B: the current model cannot honestly support arbitrary
  multi-origin non-bootstrap import because local event sequences and local
  entity IDs collide across peers.

Surviving alternatives:

- Alternative A: add a narrow same-lineage continuation import now
- Alternative C: first design peer-stable identity and ordering, then add
  non-bootstrap exchange

Unresolved question that still requires user choice:

- whether `ex5` should take the narrower same-lineage continuation step now,
  or avoid that stepping stone and solve the identity/order model first

Implications and future work:

- If Alternative A is chosen, `103` can implement one bounded non-bootstrap
  mode now, but it must be documented as lineage-only rather than general
  peer merge.
- If Alternative C is chosen, `103` should stay open and the next PromiseGrid
  TODO should define peer-stable event/artifact identity and ordering before
  non-bootstrap import lands.
- `104` and `105` remain valid follow-ons either way, but their final shape is
  cleaner if non-bootstrap exchange semantics are honest first.

## Decision status

Needs DF between Alternative A (same-lineage continuation import now) and
Alternative C (introduce peer-stable identity and ordering before any
non-bootstrap import implementation).
