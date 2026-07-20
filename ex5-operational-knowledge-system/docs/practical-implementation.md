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

## Current browser and CLI shape

Both embodiments talk to the same local HTTP surface.

Browser:

- creates places and resources
- creates responsibilities
- creates knowledge items
- edits shared live drafts for knowledge items
- snapshots the current working draft into a new revision
- approves or supersedes the current item
- inspects records and follows contextual links between them
- filters search by kind, status, place, resource, and responsibility
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

This means the browser can collaborate on the in-progress text of a procedure,
training doc, maintenance doc, or inventory audit without confusing that live
session state with the durable operational record.

## Honest current limitation

The foundation still does not yet include:

- websocket transport
- relay-visible peer exchange
- signed grid envelopes on the wire

Those are still important future steps. The current pass focuses on a runnable
standalone operational-memory tool with one local runtime, equal browser/CLI
operational embodiments, and a browser-only shared draft surface.

## Current verification shape

The current code is covered at four levels:

- app/service tests for projection, search, lifecycle, and live draft behavior
- HTTP server tests for routes, conflict handling, and mixed workflow flows
- CLI tests for command argument mapping into the HTTP adapter
- embedded web asset tests that assert the shipped UI still exposes the
  expected operational workflow sections

The browser UI is still lightweight enough that the most useful regression
coverage today is:

- asset-structure checks
- API-level tests for the data it depends on

That should be read honestly: it is stronger than no browser coverage, but it
is not yet a full browser automation suite.
