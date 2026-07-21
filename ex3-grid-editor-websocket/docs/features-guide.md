# Ex3 Features Guide

`ex3-grid-editor-websocket` is the websocket-backed Grid Editor example.
This guide explains what features exist on the page, what they do, and how
each feature is enabled in the current PromiseGrid-shaped design.

## Feature Model

This example is easiest to read in three buckets:

1. **Embodiment UI features**
   - browser or Neovim behavior presented to the user
2. **Relay-backed PromiseGrid features**
   - signed envelopes, pCID-selected protocol meaning, publish/metadata flows
3. **Browser-local workflow features**
   - convenience surfaces that are useful for the demo but not yet durable
     shared PromiseGrid artifacts

The key distinction is that `ex3` is not claiming every visible UI element is
already a stable shared protocol. Some features are intentionally local shells
over a smaller set of relay-backed protocol flows.

## Neovim Embodiment Notes

Current Neovim peer rendering is intentionally split into two visible pieces:

- a colored sign plus highlighted cursor cell at the exact remote position
- a peer name label rendered at the **end of the same line**

Why this shape exists:

- earlier overlay labels sat directly on top of document text and were hard to
  read in live demos
- the current shape keeps the remote location visible without hiding the file
  contents underneath it

What this means:

- Neovim should show the same document content as the browser
- browser peers should be visible in Neovim without opening `:GridEditorPeers`
- the full roster still lives in `:GridEditorPeers`, but the in-buffer marker
  is optimized for active remote cursor visibility rather than roster display

## Page Sections And Features

### Document

Features:
- open a document by id
- rename the current document title
- create a new shared document
- duplicate the current document into a new local/shared doc id
- paste, copy, and email a tokenized share link
- show created/viewed/edited/exported timestamps

PromiseGrid enablement:
- the active document is joined through the relay and bound to the live
  document and awareness protocols
- the share link carries the document id and demo bootstrap token used to
  mint short-lived remote capabilities
- document title can be saved durably through the separate
  `document-metadata` protocol
- browser startup now primes from the relay's current document snapshot before
  websocket catch-up begins, so an old local demo seed cannot impersonate an
  existing shared doc
- once a browser catches up to an older pre-snapshot doc, it backfills a relay
  snapshot so later joiners can start from the current document state

Current persistence shape:
- title and relay metadata can be relay-backed
- recent timestamps and duplicate/new-doc workflow are currently maintained
  in browser-local registry state

### You

Features:
- choose display name
- choose presence color
- open settings
- inspect the local participant id

PromiseGrid enablement:
- display name and color flow over the `live-awareness` protocol as signed
  awareness state
- the participant id is the relay/client identity used to label presence and
  message authorship

Current persistence shape:
- browser preferences are local
- peer-visible awareness state is live relay traffic

### Workspace

Features:
- open tabs list
- recent docs list
- starter templates
- generate demo document
- sample document

PromiseGrid enablement:
- these are embodiment workflow helpers around document ids and seeded content
- when a generated/template doc is opened, it still joins the same signed
  relay-backed live document flow

Current persistence shape:
- local browser registry only
- local seed content is only applied when a doc has no relay history yet; once
  the relay has a snapshot or message history, the relay state wins

### PromiseGrid Flow

Features:
- transport badge showing browser sync and awareness transport modes
- architecture/data-flow strip:
  - Browser
  - signed grid message
  - Relay
  - peer feed
  - Peer relay
  - websocket fanout
  - Other editor
- live message trace
- clickable message inspection via PromiseGrid Inspector

Layout note:
- this card now sits at the top of the **left sidebar**
- it is **not** at the literal top of the whole page, because the page header
  still sits above the sidebar cards

PromiseGrid enablement:
- this section is backed by relay-observed signed envelopes, not synthetic UI
  events
- document and awareness traffic are observed after protocol ingest and shown
  back to the page through a trace endpoint
- the browser transport badge verifies that the browser is using websocket
  transport to the relay while relay-to-relay propagation still goes through
  the signed peer-message feed

