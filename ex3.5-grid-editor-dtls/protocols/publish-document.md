# publish-document

Status: repo-local draft for `grid-editor`

## Purpose

`publish-document` carries signed current-time publish manifests for durable
document handoff and later import/exchange.

`publish-document` is not the live editing protocol. Live editing still uses
`live-document`. Publish is a separate action that references a chosen current
state or saved version without rewriting the past.

## Envelope

Messages are carried as:

`grid([42(pCID), payload, proof])`

Slot `0` names this exact spec by pCID. Slot `1` carries the payload as a CBOR
item. Slot `2` carries the signing proof item.

## Payload

The current payload is a CBOR map with these fields:

- `kind`: currently `publish`
- `document_id`: logical source document ID
- `author`: stable relay key ID derived from the Ed25519 public key
- `participant_id`: local participant asking the relay to publish
- `source_kind`: either `current` or `saved_version`
- `source_version_id`: optional saved-version ID when `source_kind` is
  `saved_version`
- `source_version_name`: optional saved-version label
- `title`: publish title
- `summary`: short publish summary
- `text_cid`: CAS address for the published markdown/text bytes
- `replica_cid`: CAS address for the published Automerge replica bytes
- `published_at`: RFC3339Nano publish timestamp
- `lamport`: relay-local logical clock used for stable append ordering
- `embodiment`: optional embodiment hint such as `browser` or `nvim`

## Exchange rule

The signed publish manifest is the durable PromiseGrid-facing object for this
slice. Import/exchange consumers resolve the manifest, then fetch the
referenced bytes from CAS-backed relay endpoints.

In this slice, import may materialize a new local document from the published
artifact without preserving the exact prior live relay history. That keeps
publish/import separate from restore semantics. Source: `DI-tavul`;
`DI-gosaf`.

## Verification

- reject malformed CBOR
- reject invalid signature proof
- reject payloads whose slot `0` pCID does not match this spec
- reject unsupported `source_kind`
- reject manifests whose referenced CAS objects are missing when resolved
