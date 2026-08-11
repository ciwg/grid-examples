# Package pCID Migration to Canonical Grid Carriage

TE ID: TE-ralam

## Status

superseded by TE-hofor / DI-sidoh

## Decision under test

Bring Ex6's package-defined durable records into PromiseGrid alignment without
silently reinterpreting historical JSON records or symbolic `pcid:<name>`
labels. This concerns TODO `puvok` and the package/runtime boundary.

## Assumptions and scope

- A pCID is the CID of an immutable protocol specification, not a family name,
  version label, message body hash, or implementation claim.
- `records.GridEnvelope` already provides canonical CBOR carriage with a CIDv1
  protocol selector and is the only current Ex6 carriage that meets that
  requirement.
- Existing package records use `records.Envelope`, JSON, and symbolic
  `pcid:<name>` labels. Those labels have no frozen specification object whose
  CID they name.
- The migration applies to new package-defined durable records, package
  manifests, claims, routing metadata, relay metadata, tests, examples, and
  operator documentation. It does not redefine the existing workflow or
  route-promise protocol selectors.
- Alice and Bob may run mixed Ex6 versions. Mallory may replay, alter, or
  advertise bytes but cannot forge a valid signature for a trusted key.

## Alternatives

### A. Keep symbolic labels and describe them as pCIDs

Retain the JSON envelope and treat `pcid:moks.*` names as identifiers.

### B. Breaking replacement everywhere

Replace every symbolic label and JSON envelope at once with canonical Grid
records and real pCIDs, with no reader for historical records.

### C. Compatibility-first migration

Freeze a protocol specification for each package record family, derive and
publish its CIDv1 pCID, write all new durable package records as canonical Grid
envelopes, and keep old symbolic JSON records readable only as explicitly
legacy implementation-local evidence. New routes, claims, and relay metadata
use real pCIDs. Legacy records never gain a new semantic interpretation.

## Scenario analysis

### Normal operation

Under A, Alice can name a family consistently but cannot prove that the name
selects a specific protocol specification. A changes no code but preserves the
alignment error.

Under B, Alice gets one unambiguous format immediately, but current package
data and operator examples stop working at the boundary.

Under C, Alice writes a canonical Grid envelope whose selector is the CID of a
published family specification. The implementation can still inspect existing
JSON history as historical local evidence, while new claims and routes have an
interoperable protocol identity.

### Failure, corruption, and incomplete writes

Under A, JSON parsing can detect malformed JSON but not canonical encoding, and
the symbolic selector cannot bind the bytes to a frozen specification.

Under B and C, canonical CBOR decoding rejects noncanonical carriage and the
CID selector is independently parseable. C additionally avoids turning an
interrupted rollout into data loss: immutable old JSON stays inspectable and a
failed new write is simply absent from durable history.

### Concurrent and mixed-version nodes

Under B, Bob's unupgraded runtime cannot read Alice's newly written records;
Alice cannot read Bob's retained history after an abrupt switch.

Under C, upgraded nodes preserve and identify legacy JSON as legacy rather than
mistaking it for a PromiseGrid record. Interoperation of new records requires a
real pCID-capable peer, and that boundary is visible in relay/import policy.
The runtime must never map a symbolic label to a new pCID implicitly.

### Long-horizon evolution and migration

Under A, a prose edit can change the effective meaning of a stable string.
Under B, every protocol edit and deployment must be coordinated globally.

Under C, a changed specification creates a new CID and therefore a new pCID.
The old spec and bytes remain independently interpretable. Compatibility code
can be retired later only by a separately recorded decision after retained
history and supported peers no longer need it.

### Trust-boundary changes

Under A, symbolic labels encourage trusting package vocabulary rather than the
selected protocol. Under B and C, policy can reason over an exact pCID.

C keeps local identity limits explicit: author signatures and relay carriage
signatures are evidence over the exact durable bytes, while policy continues to
decide what those signatures authorize. It does not turn package ownership,
workflow approval, or relay receipt into universal authority.

### Scale and operations

B has a lower steady-state implementation surface but moves all migration cost
into one risky rollout. C temporarily maintains legacy decoding plus canonical
writing, but it makes storage growth append-only, permits staged peer upgrades,
and leaves reproducible evidence for operators. Protocol-spec files add modest
repository overhead but make pCID provenance inspectable.

## Conclusion

Alternative A is rejected because it preserves a false pCID claim. Alternative
B is rejected for this slice because it makes existing durable records and
mixed-version operation needlessly brittle. Alternative C survives and is
locked by DI-salaf: compatibility-first migration, real CIDv1 pCIDs and
canonical Grid carriage for new durable package records, and explicit
non-PromiseGrid treatment of historical symbolic JSON records.

## Implications for TODOs and pending DIs

- TODO `puvok` owns the runtime migration because it governs the agent/kernel
  protocol boundary.
- A package-protocol specification registry, canonical record encoder/decoder,
  migration policy, package/relay/routing updates, and deterministic tests are
  required before Ex6 may claim package-level PromiseGrid alignment.
- `docs/testing.md`, a README link, implementation claims, and a CHANGELOG are
  required to make the resulting scope and evidence auditable.
