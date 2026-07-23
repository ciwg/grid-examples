# ex5 HTTP API guide

This guide documents the local HTTP adapter that the browser still uses as its
primary embodiment surface, and that CLI/Neovim still keep as explicit
compatibility transport.

The dedicated remote relay service is documented separately in
[Relay API Guide](./relay-api-guide.md).

It is intentionally a local embodiment surface, not the final PromiseGrid wire
contract. The durable history still lives in the ex5 runtime model and
protocol-family seams described elsewhere.

In current ex5, these HTTP routes are the shipped browser embodiment contract
and the explicit compatibility transport for CLI and Neovim. They are not the signed
PromiseGrid peer contract, and route names should not be read as frozen
`pCID`-selected public wire meaning. Source: `DI-sobek`; `DI-favel`.

The adapter is served by the same Go 1.24.13 runtime pinned in this module's
`go.mod`, matching the current patch-level default used across the other
`grid-examples` modules.

The browser still speaks this HTTP adapter directly. CLI and Neovim still keep
it as compatibility transport, but the direct local Unix-socket contract now
uses typed runtime `operation` messages for the first inspect/read slice
instead of forwarding those reads as generic `GET /api/*` socket requests.
Source: `DI-monuv`.

## Core shape

The server defaults to:

```text
http://127.0.0.1:7045
```

Responses are JSON except:

- `/` -> `index.html`
- `/app.js` -> browser module
- `/style.css` -> browser stylesheet

## Metadata and dashboard

### `GET /api/meta`

Returns runtime capability metadata plus the supported kinds and lifecycle
values.

The response now includes, for example:

- `knowledge_kinds`
- `run_kinds`
- `approval_decisions`
- `item_statuses`
- `peer_exchange_format`
- `peer_exchange_families`
- `relay_feed_format`
- `relay_feed_families`
- `relay_blob_transfer_enabled`
- `cas_objects_enabled`
- `cas_attachment_blobs_enabled`
- `cas_draft_bodies_enabled`
- `live_draft_websocket_enabled`
- `local_unix_socket_enabled`
- `local_unix_socket_path`
- `embodiments`

Those capability fields are the current shipped way for embodiments to confirm
which local adapter/runtime features are present without treating the HTTP
route names themselves as the PromiseGrid peer contract. They also let CLI and
Neovim discover the runtime's canonical socket path before relying on local
filesystem guesses. In Neovim, that runtime-first discovery is bounded by a
short timeout and falls back immediately to local repo-root inference if the
HTTP capability path is unavailable. Source: `DI-bavuk`; `DI-zunep`;
`DI-sorek`; `DI-batov`.

In the current runtime, `embodiments.browser` declares `local_http` as its
primary adapter and `websocket_over_local_http` as its live-draft transport.
`embodiments.cli` declares `local_unix_socket` as its primary adapter and marks
HTTP compatibility as `explicit_opt_in`. `embodiments.neovim` also declares
`local_unix_socket` as both its primary adapter and live-draft transport, with
websocket and HTTP listed as compatibility transports that are only available
through explicit opt-in mode (`oks-nvim --socket=off`). Terminal embodiments
also repeat the canonical `local_unix_socket_path` inside their own embodiment
records. Source: `DI-vurak`; `DI-zorav`; `DI-fonuv`.

For the currently shipped terminal runtime contract, the direct socket now
uses typed `operation` messages for:

- `inspect_item`
- `inspect_run`
- `inspect_entity`
- `search`
- `pending_review`
- `problem_review`

HTTP route names remain the browser adapter and explicit compatibility surface,
not the primary local terminal contract for those reads. Source: `DI-monuv`.

## Peer exchange

### `GET /api/peer-exchange/export`

Exports the current shipped peer-visible PromiseGrid bundle over the local HTTP
adapter.

The response includes:

- `format`
- `exported_at`
- `exporting_peer_id`
- family `pCID` fields for the shipped peer-visible families
- `events`
- whole-family signed record arrays
- `cas_blob_objects` for inline CID-keyed evidence blobs

This is a local embodiment/export surface, not a claim that `/api/*` route
names are the PromiseGrid wire contract. Source: `DI-voruk`; `DI-rumek`;
`DI-faruv`; `DI-loruk`.

