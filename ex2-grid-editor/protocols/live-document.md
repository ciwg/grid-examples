# live-document

Status: repo-local draft for `grid-editor`

## Purpose

`live-document` carries signed document-state updates between `grid-editor`
services.

## Envelope

Messages are carried as:

`grid([42(pCID), payload, proof])`

Slot `0` names this exact spec by pCID. Slot `1` carries the payload as a CBOR
item. Slot `2` carries the signing proof item.

## Payload

The current payload is a CBOR map with these fields:

- `kind`: currently `replace`
- `document_id`: shared logical document ID
- `content`: full current document text
- `content_cid`: CID over the exact UTF-8 content bytes
- `lamport`: per-author logical clock used for deterministic convergence
- `author`: stable author key ID derived from the Ed25519 public key
- `embodiment`: optional local embodiment hint such as `browser` or `nvim`
- `previous_cid`: optional previous accepted message CID known to the sender

## Convergence

Receivers treat the payload as a full-document replacement and project the
current document using deterministic last-writer-wins ordering:

1. larger `lamport`
2. lexical `author`
3. lexical accepted `message_cid`

This is a draft example-app rule, not a claim that upstream PromiseGrid has
frozen the same document model.

## Verification

- reject malformed CBOR
- reject invalid `content_cid`
- reject invalid signature proof
- reject payloads whose slot `0` pCID does not match this spec

