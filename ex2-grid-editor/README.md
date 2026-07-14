# grid-editor

`grid-editor` is a PromiseGrid example app for shared document editing.
It shows how Grid Editor can be embodied in both a browser and Neovim
while still speaking explicit grid-facing protocols instead of hiding the
contract inside UI code or editor-specific glue. Source: `DI-lodug`;
`DI-tofug`; `DI-jilin`; `DI-ramuv`; `DI-zegov`; `DI-lumek`; `DI-samuv`.

This repo is not just "an editor that happens to sync." It is an example of a
grid-shaped tool:

- the peer-visible contract is identified by protocol CIDs (`pCID`s)
- messages are signed envelopes, not anonymous websocket frames
- document collaboration and awareness are split into separate protocol
  families
- signed objects are persisted in content-addressed storage
- browser and Neovim are embodiments of the same app contract

## What This Example Shows

Grid Editor is built around three layers:

1. a local Go relay that signs, verifies, persists, and relays grid messages
2. embodiment-local CRDT replicas that own collaborative editing convergence
3. browser and Neovim UIs that present the shared document to humans

The current browser path uses CodeMirror plus Automerge. The current Neovim
path uses the repo-local plugin plus a sidecar helper. In both cases, the app
is meant to demonstrate how a PromiseGrid-style tool can separate:

- durable document collaboration
- ephemeral presence and cursor awareness
- local embodiment plumbing

## Current Feature Set

Grid Editor currently demonstrates four feature slices:

### Phase 1: shared editing UX

- remote cursors and selections
- peer list, peer count, and presence aging
- settings, theme, line numbers, and accessibility controls
- search and quick formatting
- new/open/paste-link document entry flows

### Phase 2: document workflow

- title and local document metadata
- recent docs and open tabs
- markdown preview and split view
- find/replace and go-to-line
- import/export, snapshots, bookmarks, and share-link helpers

### Phase 3: review and history

- inline comments and annotations
- saved versions
- recent participant history
- activity feed
- outline navigation, focus mode, summary, and diagnostics

### Phase 4: publish, import, and document metadata

- publish the current state or a named saved version
- relay-signed publish manifests
- CAS-backed text and replica objects
- published exchange URL resolution
- import a published exchange as a new local document
- relay-signed document metadata
- document description, summary, tags, collections, favorites, and archive
  state
- relay-backed catalog search over known document metadata

## How It Uses Grid Technology

Grid Editor is "grid-based" in a specific way.

### 1. Protocol selection is explicit

This repo carries two draft, repo-local protocol specs:

- [live-document](protocols/live-document.md)
- [live-awareness](protocols/live-awareness.md)
- [document-metadata](protocols/document-metadata.md)
- [publish-document](protocols/publish-document.md)

Each protocol is identified by the content hash of its own spec document. That
`pCID` acts as the selector for the wire meaning being used.

### 2. Messages are signed envelopes

Peer-visible messages are carried as:

`grid([42(pCID), payload, proof])`

In this example app:

- slot `0` identifies the exact protocol spec
- slot `1` carries the protocol payload
- slot `2` carries the signing proof

That means the browser and Neovim are not inventing their own public wire
contracts. They both send and receive messages through the same signed,
protocol-addressed grid envelope model. Source: `DI-tofug`; `DI-ramuv`.

### 3. Document sync and awareness are different things

`live-document` carries durable Automerge change packets for the shared
document. `live-awareness` carries human-facing presence state such as display
name, color, and cursor position. Keeping them separate makes the example
closer to real collaborative systems, where document truth and cursor presence
have different durability and cadence requirements. Source: `DI-tofug`.

### 4. The relay is not the editor of record

The Go relay verifies signatures, persists signed bytes, and relays change
history, but it does not own the canonical merged document text. Convergence
happens in embodiment-local CRDT replicas. That is important for this example:
the shared contract lives at the grid/protocol boundary, not inside one
authoritative app server. Source: `DI-ramuv`; `DI-lumek`; `DI-larok`.

### 5. Metadata is separate from live editing

`document-metadata` carries relay-signed latest-state document labels and
descriptions. That keeps document-management features such as favorites,
archive state, and catalog search durable and shareable without pretending
they are part of the live CRDT text stream. Source: `DI-loruk`; `DI-sukip`.

### 6. Publish/import is separate from live editing

`publish-document` is used for durable handoff artifacts, not for typing.
That means a publish manifest can point at a chosen current state or saved
version without pretending it is the same thing as joining the live sync
stream. Source: `DI-tavul`; `DI-gosaf`.

## Where Data Lives

Grid Editor intentionally splits data by role.

### Relay and CAS

The Go relay persists durable signed objects under its data root:

- relay identity seed
- append-only relay log
- CAS-backed signed envelopes
- CAS-backed signed metadata envelopes
- CAS-backed published text bytes
- CAS-backed published Automerge replica bytes

These are relay-owned and are the durable PromiseGrid-facing artifacts in the
current slices.

### Browser-local workflow and review state

Some workflow and review features are still local browser state in this repo:

- recent docs and tabs
- local timestamps
- bookmarks and local snapshots
- comments and reactions
- saved versions used by the current publish flow
- recent participant history and activity feed
- local preferences

Relay-backed document metadata now covers:

- title
- description
- summary
- tags
- collections
- favorite
- archived

Those values are relay-owned and search-visible in Phase 4, while the
remaining workflow/review items above are still local browser state. Source:
`DI-dovoz`; `DI-safor`; `DI-gosaf`; `DI-loruk`; `DI-sukip`.