### `POST /api/peer-exchange/import`

Imports one peer-exchange bundle into the current runtime.

The response includes:

- `imported_events`
- per-family import counts
- `unresolved_references`

The shipped importer accepts unseen origin-aware history into non-empty
runtimes, dedupes by `(origin_peer_id, origin_sequence)`, and preserves the
current local HTTP adapter as the embodiment surface for that exchange.
Source: `DI-rumek`; `DI-ruzok`; `DI-bavuk`.

## Relay feed

### `POST /api/relay/feed/export`

Exports an incremental relay batch over the same local HTTP adapter.

The request body includes:

- `known_origins`

`known_origins` is a per-peer cursor map. Each key is an `origin_peer_id`, and
each value is the highest `origin_sequence` the caller already has for that
peer. The server replies only with signed family records and compatibility
events whose origin tuples are newer than those cursors. Source: `DI-pazek`.

The response includes:

- `format`
- `exported_at`
- `exporting_peer_id`
- family `pCID` fields for the eight shipped peer-visible families
- `events`
- whole-family signed record arrays for just the unseen origin tuples
- `required_blob_cids`

Unlike `/api/peer-exchange/export`, this route does not inline blob bytes.
Evidence blobs are named only by CID here so the relay feed stays
record-oriented and blob transfer stays explicitly content-addressed. Source:
`DI-pazek`.

### `POST /api/relay/feed/import`

Imports one incremental relay batch into the current runtime.

The response includes:

- `imported_events`
- per-family import counts
- `missing_blob_cids`
- `unresolved_references`

If the batch references evidence blobs whose CIDs are not already present in
local CAS, the route returns `409` plus `missing_blob_cids` and does not apply
the feed. Clients must stage those blobs first through `/api/relay/blobs/{cid}`
and then retry the feed import. Source: `DI-pazek`.

### `GET /api/relay/blobs/{cid}`

Returns one raw CAS object by CID as `application/octet-stream`.

This is the blob side of the relay split. It is intentionally separate from
the signed record feed so relay clients can fetch only the missing objects they
actually need. Source: `DI-pazek`.

### `PUT /api/relay/blobs/{cid}`

Stores one raw CAS object by CID.

The request body is the object bytes. The runtime verifies that the uploaded
bytes hash to the CID named in the route before storing them. Source:
`DI-pazek`.

### `GET /api/dashboard`

Returns the current projected counts for:

- responsibilities
- places
- resources
- each knowledge-item kind
- each run kind
- approvals
- evidence
- links

### `GET /api/problem-review`

Returns grouped receiving/count problem hotspots for the browser review panel.

The CLI now reuses this same route directly for a terminal hotspot summary:

- `oks-cli problem-review`

Neovim now reuses it too:

- `:OksProblemReview`

Source: `DI-nuvaz`; `DI-ravum`; `DI-sivok`.

Terminal review surfaces also assume each run record that participates in
pending/problem review carries an explicit `approvals` array. The CLI and
Neovim pending-review queues treat an omitted `approvals` field as shared
projection contract drift, not as “unreviewed.” Source: `DI-davur`.

The response includes:

- `problem_runs`
- `place_groups`
- `resource_groups`

Each group includes:

- `group_type`
- `group_id`
- `kind`
- `name`
- `problem_count`
- `receiving_problems`
- `inventory_problems`
- `highlights`
- `runs`

## Typed links

### `POST /api/links`

Creates one typed link between two projected records.

Payload fields:

- `actor`
- `from_type`
- `from_id`
- `to_type`
- `to_id`
- `relation`
- `notes`

Supported endpoint types today:

- `place`
- `resource`
- `responsibility`
- `knowledge_item`
- `run`

The write is now validated on both ends. Unsupported endpoint types or missing
record IDs return `400` instead of entering the append-only history as dangling
graph edges. Source: `DI-luzaf`.

The CLI now reuses this same JSON route directly:

- `oks-cli add-link ACTOR FROM_TYPE FROM_ID TO_TYPE TO_ID RELATION`
- `oks-cli add-link ACTOR FROM_TYPE FROM_ID TO_TYPE TO_ID RELATION NOTES...`

Source: `DI-vuteg`.