Current persistence shape:
- live relay-observed traffic only
- this is a demo/inspection surface, not a durable store of its own

### Metadata

Features:
- description
- summary
- tags
- collections
- favorite
- archived
- save metadata

PromiseGrid enablement:
- saved through the dedicated `document-metadata` protocol
- metadata search uses relay-backed metadata records rather than live editor
  text

Current persistence shape:
- relay-backed and durable

### Relay

Features:
- local author / relay id
- active pCIDs for:
  - `live-document`
  - `live-awareness`
  - `document-metadata`
  - `publish-document`
- presence profile indicator

PromiseGrid enablement:
- this is the clearest direct proof that the demo is protocol-addressed
  rather than hiding behavior in unnamed app routes
- pCIDs expose which protocol documents are selecting the message meanings in
  the current run

Current persistence shape:
- relay runtime state

### Peers

Features:
- peer list
- presence aging legend
- presence states:
  - live
  - stale
  - offline

PromiseGrid enablement:
- driven by `live-awareness` state rather than text sync itself
- cursor position, name, color, and timing are distinct from document changes

Current persistence shape:
- live relay traffic plus local display logic for presence aging windows

### Review

Features:
- outline
- saved versions
- recent participants
- activity
- published exchanges
- catalog search

PromiseGrid enablement:
- published exchanges are relay-backed publish records
- catalog search uses relay-backed metadata
- outline/history/comments versions are currently local review helpers around
  the live document shell

Current persistence shape:
- mixed
- published exchanges and metadata search are relay-backed
- comments, saved versions, bookmarks, and most activity history are local

## Toolbar Features

### Editing Actions

Features:
- Bold
- Italic
- Underline
- Find / Replace
- Preview
- Split View
- Import
- Export / Exchange
- Snapshot
- Bookmark
- Comment
- Save Version
- Summary
- Focus
- Inspect

PromiseGrid enablement:
- live typing itself travels through `live-document`
- remote name/color/cursor state travels through `live-awareness`
- export/publish features use `publish-document`
- metadata save/search uses `document-metadata`

Current persistence shape:
- formatting and typing are relay-backed through the live doc path
- comments, bookmarks, snapshots, and local saved versions are browser-local
- preview and focus mode are embodiment-local

### Status Cluster

Features:
- connection state
- auto-save state
- message count
- local replica cid

PromiseGrid enablement:
- connection state reflects relay and websocket session status
- message count reflects relay state endpoint data
- local replica cid reflects the browser’s current Automerge replica snapshot

Current persistence shape:
- live runtime display only

## Hidden Panels And Their Meaning

### Settings

Features:
- theme
- line numbers
- font size
- dyslexia-friendly spacing
- presence profile
- editable shortcuts

PromiseGrid enablement:
- none required for durability
- affects how the embodiment presents the same shared relay flows
- browser-local preferences and workflow state now fall back to in-memory
  storage if local/session storage are blocked, so private browser sessions do
  not lose core editor startup behavior

Private/incognito note:
- the browser startup path is hardened against blocked storage and late-join
  text loss
- a final real private-browser manual verification pass is still tracked
  separately in TODO 016

### Help

Features:
- shortcut help

PromiseGrid enablement:
- none

### Find / Replace

Features:
- find
- replace all
- case sensitive
- regex
- go to line
- a short helper sentence describing the current workflow

PromiseGrid enablement:
- local editor operation
- resulting document edits still enter the live document protocol normally

### Export / Exchange

Features:
- export as markdown
- export as html
- export as plain text
- export as Automerge bytes
- copy markdown
- copy html
- publish snapshot
- publish exchange
- import exchange
- audit report

PromiseGrid enablement:
- publish/import uses the `publish-document` protocol
- exports are embodiment-level artifact generation over the current document
  state
- Automerge export is the local replica artifact, not the PromiseGrid wire
  protocol itself

### Comments

Features:
- selected-text comment capture
- comment list
- resolve selected comments

PromiseGrid enablement:
- currently local review behavior
- comments are not yet a separate shared protocol family in `ex3`

### Document Summary

