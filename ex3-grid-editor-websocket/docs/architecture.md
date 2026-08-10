# grid-editor architecture

`grid-editor` centers the PromiseGrid-facing contract in one place and keeps
the embodiment plumbing local.

## Current local-draft protocol inventory

These pCIDs are derived from the exact local profile bytes in `protocols/`.
They identify Ex3's current example contracts, not frozen upstream PromiseGrid
specifications or automatic independent-peer interoperability claims. The
complete implemented scope and non-claims are recorded in the
[implementation scope declaration](../CHANGELOG.md). Source: `DI-hadil`.

| Local draft profile | Current pCID | Canonical local source | Implemented role |
| --- | --- | --- |
| `live-document` | `bafkreidoqtzj76jlzzbytfrtni655k6zk4ymo2zuhgwnzqfiz6wtrgzwje` | [`protocols/live-document.md`](../protocols/live-document.md) | Relay-signed CRDT change carriage and embodiment-local replica convergence. |
| `live-awareness` | `bafkreidowbo76hmjrcqfa6l5ahmftlejqtjuzqgzusqjrxjnkngksntsl4` | [`protocols/live-awareness.md`](../protocols/live-awareness.md) | Relay-signed latest-state presence carriage and local presentation. |
| `document-metadata` | `bafkreih7sut3exri37qperlyvohk74vmreppzmcoxenjzozj6zkblanday` | [`protocols/document-metadata.md`](../protocols/document-metadata.md) | Relay-signed current metadata and catalog search. |
| `publish-document` | `bafkreibhvpjr5uddw5z5qmkkowfucawuvzr7hafvsc5djetkax3s6amt3a` | [`protocols/publish-document.md`](../protocols/publish-document.md) | Relay-signed publish manifests and CAS-backed local import/exchange resolution. |

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

For a non-loopback session, a browser or Neovim embodiment can bootstrap with
the relay and present its short-lived, relay-signed capability for the selected
document and profile mutation path. This is local relay admission mechanics;
it does not make the embodiment a separate cryptographic identity or add a
new peer-visible protocol. Source: `DI-hadil`; `DI-povip`.

## Public versus internal boundaries

Public, PromiseGrid-facing:

- `protocols/live-document.md`
- `protocols/live-awareness.md`
- `protocols/document-metadata.md`
- `protocols/publish-document.md`
- signed `grid([42(pCID), payload, proof])` envelopes

Internal-only:

- local HTTP endpoints
- WebSocket carriage for live document and awareness traffic
- bootstrap-secret and short-lived remote mutation capability handling
- browser polling loop and local UI state
- Neovim helper transport and `vim.system` plumbing
- browser-local review/workflow registry

## Remote mutation admission

The operator may configure a relay-local bootstrap secret. For remote HTTP or
WebSocket mutation, the relay uses it only to issue short-lived signed
capabilities scoped to an audience, document ID, selected profile pCID, and
`mutate` action; the relay verifies those claims and expiry on use. Loopback
mutation remains a local fast path.

This mechanism is not a fifth profile, a frozen PromiseGrid application-auth
API, a general person-identity system, or a cross-relay delegation/role
protocol. The relay key is the current app-node identity; participant and
capability audience fields are scoped local inputs. Source: `DI-hadil`;
`DI-povip`.

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
round-trips through CRDT sync, export, and publish flows. The browser editor
now renders that underline inline with inherited text color instead of
link-like styling.

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