## Search

### `GET /api/search`

Searches across places, resources, responsibilities, items, and runs.

Supported query params:

- `q`
- `kind`
- `status`
- `outcome`
- `place_id`
- `resource_id`
- `responsibility_id`
- `problem`

The response includes:

- `filters`
- `places`
- `resources`
- `responsibilities`
- `items`
- `runs`

This is also the browser's history-drilldown surface. The record inspector uses
structured `kind`, `outcome`, `place_id`, `resource_id`, and
`responsibility_id` filters to answer questions like:

- show me all receiving runs here
- show me all counts for this bin
- show me receiving problems in this area

When `problem=true`, the run slice is filtered by the same receiving/inventory
problem classification used by `/api/problem-review`, instead of by one
hardcoded receiving outcome. Source: `DI-vemur`.

This is now also the Neovim search/browse surface. `:OksSearch QUERY` reads the
same response and renders grouped read-only result sections with inspect hints
for the existing item, run, and generic entity inspectors. Source: `DI-givot`.

Neovim now also reuses the shared structured filters on that same route:

- `:OksSearch QUERY kind=... status=... outcome=...`
- `:OksSearch QUERY place_id=... resource_id=... responsibility_id=...`
- `:OksSearch QUERY problem=true`

Unsupported filter keys are rejected in-editor instead of being silently
dropped. Source: `DI-fanub`.

The Neovim pending-review view also reuses this route. `:OksPending` combines
`/api/search?status=draft`, `/api/search`, and `/api/search?problem=true` to
assemble draft-item, unreviewed-run, and problem-run queues without requiring
a new terminal-specific API. Source: `DI-lorav`.

The CLI now reuses this same route too:

- `oks-cli search QUERY`
- `oks-cli search QUERY kind=... status=... outcome=...`
- `oks-cli search QUERY place_id=... resource_id=... responsibility_id=...`
- `oks-cli search QUERY problem=true`
- `oks-cli pending-review`

`oks-cli pending-review` reuses the same three route reads as `:OksPending`:
`/api/search?status=draft`, `/api/search`, and `/api/search?problem=true`. It
renders one shell-facing pending-review summary instead of requiring a new
terminal-specific aggregation API. Source: `DI-vabok`; `DI-ravum`.

Unsupported or malformed filter tokens are rejected locally instead of being
silently dropped. Source: `DI-mifot`.

Neovim item approval reuses the item detail and item approval routes together:
it reads `GET /api/items/{id}` to resolve the current revision, then posts to
`POST /api/items/{id}/approvals` with the configured Neovim display name as
the approval actor. Source: `DI-vamor`.

Neovim run approval reuses `POST /api/runs/{id}/approvals` directly. It uses
the configured Neovim display name as the approval actor and refreshes the
current run or pending-review view after the write returns. Source:
`DI-bafor`.

Neovim item supersede reuses `POST /api/items/{id}/supersede` directly. It
uses the configured Neovim display name as the lifecycle actor and refreshes
the current live, inspect, or pending-review view after the write returns.
Source: `DI-pudor`.

## Places and resources

### `GET /api/places`

Lists all known places.

### `POST /api/places`

Creates a place.

Payload fields:

- `actor`
- `kind`
- `name`
- `summary`
- `parent_id`
- `tags`

### `GET /api/places/{id}`

Returns one place with:

- hierarchy context
- linked resources
- related runs
- typed links
- timeline

The CLI now reuses this route as a terminal drilldown summary:

- `oks-cli show-place PLACE_ID`

It renders hierarchy, related runs, and link context directly instead of
dumping an undifferentiated JSON blob. Source: `DI-luzom`.

### `GET /api/resources`

Lists all known resources.

### `POST /api/resources`

Creates a resource.

Payload fields:

- `actor`
- `kind`
- `name`
- `summary`
- `place_id`
- `tags`

### `GET /api/resources/{id}`

Returns one resource with:

- place context
- related runs
- typed links
- timeline

The CLI now reuses this route as a terminal drilldown summary:

- `oks-cli show-resource RESOURCE_ID`

It renders place context, related runs, and link context directly instead of
dumping an undifferentiated JSON blob. Source: `DI-luzom`.

