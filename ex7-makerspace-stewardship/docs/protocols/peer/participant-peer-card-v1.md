# Participant Peer Card v1

## Status

Frozen specification. Its pCID is the CIDv1 of these exact file bytes.

## Family

`ex7.peer.participant-card.v1`

## Purpose

This family is a signed, portable peer-card `promise` that lets another agent
find a claimed participant root-history reference, active-device references,
and optional contact hints. It is discovery evidence, not an account
credential, membership statement, or proof that a contact hint is reachable.

## Record requirements

The record uses `makerspace-record-v1`, selects this pCID, and carries
canonical JSON with exactly these fields:

```json
{
  "root_record_id": "participant-root-history-record-id",
  "active_device_record_ids": ["device-authorization-record-id"],
  "contact_hints": ["optional-opaque-contact-hint"]
}
```

All fields are required; either array may be empty. Each value is a non-empty
string, and arrays preserve their signed order. The envelope signer must be
the active root or a locally active authorized device in the referenced
history. Contact hints are opaque strings and must not be interpreted as
identity, authorization, or a promise of availability.

## Local assessment

An agent retains a valid card even when it cannot resolve its history. It may
attempt direct contact according to local safety and privacy policy. A card
cannot alter trust, authorization, recovery, or makerspace projection by
itself.

## Evolution

Any field or discovery interpretation change requires a new spec and pCID.
