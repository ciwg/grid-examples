# grid-editor architecture

`grid-editor` centers the PromiseGrid-facing contract in one place and keeps
the embodiment plumbing local.

## Topology

```text
Browser UI ----\
                \
                 local HTTP adapter -> grid-editor Go service -> peer services
                /
Neovim plugin --/
```

The local Go service owns:

- `Ed25519` identity
- repo-local pCID discovery from exact local spec bytes
- signed grid envelope creation and verification
- append-only message logging
- document and awareness state projection
- optional peer polling

The browser UI and the Neovim plugin own:

- local editing UX
- local cursor/UI wiring
- local HTTP calls into the service

They do not define the peer-visible protocol truth. Source: `DI-lodug`;
`DI-tofug`; `DI-jilin`.

## Public versus internal boundaries

Public, PromiseGrid-facing:

- `protocols/live-document.md`
- `protocols/live-awareness.md`
- signed `grid([42(pCID), payload, proof])` envelopes

Internal-only:

- local HTTP endpoints
- browser polling loop
- Neovim `curl`/`vim.system` plumbing
- local file layout under `.grid-editor/`

## Convergence model

The first slice uses a deterministic last-writer-wins rule for document and
awareness updates:

- primary order: `lamport`
- tie-break: `author`
- final tie-break: `message_cid`

Intent: make multi-host replay and mixed-version comparisons deterministic
without hiding the ordering rule inside UI code. Source: `DI-jilin`.

