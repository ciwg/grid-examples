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
example with equal browser and CLI surfaces over one local Go runtime, plus a
first-phase Neovim live-draft embodiment for teams that work there. Source:
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
- compact item lifecycle:
  - `draft`
  - `approved`
  - `superseded`
- append-only performed run records linked to exact revisions
- structured evidence with optional attachment upload
- named-role approvals with local team policy left outside the durable record
- typed links across responsibilities, items, and runs
- browser dashboard and forms
- browser search filters by kind, status, place, resource, and responsibility
- browser search filters by kind, status, outcome, place, resource, and responsibility
- browser record inspector with summary cards and timelines
- browser review panels for item revisions, run evidence, and approvals
- item detail drilldown into the runs that used that item
- place, resource, and responsibility drilldown into related run history
- receiving check review panels for inbound inspection evidence and receiving history
- inventory audit review panels for discrepancy/count facts and audit history
- place/resource/responsibility context review panels that now surface receiving facts and inventory count/discrepancy facts from related runs
- one-click place/resource/responsibility drilldowns into filtered receiving history, count history, and receiving-problem history
- grouped problem review that highlights repeated receiving and count issues by place and resource
- CLI inspection and creation commands
- first-phase Neovim live-draft commands for opening, refreshing, inspecting,
  and pushing a knowledge item draft through the same local runtime
- Neovim item inspector for revisions, approvals, and related runs
- headless browser smoke coverage for the shipped UI

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
  - copied evidence attachments grouped under per-run paths

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
go run ./cmd/oks-cli search startup
go run ./cmd/oks-cli runs
```

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
- `:OksClose`
- `:write` pushes the current buffer body through the live-draft API

The inspector phase adds a read-only split showing:

- item status and summary
- revision history
- approvals
- related run history

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
- full run creation or typed-link navigation inside Neovim

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
