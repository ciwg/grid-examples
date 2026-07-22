# Ex5 Operational Knowledge System

The spec examples can look like separate problems:

- procedures handed from one operator to the next
- training records and micro-certifications
- maintenance logs and machine history
- inventory counts, bin audits, and receiving checks
- approvals, evidence, and review
- search across operational knowledge later

They are all the same problem.

People need shared operational knowledge that stays attached to the actual work:
which procedure revision existed, who was responsible, what was done, what
evidence was captured, and what later approvals or follow-up happened.

That includes inventory-shaped work when the real problem is operational memory:
how a bin was counted, which procedure was followed, who signed off on a
discrepancy, what evidence was captured, and how a later operator can find that
history again.

`ex5-operational-knowledge-system` is the first example in this repo that tries
to solve that whole problem in one place. It is a durable operational memory
example with shared browser and CLI embodiments over one local Go runtime, plus
a first-phase Neovim live-draft embodiment for teams that work there. Source:
`DI-fudok`.

The current implementation keeps procedures, training content, maintenance
content, receiving-check content, inventory-audit content, responsibilities,
places, resources, approvals, performed runs, evidence, live working drafts,
and typed links as append-only operational events plus local draft state
projected into query views. Source: `DI-radok`; `DI-kovup`; `DI-zuvob`;
`DI-foluk`; `DI-lusov`; `DI-zoruk`; `DI-vemok`.

## Features

- first-class responsibilities
- first-class places and resources
- versioned knowledge items for:
  - procedures
  - training
  - maintenance
  - receiving checks
  - inventory audits
- browser-shared live working drafts with participant presence and explicit
  revision snapshots
- browser startup falls back to ephemeral in-memory participant identity when
  storage or UUID helpers are restricted
- compact item lifecycle:
  - `draft`
  - `approved`
  - `superseded`
- append-only performed run records linked to exact revisions
- structured evidence with optional immutable attachment upload up to 8 MiB
- named-role approvals with local team policy left outside the durable record
- five frozen PromiseGrid-native runtime families so far:
  - `knowledge-item`
  - `knowledge-approval`
  - `knowledge-evidence`
  - `knowledge-link`
  - `knowledge-responsibility`
- typed links across responsibilities, items, and runs
- validated typed links across responsibilities, items, runs, places, and resources
- browser dashboard and forms
- browser search filters by kind, status, outcome, place, resource, and responsibility
- browser search now reaches record IDs directly across places, resources, responsibilities, items, and runs
- browser free-text search now reaches run evidence facts and approval notes
- browser problem-focused search now constrains context groups to records that are actually tied to matching problem runs
- browser request failures stay inside the UI through the shared error path
- browser now defaults to a review-first workspace instead of a flat wall of equally weighted forms
- browser review queue that switches between draft queue, problem hotspots, and known-record search
- browser current-record review surface with summary cards and timelines
- browser review panels for item revisions, run evidence, and approvals
- item detail drilldown into the runs that used that item
- place, resource, and responsibility drilldown into related run history
- receiving check review panels for inbound inspection evidence and receiving history
- inventory audit review panels for discrepancy/count facts and audit history
- place/resource/responsibility context review panels that now surface receiving facts and inventory count/discrepancy facts from related runs
- one-click place/resource/responsibility drilldowns into filtered receiving history, count history, and receiving-problem history
- one-click place/resource/responsibility problem drilldowns that now use the same classification logic as grouped problem review
- grouped problem review that highlights repeated receiving and count issues by place and resource
- context-driven browser operate actions that launch run recording, evidence capture, and approvals from the current item, run, place, resource, or responsibility
- a browser primary-flow layer that now points operators toward hotspots, draft review, and record-specific next steps instead of giving every surface equal weight
- browser task-oriented labels that reduce visible ontology in the main forms while keeping the same underlying records and IDs
- a sticky browser mode rail that keeps Review, Author, Operate, Create, and Browse behaviorally distinct without splitting the single-page shell
- a richer browser authoring surface with live draft metrics, writing focus, and quieter support disclosures around the editor
- browser operate-from-context guidance that launches work, evidence, and review from the current record instead of making operators start from blank forms
- task-first browser search presets for draft review, receiving problems, inventory counts, and broad run discovery, while keeping advanced filters available
- browser search and operate advanced panels now prefer helper selects and keep manual ID overrides behind a second disclosure
- stronger browser mode behavior that compresses inactive workspaces instead of leaving every surface equally expanded
- a calmer browser writing environment with collaboration settings moved behind disclosure and the writing surface visually prioritized
- a draft-first browser review queue that now serves as the clearest home path before operators branch into hotspots or known-record search
- one active browser review lane at a time so draft review, hotspots, and known-record search do not all compete at once
- staged browser operate flows that start from action choice buttons before revealing the heavier transaction form
- quieter browser authoring that keeps lifecycle and writing-context support behind disclosure until needed
- review-mode item inspection no longer silently loads the live draft until the operator explicitly enters Author mode
- CLI inspection and creation commands
- CLI typed-link creation over the validated operational graph
- CLI evidence upload for runs, with optional facts JSON and optional file attachment
- CLI structured search filters and problem-only review over the shared search projection
- CLI pending-review and problem-review terminal summaries for draft items, unreviewed runs, problem runs, and grouped hotspots
- CLI place/resource drilldown summaries over the shared context detail routes
- CLI run drilldowns with direct handoffs into related item, place, resource, and responsibility context
- first-phase Neovim live-draft commands for opening, refreshing, inspecting,
  and pushing a knowledge item draft through the same local runtime
