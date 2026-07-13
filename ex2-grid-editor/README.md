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

## How It Uses Grid Technology

Grid Editor is "grid-based" in a specific way.

### 1. Protocol selection is explicit

This repo carries two draft, repo-local protocol specs:

- [live-document](protocols/live-document.md)
- [live-awareness](protocols/live-awareness.md)

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

## PromiseGrid Setup References

This repo was shaped against the PromiseGrid dev guide and related internal
workspace materials, but the README should stay focused on the public entry
points inside this example app.

Use this README as the editor-specific entry point, then follow the repo-local
docs below for architecture, protocols, browser usage, and Neovim usage.

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

Useful browser docs:

- [Browser UI walkthrough](docs/grid-editor-ui-example.md)
- [Architecture overview](docs/architecture.md)

Browser build notes:

```bash
cd web
npm install
npm run build
```

The browser source lives under `web/src/` and is bundled into `web/app.js`
with `esbuild`. Source: `DI-zegov`.

Phase 2 browser workflow surfaces now include:

- local document title/metadata and recent-doc tracking
- template and sample-doc creation
- markdown preview and split view
- find/replace and go-to-line tools
- import, export, snapshot, and audit-report actions
- copy/share link flows and bookmark support

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
- [CRDT relay thought experiment](docs/thought-experiments/TE-satuf-grid-editor-crdt-relay-slice.md)
- [Neovim sidecar thought experiment](docs/thought-experiments/TE-zorud-grid-editor-nvim-sidecar-hybrid-helper.md)

## Tests

```bash
go test ./...
```
