# Participant Terminal Approval v1

## Status

Frozen specification. Its pCID is the CIDv1 of these exact file bytes.

## Family

`ex7.participant.terminal-approval.v1`

## Purpose

This is an unsigned, bounded embodiment request from a shared terminal to a
participant-owned agent. It is not a PromiseGrid evidence record, does not
create an author identity, and cannot be retained as a promise. Its only result
is either expiry or exact signed `makerspace-record-v1` bytes made by an active
participant root or device.

## Request requirements

The request uses canonical JSON with exactly these fields:

```json
{
  "request_id": "opaque-request-id",
  "target_pcid": "CIDv1-pCID",
  "payload_base64": "base64-canonical-JSON-payload-bytes",
  "created_at": "RFC3339-timestamp",
  "expires_at": "RFC3339-timestamp",
  "approval_token": "base64-32-random-bytes",
  "state": "pending-or-approved-or-expired",
  "signed_record_base64": "present-only-when-approved"
}
```

All fields except `signed_record_base64` are required. The lifetime is at most
ten minutes. `target_pcid` must be an Ex7 pCID the approving agent implements.
`payload_base64` decodes to the exact canonical JSON bytes to be signed; it
must not be normalized. An agent must display the target pCID and exact payload
before local approval. The approval token is a one-time capability for polling
the request status; it is neither a signature nor identity evidence.

## Approval requirements

An approving participant agent verifies its active root/device history and
locally signs a fresh `makerspace-record-v1` record using `target_pcid` and the
exact decoded payload bytes. It marks the request approved only after creating
that exact record. The terminal may submit returned bytes to ordinary record
ingress, where author verification is independent of this request.

## Local assessment and retention

Requests reside in bounded local draft storage and are removed on expiry or
after terminal retrieval. An account session, terminal process, carrier, or
request token cannot approve or alter participant history. A remote recipient
may ignore a request entirely.

## Evolution

Any request field, capability, expiry, or approval semantic change requires a
new spec and pCID.
