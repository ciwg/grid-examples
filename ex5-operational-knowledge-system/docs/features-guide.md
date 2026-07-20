# ex5 feature guide

`ex5-operational-knowledge-system` takes a cluster of operational problems that
usually get split across wikis, ticket systems, spreadsheets, chat threads, and
local memory, and treats them as one problem:

- people need the same procedure text
- people need to know which revision was used
- people need repeatable inventory counts, audits, and discrepancy history
- people need evidence of what happened
- people need approvals and responsibility boundaries
- people need to find the same history later from either the browser or the CLI

This foundation solves that shared problem by keeping the durable record in one
append-only operational history and projecting it into equal browser and CLI
views. Source: `DI-radok`; `DI-kovup`; `DI-zuvob`.

It now also keeps a browser-shared working draft for each knowledge item,
separate from the durable revision history. That lets operators collaborate on
the current body text without pretending that live cursor state is itself the
auditable historical record. Source: `DI-lusov`; `DI-zoruk`.

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

### Places and resources

Places and resources provide reusable operational context.

Places are generic on purpose. The same model can represent:

- a site
- a room
- a receiving area
- a bench
- a rack
- a bin

Resources are the operational things that live in or move through those places:

- machines
- tools
- parts containers
- stock bins
- fixtures

Current implementation:

- create and list places from browser or CLI
- create and list resources from browser or CLI
- nest places by parent ID
- link resources to a place
- include place/resource context in run records and search results

### Knowledge items

Knowledge items are the shared operating documents.

Current knowledge item kinds:

- procedures
- training content
- maintenance content
- inventory-audit content

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

### Live draft studio

The browser now includes a live draft studio for knowledge items.

It supports:

- selecting a knowledge item
- editing the current shared working body
- seeing current participants
- refreshing the shared draft state
- snapshotting the current working body into a durable revision
- approving or superseding the current item

This is intentionally not described as a full CRDT or websocket editor. The
current implementation uses the shared local runtime and a live draft endpoint
with version checks and participant presence. That gives the browser a real
collaborative drafting surface while keeping the durable revision workflow
explicit and auditable.

### Record inspector and contextual navigation

The browser now also includes a record inspector.

It supports:

- opening a place, resource, responsibility, item, or run from the visible lists
- opening mixed search results directly into the inspector
- jumping from one record to related records, such as:
  - run -> item
  - run -> place
  - run -> resource
  - responsibility -> linked items and runs
  - place -> child places and resources

This does not replace richer future navigation, but it removes the need to
copy raw IDs manually just to understand the current operational graph.

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

Inventory-oriented work fits the same pattern when the load-bearing need is:

- which counting or receiving procedure was used
- which location, bin, or stock area was involved
- who performed the check
- what discrepancy or outcome was recorded
- what evidence and approvals followed

Current implementation also supports explicit place and resource context on a
run, so an inventory audit can say not just "a count happened" but "this count
happened in this receiving area against this bin/container context."

That does **not** make the current foundation a full inventory or MRP system. It
means inventory audits and related operational history belong to the same
operational-memory family as procedures, training, and maintenance.

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
- create places
- create resources
- create responsibility
- create knowledge item
- live draft studio for knowledge items
- record inspector with contextual navigation
- record run
- upload evidence
- record approval
- browse places and resources
- search

Current CLI surface:

- dashboard
- list and create places
- list and create resources
- list and create responsibilities
- list and create knowledge items
- list and record runs
- supersede items
- approve items and runs
- search
- show individual items and runs

## What is intentionally not in this first foundation

The current implementation does **not** yet include:

- websocket-based awareness or presence transport
- relay-to-relay peer exchange
- signed grid envelopes on the wire
- ERP-style quantity ledgers, reservations, purchasing, or planning logic

Those are still important, but this first `ex5` pass is deliberately centered on
the durable operational memory layer first.

## Why this matters

The spec examples can sound different on the surface:

- procedures handed from one operator to the next
- training records
- maintenance history
- inventory counts and bin audits
- approvals and reviews
- search across operational knowledge

But they are all the same failure mode in disguise:

information about real work gets split across text, memory, tools, and time.

`ex5` is the first runnable example in this repo that treats that as one
problem and gives it one durable local history.

## Still open

These product and architecture questions are still open:

- whether to fully port the `ex3` websocket collaboration model
- whether collaborative editing is truly core or optional
- whether `ex5` should eventually include another editor embodiment like Neovim
