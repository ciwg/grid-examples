# Unreleased Canonical Package Record Re-creation

TE ID: TE-hofor

## Status

decided

## Decision under test

Determine whether Ex6 must retain a compatibility path for its symbolic JSON
package records when correcting its pCID and canonical-carriage boundary.

## Assumptions and scope

- Ex6 has not been released and has no external compatibility or retained-user
  data commitment.
- The existing package record format is a development artifact that uses
  symbolic `pcid:<name>` labels rather than specification CIDs.
- A pCID is a CIDv1 of a frozen protocol specification, and new package
  records must use canonical Grid carriage selected by that pCID.
- Alice and Bob represent developers running a clean rebuilt version. Mallory
  represents malformed or altered input, not a historical deployed peer.

## Alternatives

### A. Preserve a legacy compatibility reader

Keep `records.Envelope` and symbolic JSON parsing alongside a new canonical
package-record model.

### B. Re-create the package record model cleanly

Remove the symbolic JSON durable-package-record contract and replace it with
one canonical Grid package-record contract, real frozen pCIDs, and updated
package/runtime/relay tests and documentation.

## Scenario analysis

### Normal development

A adds two record models that every contributor must distinguish, despite no
released user requiring it. B makes all active code, examples, and manifests
exercise the one intended protocol contract.

### Failure and malformed input

A must decide whether a malformed JSON object is a legacy record, unknown
bytes, or an invalid new record. B has one package-record decoder: canonical
Grid bytes with a real selector; malformed bytes are rejected or retained only
under the existing unknown-byte policy without a false package-record claim.

### Concurrent developers and clean environments

A leaves mixed-format fixtures and behavior for developers to maintain. B
forces Alice and Bob to rebuild from the same source contract, which is
appropriate before release and exposes missed conversions immediately in tests.

### Long-horizon evolution and trust

A preserves an incorrect public-looking pCID vocabulary indefinitely. B binds
each active package protocol to its immutable specification CID; a semantic
change creates a new spec and pCID. Signatures remain evidence over exact
canonical bytes and do not create authority.

### Operational scale

A costs less in the immediate patch but permanently expands decoding, relay,
and documentation scope. B costs one bounded re-creation now and reduces the
steady-state system to one auditable record contract.

## Conclusion

Alternative A is rejected because no release boundary requires it and it would
preserve the very ambiguity being corrected. Alternative B is locked by
DI-sidoh: Ex6 is re-created with one canonical package-record contract; no
symbolic JSON compatibility path is retained.

## Implications for TODOs and pending DIs

- This supersedes the compatibility-first conclusion in TE-ralam and the
  pending boundary choice in TE-jitam.
- TODO `puvok` now owns a clean replacement of the package record codec,
  manifests, routing, relay evidence, package writers, tests, and alignment
  documentation.
