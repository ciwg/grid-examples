# Ex4 local-draft issue-report promise profile

## Status

Local draft. This document is the exact source whose content-derived pCID
selects this profile; it is not a frozen upstream PromiseGrid specification.
Source: `DI-gonok`.

## Purpose

This profile carries one agent's signed promise that it reports a new issue
with the stated title, description, severity, and local team context.

## Envelope

The artifact is a CBOR `grid([42(pCID), payload, proof])` envelope. Slot 0 is
this document's pCID, slot 1 is the profile payload, and slot 2 is the
signer-owned proof over slots 0 and 1. The outer shape is a provisional local
draft adopted consistently with Ex3; it is not a frozen universal API.

## Payload

The payload is a CBOR map with these required fields:

- `agent_id`: text identifier derived from the enrolled public key of the
  signing local embodiment.
- `issued_at`: RFC 3339 UTC timestamp supplied by the signing embodiment.
- `team`: local team context, initially `CORE`.
- `title`: non-empty issue title, at most 120 bytes after validation.
- `description`: non-empty issue description, at most 16,000 bytes after
  validation.
- `severity`: one of `Low`, `Medium`, `High`, or `Critical`.

The local server assigns the human-facing issue ID only after accepting this
promise. The assigned ID is a local projection result, not part of the promise
payload or a portable global identifier.

## Local acceptance and non-claims

The local server verifies the proof and its local public-key/role admission
binding before it accepts an artifact. It may require local reporter policy,
but that policy does not establish general role continuity or authorization.
An accepted artifact may update the local issue projection; rejection is a
local observation/diagnostic, not a statement of global invalidity or intent.

## Evolution

Changes to any field, validation rule, or meaning require a new profile source
and therefore a new pCID. Future cross-tracker exchange, delegation,
revocation, and recognized-role semantics are outside this profile.
