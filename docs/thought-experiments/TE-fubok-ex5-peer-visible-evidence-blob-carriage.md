# Ex5 Peer-Visible Evidence Blob Carriage

TE ID: `TE-fubok`
## Status
decided

## Decision under test

How `ex5` should carry evidence blobs when it makes `knowledge-evidence`
peer-visible.

Related TODO:

- `105` - `ex5-operational-knowledge-system/TODO/TODO-faruv-ex5-peer-visible-evidence-exchange.md`

## Assumptions

- `knowledge-evidence` is already a frozen signed family in the local runtime.
- Evidence metadata already carries `attachment_cid` in the durable event and
  signed envelope payload.
- CAS already stores copied evidence blobs by blob CID.
- Current peer exchange already transports JSON bundles over the local HTTP
  adapter.
- The next evidence step should be honest and portable, not merely expose local
  attachment paths to another host.

## Alternatives

### Alternative A

Extend peer exchange so the evidence bundle is self-contained:

- export signed `knowledge-evidence` records
- include the required CAS blobs inline in the bundle, keyed by blob CID
- import the records plus blobs together

This keeps one exchange artifact sufficient to reproduce evidence on another
host.

### Alternative B

Exchange evidence metadata first, then fetch blobs separately by blob CID from
a CAS object surface.

This keeps the primary peer-exchange bundle smaller and treats blob carriage as
a separate content-addressed fetch path.

### Alternative C

Defer peer-visible evidence exchange until a later full relay/blob transport
layer exists.

## Scope and systems affected

- `ex5-operational-knowledge-system/TODO/TODO-faruv-ex5-peer-visible-evidence-exchange.md`
- `ex5-operational-knowledge-system/TODO/TODO.md`
- `ex5-operational-knowledge-system/service/peer_exchange.go`
- `ex5-operational-knowledge-system/service/server.go`
- `ex5-operational-knowledge-system/service/types.go`
- evidence/CAS storage and import tests
- PromiseGrid claims and peer-exchange/CAS staging docs

## Scenario analysis

### Scenario 1: self-contained bootstrap exchange to a fresh host

Alice exports evidence-bearing history and Bob imports it on a fresh host.

Alternative A:

- succeeds with one exchange artifact
- keeps the imported evidence immediately resolvable
- uses the existing bundle-oriented exchange surface naturally

Alternative B:

- requires Bob to perform a second fetch step for each missing blob or blob set
- makes the first evidence slice more operationally complex

Alternative C:

- still cannot exchange the evidence at all

Result:

- A is the cleanest first portable evidence step.

### Scenario 2: repeated blob reuse across multiple evidence records

Carol references the same blob from multiple evidence records or across
multiple imports.

Alternative A:

- can still dedupe naturally by carrying blobs keyed by CID and skipping
  already-present CIDs on import
- may duplicate blob bytes across separate exported bundles, but not within one
  bundle

Alternative B:

- is more bandwidth-efficient over time because blobs can be fetched on demand
  and cached independently

Alternative C:

- offers no peer-visible evidence path

Result:

- B is better for long-run transport efficiency, but A remains acceptable for a
  first self-contained exchange slice.

### Scenario 3: missing or corrupted blobs during import

Dave receives evidence metadata, but one blob is missing or tampered with.

Alternative A:

- can validate presence and integrity during one import transaction
- gives a straightforward pass/fail result for the whole evidence payload

Alternative B:

- spreads integrity handling across bundle import plus later blob fetches
- needs a second error surface and recovery path immediately

Alternative C:

- avoids the problem only by deferring the feature

Result:

- A gives the simpler first integrity story.

### Scenario 4: later relay-visible transport growth

Ellen wants eventual larger-scale multi-peer blob exchange.

Alternative A:

- is not the final most bandwidth-efficient shape
- but keeps blob identity explicit by CID, so the transport can later evolve
  without invalidating stored evidence semantics

Alternative B:

- more closely matches a future standalone CAS object transport
- but requires that extra object surface before the first evidence step can
  ship

Alternative C:

- delays real learning from peer-visible evidence exchange

Result:

- A is the best first slice if the repo wants implementable progress now while
  preserving CID-based blob identity.

### Scenario 5: user priority on “fully on the grid”

Steve wants `ex5` to keep moving rather than waiting for a perfect blob
transport stack.

Alternative A:

- advances `knowledge-evidence` into the peer-visible set now
- stays honest because the bundle carries both the metadata and the actual blob
  bytes referenced by CID

Alternative B:

- is arguably more elegant
- but asks the repo to design and implement a second object-exchange surface
  before evidence itself becomes peer-visible

Alternative C:

- is too conservative for the current goal

Result:

- A best balances implementability with on-grid progress.

## Conclusions

Rejected alternatives:

- Alternative C: too conservative; it leaves the evidence family off-grid
- Alternative B: a plausible later refinement, but it adds an extra transport
  surface before the first peer-visible evidence slice exists

Surviving alternative:

- Alternative A: make the first evidence exchange self-contained by carrying
  signed evidence records plus the referenced CAS blobs inline, keyed by CID

Implications and future work:

- `105` should extend the peer-exchange bundle to include evidence records and
  blob payloads together
- import can dedupe blobs by CID while still remaining self-contained
- a later standalone CAS object exchange surface can still be added without
  changing the evidence family semantics

## Decision status

Alternative A locked by `DI-faruv`: the first peer-visible
`knowledge-evidence` slice carries signed evidence records plus inline
CID-keyed CAS blobs inside the bootstrap bundle.
