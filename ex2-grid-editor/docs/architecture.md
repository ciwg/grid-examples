# grid-editor architecture

`grid-editor` centers the PromiseGrid-facing contract in one place and keeps
the embodiment plumbing local.

## Topology

```text
Browser UI -----------------------------\
                                         \
Browser-local Automerge replica ----------> local HTTP adapter -> grid-relay -> peer relays
                                         /
Neovim plugin -> nvim sidecar ----------/
                    local Automerge replica
```

The relay owns:

- `Ed25519` identity
- repo-local pCID discovery from exact local spec bytes
- signed grid envelope creation and verification
- append-only message logging
- relay-visible sync and awareness feed projection
- publish manifest signing and resolution
- CAS object persistence
- optional peer polling

The browser UI and the Neovim sidecar own:

- local Automerge replicas
- local editing UX
- local cursor and selection wiring
- local HTTP calls into the service

They do not define the peer-visible protocol truth. Source: `DI-lodug`;
`DI-tofug`; `DI-jilin`; `DI-ramuv`; `DI-tavul`.

## Public versus internal boundaries

Public, PromiseGrid-facing:

- `protocols/live-document.md`
- `protocols/live-awareness.md`
- `protocols/document-metadata.md`
- `protocols/publish-document.md`
- signed `grid([42(pCID), payload, proof])` envelopes

Internal-only:

- local HTTP endpoints
- browser polling loop and local UI state
- Neovim helper transport and `vim.system` plumbing
- browser-local review/workflow registry

## Protocol roles

### `live-document`

- carries Automerge change packets for collaborative editing
- drives browser and Neovim replica convergence
- is durable and replayable

### `live-awareness`

- carries cursor, selection, display name, color, and typing presence
- is human-facing and ephemeral
- is kept separate from document truth

### `publish-document`

- carries relay-signed publish manifests
- references CAS-backed text and replica objects
- supports current-state or saved-version handoff
- is separate from restore semantics and separate from live sync

### `document-metadata`

- carries relay-signed latest-state document metadata
- covers title, description, summary, tags, collections, favorite, and
  archived state
- powers relay-backed catalog search over known documents
- is separate from both live CRDT typing and publish/import exchange

Source: `DI-tofug`; `DI-ramuv`; `DI-tavul`; `DI-gosaf`; `DI-loruk`;
`DI-sukip`.

## Storage model

### Durable relay storage

The relay data root stores:

- the relay signing identity
- the append-only message log
- CAS-backed signed envelopes
- CAS-backed signed metadata envelopes
- CAS-backed published text bytes
- CAS-backed published replica bytes

These are the durable artifacts the relay can verify and serve back later.

### Embodiment-local state

The browser currently keeps some product-facing metadata locally:

- preferences
- recent docs
- local timestamps, bookmarks, and snapshots
- comments and review metadata
- saved versions used by the current publish flow

Relay-backed document metadata now covers title, description, summary, tags,
collections, favorites, archive state, and relay-backed catalog search.

The Neovim embodiment keeps its own local editor/session state and relies on
the relay for shared collaboration artifacts.

## Current browser UX note

The browser supports two distinct markdown view modes:

- `Preview`
  - opens the preview pane for the same document
- `Split View`
  - keeps editor and preview visible together

Underline is stored as raw `<u>...</u>` markup in the shared text stream so it
round-trips through CRDT sync, export, and publish flows. The visible browser
underline rendering is still being polished and should not be treated as fully
closed UX work yet.

## Convergence model

Live editing convergence is CRDT-based:

- browser uses Automerge locally
- Neovim sidecar uses Automerge locally
- relay stores and relays signed change traffic but is not the canonical text
  owner

Relay indexing still uses stable append ordering and lamport metadata where
needed for deterministic feed handling, but the shared document truth is no
longer a last-writer-wins text projection. Source: `DI-ramuv`; `DI-lumek`;
`DI-larok`.