- Neovim item inspector for revisions, approvals, and related runs
- Neovim run inspector for direct evidence and approval review
- Neovim typed-link browsing over linked items, runs, places, resources, and responsibilities
- Neovim search/browse over grouped `/api/search` results with direct inspect hints
- Neovim pending-review view for draft items, unreviewed runs, and problem runs
- Neovim run inspector handoffs into related item, place, resource, and responsibility inspectors
- Neovim item approval action over the existing item approval API
- Neovim run approval action over the existing run approval API
- Neovim item supersede action over the existing item supersede API
- stub-backed headless browser smoke coverage for the shipped UI, including explicit Author handoff, record-ID search, context-driven run logging, and durable browser snapshots

For the longer feature walkthrough, see
[features guide](docs/features-guide.md).

For operator-facing docs that describe the current product directly, see:

- [Product Overview](docs/product-overview.md)
- [User Guide](docs/user-guide.md)
- [Browser UI Guide](docs/browser-ui-guide.md)
- [Terminal Capability Matrix](docs/terminal-capability-matrix.md)

## Current PromiseGrid Shape

The implemented foundation ships with the PromiseGrid examples and dev guide,
and it follows that development model. The current ex5 runtime is still the
local-runtime layer of that model rather than the full signed-envelope / relay
layer. Source: `DI-sobek`.

What is already true today:

- the durable contract is intended to live in protocol docs under
  `protocols/`
- the Go runtime owns append-only storage, projections, and the local adapter
  surfaces
- the browser and CLI are embodiments over the same shared operational state

What the current runtime actually implements today:

- durable event history
- generic place hierarchy and resource records
- local HTTP API
- browser UI over that API, including shared working drafts for knowledge items
- CLI over that API, including place/resource commands
- versioned text bodies inside knowledge-item revisions
- local shared-draft persistence and live participant presence for browser
  editing
- one frozen `knowledge-item` PromiseGrid profile computed from the shipped
  protocol bytes
- one local signed-envelope runtime slice for durable knowledge-item
  create/revision/lifecycle events

What is not yet implemented in the shipped runtime:

- websocket-based collaboration transport
- peer-to-peer relay exchange
- signed grid envelopes for the remaining ex5 families
- frozen runtime behavior selected by shipped `pCID`s for the remaining ex5
  families
- ERP-style inventory quantities, reservations, or planning logic

That distinction is intentional in the docs: this README describes the actual
implemented ex5 runtime layer, not runtime behavior that has not shipped yet.

For the explicit technical statement of the current PromiseGrid boundary, see
[PromiseGrid Implementation Claims](docs/promisegrid-implementation-claims.md).

## Runtime