## PromiseGrid Setup References

This repo was shaped against the PromiseGrid dev guide and related internal
workspace materials, but the README should stay focused on the public entry
points inside this example app.

Use this README as the editor-specific entry point, then follow the repo-local
docs below for architecture, protocols, browser usage, and Neovim usage.

Useful docs:

- [Browser UI walkthrough](docs/grid-editor-ui-example.md)
- [Architecture overview](docs/architecture.md)
- [Practical implementation notes](docs/practical-implementation.md)
- [Docker simulation guide](docs/docker-simulation.md)

## Quick Start

Start the local relay:

```bash
go run ./cmd/grid-relay --listen 127.0.0.1:7015
```

Then pick an embodiment:

- browser: open `http://127.0.0.1:7015/?doc=demo`
- Neovim: run `./scripts/grid-editor-nvim demo`

Local browser and Neovim mutation endpoints are loopback-only in this slice.
For multi-machine collaboration, run one relay per machine and connect the
relays with `--peer` instead of pointing remote editors at one shared relay.

If you want a second relay to poll the first one for peer exchange:

```bash
go run ./cmd/grid-relay --listen 127.0.0.1:7016 --peer http://127.0.0.1:7015
```

## Browser Version

The browser embodiment is the easiest way to see the current CRDT slice in
action.

Run the relay, then open:

```text
http://127.0.0.1:7015/?doc=demo
```

Browser build notes:

```bash
cd web
npm install
npm run build
```

The browser source lives under `web/src/` and is bundled into `web/app.js`
with `esbuild`. Source: `DI-zegov`.

## Docker simulation

If you want to simulate two separate relay machines quickly on Linux, use the
repo-local Docker setup:

```bash
docker-compose up --build
```

Then open:

- `http://127.0.0.1:7015/?doc=demo`
- `http://127.0.0.1:7016/?doc=demo`

See [Docker simulation guide](docs/docker-simulation.md) for the short
workflow and caveats.

Phase 2 browser workflow surfaces now include:

- local document title/metadata and recent-doc tracking
- template and sample-doc creation
- markdown preview pane and split view
- find/replace and go-to-line tools
- import, export, snapshot, and audit-report actions
- copy/share link flows and bookmark support

Phase 3 browser review surfaces now include:

- inline comments and annotation ranges
- saved versions and local review history
- activity feed and recent participant history
- outline navigation and focus mode
- summary and diagnostics overlays

Phase 4 browser exchange surfaces now include:

- relay-signed publish manifests for either the current state or a named saved
  version

## Current demo caveats

- Browser underline currently stores raw `<u>...</u>` markup correctly, but
  the visible underline rendering in the browser editor is still under active
  polish for the demo.
- `Preview` opens the markdown preview pane below the editor.
- `Split View` shows the editor and preview together.
- import/exchange from a published manifest URL
- a published-exchanges list for the current document
- a separate `publish-document` protocol instead of overloading
  `live-document`

In this slice, the publish manifest plus its referenced text and replica bytes
are durable CAS-backed relay objects, while the browser still chooses named
saved versions from its local review metadata. Source: `DI-tavul`; `DI-gosaf`.

Phase 4 browser metadata surfaces now include:

- relay-backed document description and summary
- tags and collections
- favorite and archived flags
- catalog search across relay-known document metadata
- a separate `document-metadata` protocol instead of folding document
  management into `live-document`

This slice keeps document management durable and shareable without treating it
as live typing traffic. Source: `DI-loruk`; `DI-sukip`.

## Neovim Version

The easiest Neovim path is the launcher script:

```bash
./scripts/grid-editor-nvim demo
```

Optional environment overrides:

```bash
GRID_EDITOR_RELAY_URL=http://127.0.0.1:7015 \
GRID_EDITOR_DISPLAY_NAME="Alice" \
GRID_EDITOR_COLOR="#1d6fd6" \
./scripts/grid-editor-nvim demo
```

Inside Neovim, the main commands are:

```vim
:GridEditorInfo
:GridEditorPeers
:GridEditorClose
```

If you want the manual path, load the repo-local plugin yourself:

```vim
:set runtimepath+=/home/jj/lab/cswg/grid-examples/ex2-grid-editor/nvim
:runtime plugin/grid_editor.vim
:GridEditorOpen demo
```

The launcher script and Neovim plugin live here:

- [Neovim launcher](scripts/grid-editor-nvim)
- [Neovim plugin](nvim/lua/grid_editor/init.lua)

Source: `DI-samuv`; `DI-gafit`.

## Important Docs

Use these as the main reading path after this README:

- [Architecture overview](docs/architecture.md)
- [Browser UI example](docs/grid-editor-ui-example.md)
- [live-document protocol](protocols/live-document.md)
- [live-awareness protocol](protocols/live-awareness.md)
- [document-metadata protocol](protocols/document-metadata.md)
- [publish-document protocol](protocols/publish-document.md)
- [CRDT relay thought experiment](docs/thought-experiments/TE-satuf-grid-editor-crdt-relay-slice.md)
- [Neovim sidecar thought experiment](docs/thought-experiments/TE-zorud-grid-editor-nvim-sidecar-hybrid-helper.md)
- [Publish/import thought experiment](docs/thought-experiments/TE-vafor-grid-editor-publish-exchange-slice.md)
- [Document metadata thought experiment](docs/thought-experiments/TE-mifud-grid-editor-document-metadata-slice.md)
- [Practical implementation notes](docs/practical-implementation.md)

## Tests

```bash
go test ./...
```
