# document-metadata

Status: repo-local draft for `grid-editor`

## Purpose

`document-metadata` carries relay-signed latest-state metadata for a document.

This protocol is for durable metadata such as description, summary, tags,
collections, favorite state, and archive state. It is separate from
`live-document`, `live-awareness`, and `publish-document`.

## Envelope

Messages are carried as:

`grid([42(pCID), payload, proof])`

Slot `0` names this exact spec by pCID. Slot `1` carries the payload as a CBOR
item. Slot `2` carries the signing proof item.

## Payload

The current payload is a CBOR map with these fields:

- `kind`: currently `metadata`
- `document_id`: logical document ID
- `author`: stable relay key ID derived from the Ed25519 public key
- `participant_id`: local participant asking the relay to update metadata
- `title`: optional compatibility title hint
- `description`: human-written document description
- `summary`: short summary used for search/result display
- `tags`: ordered list of document tags
- `collections`: ordered list of collection or folder labels
- `favorite`: boolean pin/favorite state
- `archived`: boolean archive state
- `updated_at`: RFC3339Nano metadata update timestamp
- `lamport`: relay-local logical clock used for deterministic latest-state
  ordering
- `embodiment`: optional embodiment hint such as `browser` or `nvim`

## Semantics

`document-metadata` uses latest-state semantics per document in this slice.

That means:

- every metadata update is a signed current-time action
- relays keep the latest accepted metadata state for each document
- search and document-management views operate over that latest state

This protocol does not define restore semantics and does not replace the live
editing history. Source: `DI-loruk`; `DI-sukip`.

## Verification

- reject malformed CBOR
- reject invalid signature proof
- reject payloads whose slot `0` pCID does not match this spec
- reject payloads whose `author` does not match the proof key
- reject oversize description, summary, tags, or collections
