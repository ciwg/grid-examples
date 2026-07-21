# ex5 HTTP API guide

This guide documents the local HTTP adapter that the browser, CLI, and
first-phase Neovim embodiment all use.

It is intentionally a local embodiment surface, not the final PromiseGrid wire
contract. The durable history still lives in the ex5 runtime model and
protocol-family seams described elsewhere.

The adapter is served by the same Go 1.24.13 runtime pinned in this module's
`go.mod`, matching the current patch-level default used across the other
`grid-examples` modules.

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

Returns the supported kinds and lifecycle values, for example:

- `knowledge_kinds`
- `run_kinds`
- `approval_decisions`
- `item_statuses`

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
- timeline

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

Returns one resource and its timeline.

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

### `POST /api/items/{id}/revisions`

Creates a durable new revision snapshot.

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

### `POST /api/items/{id}/supersede`

Marks a knowledge item as `superseded`.

Payload fields:

- `actor`
- `notes`

## Live draft endpoints

The live draft surface is browser-oriented, but it is also reused directly by
the first Neovim phase. It shares the current working body for a knowledge item
without turning that live state into a durable revision automatically. Source:
`DI-fudok`.

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
- `body`

Behavior:

- presence is refreshed on every call
- if `body` differs and `base_version` matches the current live version, the
  working draft is updated
- if `body` is empty, the call still refreshes participant presence and returns
  the current shared state without advancing the version
- if `base_version` is stale, the server returns `409`

Conflict response shape:

- `conflict: true`
- `state: <current live state>`

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

For a `receiving_check` run, the browser uses the evidence facts from this
response to render the `Receiving review` panel.

## Evidence

### `POST /api/runs/{id}/evidence`

Adds evidence to a run using multipart form upload.

Fields:

- `actor`
- `summary`
- `facts_json`
- optional `attachment`

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
