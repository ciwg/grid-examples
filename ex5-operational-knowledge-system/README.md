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
- typed links across responsibilities, items, and runs
- validated typed links across responsibilities, items, runs, places, and resources
- browser dashboard and forms
- browser search filters by kind, status, place, resource, and responsibility
- browser search filters by kind, status, outcome, place, resource, and responsibility
- browser free-text search now reaches run evidence facts and approval notes
- browser request failures stay inside the UI through the shared error path
- browser record inspector with summary cards and timelines
- browser review panels for item revisions, run evidence, and approvals
- item detail drilldown into the runs that used that item
- place, resource, and responsibility drilldown into related run history
- receiving check review panels for inbound inspection evidence and receiving history
- inventory audit review panels for discrepancy/count facts and audit history
- place/resource/responsibility context review panels that now surface receiving facts and inventory count/discrepancy facts from related runs
- one-click place/resource/responsibility drilldowns into filtered receiving history, count history, and receiving-problem history
- one-click place/resource/responsibility problem drilldowns that now use the same classification logic as grouped problem review
- grouped problem review that highlights repeated receiving and count issues by place and resource
- CLI inspection and creation commands
- CLI evidence upload for runs, with optional facts JSON and optional file attachment
- first-phase Neovim live-draft commands for opening, refreshing, inspecting,
  and pushing a knowledge item draft through the same local runtime
- Neovim item inspector for revisions, approvals, and related runs
- Neovim run inspector for direct evidence and approval review
- Neovim typed-link browsing over linked items, runs, places, resources, and responsibilities
- Neovim search/browse over grouped `/api/search` results with direct inspect hints
- Neovim pending-review view for draft items, unreviewed runs, and problem runs
- Neovim item approval action over the existing item approval API
- Neovim run approval action over the existing run approval API
- Neovim item supersede action over the existing item supersede API
- stub-backed headless browser smoke coverage for the shipped UI

For the longer feature walkthrough, see
[features guide](docs/features-guide.md).

## Current PromiseGrid Shape

The implemented foundation already follows the repo's actual grid direction:

- the durable contract is intended to live in protocol docs under
  `protocols/`
- the Go runtime owns append-only storage, projections, and the local adapter
  surfaces
- the browser and CLI are embodiments over the same shared operational state

What is implemented today:

- durable event history
- generic place hierarchy and resource records
- local HTTP API
- browser UI over that API, including shared working drafts for knowledge items
- CLI over that API, including place/resource commands
- versioned text bodies inside knowledge-item revisions
- local shared-draft persistence and live participant presence for browser
  editing

What is not yet implemented in this foundation:

- websocket-based collaboration transport
- peer-to-peer relay exchange
- signed grid envelopes on the wire
- ERP-style inventory quantities, reservations, or planning logic

