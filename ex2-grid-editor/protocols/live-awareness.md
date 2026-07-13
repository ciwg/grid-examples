# live-awareness

Status: repo-local draft for `grid-editor`

## Purpose

`live-awareness` carries signed latest-state presence updates between
`grid-editor` services.

## Envelope

Messages are carried as:

`grid([42(pCID), payload, proof])`

Slot `0` names this exact spec by pCID. Slot `1` carries the payload as a CBOR
item. Slot `2` carries the signing proof item.

## Payload

The current payload is a CBOR map with these fields:

- `kind`: currently `state`
- `document_id`: logical document ID
- `author`: stable author key ID derived from the Ed25519 public key
- `display_name`: presentation label
- `color`: presentation color
- `cursor`: primary cursor offset
- `head`: optional selection head offset
- `typing`: transient typing hint
- `lamport`: per-author logical clock used for deterministic latest-state choice
- `embodiment`: optional embodiment hint such as `browser` or `nvim`

## Projection

Receivers keep the latest accepted state per `author` using the same ordering
rule as `live-document`:

1. larger `lamport`
2. lexical `author`
3. lexical accepted `message_cid`

Awareness is a latest-state projection rather than a durable shared document
truth.

## Verification

- reject malformed CBOR
- reject invalid signature proof
- reject payloads whose slot `0` pCID does not match this spec

