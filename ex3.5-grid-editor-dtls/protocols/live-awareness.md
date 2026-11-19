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

## Presence lifecycle

`live-awareness` is for live presence, not historical membership. The current
repo policy is to support two presence profiles:

- `demo`:
  - `0-5 minutes`: treat a peer as live
  - `5-15 minutes`: treat a peer as stale or dimmed
  - `15-30 minutes`: treat a peer as offline
  - `30+ minutes`: remove the peer from the main live `Peers` list
- `normal`:
  - `0-1 minute`: treat a peer as live
  - `1-5 minutes`: treat a peer as stale or dimmed
  - `5-15 minutes`: treat a peer as offline
  - `15+ minutes`: remove the peer from the main live `Peers` list

The demo profile exists because multi-machine demos, restarts, and walking
between embodiments take much longer than strict real-time collaboration
timeouts. The normal profile remains the tighter everyday collaboration model.
Source: `DI-mivor`.

## Live presence versus historical activity

Historical collaboration signals should not be overloaded into the main
`Peers` list.

- document activity belongs in a durable activity stream or change history
- comments belong in a comment or annotation surface
- version history belongs in document/history tooling
- `last viewed` and `last edited` belong in separate historical or audit views

The main live `Peers` list should answer "who is here now?" rather than "who
has ever been here?" Source: `DI-mivor`.

## Verification

- reject malformed CBOR
- reject invalid signature proof
- reject payloads whose slot `0` pCID does not match this spec
