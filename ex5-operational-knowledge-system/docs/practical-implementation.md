# practical implementation notes

## Current storage and projection model

The `ex5` foundation uses one append-only `events.jsonl` file plus an
`attachments/` directory. The service replays the event log at startup and
projects current query views for:

- responsibilities
- knowledge items
- runs
- approvals
- links

This is closer to the `ex4` durable workflow pattern than the `ex3`
collaborative runtime pattern.

## Current browser and CLI shape

Both embodiments talk to the same local HTTP surface.

Browser:

- creates responsibilities
- creates knowledge items
- records runs
- records approvals
- uploads evidence
- searches responsibilities, items, and runs

CLI:

- prints dashboard counts
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

## Honest current limitation

The foundation does not yet include:

- live collaborative editing
- websocket transport
- relay-visible peer exchange

Those were deliberately left out of this first implementation pass so the
durable operational record and equal browser/CLI embodiment model could exist
as a runnable example first.
