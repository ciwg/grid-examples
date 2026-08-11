# Participant Key Revocation v1

## Status

Frozen specification. Its pCID is the CIDv1 of these exact file bytes.

## Family

`ex7.participant.key-revocation.v1`

## Purpose

This family is an active-root `promise` to stop accepting a named root or
device key at a declared time. It is evidence for each receiving agent's local
assessment; it cannot erase or rewrite earlier exact records.

## Record requirements

The record uses `makerspace-record-v1`, selects this pCID, and carries
canonical JSON with exactly these fields:

```json
{
  "subject_key_id": "ed25519:lower-case-sha256-hex",
  "subject_kind": "root-or-device",
  "effective_at": "RFC3339-timestamp",
  "reason": "non-empty-human-readable-reason"
}
```

Every field is required. `subject_kind` is exactly `root` or `device`. The
envelope signer must be an active root in the participant history that locally
binds `subject_key_id`. A revocation record does not by itself select a
replacement key; a root-history continuation or completed recovery does that.

## Local assessment

Agents retain revocation records even when their root history is incomplete.
For an accepted history, records signed by the subject after `effective_at` do
not support local projection. Earlier records remain evidence and are not
silently discarded.

## Evolution

Any field, signer rule, or temporal interpretation change requires a new spec
and pCID.
