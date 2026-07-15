# grid-editor publish and exchange slice

TE ID: TE-vafor
## Status
needs DF

## Decision under test

How should `grid-editor` add its first PromiseGrid-native publish / document
exchange slice without breaking the existing CRDT relay model?

## Assumptions

- The current relay stays the live CRDT exchange surface for editing.
- CAS already stores exact signed envelope bytes under `<data-root>/cas/**`.
- Phase 4 should add a PromiseGrid-native publish/exchange path, not permissions
  or restore semantics yet.
- Alice, Bob, and Carol are cooperative users; Mallory is the adversary.

## Alternatives

### Alternative 1

Publish the current browser-local snapshot metadata as-is from the browser.

Meaning:
- browser writes snapshot/export metadata directly
- relay remains unaware of publish objects
- publish links point to browser-managed local artifacts

### Alternative 2

Add a relay-owned signed publish manifest object that references existing CAS
document change history and browser-exported presentation bytes.

Meaning:
- relay accepts a local publish request
- relay signs a new publish manifest envelope
- manifest references:
  - document ID
  - selected version marker or current feed offset
  - optional exported bytes CID
  - summary/title metadata
- CAS is the durable publish store

### Alternative 3

Make publish a full new top-level server-side document registry with canonical
document metadata and document ownership now.

Meaning:
- relay becomes authoritative for publish registry state
- publish/exchange, owner/admin, favorites, tags, archive, and collections
  arrive together

## Scenario analysis

### S1 — normal publish from current document

Alice edits a document in the browser and wants a stable artifact she can hand
to Bob.

Alternative 1 is easy but weak. The published artifact depends on browser-local
state and does not create a relay-verifiable PromiseGrid-shaped object. Bob
cannot inspect a signed publish record from the relay.

Alternative 2 creates a stable signed manifest at publish time. Alice still
edits through Automerge, but publish becomes a separate current-time action.
Bob can fetch the manifest and then fetch the referenced content from CAS or
the relay. This fits the earlier decision that publish should be a new action,
not an invisible mutation of live state.

Alternative 3 also works, but it drags in too much unneeded registry behavior
for the first slice.

### S2 — multi-machine exchange

Alice publishes on one relay and Bob wants to import or view the result from a
different relay.

Alternative 1 is poor here because the publish artifact is browser-local and
not clearly replayable by another relay.

Alternative 2 fits well. The manifest is a signed object with explicit fields
and CAS references. Bob's relay can ingest the signed manifest just as it
ingests other signed envelopes, and Bob can materialize a local imported copy
from the referenced state.

Alternative 3 can support this, but again at the cost of pulling permissions
and registry semantics into the first slice.

### S3 — failure and replay

Alice publishes, then her browser crashes.

Alternative 1 risks losing the publish metadata if it only lives in browser
storage.

Alternative 2 survives because the relay signs and persists the publish object
to CAS and the append-only log. Alice can rediscover it later.

Alternative 3 also survives, but with more surface area than needed.

### S4 — adversarial tampering

Mallory modifies a manifest or tries to claim someone else's publish.

Alternative 1 offers the weakest evidence chain.

Alternative 2 preserves the existing signed-envelope model. Mallory can only
introduce a forged publish if she controls a signing key accepted by the relay,
and manifest/object tampering breaks the CAS address or signature proof.

Alternative 3 offers similar security if done correctly, but with a larger
attack surface.

### S5 — long-horizon migration

Later phases need restore, permissions, archive, favorites, and collections.

Alternative 1 creates a dead end because publish has no durable relay-native
shape.

Alternative 2 gives a narrow durable object model now and leaves room to layer
later Phase 4 features around it. Restore can later reference publish/version
manifests without redefining publish.

Alternative 3 front-loads too much product/backend policy before those later
decisions are separately tested.

## Conclusions

- Reject Alternative 1 because browser-local publish artifacts are not
  PromiseGrid-native enough and do not survive relay-to-relay exchange well.
- Reject Alternative 3 for the first slice because it bundles too many later
  Phase 4 concerns into one large authoritative registry change.
- Keep Alternative 2 as the surviving design.

## Surviving alternative

Alternative 2:
- add a relay-owned signed publish manifest
- persist it to CAS and the append-only log
- treat publish as a current-time action that references existing document
  state, not as a rewrite of past live edits

## New decisions exposed

- whether publish should reference:
  - the current sync-feed offset
  - a named saved version
  - or either one
- whether the first slice should ship:
  - publish only
  - or publish plus import/exchange in the same slice
- where the publish manifest protocol doc should live

## Implications for TODOs and DIs

- A new Phase 4 TODO is needed for publish/exchange.
- The first Phase 4 slice should update docs and protocol notes to explain that
  publish is a relay-signed durable object, separate from live CRDT editing.
- Later permissions/restore work should build on publish manifests instead of
  replacing them.

## Decision status

needs DF