## Responsibilities

### `GET /api/responsibilities`

Lists all responsibilities.

### `POST /api/responsibilities`

Creates a responsibility.

Payload fields:

- `actor`
- `title`
- `summary`
- `role_keys`
- `tags`

### `GET /api/responsibilities/{id}`

Returns one responsibility with linked items/runs and its timeline.

The CLI now reuses this same route directly for a terminal drilldown summary of
linked items, linked runs, related runs, and typed links:

- `oks-cli show-responsibility RESPONSIBILITY_ID`

Source: `DI-jubav`; `DI-salup`.

## Knowledge items

### `GET /api/items`

Lists knowledge items. Optional query:

- `kind`

Supported kinds today:

- `procedure`
- `training`
- `maintenance`
- `receiving_check`
- `inventory_audit`

`receiving_check` is the broad inbound inspection and intake workflow kind. It
is meant for received parts, returned items, tool intake, staged kits, and
similar receipt/review work that should not be forced into either a plain
inventory count or a generic procedure label.

### `POST /api/items`

Creates a new knowledge item at revision `1`.

Payload fields:

- `actor`
- `kind`
- `title`
- `summary`
- `body`
- `tags`
- `responsibility_ids`

### `GET /api/items/{id}`

Returns one projected knowledge item, including:

- `status`
- `current_revision`
- `working_body`
- `working_version`
- revision list
- approvals
- links

For a `receiving_check` item, the browser uses the same response shape to show
revision history plus receiving-specific related-run history in the record
inspector.

The current Neovim inspector also reuses this response shape for a read-only
split that shows status, revisions, approvals, and related runs. Source:
`DI-lonuk`.

That same inspector now also reads the `links` array for typed-link browsing.
Source: `DI-zalor`.

The CLI now also reuses this same response shape for a terminal drilldown
summary of revisions, approvals, related runs, and typed links:

- `oks-cli show-item ITEM_ID`

Source: `DI-salup`.

### `GET /api/responsibilities/{id}`

Returns one projected responsibility, including:

- `linked_item_ids`
- `linked_run_ids`
- `related_runs`
- `links`
- `timeline`

That `links` array now uses the same projection shape as place, resource, item,
and run detail, so browser and Neovim inspectors see the same typed-link graph
for responsibility records too. Source: `DI-luzaf`.

### `POST /api/items/{id}/revisions`

Creates a durable new revision snapshot.

Neovim `:OksSnapshot` now reuses this route directly after flushing the
current live draft body through `/api/items/{id}/live`. That keeps durable
authoring on the shared HTTP model instead of inventing an editor-only
snapshot path. Source: `DI-jabup`.

The CLI now reuses this route too. `oks-cli snapshot-item ITEM_ID ACTOR BODY`
loads the existing item title, summary, and tags from `GET /api/items/{id}`,
then posts the supplied body through this revision endpoint. Source:
`DI-muvok`.

Payload fields:

- `actor`
- `title`
- `summary`
- `body`
- `tags`

### `POST /api/items/{id}/approvals`

Records an approval against a knowledge item revision.

Payload fields:

- `actor`
- `revision`
- `role`
- `decision`
- `notes`

If `decision == "approved"`, the item lifecycle moves to `approved`.
That lifecycle change is revision-aware: approving a stale older revision now
returns `400` instead of silently marking a newer draft as approved. Source:
`DI-dazim`.

The CLI now requires an explicit `actor` for approval commands instead of
inventing a placeholder approver name, so durable approval records preserve the
real identity that invoked the command. Source: `DI-tarok`.

### `POST /api/items/{id}/supersede`

Marks a knowledge item as `superseded`.

Payload fields:

- `actor`
- `notes`

## Live draft endpoints

The live draft surface is browser-oriented, but it is also reused directly by
the first Neovim phase. It shares the current working body for a knowledge item
without turning that live state into a durable revision automatically. Source:
`DI-fudok`; `DI-noruv`.

### `GET /api/items/{id}/live`

Returns:

- `item_id`
- `title`
- `status`
- `body`
- `version`
- `current_revision`
- `participants`

### `POST /api/items/{id}/live`

Updates participant presence and optionally updates the shared body.

Payload fields:

