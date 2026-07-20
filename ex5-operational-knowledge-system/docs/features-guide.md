# ex5 feature guide

`ex5-operational-knowledge-system` takes a cluster of operational problems that
usually get split across wikis, ticket systems, spreadsheets, chat threads, and
local memory, and treats them as one problem:

- people need the same procedure text
- people need to know which revision was used
- people need evidence of what happened
- people need approvals and responsibility boundaries
- people need to find the same history later from either the browser or the CLI

This foundation solves that shared problem by keeping the durable record in one
append-only operational history and projecting it into equal browser and CLI
views. Source: `DI-radok`; `DI-kovup`; `DI-zuvob`.

## What the current feature set covers

### Responsibilities

Responsibilities are first-class durable records.

They exist so the app can answer questions like:

- who owns this procedure
- who reviews this kind of run
- which duties belong to training versus maintenance

Current implementation:

- create responsibilities from browser or CLI
- search responsibility titles and summaries
- link responsibilities to knowledge items and runs

### Knowledge items

Knowledge items are the shared operating documents.

Current knowledge item kinds:

- procedures
- training content
- maintenance content

Each item keeps:

- title and summary
- revisioned body text
- responsibility links
- approvals
- typed links
- append-only timeline

This is the current bridge between plain collaborative documents and structured
workflow state: the text stays readable, but the durable context around it is
not lost.

### Performed runs

Performed runs are the durable anchor for completed work.

A run answers:

- which item was used
- which exact revision was used
- who performed it
- what outcome was recorded
- what notes, machine, and location context applied

That is the main reason the spec’s separate use cases collapse into the same
core problem. Procedures, training, and maintenance all need an auditable
"what was done, from which revision, by whom, with what evidence" record.

### Evidence

Evidence adds facts and optional copied attachments to a run.

Current implementation supports:

- structured fact maps
- optional copied file attachments
- browser upload flow
- durable attachment path under the local runtime root

### Approvals

Approvals record review decisions without hardcoding one organization chart.

Current implementation uses:

- named roles
- actor
- decision
- notes
- target item or run

Team policy remains local and outside the durable event schema, which keeps the
model portable across organizations.

### Typed links

Typed links connect responsibilities, knowledge items, and runs.

This is how the app begins to behave like one operational knowledge system
instead of separate record silos. Links let users move through related context
instead of reconstructing the story manually from separate tools.

### Equal embodiments

The browser and CLI are meant to be equal first-class embodiments.

Current browser surface:

- dashboard
- create responsibility
- create knowledge item
- record run
- upload evidence
- record approval
- search

Current CLI surface:

- dashboard
- list and create responsibilities
- list and create knowledge items
- list and record runs
- approve items and runs
- search
- show individual items and runs

## What is intentionally not in this first foundation

The current implementation does **not** yet include:

- live collaborative editing
- websocket awareness or presence
- relay-to-relay peer exchange
- signed grid envelopes on the wire

Those are still important, but this first `ex5` pass is deliberately centered on
the durable operational memory layer first.

## Why this matters

The spec examples can sound different on the surface:

- procedures handed from one operator to the next
- training records
- maintenance history
- approvals and reviews
- search across operational knowledge

But they are all the same failure mode in disguise:

information about real work gets split across text, memory, tools, and time.

`ex5` is the first runnable example in this repo that treats that as one
problem and gives it one durable local history.
