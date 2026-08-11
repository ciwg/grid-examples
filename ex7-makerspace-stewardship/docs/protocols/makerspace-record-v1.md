# Makerspace Record v1

Status: frozen

This document defines Ex7's common canonical evidence envelope. It is not a
family selector: a family specification selects its own pCID in the outer Grid
selector.

## Canonical carriage

Each durable record is canonical CBOR `grid(...)` with tag 42 over the CIDv1
pCID of a frozen makerspace family specification. Slots after the selector are
exactly:

```text
[ record_id, signer_label, created_at_rfc3339,
  canonical_json_payload_bytes, author_key_id,
  author_public_key_bytes, author_signature_bytes ]
```

`record_id` and `signer_label` are non-empty UTF-8 text; the former is an
opaque writer-unique identifier and the latter is local bootstrap/presentation
data, not global identity. `created_at_rfc3339` is UTC RFC3339 text.
`canonical_json_payload_bytes` are canonical JSON defined by the selected
family. `author_key_id` is `ed25519:` plus lower-case hexadecimal SHA-256 of
the 32-byte Ed25519 `author_public_key_bytes`; `author_signature_bytes` is a
64-byte Ed25519 signature.

## Signing view

The author signs the exact canonical CBOR encoding of the same `grid(...)`
record with `author_signature_bytes` represented as CBOR null. All other slots,
including selector and payload bytes, remain unchanged. A verifier recomputes
that view before checking signature and key fingerprint.

## Local admission and projection

An implementation may retain well-framed canonical bytes for an unknown pCID,
unrecognized key, or locally revoked key, but must not assign unknown-family
semantics. Ex7 projects only known-family records with valid signatures from
active, locally recognized author/key pairs. This is local policy, not global
authority or portable key continuity. Malformed frames, encodings, pCIDs, key
fingerprints, or signatures are corruption and make Ex7's persistent store fail
closed.

## Versioning

This document is immutable. A change to slots, signing view, encoding, or
meaning requires a new versioned specification.
