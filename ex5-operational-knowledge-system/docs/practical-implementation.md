# practical implementation notes

## Current storage and projection model

The `ex5` foundation uses one append-only `events.jsonl` file plus
`knowledge-item-messages.jsonl`, `knowledge-approval-messages.jsonl`,
`knowledge-evidence-messages.jsonl`, `operational-run-messages.jsonl`,
`operational-place-messages.jsonl`, `operational-resource-messages.jsonl`,
`knowledge-link-messages.jsonl`, `knowledge-responsibility-messages.jsonl`, an
`attachments/` directory, and a `drafts/` directory for per-item draft
manifests that point at browser-shared working bodies in CAS. The service
replays the event log at startup, verifies the signed
knowledge-item, knowledge-approval, knowledge-evidence, operational-run,
operational-place, operational-resource, knowledge-link, and
knowledge-responsibility logs against that replay, overlays any saved live
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

Imported evidence now also re-materializes attachment blobs from CAS into the
local compatibility attachment tree when the source-host attachment path is not
valid on the current machine. That keeps bootstrap-exchanged evidence usable in
the existing browser, CLI, and Neovim surfaces without changing the durable
signed evidence payload. Source: `DI-faruv`.

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
- creates typed links over validated existing record endpoints
- uploads run evidence with optional facts JSON and optional copied attachment
- reuses structured and `problem=true` search filters over the shared search projection
- reads grouped problem hotspots over the shared problem-review projection
- reads projected responsibility detail over the shared responsibility route
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
- extends that same search buffer with trailing shared `key=value` filters over `/api/search`
- opens a read-only pending-review buffer by combining draft-item, all-run, and problem-run search slices from the same projection layer
- opens a read-only grouped problem-review buffer over `/api/problem-review`
- posts durable item revision snapshots with `POST /api/items/{id}/revisions`
- posts limited item approvals by resolving current revision from `GET /api/items/{id}` and then using `POST /api/items/{id}/approvals`
- posts limited run approvals with `POST /api/runs/{id}/approvals`
- posts limited item supersede actions with `POST /api/items/{id}/supersede`

In practice, that gives `ex5` a real terminal-first operating mode:

- CLI for direct command execution and mutation
- Neovim for continuous text work, durable item snapshots, review, browsing,
  and pending-work triage

The CLI evidence path now reuses the same multipart run-evidence route as the
browser. That keeps terminal-side evidence entry honest: it is not a second
schema or a special shell-only shortcut, just a narrower embodiment of the
same durable evidence contract. Source: `DI-zanub`.

The CLI typed-link path follows the same rule. It reuses the existing
`POST /api/links` contract and surfaces the same server-side endpoint
validation errors instead of inventing a special shell-only graph mutation
path. Source: `DI-vuteg`.

The richer CLI search path follows that same rule too. It reuses the same
`GET /api/search` query params the browser and Neovim already depend on,
including `problem=true`, instead of inventing a second terminal-only review
surface. Source: `DI-mifot`.

The Neovim structured-search follow-on now follows the same rule too. It keeps
` :OksSearch` on the same `/api/search` route and the same filter vocabulary as
the CLI, instead of adding a second editor-only search command or backend
shape. Source: `DI-fanub`.

The CLI pending-review path follows the same rule again. It reuses
`/api/search?status=draft`, `/api/search`, and `/api/search?problem=true` to
build one shell-facing summary for draft items, review queue runs, and problem
runs instead of inventing a terminal-only aggregation endpoint. Source:
`DI-vabok`; `DI-ravum`.

The Neovim grouped problem-review path now follows the same rule too. It
reuses `/api/problem-review` to render place and resource hotspot groups with
direct run and entity handoffs, instead of inventing an editor-only grouped
review backend. Source: `DI-sivok`.

The Neovim revision-snapshot path now follows the same rule too. It reuses
`POST /api/items/{id}/revisions` directly after flushing the current live-draft
body through `/api/items/{id}/live`. That keeps durable authoring inside the
editor without inventing a separate snapshot-only transport or schema. Source:
`DI-jabup`.

That queue path now also treats an omitted run `approvals` field as a shared
projection contract error instead of silently reclassifying the run as genuine
unreviewed work. Neovim follows the same rule so the two terminal review
surfaces stay aligned. Source: `DI-davur`.

The CLI place/resource drilldown path follows the same rule too. It keeps
`show-place` and `show-resource` on the existing context-detail routes and
renders hierarchy, related runs, and typed links in a more useful terminal
layout instead of inventing a shell-only detail backend. Source: `DI-luzom`.

The matching CLI drilldown pass follows the same rule too. It keeps `show-run`,
`show-item`, and `show-responsibility` on the existing projected detail routes
and renders approvals, evidence, revisions, related runs, and typed links in
terminal-friendly summaries instead of falling back to raw JSON. Source:
`DI-salup`.

