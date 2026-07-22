# Ex5 User Guide

This guide is for operators and reviewers who want to use
`ex5-operational-knowledge-system` as it exists today. It describes the current
product surface, not planned follow-on work. Source: `DI-movar`.

## Start The System

Run the server:

```bash
go run ./cmd/operational-knowledge
```

Then choose the surface you want:

- Browser: open `http://127.0.0.1:7045/`
- CLI: run `go run ./cmd/oks-cli ...`
- Neovim: use the shipped `oks` plugin/launcher against the same local server

All three surfaces talk to the same local runtime and projected state. Source:
`DI-fudok`; `DI-givot`; `DI-ravum`.

## Core Concepts

You will see the same core record types in every surface:

- Responsibility: who owns or reviews something
- Place: where work happens
- Resource: what bin, station, tool, or container is involved
- Knowledge item: a versioned operational record such as a procedure,
  receiving check, training item, maintenance item, or inventory audit
- Run: a performed event against an exact item revision
- Evidence: facts and optional attachment captured for a run
- Approval: a named review record for an item or run
- Link: a typed connection between records

Source: `DI-radok`; `DI-zuvob`; `DI-luzaf`; `DI-zanub`.

## Common Workflow 1: Set Up Context

Create the context records first:

```bash
go run ./cmd/oks-cli new-responsibility alice "Line lead" "Owns startup and approval routing"
go run ./cmd/oks-cli new-place alice area Receiving "Inbound receiving and count area"
go run ./cmd/oks-cli new-resource alice container "RJ45 Bin" "Connectors bin" PLACE-0001
```

Then inspect them:

```bash
go run ./cmd/oks-cli show-responsibility RESP-0001
go run ./cmd/oks-cli show-place PLACE-0001
go run ./cmd/oks-cli show-resource RES-0001
```

Those terminal drilldowns now summarize hierarchy, related runs, links, and
handoff commands instead of raw JSON dumps. Source: `DI-jubav`; `DI-luzom`;
`DI-salup`.

## Common Workflow 2: Create A Knowledge Item

Create a knowledge item for the work you want to preserve:

```bash
go run ./cmd/oks-cli new-item alice procedure "Start line A" "Startup procedure" "# Start line A"
go run ./cmd/oks-cli new-item alice receiving_check "Inspect inbound pallet" "Receiving check for inbound pallet" "# Inspect inbound pallet"
go run ./cmd/oks-cli new-item alice inventory_audit "Count RJ45 bin" "Cycle count for RJ45 connectors" "# Count RJ45 bin"
```

Use the browser or Neovim when you want richer ongoing text work on the item
body. Use the CLI when you want one-shot creation from the terminal. Source:
`DI-fudok`; `DI-vamor`; `DI-pudor`.

## Common Workflow 3: Record A Run

Once work is performed, record a run against an exact item revision:

```bash
go run ./cmd/oks-cli record-run bob receiving_check RECV-0001 1 accepted_with_notes "Outer wrap torn" PLACE-0001 RES-0001
go run ./cmd/oks-cli record-run bob inventory_audit INV-0001 1 completed "Counted receiving bin" PLACE-0001 RES-0001
```

This is how the system keeps procedural history attached to the exact revision
that was used. Source: `DI-zuvob`; `DI-vemok`.

## Common Workflow 4: Add Evidence And Links

You can add evidence from the terminal:

```bash
go run ./cmd/oks-cli add-evidence RUN-0001 dave "Dock photo" '{"result":"ok"}' ./evidence.txt
```

You can also connect records with typed links:

```bash
go run ./cmd/oks-cli add-link alice responsibility RESP-0001 knowledge_item PROC-0001 owns "Primary startup procedure"
```

Evidence uploads reuse the same multipart route as the browser, and typed links
reuse the same validated graph contract as the other surfaces. Source:
`DI-zanub`; `DI-vuteg`; `DI-luzaf`.

## Common Workflow 5: Review And Approve

CLI:

```bash
go run ./cmd/oks-cli snapshot-item PROC-0001 alice "# Start line A\nAdd audited latch check"
go run ./cmd/oks-cli approve-item PROC-0001 1 carol reviewer approved "Ready for use"
go run ./cmd/oks-cli approve-run RUN-0001 dave approver noted "Shift handoff recorded"
```

Neovim:

- `:OksSnapshot`
- `:OksApproveItem`
- `:OksApproveRun`
- `:OksSupersedeItem`

Snapshot, approval, and supersede actions stay on the existing shared HTTP API
and refresh the relevant terminal view afterward. `snapshot-item` lets a
CLI-only operator cut a durable revision by supplying the new body directly,
while `:OksSnapshot` requires an open live draft and snapshots the current
editor body using the item's existing title, summary, and tags. Source:
`DI-muvok`; `DI-jabup`; `DI-vamor`; `DI-bafor`; `DI-pudor`; `DI-dazim`.

## Common Workflow 6: Search And Triage

Use CLI search and queue views:

```bash
go run ./cmd/oks-cli search "supplier: Acme Parts & variance=-2" kind=receiving_check problem=true place_id=PLACE-0001
go run ./cmd/oks-cli pending-review
go run ./cmd/oks-cli problem-review
```

Use Neovim for the same staged review shape:

- `:OksSearch <query> [kind=...] [status=...] [outcome=...] [place_id=...] [resource_id=...] [responsibility_id=...] [problem=true]`
- `:OksPending`
- `:OksProblemReview`
- `:OksInspect`
- `:OksInspectRun`
- `:OksInspectEntity`

Pending-review queues in both CLI and Neovim now require an explicit
`approvals` array in the shared run payloads. If the shared projection omits
that field, the queue fails loudly instead of inventing fake unreviewed work.
Source: `DI-givot`; `DI-lorav`; `DI-ravum`; `DI-davur`.

## When To Use Which Surface

Use the browser when you want:

- the broadest operational surface
- shared live browser drafting
- contextual review panels
- grouped hotspot review

Use the CLI when you want:

- fast shell-first creation or mutation
- shell-only durable revision snapshots
- evidence upload
- terminal search and queue review
- one-shot drilldown commands

Use Neovim when you want:

- live draft editing
- durable revision snapshots from the current draft
- in-editor record inspection
- pending-review browsing
- approval or supersede actions without leaving the editor

Neovim is now enough for continuous item authoring and review, but browser/CLI
still own the broader create-run-evidence workflows. Source: `DI-jabup`;
`DI-vogar`.

Source: `DI-fudok`; `DI-givot`; `DI-lorav`; `DI-ravum`.

## Where To Read Next

- [Product Overview](./product-overview.md)
- [Terminal Capability Matrix](./terminal-capability-matrix.md)
- [Features Guide](./features-guide.md)
- [HTTP API Guide](./http-api-guide.md)
