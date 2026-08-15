# restore-published-version

Status: repo-local draft for `grid-editor`

## Purpose

`restore-published-version` carries one current-time promise that names an
immutable `publish-document` manifest and the exact Automerge change used to
continue the same live document from that selected publication.

It does not rewrite history, replace a snapshot, make a prior publication
authoritative, or create a new imported document. Concurrent edits remain CRDT
inputs, so the resulting working state can differ from the historical bytes.

## Envelope

Messages are carried as `grid([42(pCID), payload, proof])`.

## Payload

- `kind`: `restore`
- `document_id`: existing live document ID
- `author`: stable relay signing-key ID
- `participant_id`: admitted requester
- `source_manifest_cid`: exact resolved `publish-document` envelope CID
- `live_change_base64`: exact Automerge change bytes carried as a CBOR byte
  string despite the historical field name
- `restored_at`: RFC3339Nano timestamp
- `lamport`: relay-local append ordering value
- `embodiment`: optional local embodiment hint

## Verification

- require the referenced publish manifest to exist for the same document
- reject empty or oversized change bytes
- verify the signed envelope and signer/author binding
- replay the one accepted artifact into the live sync feed

This is a repo-local draft, not a frozen PromiseGrid specification or a claim
of generalized authority, delegation, or interoperability. Source: DI-hihok;
DI-tibum; DI-nihiz.