The next drilldown refinement follows that same rule too. It keeps run review
on the existing run-detail projection, but adds terminal handoff hints from the
run back into related item, place, resource, and responsibility context so
queue-driven review can keep walking the same projected model instead of
stopping at one record. Source: `DI-vunep`.

The next Neovim refinement applies the same idea inside the older inspectors:
whenever item, place, resource, or responsibility detail lists related runs,
the editor now emits direct `:OksInspectRun` hints so the user can keep walking
projected review context without leaving the current embodiment. Source:
`DI-josav`.

The matching CLI queue pass follows the same rule too. It keeps
`problem-review` on the existing grouped hotspot route and keeps
`pending-review` on the existing shared search projections, but renders both as
terminal-friendly summaries instead of raw JSON. Source: `DI-ravum`.

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
- `operational-run`
- `operational-place`
- `operational-resource`
- `knowledge-link`
- `knowledge-responsibility`
- `knowledge-search-metadata`

The current Go code now emits and verifies signed grid envelopes for the
`knowledge-item`, `knowledge-approval`, `knowledge-evidence`,
`operational-run`, `operational-place`, `operational-resource`,
`knowledge-link`, and `knowledge-responsibility` families, and it freezes
those eight families' runtime behavior by their shipped `pCID`s.
`knowledge-search-metadata` remains derived projection state rather than a
sixth signed family, and the local HTTP adapter is still not the full
PromiseGrid peer contract. Source: `DI-sobek`; `DI-mibor`; `DI-vosul`;
`DI-kavup`; `DI-votek`; `DI-sarib`; `DI-fusok`.

For the explicit statement of current shipped promises versus future
PromiseGrid-facing work, see
[PromiseGrid Implementation Claims](./promisegrid-implementation-claims.md).

The next staged PromiseGrid boundary is now also explicit: the first
relay-visible exchange slice carries `knowledge-item`, `knowledge-approval`,
`knowledge-evidence`, `knowledge-link`, `knowledge-responsibility`,
`operational-run`, `operational-place`, and `operational-resource`. The
shipped importer accepts origin-aware unseen history for those families into
non-empty runtimes, carries inline CID-keyed evidence blobs, and treats the
create-envelope CID as the durable peer-visible entity ID while preserving the
older short ID only as an alias. Source: `DI-guzab`; `DI-voruk`; `DI-vamok`;
`DI-faruv`; `DI-ruzok`; `DI-rumek`; `DI-loruk`; `DI-pivul`.

That later storage decision is now also staged: the first CAS pass dual-writes
signed family envelopes and copied evidence blobs into content-addressed
storage while keeping the current family logs and attachment tree during
migration. Source: `DI-ribek`.

That CAS pass now ships. Evidence records keep their compatibility attachment
paths, and they now also carry blob CIDs so later peer-visible evidence
exchange can bind to portable blob identity without losing the current read
model. The five frozen family envelopes also now replay/export authoritatively
from CAS, with one-time manifest backfill kept only for older runtime
migration. Source: `DI-lavuz`; `DI-rovud`.

Embodiment wording now tightens one step as well. Browser, CLI, and Neovim
remain adapter-delivered, but the adapter now exposes the shipped peer-exchange
format and CAS capability through runtime metadata and is described as an
adapter over that broader runtime. Source: `DI-vabek`; `DI-rovuz`.

What remains after that shipped runtime wave is narrower: CAS is still only
authoritative for the eight frozen families rather than every local
projection/runtime artifact. Source: `DI-lavuz`; `DI-tivor`; `DI-loruk`;
`DI-pivul`.

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
- embodiment/product follow-on work beyond the current local HTTP adapter
  contract

Those are still important future steps. The current pass focuses on a runnable
standalone operational-memory tool with one local runtime, a richer browser
embodiment, a thinner CLI embodiment, and a browser-only shared draft surface.

The current product direction is to keep that live-draft surface optional,
rather than making collaborative editing the core of the tool, and to revisit a
future richer Neovim embodiment later without porting the full `ex3`
websocket stack into `ex5` now. The current phase is intentionally a thin HTTP
live-draft client plus read-only item/run/entity inspection and search/browse
over projected detail. It now also includes a read-only pending-review queue
over the same shared search projections plus a narrow set of revision-safe
item/run approval and item supersede actions over the same local HTTP runtime
the CLI and browser already use. Source:
`DI-tabiv`; `DI-fudok`; `DI-lonuk`; `DI-ravok`; `DI-zalor`; `DI-givot`;
`DI-lorav`; `DI-vamor`; `DI-bafor`; `DI-pudor`; `DI-tivor`.

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
