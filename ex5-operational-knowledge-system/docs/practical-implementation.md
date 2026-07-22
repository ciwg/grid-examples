# practical implementation notes

## Current storage and projection model

The `ex5` foundation uses one append-only `events.jsonl` file plus an
`attachments/` directory, plus a `drafts/` directory for browser-shared working
bodies. The service replays the event log at startup, overlays any saved live
draft bodies, and projects current query views for:

- places
- resources
- responsibilities
- knowledge items
- runs
- approvals
- links

This is closer to the `ex4` durable workflow pattern than the `ex3`
collaborative runtime pattern, but it now has a lightweight live-draft layer on
top for browser collaboration.

The durability rules for that storage are now stricter:

- evidence attachments are written to immutable per-upload paths, even when two
  uploads reuse the same filename
- event replay is buffered for large stored item and revision bodies inside the
  server's current request-size envelope

That keeps old evidence bytes and larger saved knowledge text readable after a
restart instead of depending on filename uniqueness or the scanner default
token size. Source: `DI-busor`.

## Current browser, CLI, and Neovim shape

All current embodiments talk to the same local HTTP surface.

The current `ex5` module now pins Go 1.24.13, which matches the other
`grid-examples` modules and avoids a separate patch-level Go requirement when
switching between examples.

Browser:

- creates places and resources
- creates responsibilities
- creates knowledge items
- edits shared live drafts for knowledge items
- snapshots the current working draft into a new revision
- approves or supersedes the current item
- inspects records, reads timelines, and follows contextual links between them
- reviews item revisions plus run evidence/approvals in the record inspector
- shows related run history from the selected item detail view
- shows related run history from place, resource, and responsibility detail views
- shows receiving-check evidence and receiving history in dedicated review sections
- shows inventory audit discrepancy/count facts and inventory audit history in dedicated review sections
- shows receiving and inventory fact history directly from place/resource/responsibility context views, instead of only bare related-run summaries
- filters search by kind, status, outcome, place, resource, and responsibility
- uses those same structured filters for one-click context drilldowns into receiving/count/problem history
- uses an explicit shared `problem=true` search filter so browser drilldowns and grouped hotspot review classify problems the same way
- summarizes repeated receiving and inventory problems by place and resource in a dedicated browser review panel
- falls back to ephemeral in-memory participant identity when browser storage or UUID helpers are restricted, so the live-draft UI still boots in private/policy-limited contexts
- routes form/search failures through one shared in-app error path instead of relying on unhandled async request failures
- records runs
- records approvals
- uploads evidence
- searches places, resources, responsibilities, items, and runs

CLI:

- prints dashboard counts
- creates and lists places and resources
- creates responsibilities and knowledge items
- records runs
- records approvals with explicit actor identity
- shows individual items and runs

Neovim phase 1:

- opens a knowledge item live draft by item ID
- polls the same live-draft HTTP endpoint for refresh
- pushes the current body with `:write`
- sends presence/typing heartbeats through the live-draft HTTP endpoint
- exposes local status/participant inspection commands
- opens a read-only projected item inspector for revisions, approvals, and related runs
- opens a read-only projected run inspector for evidence and approvals
- opens a generic read-only entity inspector for linked places, resources, responsibilities, items, and runs
- opens a read-only grouped search buffer over `/api/search` for discovery and browse handoff into the existing inspectors
- opens a read-only pending-review buffer by combining draft-item, all-run, and problem-run search slices from the same projection layer
- posts limited item approvals by resolving current revision from `GET /api/items/{id}` and then using `POST /api/items/{id}/approvals`
- posts limited run approvals with `POST /api/runs/{id}/approvals`

In practice, that gives `ex5` a real terminal-first operating mode:

- CLI for direct command execution and mutation
- Neovim for continuous text work, review, browsing, and pending-work triage

The important behavior point is that these are not separate backends. They are
two terminal-facing views over the same local runtime and projected state, so a
user can mix shell commands and Neovim inspection without crossing embodiment
boundaries or losing context. Source: `DI-fudok`; `DI-givot`; `DI-lorav`.

The plugin now tracks the live-draft window explicitly so Neovim presence and
body pushes continue to report cursor/head against the draft buffer even after
the user moves focus into a read-only inspector split. Source: `DI-pazud`.

The session close path now wipes both the live-draft buffer and the current
read-only inspector buffer. That keeps `:OksClose` honest as a teardown action
instead of leaving a visible buffer whose live session hooks are already gone.
Source: `DI-mabek`.

## Why the docs mention protocol families

The implementation already organizes the model around protocol-family seams:

- `knowledge-item`
- `knowledge-approval`
- `knowledge-evidence`
- `knowledge-link`
- `knowledge-responsibility`
- `knowledge-search-metadata`

The current Go code does not yet emit signed grid envelopes for those families,
but the seams are intentionally visible in the data model so the example can
move there without a total rewrite.

The typed-link path is also now stricter at write time. The runtime validates
both endpoint types and endpoint IDs before appending a link event, and
responsibility records now project their `links` array alongside the older
linked-item and linked-run convenience fields. Source: `DI-luzaf`.

## Live draft behavior

The browser draft studio is intentionally separate from durable history:

- the current working body is shared
- participant presence is shared
- a live version number guards against blind overwrite
- creating a revision snapshots the working body into append-only history

This means the browser and first-phase Neovim embodiment can collaborate on the
in-progress text of a procedure, training doc, maintenance doc, receiving
check, or inventory audit without confusing that live session state with the
durable operational record.

The live HTTP surface now also distinguishes presence-only posts from real body
writes. Browser and Neovim clients set `update_body=true` when they intend to
change the shared draft, including clearing it to empty text, and use
`update_body=false` for presence heartbeats that should not advance the live
version. Source: `DI-dazim`.

## Honest current limitation

The foundation still does not yet include:

- websocket transport
- relay-visible peer exchange
- signed grid envelopes on the wire

Those are still important future steps. The current pass focuses on a runnable
standalone operational-memory tool with one local runtime, a richer browser
embodiment, a thinner CLI embodiment, and a browser-only shared draft surface.

The current product direction is to keep that live-draft surface optional,
rather than making collaborative editing the core of the tool, and to revisit a
future richer Neovim embodiment later without porting the full `ex3`
websocket stack into `ex5` now. The current phase is intentionally a thin HTTP
live-draft client plus read-only item/run/entity inspection and search/browse
over projected detail. The next terminal-first step now adds a read-only
pending-review queue over the same search projections, still before any
write-side approval actions are added to the editor. The next small step now
adds one revision-safe item approval action, still using the same local HTTP
runtime and item approval semantics as the CLI and browser. Source:
`DI-tabiv`; `DI-fudok`; `DI-lonuk`; `DI-ravok`; `DI-zalor`; `DI-givot`;
`DI-lorav`; `DI-vamor`; `DI-bafor`.

## Current verification shape

The current code is covered at four levels:

- app/service tests for projection, search, lifecycle, and live draft behavior
- HTTP server tests for routes, conflict handling, and mixed workflow flows
- CLI tests for command argument mapping into the HTTP adapter
- Neovim asset tests for the shipped launcher and command surface
- embedded web asset tests that assert the shipped UI still exposes the
  expected operational workflow sections
- a headless browser smoke layer that loads the real UI shell and checks
  rendered operational detail content against mostly stubbed API responses

The browser UI is still not covered by a deep end-to-end interaction suite, but
it is no longer limited to static asset checks. The current baseline now
includes:

- asset-structure checks
- API-level tests for the data the browser depends on
- a real headless browser render smoke over the live app shell
