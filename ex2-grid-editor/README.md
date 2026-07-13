# grid-editor

`grid-editor` is a PromiseGrid example app for collaborative document editing.
The current slice combines:

- a non-canonical Go relay that signs and relays CRDT messages
- CAS-backed persistence for signed envelope bytes
- a browser editor surface with CodeMirror and Automerge
- transitional Neovim compatibility code
- repo-local draft PromiseGrid-facing live protocols for document sync and awareness

The local service is embodiment plumbing. The peer-visible contract lives in:

- `protocols/live-document.md`
- `protocols/live-awareness.md`

Both browser and Neovim are embodiments of the same app contract. Source:
`DI-lodug`; `DI-tofug`; `DI-jilin`.

## Current shape

The current CRDT slice is centered on:

- signed `grid([42(pCID), payload, proof])` envelopes
- durable `Ed25519` identity persisted locally
- CAS-backed storage under `<data-root>/cas/**`
- append-only local message log kept as transitional debug/monitoring output
- repo-local draft `live-document` and `live-awareness` specs
- a browser-local Automerge replica that syncs through explicit relay HTTP endpoints
- real remote cursor rendering rooted in `collab-awareness`

This is a working example repo, not a frozen upstream PromiseGrid protocol.

## Run

```bash
go run ./cmd/grid-relay --listen 127.0.0.1:7001
```

Open the browser UI:

```text
http://127.0.0.1:7001/?doc=demo
```

To point another local service at the first one for peer sync:

```bash
go run ./cmd/grid-relay --listen 127.0.0.1:7002 --peer http://127.0.0.1:7001
```

## Browser build

The browser source now lives under `web/src/` and is bundled into `web/app.js`
with `esbuild`:

```bash
cd web
npm install
npm run build
```

## Neovim

The Neovim path is still transitional in this slice. The compatibility plugin
can still be loaded from the local plugin directory:

```vim
:set runtimepath+=/home/jj/lab/cswg/grid-examples/ex2-grid-editor/nvim
:GridEditorOpen demo
```

The plugin expects Neovim with `vim.system()` available and a running local
`grid-editor` service at `http://127.0.0.1:7001` by default.

## Tests

```bash
go test ./...
```
