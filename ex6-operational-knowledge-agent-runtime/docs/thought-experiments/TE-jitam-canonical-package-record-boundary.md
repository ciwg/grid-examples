# Canonical Package Record Boundary

TE ID: TE-jitam

## Status

superseded by TE-hofor / DI-sidoh

## Decision under test

Choose the implementation boundary for DI-salaf's compatibility-first migration
of Ex6 package records to real pCIDs and canonical Grid carriage.

## Assumptions and scope

- Legacy JSON `records.Envelope` objects remain readable only as historical
  implementation-local evidence.
- New package records must be canonical CBOR Grid envelopes selected by a real
  CIDv1 pCID and must retain current local author-signature evidence.
- Alice and Bob may retain mixed legacy and canonical history. Mallory may
  submit malformed, noncanonical, replayed, or incorrectly selected bytes.
- Existing workflow and route-promise Grid formats remain separately specified.

## Alternatives

### A. Replace `records.Envelope` in place

Change the existing public type and `AppendRecord` contract so every caller
uses canonical Grid bytes directly.

### B. Add an explicit canonical package-record type

Keep `records.Envelope` as the named legacy JSON compatibility type. Add a
separate canonical package-record type and append/import path for all new
package writes; dispatch and validation determine the format from exact bytes.

### C. Wrap the legacy JSON envelope inside a Grid slot

Use one generic real pCID for a Grid wrapper whose payload remains the old JSON
shape and symbolic protocol label.

## Scenario analysis

### Normal operation

A has one API but makes every existing package and external package source
change immediately. B makes the boundary visible: new package writes cannot
accidentally use the legacy encoder. C preserves most code but leaves protocol
selection inside a symbolic JSON field rather than at the Grid selector.

### Failure and incomplete migration

A turns any incomplete caller migration into a build or runtime outage. B
allows exact-byte detection: malformed canonical records are rejected, valid
legacy records remain quarantined to the compatibility reader, and a partial
rollout cannot reinterpret stored JSON. C admits a valid wrapper around an
invalid protocol identity.

### Mixed-version and concurrent nodes

A forces Alice and Bob to upgrade together. B permits both histories to remain
inspectable while routing and new claims advertise only real pCIDs. C keeps the
old ambiguity at the wire boundary, so peers cannot safely decide whether two
symbolic labels mean the same protocol.

### Evolution and trust

A and B can require a new specification CID whenever protocol meaning changes.
B isolates the retired JSON shape, making later deletion a conscious policy
change. C binds a wrapper spec but not the actual package-record semantics,
weakening pCID-selected interpretation. All alternatives retain signatures as
evidence, not authority; B makes the signed bytes and selected spec easiest to
inspect together.

### Scale and operational complexity

A minimizes steady-state code but concentrates risk in one migration. B carries
a bounded compatibility reader and two clearly named paths during transition.
C appears smallest but creates ongoing ambiguity and support burden. B's extra
code is justified because local history, relay inputs, and installed packages
must be upgraded independently.

## Conclusion

C is rejected because it preserves symbolic protocol semantics inside the new
carriage. A is rejected for this slice because it makes the compatibility-first
decision ineffective. B survives: retain `records.Envelope` as explicit legacy
JSON compatibility evidence and add a distinct canonical package-record
boundary for all new writes and imports.

## Decision status

superseded by TE-hofor / DI-sidoh. The unreleased-product decision removes the
compatibility premise and selects a clean canonical re-creation instead.

## Implications for TODOs and pending DIs

- TODO `puvok` needs DI-jitod before runtime code changes.
- The implementation will touch the canonical record codec, runtime append and
  import dispatch, manifests/claims/routes, package writers, tests, protocol
  specifications, and alignment documentation.
