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
- settings, theme, line numbers, accessibility controls, and a visible color
  picker button
- find/replace and quick formatting
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

The browser startup path now primes from the relay's current snapshot before
it considers any browser-local demo seed. That keeps old local sample content
from impersonating an existing shared document during reconnect or late join.

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

### 7. Restore continues; it does not roll back

`restore-published-version` lets an admitted participant select one immutable
publish manifest and create one new, provenance-bearing CRDT change in the
existing document. It does not rewrite live history or make the old publication
authoritative. Concurrent edits can merge with the restored content. Import
remains separate: it creates a new local document from a publish artifact.
Source: `DI-hihok`; `DI-tibum`; `DI-nihiz`.

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

When local or session storage is blocked, the browser now falls back to
in-memory workflow/session state so fresh or private browser sessions do not
lose the shared relay-backed document path just because persistence APIs are
restricted. Source: `DI-ribaf`; `DI-bonuv`.

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

- [Implementation scope declaration](CHANGELOG.md)
- [Current local-draft protocol inventory](docs/architecture.md#current-local-draft-protocol-inventory)
- [Testing guide](docs/testing.md)
- [Browser UI walkthrough](docs/grid-editor-ui-example.md)
- [Feature guide](docs/features-guide.md)
- [Architecture overview](docs/architecture.md)
- [Practical implementation notes](docs/practical-implementation.md)
- [Docker simulation guide](docs/docker-simulation.md)

## What You Need To Run

For the basic local relay plus browser flow:

- Go
- a modern browser

Optional tools:

- Node and npm, only if you want to rebuild or test the browser bundle under
  `web/`
- Neovim, only if you want to run the Neovim embodiment
- Docker plus either `docker compose` or `docker-compose`, only if you want the
  two-relay simulation

By default the relay stores runtime data under `.grid-editor/`. If you do not
pass `--data-root`, that directory is created next to this README.

## Quick Start

This copied example defaults to `127.0.0.1:7025` and `127.0.0.1:7026` so it
can run beside `ex2-grid-editor` on the same machine. Source: `DI-vatub`.

Start the local relay:

```bash
go run ./cmd/grid-relay --listen 127.0.0.1:7025 --data-root .grid-editor/alice
```

Then pick an embodiment:

- browser: open `http://127.0.0.1:7025/?doc=demo`
- Neovim: run `./scripts/grid-editor-nvim demo`

Local loopback clients still work with no extra setup. For multi-machine
browser or Neovim collaboration, start the relay with a bootstrap token and
share that token in the document link or sidecar environment:

```bash
go run ./cmd/grid-relay --listen 0.0.0.0:7025 --remote-access-token ex3-demo-access
```

`ex3-demo-access` is a checked-in demo-only value for this example. Operators
must choose their own local bootstrap secret for a real deployment; it remains
relay-local configuration, not a general PromiseGrid identity or auth API.

Then open the browser with:

```text
http://127.0.0.1:7025/?doc=demo&access_token=ex3-demo-access
```

Or launch Neovim with:

```bash
GRID_EDITOR_RELAY_URL=http://127.0.0.1:7025 GRID_EDITOR_ACCESS_TOKEN=ex3-demo-access ./scripts/grid-editor-nvim demo
```

That bootstrap token is only used to mint short-lived relay-signed mutation
capabilities for the current document; steady-state live traffic still uses the
existing `live-document`, `live-awareness`, metadata, and publish surfaces.
Source: `DI-povip`.

If you want a second relay to poll the first one for peer exchange:

```bash
go run ./cmd/grid-relay --listen 127.0.0.1:7026 --data-root .grid-editor/bob \
  --peer http://127.0.0.1:7025
```

## Scope, topology, admission, and evidence

Ex3 implements four source-derived, repo-local draft profiles. They are not
frozen upstream PromiseGrid specifications or automatic interoperability
claims. The exact pCIDs, implemented scope, and explicit non-claims are in the
[implementation scope declaration](CHANGELOG.md).

The normal decentralized collaboration model is one relay per logical node.
Browser and Neovim are local embodiments; their participant fields do not make
them separate cryptographic identities in this slice. The two commands above
can simulate Alice and Bob on one host because they use distinct relay
processes, keys, and data roots. On separate machines, use the same arrangement
with each relay's own local root and reachable peer URL.

For remote mutation, the relay uses its operator-configured bootstrap secret
only to mint short-lived capabilities for the selected document and profile.
WebSocket carries live traffic but does not define a fifth profile or replace
the pCID-selected protocol meanings.

Each relay root contains that relay's durable local artifacts. To inspect a
run, use the root selected by its `--data-root` value, then follow the
[architecture overview](docs/architecture.md) and [testing guide](docs/testing.md)
for accepted replay, observations, admission diagnostics, and verification
semantics. Those relay-local records are not shared proof and do not settle
another participant's intent.

Automated hardening covers likely private/incognito startup problems. Browser-level
verification has also passed with isolated normal and incognito Chrome sessions
against one isolated relay, confirming bidirectional document-text convergence.
The check used local DevTools and native browser input rather than a human-driven
usability review; see the completed
[TODO tamuk](TODO/TODO-tamuk-grid-editor-private-browser-document-sync.md).
Source: `DI-tajin`; `DI-hadil`; `DI-pazis`; `DI-hosit`; `DI-sodoj`.

## Docker Container Start

If you want the two-relay Docker demo instead of running `grid-relay`
directly, use the checked-in `compose.yaml`. The container setup publishes
`7025` and `7026` on the host and enables the demo bootstrap token
`ex3-demo-access`. Source: `DI-vatub`; `DI-povip`.

From this directory:

```bash
docker compose up -d --build
```

If your machine does not have the `docker compose` plugin, use the standalone
binary instead:

```bash
docker-compose up -d --build
```

Check that both relays are up:

```bash
docker compose ps
```

or:

```bash
docker-compose ps
```

Then open either browser URL:

- `http://127.0.0.1:7025/?doc=demo&access_token=ex3-demo-access`
- `http://127.0.0.1:7026/?doc=demo&access_token=ex3-demo-access`

The browser demo now shows a `PromiseGrid Flow` card while you edit. That
surface puts the active browser transport modes on screen, draws the
browser-to-relay data path, and lists the most recent relay-observed signed
messages for the current document. Click any listed message to inspect its
decoded payload plus raw envelope/payload base64 in the `PromiseGrid
Inspector`. Source: `DI-holoz`; `DI-dogub`.

For the current layout and a fuller explanation of the PromiseGrid Flow and
PromiseGrid Inspector surfaces, see:

- [Updated UI notes](docs/updated-ui.md)
- [Feature guide](docs/features-guide.md)

To attach Neovim to the same shared document through relay `7026`:

```bash
GRID_EDITOR_RELAY_URL=http://127.0.0.1:7026 \
GRID_EDITOR_ACCESS_TOKEN=ex3-demo-access \
GRID_EDITOR_DISPLAY_NAME="Neovim" \
GRID_EDITOR_COLOR="#d66f1d" \
./scripts/grid-editor-nvim demo
```

To stop the containers:

```bash
docker compose down
```

or:

```bash
docker-compose down
```

To fully reset the Docker runtime state and bring the demo back up, use the
repo-local helper:

```bash
./scripts/run-clean.sh
```

That script detects either `docker compose` or `docker-compose` on the local
host before resetting and rebuilding the demo stack. Source: `DI-samuv`.

## Browser Version

In `ex3-grid-editor-websocket`, the browser and Neovim live `sync` and
`awareness` flows prefer websocket transport when the runtime supports it,
while metadata and publish remain on the existing HTTP endpoints. Source:
`DI-vubih`; `DI-bitus`.

This transport split is aligned with the current PromiseGrid dev-guide
direction: websocket is only the carriage for live traffic here, not the
protocol meaning, and the existing `live-document` / `live-awareness` /
metadata / publish boundary remains explicit. It is still a repo-local
provisional implementation, not a frozen upstream PromiseGrid app API. Source:
`DI-vipat`; `DI-bitus`.

The latest upstream guide refresh on July 13-14, 2026 kept that caution in
place explicitly: app-facing auth/API guidance remains provisional under
`DR-tuhaz`, and `POC20` semantic-model work is now tracked separately from the
`POC21` DevOps/bootstrap track. This example follows that direction, but does
not claim that its remote bootstrap or capability shape is a frozen upstream
PromiseGrid contract. Source: `DI-talih`; `DI-vipat`.

The browser embodiment is the easiest way to see the current CRDT slice in
action.

Run the relay, then open:

```text
http://127.0.0.1:7025/?doc=demo
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

- `http://127.0.0.1:7025/?doc=demo&access_token=ex3-demo-access`
- `http://127.0.0.1:7026/?doc=demo&access_token=ex3-demo-access`

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

- Browser underline stores raw `<u>...</u>` markup and now renders inline in
  the browser editor using normal text color rather than link-like styling.
- Browser color selection now uses a visible in-app picker button so Chrome
  does not depend on the tiny native color preview line.
- `Find / Replace` names the search tool directly and the overlay explains the
  find/replace/go-to-line flow more plainly.
- Private/incognito browser startup is hardened against blocked storage and
  fresh late-join text loss. If websocket startup still opens blank while the
  relay reports shared history, the browser now forces one bounded HTTP sync
  catch-up before staying empty. If the relay state advertises a stale blank
  snapshot, that recovery path now distrusts the snapshot offset and replays
  from full history instead of skipping the real shared text forever. Isolated
  normal and incognito Chrome sessions now have browser-level bidirectional
  document-sync verification through one relay; the DevTools/native-input
  method is recorded in TODO 016. Source: `DI-sodoj`.
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

The Neovim embodiment needs Neovim in addition to the relay above. For
multi-machine use, also provide `GRID_EDITOR_ACCESS_TOKEN` so the sidecar can
bootstrap a remote session.

The easiest Neovim path is the launcher script:

```bash
./scripts/grid-editor-nvim demo
```

Optional environment overrides:

```bash
GRID_EDITOR_RELAY_URL=http://127.0.0.1:7025 \
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

Current remote-peer rendering in Neovim:

- the shared document text should match the browser
- active remote peers should show as a colored sign/cursor cell at the exact
  position
- the peer name label renders at the **end of that line** instead of on top of
  the document text, to keep the file readable during live demos

If you want the manual path, load the repo-local plugin yourself:

```vim
:set runtimepath+=/home/jj/lab/cswg/grid-examples/ex3-grid-editor-websocket/nvim
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

See the [testing guide](docs/testing.md) for the complete verification command
set and the distinct claims covered by each test layer. Source: `DI-tajin`.