That omission is intentional in the docs: this README describes the actual
implemented state, not the longer-term aspiration.

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
go run ./cmd/oks-cli resources
go run ./cmd/oks-cli new-resource alice container "RJ45 Bin" "Connectors bin" PLACE-0001
go run ./cmd/oks-cli responsibilities
go run ./cmd/oks-cli new-responsibility alice "Line lead" "Owns startup and approval routing"
go run ./cmd/oks-cli items
go run ./cmd/oks-cli new-item alice procedure "Start line A" "Startup procedure" "# Start line A"
go run ./cmd/oks-cli new-item alice receiving_check "Inspect inbound pallet" "Receiving check for inbound pallet" "# Inspect inbound pallet"
go run ./cmd/oks-cli new-item alice inventory_audit "Count RJ45 bin" "Cycle count for RJ45 connectors" "# Count RJ45 bin"
go run ./cmd/oks-cli record-run bob receiving_check RECV-0001 1 accepted_with_notes "Outer wrap torn" PLACE-0001 RES-0001
go run ./cmd/oks-cli record-run bob inventory_audit INV-0001 1 completed "Counted receiving bin" PLACE-0001 RES-0001
go run ./cmd/oks-cli approve-item PROC-0001 1 carol reviewer approved "Ready for use"
go run ./cmd/oks-cli approve-run RUN-0001 dave approver noted "Shift handoff recorded"
go run ./cmd/oks-cli add-evidence RUN-0001 dave "Dock photo" '{"result":"ok"}' ./evidence.txt
go run ./cmd/oks-cli search "supplier: Acme Parts & variance=-2"
go run ./cmd/oks-cli runs
```

## Terminal Behavior

`ex5` now has two terminal-facing embodiments over the same local runtime:

- the CLI for direct create/list/show/approve/search commands
- Neovim for live draft editing plus read-only review and browse flows

The intended terminal behavior today is:

- use the CLI for fast shell-oriented creation and direct inspection
- use the CLI when you need one-shot run evidence upload with optional facts
  JSON and optional copied file attachments
- use Neovim when you want to stay inside one editor session while:
  - editing a live draft
  - reviewing item and run detail
  - browsing linked entities
  - searching across the operational graph
  - opening a pending-review queue for draft items and runs that need attention

The terminal surface is intentionally staged, not fully symmetric yet:

- the CLI is better for direct command-style mutation
- Neovim is better for text work, read-only browsing, and reviewer flow
- both talk to the same `ex5` runtime and see the same projected state

That means a terminal-heavy operator can already do a large amount of real work
without opening the browser, while later follow-ons can still add narrower
workflow actions instead of trying to duplicate the whole browser at once.
Source: `DI-fudok`; `DI-givot`; `DI-lorav`.

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
- `:OksPending`
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

The run inspector adds a read-only split showing:

- run outcome and context
- evidence fact summaries
- run approvals

The typed-link phase adds:

- link sections inside item and run inspectors
- generic read-only inspection of linked `place`, `resource`, `responsibility`,
  `item`, and `run` records

The next read-only browse phase adds:

- grouped Neovim search results over the shared `/api/search` projection
- direct inspect hints for places, resources, responsibilities, items, and runs
- a read-only `oks-search://...` buffer so discovery can stay inside the editor

It still does not add write-side review or approval actions to Neovim. Source:
`DI-givot`.

The next terminal-first review phase adds:

- a read-only pending-review buffer for draft items
- a read-only pending-review buffer for unreviewed runs
- a read-only pending-review buffer for problem runs
- direct inspect hints so the next step stays inside the existing Neovim inspectors

It still keeps approval actions themselves out of Neovim for now. Source:
`DI-lorav`.

The next small write-side phase adds:

- one item approval action over the existing item approval API
- current-revision lookup before the approval is posted
- refresh of the current live, inspector, or pending-review context afterward

`OksApproveItem` uses the configured Neovim display name as the approval
`actor`. If you are on a live draft or item inspector, you can omit the item
ID and approve the current item directly. Source: `DI-vamor`.

The next matching phase adds:

- one run approval action over the existing run approval API
- direct use from the current run inspector or from an explicit run ID
- refresh of the current run or pending-review view afterward

`OksApproveRun` also uses the configured Neovim display name as the approval
`actor`. If you are on a run inspector, you can omit the run ID and approve
the current run directly. Source: `DI-bafor`.

The next lifecycle phase adds:

- one item supersede action over the existing item supersede API
- direct use from the current live draft or item inspector, or from an explicit
  item ID
- refresh of the current live, inspector, or pending-review view afterward

`OksSupersedeItem` uses the configured Neovim display name as the lifecycle
`actor`. If you are on a live draft or item inspector, you can omit the item
ID and supersede the current item directly. Source: `DI-pudor`.

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
- [Knowledge item protocol](protocols/knowledge-item.md)
- [Knowledge approval protocol](protocols/knowledge-approval.md)
- [Knowledge evidence protocol](protocols/knowledge-evidence.md)
- [Knowledge link protocol](protocols/knowledge-link.md)
- [Knowledge responsibility protocol](protocols/knowledge-responsibility.md)
- [Knowledge search metadata protocol](protocols/knowledge-search-metadata.md)

## Current direction

- keep the current local HTTP live-draft model instead of porting the full `ex3` websocket collaboration stack
- treat collaborative editing as optional rather than core to the product
- keep a richer future Neovim embodiment on the roadmap because it fits real team and customer workflows
