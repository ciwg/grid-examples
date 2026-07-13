# grid-editor

`grid-editor` is a PromiseGrid example app for collaborative document editing.
It combines:

- a shared Go runtime and local service
- a browser editor surface
- a Neovim editor surface
- repo-local draft PromiseGrid-facing live protocols for document sync and awareness

The local service is embodiment plumbing. The peer-visible contract lives in:

- `protocols/live-document.md`
- `protocols/live-awareness.md`

Both browser and Neovim are embodiments of the same app contract. Source:
`DI-lodug`; `DI-tofug`; `DI-jilin`.

## Current shape

The first slice is intentionally small:

- signed `grid([42(pCID), payload, proof])` envelopes
- durable `Ed25519` identity persisted locally
- append-only local message log
- repo-local draft `live-document` and `live-awareness` specs
- internal HTTP adapter between the local service and the browser/Neovim embodiments
- deterministic last-writer-wins document convergence based on `(lamport, author, message_cid)`

This is a working example repo, not a frozen upstream PromiseGrid protocol.

## Run

```bash
go run ./cmd/grid-editor --listen 127.0.0.1:7001
```

Open the browser UI:

```text
http://127.0.0.1:7001/?doc=demo
```

To point another local service at the first one for peer sync:

```bash
go run ./cmd/grid-editor --listen 127.0.0.1:7002 --peer http://127.0.0.1:7001
```

## Neovim

Add the local plugin directory to your runtime path, then open a shared
document:

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