By default the server stores runtime data under `.operational-knowledge-system/`:

- `events.jsonl`
  - append-only operational event log
- `drafts/`
  - per-item shared working drafts used by the browser collaboration surface
- `attachments/`
  - immutable copied evidence attachments grouped under per-run paths

The runtime now keeps those attachment paths immutable and replays stored large
item/revision bodies within the current request envelope, so restart does not
silently invalidate either old evidence bytes or larger saved knowledge text.
Source: `DI-busor`.

The workflow-correctness pass also tightened two behavior edges:

- a knowledge item only moves to `approved` when the approval targets its
  current revision
- a live draft can now be intentionally cleared to empty text without
  confusing that edit with a presence-only heartbeat

Source: `DI-dazim`.

The browser request path is also now hardened for routine operator errors and
server validation failures. Create/search/review actions route failures through
the shared in-app error path instead of leaking them as unhandled async
rejections. Source: `DI-ruvot`.

## What You Need To Run

- Go 1.24.13
- a modern browser for the browser surface

Optional:

- a shell for CLI use

You do not need Node, npm, or Docker to run the current `ex5` foundation.

Neovim is optional. If you want the first-phase Neovim embodiment, you need a
local `nvim` binary in addition to the Go runtime.

This now matches the Go version pinned by the other `grid-examples` modules, so
you should not need a separate patch-level Go toolchain just to switch between
examples in this repo.

## Run

Start the server:

```bash
go run ./cmd/operational-knowledge
```

Then open:

```text
http://127.0.0.1:7045/
```

## CLI

The CLI talks to the same server:

```bash
go run ./cmd/oks-cli dashboard
go run ./cmd/oks-cli places
go run ./cmd/oks-cli new-place alice area Receiving "Inbound receiving and count area"
go run ./cmd/oks-cli show-place PLACE-0001
go run ./cmd/oks-cli resources
go run ./cmd/oks-cli new-resource alice container "RJ45 Bin" "Connectors bin" PLACE-0001
go run ./cmd/oks-cli show-resource RES-0001
go run ./cmd/oks-cli responsibilities
go run ./cmd/oks-cli show-responsibility RESP-0001
go run ./cmd/oks-cli new-responsibility alice "Line lead" "Owns startup and approval routing"
go run ./cmd/oks-cli items
go run ./cmd/oks-cli new-item alice procedure "Start line A" "Startup procedure" "# Start line A"
go run ./cmd/oks-cli new-item alice receiving_check "Inspect inbound pallet" "Receiving check for inbound pallet" "# Inspect inbound pallet"
go run ./cmd/oks-cli new-item alice inventory_audit "Count RJ45 bin" "Cycle count for RJ45 connectors" "# Count RJ45 bin"
go run ./cmd/oks-cli snapshot-item PROC-0001 alice "# Start line A\nAdd audited latch check"
go run ./cmd/oks-cli record-run bob receiving_check RECV-0001 1 accepted_with_notes "Outer wrap torn" PLACE-0001 RES-0001
go run ./cmd/oks-cli record-run bob inventory_audit INV-0001 1 completed "Counted receiving bin" PLACE-0001 RES-0001
go run ./cmd/oks-cli approve-item PROC-0001 1 carol reviewer approved "Ready for use"
go run ./cmd/oks-cli approve-run RUN-0001 dave approver noted "Shift handoff recorded"
go run ./cmd/oks-cli add-link alice responsibility RESP-0001 knowledge_item PROC-0001 owns "Primary startup procedure"
go run ./cmd/oks-cli add-evidence RUN-0001 dave "Dock photo" '{"result":"ok"}' ./evidence.txt
go run ./cmd/oks-cli search "supplier: Acme Parts & variance=-2" kind=receiving_check problem=true place_id=PLACE-0001
go run ./cmd/oks-cli problem-review
go run ./cmd/oks-cli pending-review
go run ./cmd/oks-cli runs
```

## Terminal Behavior

`ex5` now has two terminal-facing embodiments over the same local runtime:

- the CLI for direct create/list/show/approve/search commands
- Neovim for live draft editing, staged review/browse surfaces, and a small set
  of direct authoring and review actions

