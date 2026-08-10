# Ex2 relay exception evidence policy

TE ID: TE-bulaj

## Status

decided

## Decision status

Locked by DI-todav: retain bounded raw exception bytes, append one separate
local observation per receipt, use `no_supported_handler`, and exclude
exception observations from accepted-message replay.

## Decision under test

What exact local evidence should the Ex2 relay retain when it receives bytes
that are malformed, carry an invalid proof, or select a pCID for which this
relay has no handler, without treating rejection as a statement about global
pCID validity or another participant's intent?

## Assumptions and trust model

- The relay's current accepted-message path verifies a signed envelope, stores
  its exact bytes in CAS, appends a normal `message-log.jsonl` entry, and then
  projects known pCID payloads into feeds or latest-state indices.
- Current malformed or invalid-proof bytes are rejected before CAS and normal
  log retention. A valid signed envelope with an unsupported pCID currently
  reaches CAS and the normal message log before the handler switch rejects it.
- The relay is an observer and dispatcher for its local supported pCIDs. It is
  not a global pCID registry, a canonical document host, or a business/promise
  assessment authority.
- Alice and Bob run cooperative relays. Mallory may send malformed bytes,
  signature failures, replayed invalid traffic, valid future pCIDs, or traffic
  that this checkout does not support.
- Retaining arbitrary inbound bytes has storage and abuse cost. Any selected
  policy must remain bounded by existing request-size limits and must keep
  exception evidence separate from accepted-message records.
- The current PromiseGrid guide treats timeout, silence, and local observations
  as observer-local evidence, while signed artifacts can later be reassessed
  without converting the relay's local rejection into a global verdict.

## Alternatives

### A. Retain bounded raw bytes and append distinct local observations

Compute a raw CID for every received byte sequence within the existing size
bound. Store the exact bytes in CAS when possible and append an
`observations.jsonl` record for `malformed_input`, `invalid_proof`, or
`no_supported_handler`. Keep accepted valid envelopes in `message-log.jsonl`;
exception observations are a separate append-only local record. For a valid
unsupported pCID, say only that this relay has no supported handler, not that
the pCID is globally unknown or invalid.

### B. Preserve the current asymmetric behavior

Reject malformed and invalid-proof bytes without durable evidence; retain valid
unsupported-pCID envelopes in the normal CAS and message log, then return an
unknown-pCID error.

### C. Discard all rejected bytes

Keep only operational error responses and retain neither invalid nor
unsupported inbound traffic.

## Scenario analysis

### Malformed bytes

Mallory sends bytes that cannot parse as a grid envelope. A lets Alice retain
the bounded exact bytes and a local `malformed_input` observation; the record
does not say Mallory promised anything. B and C leave Alice with only a
transient error, making later diagnosis and abuse review impossible. A creates
a local storage obligation, bounded by the already accepted request limit.

### Invalid proof

Mallory sends a syntactically valid envelope whose proof fails verification.
Under A, the relay retains raw bytes and records `invalid_proof` with the raw
CID and locally assessed reason. This is not a valid signed message and must
not enter the accepted-message log or feed. B and C reject it without durable
evidence. A makes it possible to distinguish a proof failure from a malformed
frame later without treating either as a broken promise.

### Valid future or unsupported pCID

Bob sends a valid signed envelope for a pCID introduced by a newer checkout.
A retains bytes and records `no_supported_handler`, a fact about Alice's local
handler set. It neither rejects the pCID's meaning globally nor appends it to
the accepted-message log. B currently retains it as though it were a normal
accepted message, then reports `unknown pCID`, conflating carriage evidence
with accepted relay history. C loses the only evidence of local non-support.

### Replay, restart, and mixed versions

Mallory repeats the same invalid bytes; Bob restarts after receiving a future
pCID. A's CAS identity deduplicates raw bytes, while append-only observations
may either record each receipt or use a local deduplication rule selected with
the record schema. On restart, accepted-message replay continues to rebuild
only accepted projections; exception observations remain evidence rather than
replay input. B's normal log can fail replay when it encounters the previously
logged unsupported pCID. C has no durable record of the event.

### Scale and denial-of-service pressure

A must cap ingress before retention, avoid relay fan-out for rejected bytes,
and expose local retention as an operator policy rather than a global service
guarantee. B consumes CAS/log storage for unsupported pCIDs but provides no
consistent malformed/proof evidence. C is cheapest but abandons auditability.
A has the clearest operational cost and can later add retention limits without
changing the semantic boundary.

## Conclusions

C is rejected because it discards the local evidence needed to distinguish
malformed input, proof failure, and absent local handler support. B is rejected
because it makes unsupported valid traffic appear in the accepted-message log,
leaves other rejection classes unrecorded, and risks replay failure across
mixed versions.

A survives and is recommended: retain bounded exact bytes and append distinct
relay-local observation records; preserve `message-log.jsonl` for accepted
known-profile envelopes only; call unsupported traffic `no_supported_handler`;
and make exception observations non-replay input and non-global claims.

## Decisions still requiring DF

1. **Evidence form:** add bounded raw-byte retention plus separate append-only
   local observation records (recommended), preserve current asymmetry, or
   discard rejected bytes?
2. **Unsupported-pCID name:** use `no_supported_handler` (recommended), or
   retain the current `unknown_pcid` wording?
3. **Observation multiplicity:** append one observation for every received
   exception (recommended), or deduplicate observations by raw CID within one
   running process?
4. **Storage path:** use `observations.jsonl` beside the relay's existing
   `message-log.jsonl` and `cas/` data (recommended), or fold exceptions into
   the normal message log?

## Implications for open work

- `sojot.2` can proceed after these DFs are locked and recorded in
  `TODO-sojot`.
- The selected policy requires schema, CAS/ingress ordering, replay, and
  regression-test decisions before implementation.
- `sojot.3` must cover each selected observation class and confirm rejected
  traffic does not enter accepted-message replay.
