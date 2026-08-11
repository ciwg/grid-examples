# Participant Root History v1

## Status

Frozen specification. Its pCID is the CIDv1 of these exact file bytes.

## Family

`ex7.participant.root-history.v1`

## Purpose

This family establishes a participant root signing key or continues that
participant's signed root-key history. It is a `promise` whose meaning is
defined only by this pCID. It does not create a global identity or compel any
agent to recognize the history.

## Record requirements

The record uses `makerspace-record-v1`, selects this pCID, and carries
canonical JSON with exactly these fields:

```json
{
  "root_key": "base64-ed25519-public-key",
  "previous_root_record_id": "optional-record-id",
  "history_note": "human-readable-continuity-note",
  "recovery_set": [
    "base64-ed25519-public-key",
    "base64-ed25519-public-key",
    "base64-ed25519-public-key"
  ]
}
```

`root_key` is required and is exactly 32 decoded bytes. `history_note` is
required and non-empty. `previous_root_record_id` is omitted for establishment
and required for continuation. An establishment record's envelope signer key
must equal `root_key`. A continuation record's signer must be the active root
identified by the referenced predecessor under the receiving agent's local
assessment. `recovery_set` is required, ordered, and contains exactly three
distinct 32-byte public keys. It is the participant-declared 2-of-3 recovery
set for this root-history point. The `Signer` label is not identity evidence.

## Local assessment

An agent retains every valid envelope as exact bytes. It may treat a history as
active only after verifying this pCID's linkage and applying its own trust
policy. Conflicting histories, unknown predecessors, and unrecognized roots
remain retained evidence without an active identity conclusion.

## Evolution

Any field, signer rule, or semantic change requires a new spec and pCID.
