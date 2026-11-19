# Practical implementation notes

This is a practical guide for getting `grid-editor` working in a more
server-shaped environment where grid agents serve only a thin browser shell,
bootstrap the grid kernel bits, and then move the rest of the live session
over WebSocket.

## Keep the browser shell thin

- Serve only enough HTML, CSS, and JS to:
  - create the editor shell
  - bootstrap identity/session state
  - learn the active pCIDs
  - open the live transport
- Keep protocol meaning in the signed payloads and pCIDs, not in ad hoc
  front-end route behavior.

## Treat WebSocket as transport, not protocol meaning

- The current repo uses HTTP polling and loopback mutation endpoints for local
  relay clients.
- In the server experiment, WebSocket can replace that transport layer for
  lower latency.
- Do not let WebSocket frames become the new undocumented protocol.
- Keep the same protocol split:
  - `live-document`
  - `live-awareness`
  - `document-metadata`
  - `publish-document`

## Suggested bootstrap flow

1. Serve the thin browser shell.
2. Browser fetches a small bootstrap payload:
   - local/session identity hints
   - active relay URL or WebSocket URL
   - current pCIDs
   - document ID
3. Browser opens the live transport.
4. Browser joins the live document flow.
5. Browser starts separate awareness and metadata flows.

That keeps live typing, awareness, and document-management features explicit
instead of smearing them together in one opaque app channel.

## Map the current repo to that model

Current repo pieces already line up well:

- browser UI:
  - CodeMirror shell
  - local Automerge replica
  - local workflow/review state
- relay:
  - signer/verifier
  - CAS object store
  - append-only log
  - protocol-aware ingest/search/resolve behavior
- protocols:
  - pCID-selected meaning
  - separate live typing, awareness, metadata, and publish flows

The main thing changing in tomorrow's experiment is transport shape, not the
high-level contract boundaries.

## Keep durable and ephemeral flows separate

- `live-document` is durable CRDT change traffic.
- `live-awareness` is ephemeral social state.
- `document-metadata` is durable latest-state document labeling.
- `publish-document` is durable handoff/exchange.

Do not collapse those into one "session state" blob just because WebSocket is
available.

## Current codebase cautions

- Local mutation started as loopback-only hardening for the first single-host
  relay slice.
- `ex3` now adds a repo-local authenticated remote mutation mode: a bootstrap
  token mints short-lived document-scoped capabilities for browser or Neovim
  clients.
- That is PromiseGrid-aligned in direction, but still provisional rather than a
  frozen upstream app-auth API.
- Browser underline round-trips as raw `<u>...</u>` text and now renders in
  the browser editor with inherited text color so it no longer presents like a
  hyperlink.
- The browser markdown UI currently exposes a normal preview pane plus a split
  view that keeps editor and preview visible together. Keep that behavior in
  mind when remapping the shell onto a server-served bootstrap.

## Practical suggestion for tomorrow

If the goal is to get something working quickly:

1. keep the relay-side protocol model
2. add a thin authenticated WebSocket adapter
3. reuse the current document / awareness / metadata / publish boundaries
4. delay bigger server-only redesigns until after the first proof run

That remains the right shape after the July 2026 upstream guide refresh as
well: `POC20` semantic-model work and `POC21` DevOps/bootstrap work are now
tracked separately, and app-facing auth/API guidance remains provisional.
Source: `DI-talih`.

That will give you a shorter path to a working demo without losing the
PromiseGrid structure that is already paying off here.

Current ex3 status:

- browser and Neovim live `sync` and `awareness` now prefer websocket transport
- metadata and publish remain on the existing HTTP endpoints
- loopback-only mutation rules still apply until a separate authenticated
  remote-client design is locked
