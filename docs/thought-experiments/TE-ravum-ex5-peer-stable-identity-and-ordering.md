# Ex5 Peer-Stable Identity And Ordering

TE ID: `TE-ravum`
## Status
decided

## Decision under test

What first peer-stable identity and ordering layer `ex5` should introduce so
non-bootstrap multi-origin peer exchange can be implemented honestly.

Related TODOs:

- `107` - `ex5-operational-knowledge-system/TODO/TODO-navud-ex5-peer-stable-identity-and-ordering.md`
- `103` - `ex5-operational-knowledge-system/TODO/TODO-rumek-ex5-peer-exchange-beyond-bootstrap.md`

## Assumptions

- `ex5` already ships six frozen signed families:
  `knowledge-item`, `knowledge-approval`, `knowledge-evidence`,
  `operational-run`, `knowledge-link`, and `knowledge-responsibility`.
- Bootstrap peer exchange already exports and imports whole-family signed
  records plus compatibility events, but only into an empty runtime.
- The current compatibility event stream still uses runtime-local
  `OperationalEvent.Sequence`.
- Runtime-local entity IDs such as `ITEM-*`, `RUN-*`, `RESP-*`, and `LINK-*`
  are still minted from local counters.
- Signed family records already have stable content identities in practice via
  their envelope bytes and envelope CIDs, but compatibility replay and
  projections still rely on local event sequence order.
- The user wants the strongest practical PromiseGrid direction, not a
  convenience-only stepping stone.

## Threat and trust model

- Alice and Bob are honest peers exchanging valid signed bundles.
- Carol is an honest peer who may run a different history depth or import the
  same artifacts more than once.
- Dave runs a mixed-version node during migration.
- Mallory may tamper with transport payloads, replay previously seen bundles,
  or try to confuse ordering by resubmitting artifacts in a different order.
- The runtime should not silently merge unrelated local histories by assuming
  one shared local sequence or one shared local ID namespace.

## Alternatives

### Alternative A

Use signed envelope CID as the only new durable identity and derive ordering
from existing timestamps plus current local append order.

Under this model:

- the envelope CID identifies whether a signed record is already known
- duplicate detection uses envelope CID only
- ordering across peers uses event timestamps, then local import order as a tie
  breaker
- compatibility event sequence remains local and is not promoted into a
  peer-stable origin-aware tuple

### Alternative B

Introduce explicit peer-stable origin identity plus per-origin order, and bind
compatibility replay to that tuple while keeping envelope CID as content
    identity.

Under this model:

- each emitted operational event and signed family record carries:
  - `origin_peer_id`
  - `origin_sequence`
- the pair `(origin_peer_id, origin_sequence)` is the first-class durable
  artifact identity for replay, dedupe, and lineage-aware ordering
- envelope CID remains the stable content identity and tamper detector
- local `Sequence` can remain as a compatibility projection, but it is no
  longer the canonical cross-peer identity

### Alternative C

Introduce a separate globally unique event ID for each compatibility event and
order imported history primarily by wall-clock timestamp plus that event ID,
without introducing per-origin monotonic order.

Under this model:

- each event gets one UUID-like durable identifier
- duplicate detection uses that event ID
- imported ordering uses timestamps first, with event ID only as a tie breaker
- per-origin append order is not explicitly preserved

## Scope and systems affected

- `docs/thought-experiments/TE-ravum-ex5-peer-stable-identity-and-ordering.md`
- `ex5-operational-knowledge-system/TODO/TODO-navud-ex5-peer-stable-identity-and-ordering.md`
- `ex5-operational-knowledge-system/TODO/TODO-rumek-ex5-peer-exchange-beyond-bootstrap.md`
- `ex5-operational-knowledge-system/service/types.go`
- `ex5-operational-knowledge-system/service/app.go`
- `ex5-operational-knowledge-system/service/persistence.go`
- `ex5-operational-knowledge-system/service/peer_exchange.go`
- signed-family helper files if event/record metadata changes
- service replay/import tests
- PromiseGrid claims and peer-exchange staging docs

## Scenario analysis

### Scenario 1: same artifact delivered twice from one peer

Alice exports a bundle to Bob, then Bob receives the same bundle again later.

Alternative A:

- envelope CID dedupe works for signed family records
- but compatibility replay still lacks one peer-stable event identity
- if replay still depends on local `Sequence`, Bob must either reject whole
  imports based on local sequence mismatch or invent secondary heuristics

Alternative B:

- duplicate detection works at two levels:
  - same envelope CID means same artifact bytes
  - same `(origin_peer_id, origin_sequence)` means same logical event slot
