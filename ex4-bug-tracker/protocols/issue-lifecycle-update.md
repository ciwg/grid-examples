# Ex4 local-draft issue-lifecycle-update promise profile

## Status

Local draft. This document is the exact source whose content-derived pCID
selects this profile; it is not a frozen upstream PromiseGrid specification.
Source: `DI-gonok`.

## Purpose

This profile carries one agent's signed promise about an existing locally
projected issue: a comment, an assignment, or a status transition.

## Envelope

The artifact is a CBOR `grid([42(pCID), payload, proof])` envelope. Slot 0 is
this document's pCID, slot 1 is the profile payload, and slot 2 is the
signer-owned proof over slots 0 and 1. The outer shape is a provisional local
draft adopted consistently with Ex3; it is not a frozen universal API.

## Payload

The payload is a CBOR map with these required common fields:

- `agent_id`: text identifier derived from the enrolled public key of the
  signing local embodiment.
- `issued_at`: RFC 3339 UTC timestamp supplied by the signing embodiment.
- `issue_id`: locally projected issue identifier.
- `kind`: one of `comment`, `assignment`, or `status`.

`kind=comment` additionally requires non-empty `comment` text no longer than
8,000 bytes after validation.

`kind=assignment` additionally requires `assignee_agent_id`, the locally
enrolled agent identifier selected as assignee. The local service may map this
to a current presentation label for its projection.

`kind=status` additionally requires `status`, one of `Triaged`, `In Progress`,
or `Resolved`. Reopening is represented by `status=Triaged` from locally
projected `Resolved` state.

## Local acceptance and non-claims

The local service verifies proof, public-key enrollment, issue existence, and
its own transition/role policy before accepting an artifact. That local policy
does not establish general identity, delegation, role continuity, or global
authorization. An accepted artifact may update the local projection; rejection
is a local observation/diagnostic, not a statement of global invalidity or
intent.

## Evolution

Changes to any field, validation rule, or meaning require a new profile source
and therefore a new pCID. Future cross-tracker exchange, delegation,
revocation, and recognized-role semantics are outside this profile.