The intended terminal behavior today is:

- use the CLI for fast shell-oriented creation and direct inspection
- use the CLI when you need one-shot typed-link creation over validated
  existing records
- use the CLI when you need one-shot run evidence upload with optional facts
  JSON and optional copied file attachments
- use the CLI when you need a shell-only durable revision snapshot without
  opening Neovim
- use the CLI when you need the same structured or `problem=true` search
  slices that already drive browser and Neovim review views
- use the CLI when you need grouped hotspot review or projected responsibility
  detail from the same routes the browser and Neovim embodiments already read
- use Neovim when you want the same grouped hotspot review inside an editor
  scratch buffer with direct inspect handoffs for places, resources, and runs
- use the CLI when you need pending-review and problem-review queues rendered as
  terminal summaries instead of raw JSON
- use the CLI and Neovim review queues only against the shared search payloads
  that carry an explicit `approvals` array per run; omitted approvals are
  treated as contract failure instead of “unreviewed work”
- use the CLI when you need contextual place/resource drilldowns plus
  review-oriented run/item/responsibility detail with related runs, approvals,
  evidence, and link summaries in one terminal view
- use Neovim when you want to stay inside one editor session while:
  - editing a live draft
  - reviewing item and run detail
  - browsing linked entities
  - searching across the operational graph with shared structured filters
  - reviewing grouped hotspot problems
  - opening a pending-review queue for draft items and runs that need attention
  - cutting a durable revision snapshot from the current live draft
  - approving items or runs
  - superseding an item

The terminal surface is intentionally staged, not fully symmetric yet:

- the CLI is better for direct command-style mutation
- Neovim is better for text work, staged review surfaces, and reviewer flow
- both talk to the same `ex5` runtime and see the same projected state

That means a terminal-heavy operator can already do a large amount of real work
without opening the browser, while later follow-ons can still add narrower
workflow actions instead of trying to duplicate the whole browser at once.
Source: `DI-fudok`; `DI-givot`; `DI-lorav`; `DI-vabok`; `DI-muvok`.

## Neovim

The first Neovim phase is intentionally thin.

It reuses the same `GET/POST /api/items/{id}/live` surface as the browser live
draft studio instead of inventing a separate transport or porting the full
`ex3` websocket stack into `ex5`. Source: `DI-tabiv`; `DI-fudok`.

What it supports now:

- `:OksOpen ITEM_ID`
- `:OksRefresh`
- `:OksPush`
- `:OksInfo`
- `:OksInspect`
- `:OksInspectRun`
- `:OksInspectEntity TYPE ID`
- `:OksSearch QUERY`
- `:OksSearch QUERY [kind=VALUE] [status=VALUE] [outcome=VALUE] [place_id=VALUE] [resource_id=VALUE] [responsibility_id=VALUE] [problem=true]`
- `:OksPending`
- `:OksProblemReview`
- `:OksSnapshot`
- `:OksApproveItem [ITEM_ID] ROLE DECISION [NOTES...]`
- `:OksApproveRun [RUN_ID] ROLE DECISION [NOTES...]`
- `:OksSupersedeItem [ITEM_ID] [NOTES...]`
- `:OksClose`
- `:write` pushes the current buffer body through the live-draft API

`:OksClose` now tears down the active live-draft session by wiping the live
draft buffer and any open read-only inspector buffer, instead of only stopping
timers behind a still-visible detached buffer. Source: `DI-mabek`.

The Neovim embodiment now keeps cursor and presence offsets tied to the actual
live-draft window even after opening read-only inspector splits. Source:
`DI-pazud`.

The inspector phase adds a read-only split showing:

- item status and summary
- revision history
- approvals
- related run history
- direct `:OksInspectRun` hints from related runs

The run inspector adds a read-only split showing:

- run outcome and context
- evidence fact summaries
- run approvals

The typed-link phase adds:

- link sections inside item and run inspectors
- generic read-only inspection of linked `place`, `resource`, `responsibility`,
  `item`, and `run` records
- direct `:OksInspectRun` hints inside place/resource/responsibility related-run sections

Neovim also supports read-only search and browse over the shared
`/api/search` projection:

