# live-document

Status: repo-local draft for `grid-editor`

## Purpose

`live-document` carries signed Automerge change packets between `grid-editor`
relays and embodiments.

## Envelope

Messages are carried as:

`grid([42(pCID), payload, proof])`

Slot `0` names this exact spec by pCID. Slot `1` carries the payload as a CBOR
item. Slot `2` carries the signing proof item.

## Payload

The current payload is a CBOR map with these fields:

- `kind`: currently `change`
- `document_id`: shared logical document ID
- `author`: stable relay key ID derived from the Ed25519 public key
- `participant_id`: local replica participant sending the message
- `recipient_id`: currently empty in this repo-local draft and reserved for
  future directed delivery if a later TE/DI needs it
- `change_bytes`: raw Automerge change bytes
- `lamport`: relay-local logical clock used for stable append ordering
- `embodiment`: optional local embodiment hint such as `browser` or `nvim`

## Convergence

The relay does not project canonical document text. It verifies and persists
signed change envelopes, then forwards them to embodiment-local Automerge
replicas. Late joiners reconstruct a document by replaying the append-only
change history from the relay log. Convergence happens in the CRDT replicas,
not in the relay.

This is a repo-local example-app rule, not a claim that upstream PromiseGrid
has frozen the same live-CRDT model. Source: `DI-ramuv`; `DI-lumek`;
`DI-larok`.

## Verification

- reject malformed CBOR
- reject invalid signature proof
- reject payloads whose slot `0` pCID does not match this spec
- persist the signed canonical envelope bytes to CAS addressed by the signed
  object hash
