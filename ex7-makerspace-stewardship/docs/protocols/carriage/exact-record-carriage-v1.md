# Exact Record Carriage v1

## Status

Frozen specification. Its pCID is the CIDv1 of these exact file bytes.

## Family

`ex7.carriage.exact-record.v1`

## Purpose

This family is a carrier's `promise` that it is forwarding an opaque batch of
exact record bytes at a cursor. It supports direct peer exchange and optional
relay/feed transport without granting the carrier authorship, interpretation,
trust, or makerspace authority over the enclosed records.

## Record requirements

The wrapper uses `makerspace-record-v1`, selects this pCID, and carries
canonical JSON with exactly these fields:

```json
{
  "sender_card_record_id": "participant-peer-card-record-id",
  "cursor": "opaque-monotonic-carrier-cursor",
  "records": ["base64-exact-record-bytes"]
}
```

Every field is required. `records` is non-empty. Each item decodes to one
complete exact Grid record byte string; the wrapper must not normalize,
re-sign, filter, reorder, or interpret it. A receiving agent deduplicates by
the enclosed record's exact bytes or validated record ID under its own policy.
Cursor gaps, replay, outage, and retention are carriage observations, never
evidence that a makerspace promise is true.

## Local assessment

The wrapper signature only identifies its carrier key. A receiving agent
verifies each enclosed record independently, retains unknown pCID bytes
unchanged, and projects only records whose semantics and author evidence it
locally accepts. Direct carriage and relay carriage use this same format.

## Evolution

Any wrapper field, opaque-byte handling, or duplicate rule change requires a
new spec and pCID.
