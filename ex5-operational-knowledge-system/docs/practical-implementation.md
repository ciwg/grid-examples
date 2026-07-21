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
- summarizes repeated receiving and inventory problems by place and resource in a dedicated browser review panel
- records runs
- records approvals
- uploads evidence
- searches places, resources, responsibilities, items, and runs

CLI:

- prints dashboard counts
- creates and lists places and resources
- creates responsibilities and knowledge items
- records runs
- records approvals
- shows individual items and runs

Neovim phase 1:

- opens a knowledge item live draft by item ID
- polls the same live-draft HTTP endpoint for refresh
- pushes the current body with `:write`
- sends presence/typing heartbeats through the live-draft HTTP endpoint
- exposes local status/participant inspection commands
- opens a read-only projected item inspector for revisions, approvals, and related runs
- opens a read-only projected run inspector for evidence and approvals

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

## Honest current limitation

The foundation still does not yet include:

- websocket transport
- relay-visible peer exchange
- signed grid envelopes on the wire

Those are still important future steps. The current pass focuses on a runnable
standalone operational-memory tool with one local runtime, equal browser/CLI
operational embodiments, and a browser-only shared draft surface.

The current product direction is to keep that live-draft surface optional,
rather than making collaborative editing the core of the tool, and to revisit a
future richer Neovim embodiment later without porting the full `ex3`
websocket stack into `ex5` now. The current phase is intentionally a thin HTTP
live-draft client plus read-only item/run inspection over projected detail.
Source: `DI-tabiv`; `DI-fudok`; `DI-lonuk`; `DI-ravok`.

## Current verification shape

The current code is covered at four levels:

- app/service tests for projection, search, lifecycle, and live draft behavior
- HTTP server tests for routes, conflict handling, and mixed workflow flows
- CLI tests for command argument mapping into the HTTP adapter
- Neovim asset tests for the shipped launcher and command surface
- embedded web asset tests that assert the shipped UI still exposes the
  expected operational workflow sections
- a headless browser smoke test that loads the real UI against a live test
  server and checks for rendered operational detail content

The browser UI is still not covered by a deep end-to-end interaction suite, but
it is no longer limited to static asset checks. The current baseline now
includes:

- asset-structure checks
- API-level tests for the data the browser depends on
- a real headless browser render smoke over the live app shell
