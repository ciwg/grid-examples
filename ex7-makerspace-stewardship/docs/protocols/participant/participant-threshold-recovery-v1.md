# Participant Threshold Recovery v1

## Status

Frozen specification. Its pCID is the CIDv1 of these exact file bytes.

## Family

`ex7.participant.threshold-recovery.v1`

## Purpose

This family carries one recovery-witness `promise` for replacement of a lost
or unusable participant root. A completed recovery needs two agreeing witness
records from the three public keys declared in `recovery_set`; one record is
never sufficient.

## Record requirements

The record uses `makerspace-record-v1`, selects this pCID, and carries
canonical JSON with exactly these fields:

```json
{
  "root_record_id": "participant-root-history-record-id",
  "recovery_id": "opaque-nonce-or-incident-id",
  "replacement_root_key": "base64-ed25519-public-key",
  "recovery_set": [
    "base64-ed25519-public-key",
    "base64-ed25519-public-key",
    "base64-ed25519-public-key"
  ]
}
```

Every field is required. `replacement_root_key` and each recovery-set item are
exactly 32 decoded bytes. `recovery_set` has exactly three distinct keys. Two
records complete recovery only when they have the same `root_record_id`,
`recovery_id`, `replacement_root_key`, and byte-for-byte ordered
`recovery_set`, and their envelope signers are distinct members of that set.
The `recovery_set` must exactly equal the set anchored by the referenced root
history record. The completed recovery is linked into root history by a
continuation record.

## Local assessment

Agents preserve conflicting witness sets and incomplete recoveries unchanged.
They may accept a completed recovery only if they locally trust the referenced
root history and the declared witness keys. Recovery witnesses do not gain a
standing authority over ordinary makerspace promises.

## Evolution

Changing the fixed 2-of-3 threshold or any payload meaning requires a new spec
and pCID.