Features:
- generated summary text
- read aloud
- voice input support for comment dictation

PromiseGrid enablement:
- none directly
- derived from the current local embodiment text

### PromiseGrid Inspector

Features:
- inspect current transport state
- inspect selected trace message
- inspect local comments/activity/participants/versions summary

PromiseGrid enablement:
- selected message data comes from relay-observed signed traffic
- transport section proves the page is using relay websocket transport

Why it exists:
- it turns the demo from “a collaborative editor that happens to work” into a
  visibly explainable grid example
- it gives viewers a concrete answer to:
  - what actually moved over the relay
  - which protocol family the message belonged to
  - which participant authored it
  - whether the browser is using websocket transport to the relay

What the inspector JSON means:

- `documentID`
  - the currently open shared document id
- `browser_transport`
  - the current browser transport status for the live page
- `browser_transport.sync`
  - the transport used for live document sync
- `browser_transport.awareness`
  - the transport used for presence/cursor awareness
- `browser_transport.relay_path`
  - a boolean marker showing the browser is still going through the relay path
    rather than using a direct browser-to-browser shortcut
- `selected_message`
  - the currently clicked relay-observed PromiseGrid message

Important `selected_message` fields:

- `offset`
  - the relay log offset where this message was observed
- `envelope_cid`
  - the CID of the signed outer envelope
- `protocol`
  - the human-readable protocol family, such as `live-awareness`
- `pcid`
  - the protocol CID selecting the message meaning
- `kind`
  - the message kind inside that protocol family
- `document_id`
  - which shared document this message belongs to
- `participant_id`
  - which participant emitted the message
- `author`
  - the signing author key id used on the envelope
- `embodiment`
  - which UI embodiment emitted the message, such as `browser`
- `lamport`
  - the Lamport clock value carried by the message
- `received_at`
  - when the local relay observed the message
- `summary`
  - short human-readable explanation of what happened
- `envelope_base64`
  - raw signed envelope bytes encoded as base64
- `payload_base64`
  - raw protocol payload bytes encoded as base64
- `proof_algorithm`
  - signing algorithm used on the proof
- `proof_key_id`
  - key id associated with the proof
- `decoded_payload`
  - the decoded protocol payload, shown as fields instead of raw CBOR only

How to read the example awareness payload:
- `protocol: live-awareness`
  - this is presence/cursor traffic, not durable text edits
- `kind: state`
  - this is a participant state update
- `cursor` and `head`
  - the remote selection/cursor positions
- `typing`
  - whether the participant was marked as actively typing
- `display_name`
  - the peer-facing display label
- `color`
  - the peer-facing presence color

## Feature Reliability Notes

These are the features that are genuinely PromiseGrid-shaped in the current
example:

- signed live document traffic
- signed live awareness traffic
- relay-backed metadata
- relay-backed publish/import
- visible pCID selection
- relay-observed message trace
- browser and Neovim embodiments over the same relay contract

These are still mostly embodiment-local demo/workflow helpers:

- recent docs and tab registry
- bookmarks
- local snapshots
- comments and reactions
- saved versions
- recent participant history cache
- activity feed aggregation
- most preferences

That split is intentional. `ex3` is demonstrating how a PromiseGrid-facing
protocol core can coexist with a richer local UX shell without pretending that
every convenience feature is already a hardened shared protocol.

## What “Enabled On PromiseGrid” Means In Ex3

In this repo, a feature is effectively “enabled on PromiseGrid” when at least
one of these is true:

- it is carried as a signed envelope with pCID-selected meaning
- it is persisted or resolved by the relay as a durable signed artifact
- it is surfaced from relay-observed protocol traffic
- it is shared across browser and Neovim embodiments without a private
  embodiment-only contract

That is true today for:

- live collaborative text
- live awareness and peer presence
- relay metadata
- publish/import handoff artifacts
- transport/trace inspection

It is not yet fully true for:

- comments
- bookmarks
- snapshots
- local saved versions
- local activity/history summaries

Those remain useful features, but they should be described honestly as local
workflow surfaces built around the PromiseGrid-backed editing core.
