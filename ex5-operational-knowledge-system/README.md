# Ex5 Operational Knowledge System

The spec examples can look like separate problems:

- procedures handed from one operator to the next
- training records and micro-certifications
- maintenance logs and machine history
- approvals, evidence, and review
- search across operational knowledge later

They are all the same problem.

People need shared operational knowledge that stays attached to the actual work:
which procedure revision existed, who was responsible, what was done, what
evidence was captured, and what later approvals or follow-up happened.

`ex5-operational-knowledge-system` is the first example in this repo that tries
to solve that whole problem in one place. It is a durable operational memory
example with equal browser and CLI surfaces over one local Go runtime.

The current implementation keeps procedures, training content, maintenance
content, responsibilities, approvals, performed runs, evidence, and typed
links as append-only operational events projected into query views. Source:
`DI-radok`; `DI-kovup`; `DI-zuvob`.

## Features

- first-class responsibilities
- versioned knowledge items for:
  - procedures
  - training
  - maintenance
- append-only performed run records linked to exact revisions
- structured evidence with optional attachment upload
- named-role approvals with local team policy left outside the durable record
- typed links across responsibilities, items, and runs
- browser dashboard and forms
- CLI inspection and creation commands

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
- local HTTP API
- browser UI over that API
- CLI over that API
- versioned text bodies inside knowledge-item revisions

What is not yet implemented in this foundation:

- live collaborative document editing
- websocket awareness or presence
- peer-to-peer relay exchange

That omission is intentional in the docs: this README describes the actual
implemented state, not the longer-term aspiration.

## Runtime

By default the server stores runtime data under `.operational-knowledge-system/`:

- `events.jsonl`
  - append-only operational event log
- `attachments/`
  - copied evidence attachments grouped under per-run paths

## What You Need To Run

- Go
- a modern browser for the browser surface

Optional:

- a shell for CLI use

You do not need Node, npm, Docker, or Neovim to run the current `ex5`
foundation.

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
go run ./cmd/oks-cli responsibilities
go run ./cmd/oks-cli new-responsibility alice "Line lead" "Owns startup and approval routing"
go run ./cmd/oks-cli items
go run ./cmd/oks-cli new-item alice procedure "Start line A" "Startup procedure" "# Start line A"
go run ./cmd/oks-cli record-run bob procedure PROC-0001 1 completed "Completed startup cleanly"
go run ./cmd/oks-cli search startup
go run ./cmd/oks-cli runs
```

## Docs

- [Architecture notes](docs/architecture.md)
- [Features guide](docs/features-guide.md)
- [Practical implementation notes](docs/practical-implementation.md)
- [Knowledge item protocol](protocols/knowledge-item.md)
- [Knowledge approval protocol](protocols/knowledge-approval.md)
- [Knowledge evidence protocol](protocols/knowledge-evidence.md)
- [Knowledge link protocol](protocols/knowledge-link.md)
- [Knowledge responsibility protocol](protocols/knowledge-responsibility.md)
- [Knowledge search metadata protocol](protocols/knowledge-search-metadata.md)