- grouped places, resources, responsibilities, items, and runs
- trailing shared `key=value` filters for `kind`, `status`, `outcome`,
  `place_id`, `resource_id`, `responsibility_id`, and `problem=true`
- direct inspect hints for the existing inspector commands
- a read-only `oks-search://...` buffer so discovery can stay inside the editor

It also supports a read-only pending-review buffer over the same shared search
projections:

- draft items
- unreviewed runs
- problem runs
- direct inspect hints so the next step stays inside the existing inspectors

And it now supports a read-only grouped problem-review buffer over the shared
`/api/problem-review` projection:

- grouped place hotspots
- grouped resource hotspots
- direct inspect hints for places, resources, and runs inside each hotspot

It now includes a small durable authoring action too:

- `:OksSnapshot`

That command requires an open live draft, flushes the current buffer body
through the shared live-draft endpoint, and then reuses
`POST /api/items/{id}/revisions` with the item's existing title, summary, and
tags. That lets a terminal-first author cut the durable revision without
leaving Neovim, while broader record creation, run entry, and evidence upload
still stay in the browser or CLI. Source: `DI-jabup`; `DI-vogar`.

And it now includes a narrow set of write-side authoring and review actions
over the existing HTTP runtime:

- `:OksSnapshot`
- `:OksApproveItem [ITEM_ID] ROLE DECISION [NOTES...]`
- `:OksApproveRun [RUN_ID] ROLE DECISION [NOTES...]`
- `:OksSupersedeItem [ITEM_ID] [NOTES...]`

Those actions stay intentionally small:

- revision snapshot reuses the existing item revision endpoint directly
- item approval resolves the current revision before posting approval
- run approval reuses the existing run approval endpoint directly
- item supersede reuses the existing item supersede endpoint directly
- each action refreshes the relevant live, inspector, or pending-review view afterward

All four use the configured Neovim display name as the `actor`, and the item
commands can omit the item ID when you are already on the current live draft or
item inspector. `:OksApproveRun` can omit the run ID when you are already on a
run inspector. Source: `DI-givot`; `DI-lorav`; `DI-jabup`; `DI-vamor`;
`DI-bafor`; `DI-pudor`.

Start it against a running server with:

```bash
./scripts/oks-nvim ITEM-0001
```

Optional environment variables:

- `OKS_BASE_URL`
- `OKS_DISPLAY_NAME`
- `OKS_COLOR`

What it does not try to do yet:

- websocket sync
- remote cursor rendering
- durable revision approval/review directly inside Neovim
- full run creation or write-side approval actions inside Neovim

## Docs

- [Architecture notes](docs/architecture.md)
- [Features guide](docs/features-guide.md)
- [HTTP API guide](docs/http-api-guide.md)
- [Practical implementation notes](docs/practical-implementation.md)
- [PromiseGrid peer-exchange staging](docs/promisegrid-peer-exchange-staging.md)
- [PromiseGrid CAS staging](docs/promisegrid-cas-staging.md)
- [PromiseGrid embodiment staging](docs/promisegrid-embodiment-staging.md)
- [Knowledge item protocol](protocols/knowledge-item.md)
- [Knowledge approval protocol](protocols/knowledge-approval.md)
- [Knowledge evidence protocol](protocols/knowledge-evidence.md)
- [Knowledge link protocol](protocols/knowledge-link.md)
- [Knowledge responsibility protocol](protocols/knowledge-responsibility.md)
- [Knowledge search metadata note](protocols/knowledge-search-metadata.md)

## Current direction

- keep the current local HTTP live-draft model instead of porting the full `ex3` websocket collaboration stack
- treat collaborative editing as optional rather than core to the product
- keep a richer future Neovim embodiment on the roadmap because it fits real team and customer workflows
- ship the first PromiseGrid peer exchange as bootstrap export/import for the four attachment-free signed families while deferring evidence carriage to the later CAS step
- dual-write signed envelopes and copied evidence blobs into additive CAS sidecar storage while keeping the current logs and attachment paths
- expose peer-exchange and CAS runtime capability metadata through the existing adapter so embodiments reflect the broader runtime contract without changing transport yet
