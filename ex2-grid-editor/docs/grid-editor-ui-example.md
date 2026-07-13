# grid-editor UI example

This document explains what you are looking at when you open the current
browser demo for `grid-editor`.

In plain terms, this screen is a small browser-based view into a local
`grid-editor` service. That service is the real PromiseGrid-facing part of the
example app: it owns the signing identity, the repo-local protocol selection,
the append-only message log, and the current projected document and awareness
state. The browser page is intentionally thin. It lets you choose a document,
type text, and send cursor/presence updates, but it is not the source of truth
for the protocol or the shared state.

So when this page shows values like an `Author`, a `live-document pCID`, or a
`local replica` fingerprint, it is really showing you facts about the local Go
service and the CRDT state around it, not just browser-only UI values. Source:
`DI-lodug`; `DI-tofug`; `DI-jilin`; `DI-zegov`; `DI-larok`.

![grid-editor demo](./demo.png)

## Header

### `PromiseGrid example app`

- The page is a demo UI for the `grid-editor` example repo.
- It talks to the local Go service, which owns the actual document and
  awareness state.

## Phase 1 browser shell

### Toolbar and quick actions

- The browser UI now has a small Phase 1 action bar for `Search`, `Bold`,
  `Italic`, and `Underline`.
- There is also a lightweight welcome banner and a `Help` surface so a demo
  user does not need to read the README first.
- `New Shared Doc` creates a fresh logical document ID and updates the share
  link in the sidebar.
- `Paste Link` accepts either a full grid-editor URL or a raw document ID.

### Settings

- The browser UI now has a local settings surface for:
  - theme
  - line numbers
  - font size
  - dyslexia-friendly spacing
  - presence profile
  - shortcut overrides
- In this Phase 1 slice those preferences stay local to the browser, but they
  are treated as app-level preferences rather than scattered one-off browser
  toggles. Source: `DI-vasul`.

## Document controls

### `Document ID`

- This is the logical name of the shared document.
- Opening the same document ID on the same service, or on synced peer
  services, points at the same document stream.

### `Open`

- This tells the page to switch to the selected document ID.
- After opening, the page polls the local service for the current document and
  awareness state for that document.

### `Display name`

- This is a human-facing presence label.
- It is not the durable identity used to sign messages.

### `Color`

- This is a human-facing presence color used in the UI.
- It is presentation data, not identity data.

## Local service section

### `Author`

- This is the durable local author ID derived from the service's stored
  `Ed25519` public key. Source: `DI-jilin`.
- The service uses this identity when it signs `live-document` and
  `live-awareness` messages.

### `live-document pCID`

- This is the content-addressed ID of the exact local `live-document` spec
  file.
- It identifies the current draft protocol used for signed document-update
  messages. Source: `DI-tofug`.

### `live-awareness pCID`

- This is the content-addressed ID of the exact local `live-awareness` spec
  file.
- It identifies the current draft protocol used for signed awareness messages.
  Source: `DI-tofug`.

## Peers

### `Peers`

- This section lists the currently visible awareness states for the active
  document.
- `No remote peers yet` means no other authors are currently visible for that
  document, or only the local author has written awareness state so far.
- This list is meant to show live presence, not historical participation.
- The browser now renders peer presence using the selected profile:
  - `demo`: `0-5 minutes` live, `5-15 minutes` stale, `15-30 minutes`
    offline, `30+ minutes` removed
  - `normal`: `0-1 minute` live, `1-5 minutes` stale, `5-15 minutes`
    offline, `15+ minutes` removed
- Historical information such as comments, version history, `last viewed`, or
  `last edited` should live in separate surfaces instead of staying in the
  live peer roster. Source: `DI-mivor`.

## Status bar

### `connected`

- The browser page can currently reach the local Go service over the internal
  HTTP adapter.
- In the current Phase 1 slice, this can also show more explicit states such
  as `connecting`, `syncing`, `unsynced local changes`, `disconnected`, or an
  error message when sync fails.

### `messages: 4`

- This is the current relay-visible message count for the active document.
- It is a quick status value showing how many signed document-change records
  the local relay currently has for that document.
- It is not a Git revision, file save count, or direct CRDT version number.

### `local replica: hW9K...`

- This is a short debug fingerprint of the browser's current local Automerge
  replica state.
- It is derived from the serialized local CRDT document bytes and then trimmed
  for display.
- It helps you tell whether the browser replica has moved to a new local CRDT
  state, even when the visible text is similar.
- It is not currently shown as a formal CID in the UI.

## Quick distinction

### `Author` vs `Display name`

- `Author` is the durable signing identity.
- `Display name` is just a presentation label shown to people.

### `pCID` vs `local replica`

- A `pCID` identifies the protocol spec being spoken.
- The `local replica` value is a browser-local fingerprint of the current CRDT
  state snapshot.

### `local replica` vs document text identity

- The `local replica` value describes the whole local CRDT state, not just the
  plain text you can read on screen.
- Two replicas can temporarily show the same text while still having different
  internal CRDT histories or state bytes.
- A plain-text content CID, if shown in some future debugging view, would mean
  "what exact text bytes do I have?", while `local replica` means "what exact
  local CRDT snapshot do I have?"

### Why there are two pCIDs

- `live-document` and `live-awareness` are separate protocol families in this
  repo because document state and awareness state have different cadence and
  durability pressure. Source: `DI-tofug`.