- `participant_id`
- `display_name`
- `color`
- `cursor`
- `head`
- `typing`
- `base_version`
- `update_body`
- `body`

Behavior:

- presence is refreshed on every call
- if `update_body == true` and `base_version` matches the current live version,
  the working draft body is updated
- if `update_body == true`, an empty `body` is treated as an intentional clear
  and advances the live version
- if `update_body == false`, the call is presence-only and does not change the
  shared body or version
- if `update_body == true` and `base_version` is stale, the server returns
  `409`

Conflict response shape:

- `conflict: true`
- `state: <current live state>`

### `GET /api/items/{id}/live/socket`

Upgrades the existing local adapter to websocket for shared live-draft
carriage. Browser and Neovim now prefer this route for participant presence and
shared body updates, while `GET/POST /api/items/{id}/live` stay available for
bootstrap fetch and fallback. Source: `DI-noruv`.

Client message shape:

- `type`
- `participant_id`
- `display_name`
- `color`
- `cursor`
- `head`
- `typing`
- `base_version`
- `update_body`
- `body`

Server message shape:

- `type: "live-state" | "live-conflict" | "error"`
- `state: <current live state>` when applicable
- `message` for error payloads

## Runs

### `GET /api/runs`

Lists runs. Optional query:

- `kind`

Supported run kinds today match the knowledge-item kinds:

- `procedure`
- `training`
- `maintenance`
- `receiving_check`
- `inventory_audit`

### `POST /api/runs`

Creates a run record.

Payload fields:

- `actor`
- `kind`
- `item_id`
- `revision`
- `outcome`
- `notes`
- `machine`
- `location`
- `place_id`
- `resource_ids`
- `responsibility_ids`

### `GET /api/runs/{id}`

Returns one projected run, including:

- evidence
- approvals
- links
- timeline

The CLI now also reuses this same response shape for a terminal drilldown
summary of context, evidence, approvals, and typed links:

- `oks-cli show-run RUN_ID`

The CLI and Neovim run inspectors now also use the same run detail projection
to hand terminal users back into related item, place, resource, and
responsibility context without inventing a second terminal-only route.
Source: `DI-salup`; `DI-vunep`.

For a `receiving_check` run, the browser uses the evidence facts from this
response to render the `Receiving review` panel.

The current Neovim run inspector also reuses this response shape for a
read-only split that shows run context, evidence summaries, approvals, and
direct handoff hints into related context inspectors. Source: `DI-ravok`;
`DI-vunep`.

It also reads the run `links` array for typed-link browsing. Source: `DI-zalor`.

### `GET /api/places/{id}`

The current Neovim entity inspector reuses this response shape for read-only
place browsing. Source: `DI-zalor`.

### `GET /api/resources/{id}`

The current Neovim entity inspector reuses this response shape for read-only
resource browsing. Source: `DI-zalor`.

### `GET /api/responsibilities/{id}`

The current Neovim entity inspector reuses this response shape for read-only
responsibility browsing. Source: `DI-zalor`.

## Evidence

### `POST /api/runs/{id}/evidence`

Adds evidence to a run using multipart form upload.

Fields:

- `actor`
- `summary`
- `facts_json`
- optional `attachment`

The CLI now reuses this same multipart route directly:

- `oks-cli add-evidence RUN_ID ACTOR SUMMARY`
- `oks-cli add-evidence RUN_ID ACTOR SUMMARY FACTS_JSON`
- `oks-cli add-evidence RUN_ID ACTOR SUMMARY FACTS_JSON FILE`

Source: `DI-zanub`.

Attachment uploads are limited to 8 MiB. Larger files are rejected with a
`400 Bad Request` response instead of being truncated into durable evidence.
Source: `DI-navos`.

## Run approvals

### `POST /api/runs/{id}/approvals`

Records an approval against a run.

Payload fields:

- `actor`
- `role`
- `decision`
- `notes`

## Search

### `GET /api/search?q=...`

Returns mixed projected results across:

- places
- resources
- responsibilities
- items
- runs

Run free-text matching includes:

- run outcome, notes, machine, and location
- linked place and resource names/summaries
- evidence summaries and fact key/value text
- approval actor/role/decision/notes text

Source: `DI-farun`.
