# Participant Device Authorization v1

## Status

Frozen specification. Its pCID is the CIDv1 of these exact file bytes.

## Family

`ex7.participant.device-authorization.v1`

## Purpose

This family is a root-authored `promise` that a particular public key may act
as a participant device key during a stated interval. It does not turn the
device, browser, account, or facility into the participant's identity.

## Record requirements

The record uses `makerspace-record-v1`, selects this pCID, and carries
canonical JSON with exactly these fields:

```json
{
  "root_record_id": "participant-root-history-record-id",
  "device_key": "base64-ed25519-public-key",
  "device_label": "participant-chosen-label",
  "not_before": "RFC3339-timestamp",
  "not_after": "optional-RFC3339-timestamp"
}
```

All fields except `not_after` are required. `device_key` is exactly 32 decoded
bytes. The envelope signer must be the active root key resolved from
`root_record_id` under local assessment. A device key may sign ordinary Ex7
records only while a locally trusted authorization is active and not revoked.

## Local assessment

Every valid record remains exact retained evidence. An agent must not infer
authorization from a device label, account session, peer card, or carrier. A
receiving agent independently decides whether it recognizes the root history,
time source, and device authorization.

## Evolution

Any field, time rule, or delegated authority change requires a new spec and
pCID.
