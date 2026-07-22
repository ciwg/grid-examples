# ex5 relay API guide

This guide documents the dedicated remote relay service shipped by
`cmd/operational-relay`.

The relay is not the local embodiment adapter and it is not the main
operational runtime. Its job is narrower:

- accept origin-aware incremental relay-feed publishes
- serve relay-feed pulls by per-origin cursor map
- stage and fetch CID-addressed blobs

That separation is intentional. The browser, CLI, and Neovim embodiments still
use the local HTTP adapter, while the relay carries signed durable history and
evidence blobs remotely. Source: `DI-rovik`; `DI-tasov`; `DI-nulav`.

## Core shape

The relay binary defaults to:

```text
http://127.0.0.1:7046/relay/v1
```

Its data root defaults to:

```text
.operational-relay/
```

Durable state includes:

- `events.jsonl`
- per-family signed message logs
- `origin-cursors.json`
- `cas/objects/*`

## Routes

### `GET /relay/v1/meta`

Returns relay capability metadata:

- `service_name`
- `route_prefix`
- `relay_feed_format`
- `relay_feed_families`
- `relay_blob_transfer_enabled`
- `publish_requires_staged_blobs`

### `POST /relay/v1/feed/publish`

Publishes a `RelayFeedBatch` into the remote relay.

Important behavior:

- the relay persists only unseen origin tuples
- per-origin relay history must be contiguous
- a brand-new origin must start at sequence `1`
- evidence-bearing publishes fail with `409` until all referenced blob CIDs are
  already staged into relay CAS

On success the response includes per-family publish counts.

### `POST /relay/v1/feed/pull`

Pulls unseen relay history by per-origin cursor map.

Request body:

- `known_origins`

The response is a `RelayFeedBatch` containing only origin tuples newer than
the caller's cursor map. Events are renumbered to a fresh batch-local
compatibility `sequence`, while `origin_peer_id` plus `origin_sequence`
remain the durable identity. Source: `DI-tasov`.

### `PUT /relay/v1/blobs/{cid}`

Stages a raw CAS blob into the relay by CID.

The relay verifies that the uploaded bytes hash to the CID named in the route.

### `GET /relay/v1/blobs/{cid}`

Fetches one staged CAS blob by CID.

## Relationship to the local adapter

The local ex5 adapter still ships:

- `/api/peer-exchange/*`
- `/api/relay/feed/*`
- `/api/relay/blobs/{cid}`

Those local routes are embodiment/runtime surfaces for one node. The remote
relay routes above are the separate durable transport surface for store-and-
forward deployment. Source: `DI-rovik`; `DI-pazek`.
