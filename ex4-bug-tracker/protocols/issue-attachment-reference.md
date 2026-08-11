# Ex4 local-draft issue-attachment-reference promise profile

## Status

Local draft. This document is the exact source whose content-derived pCID
selects this profile; it is not a frozen upstream PromiseGrid specification.
Source: `DI-gonok`.

## Purpose

This profile carries one agent's signed promise that it associates one
content-addressed attachment object with an existing locally projected issue.

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
- `issue_id`: locally projected issue identifier.
- `attachment_cid`: content-derived identifier of the attachment bytes.
- `name`: sanitized presentation filename.
- `content_type`: media type selected or detected by the local adapter.
- `size`: non-zero attachment size in bytes, bounded by local service policy.

The attachment bytes are retained separately in local content-addressed
storage. The original host path is not part of the payload.

## Local acceptance and non-claims

The local service verifies proof, public-key enrollment, issue existence,
attachment-object presence, and its local size/type policy before accepting an
artifact. These are local admission and projection decisions, not a portable
authorization system or proof of another actor's intent. Rejection is a local
observation/diagnostic, not global invalidity.

## Evolution

Changes to any field, validation rule, or meaning require a new profile source
and therefore a new pCID. Future cross-tracker object exchange, delegation,
revocation, and recognized-role semantics are outside this profile.