- Bob can reject or ignore already-seen logical history cleanly

Alternative C:

- duplicate detection works only if the event ID was already preserved
- ordering still depends too heavily on timestamps, which do not express
  origin-local append order

Result:

- B handles repeated delivery most cleanly.

### Scenario 2: two peers create valid local history before exchange

Bob and Carol each create items, approvals, and runs independently, then
exchange state later.

Alternative A:

- envelope CID only says whether two records have identical bytes
- it does not explain how two distinct valid histories interleave
- timestamp ordering is vulnerable to clock skew and tie ambiguity

Alternative B:

- preserves each peer's local append order explicitly
- gives the importer a stable ordering key across mixed origins
- makes it possible to talk about "Carol sequence 8 arrived after Bob sequence
  13" without pretending they are one sequence

Alternative C:

- unique IDs avoid pure duplication, but still do not model one peer's local
  monotonic history well enough
- timestamp-first ordering becomes the de facto arbiter, which is weak under
  skew or batched imports

Result:

- B is the strongest multi-origin model.

### Scenario 3: mixed-version migration

Dave upgrades one node while Ellen still runs an older build.

Alternative A:

- minimizes structural changes
- but keeps the older local-sequence assumption alive, which is exactly the
  blocker to non-bootstrap import

Alternative B:

- introduces one new durable tuple that can be projected back into older local
  sequence views
- gives the migration a clear canonical model and a compatibility layer

Alternative C:

- adds a new event ID but not a clear per-origin ordering story
- likely still needs later augmentation once real non-bootstrap merge exists

Result:

- B gives one upgrade path instead of two partial ones.

### Scenario 4: tamper, replay, and adversarial reorder

Mallory resends a valid old bundle after Alice and Bob have exchanged newer
history, or reorders records inside a delivery.

Alternative A:

- envelope CID catches byte tamper
- but replay and reorder handling still lean on timestamps or local append
  heuristics
- that makes "already known versus newly ordered" harder to defend

Alternative B:

- envelope CID still catches byte tamper
- origin tuple makes replay and ordering checks explicit
- the importer can sort by `(origin_peer_id, origin_sequence)` within each
  origin and keep imported histories stable across repeated deliveries

Alternative C:

- unique IDs help replay detection
- but timestamp-first ordering still gives Mallory more room to perturb import
  order semantically without altering content

Result:

- B is strongest under trust-boundary stress.

### Scenario 5: long-horizon PromiseGrid evolution

Steve wants `ex5` to stop relying on local-only assumptions and move toward a
real multi-peer grid node.

Alternative A:

- is easy to bolt on
- but would still leave the repo needing a second, more fundamental identity
  migration later

Alternative B:

- introduces exactly the origin-aware layer that later non-bootstrap exchange,
  richer peer sync, and stronger embodiment claims need
- can coexist with envelope CID and CAS without changing family semantics

Alternative C:

- adds one more global ID surface now
- but still leaves per-origin ordering underspecified, so it is likely another
  stepping stone

Result:

- B best matches the "fully on the grid" direction.

### Scenario 6: scale and operational complexity

Frank expects many peers, repeated imports, and larger mixed histories.

Alternative A:

- appears simpler initially
- but moves complexity into importer heuristics and ambiguity handling

Alternative B:

- adds explicit metadata to emitted events and signed records
- but simplifies import semantics, dedupe, and future peer continuity

Alternative C:

- adds a second ID namespace without eliminating timestamp-order ambiguity

Result:

- B adds the most disciplined complexity for the best operational payoff.

## Conclusions

Rejected alternatives:

- Alternative A: envelope CID plus timestamp order is not strong enough to
  become the first honest multi-origin identity/order model
- Alternative C: globally unique event IDs without per-origin monotonic order
  still leave the central ordering problem unresolved

Surviving alternative:

- Alternative B: add explicit `origin_peer_id` plus `origin_sequence`, keep
  envelope CID as content identity, and make the origin tuple the replay and
  dedupe anchor for multi-origin import

## Implications for TODOs and pending DIs

- `107` should lock the first peer-stable identity/order layer around
  `origin_peer_id` plus `origin_sequence`.
- `103` can then implement non-bootstrap import using:
  - known-origin dedupe
  - per-origin ordering
  - explicit conflict/replay handling
- local `Sequence` may remain as a compatibility projection, but it should stop
  being the cross-peer source of truth.

## Decision status

Alternative B locked by `DI-ruzok`: use `origin_peer_id` plus
`origin_sequence` as the first peer-stable replay and dedupe identity, keep
envelope CID as content identity, and demote local `Sequence` to a
compatibility projection.
