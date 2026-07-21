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
- receiving-check content
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

### Neovim live draft phase 1

`ex5` now also has a first-phase Neovim embodiment for knowledge-item live
drafts.

It is intentionally narrower than the browser surface:

- open one item into an `acwrite` buffer
- refresh the shared live draft from the runtime
- push the current buffer body with `:write` or `:OksPush`
- inspect the current item/version/participant state with `:OksInfo`
- inspect projected item metadata, revisions, approvals, and related runs with `:OksInspect`
- publish presence and typing heartbeats over the same local HTTP live endpoint

This phase deliberately reuses the existing live-draft API rather than adding a
separate websocket sidecar or remote-cursor renderer. The point is to give
Neovim-heavy teams a real operational embodiment now without reopening the
larger transport decision. Source: `DI-tabiv`; `DI-fudok`.

The first richer follow-on stays read-only on purpose. The inspector reads the
same item detail projection as the browser and CLI, so Neovim users can review
revision history, approvals, and related runs without pretending the editor
already supports every workflow action. Source: `DI-lonuk`.

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
- viewing timeline entries for the selected record
- viewing summary cards for the selected record instead of only raw JSON
- reviewing item revisions and approvals directly in the inspector
- reviewing run evidence and approvals directly in the inspector
- reviewing the runs that used a selected item directly in the inspector
- reviewing related runs from places, resources, and responsibilities directly in the inspector

This does not replace richer future navigation, but it removes the need to
copy raw IDs manually just to understand the current operational graph.

### Structured search filters

The browser search surface now supports more than plain text.

Operators can combine a free-text query with filters for:

- `kind`
- `status`
- `outcome`
- `place_id`
- `resource_id`
- `responsibility_id`

That makes it easier to ask targeted questions like:

- which approved inventory audits belong to one responsibility
- which runs and resources belong to one place
- which records are still in `draft` status
- which receiving runs ended with `accepted_with_notes`

The record inspector now also uses those same filters for one-click context
drilldown from places, resources, and responsibilities. That makes it possible
to answer questions like:

- show me all receiving history here
- show me all counts for this bin
- show me all receiving problems in this area

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

Receiving and inbound inspection fit that same family too. `receiving_check`
lets the system represent intake work like inbound parts receipts, returned
items, tool intake, and staged kit inspections without pretending every
receiving event is only an inventory count.

The browser inspector now also treats receiving checks as a first-class review
surface:

- receiving check items show `Receiving history`
- receiving check runs show `Receiving review`
- place/resource/responsibility detail views show receiving history when those
  runs are part of the surrounding context

The browser inspector now also treats inventory audits as a first-class review
surface:

- inventory audit items show `Inventory count history`
- inventory audit runs show `Inventory discrepancy`
- place/resource/responsibility detail views show inventory audit history when
  those runs are part of the surrounding context

Context anchors now go one step further for both receiving and inventory work:

- place/resource/responsibility detail views show receiving fact history from
  related runs, not just bare run ids
- place/resource/responsibility detail views show inventory count/discrepancy
  fact history from related runs, not just bare run ids

The browser now also includes a grouped `Problem Review` surface:

- repeated receiving and inventory problems are summarized by place
- repeated receiving and inventory problems are summarized by resource
- operators can see hotspot counts before drilling into any single run
- each hotspot card links back into the existing inspector flow

This keeps the feature in the operational-memory lane. It does not turn `ex5`
into a quantity ledger or planning engine. It just makes repeated receiving and
count problems visible without forcing the user to rebuild them from raw runs.

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

The browser inspector now makes that review trail easier to read for both
knowledge items and runs, rather than forcing operators to inspect raw JSON.

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
- timeline-first record detail view
- item history drilldown into related runs
- context history drilldown from places, resources, and responsibilities
- receiving review for inbound inspection runs
- inventory discrepancy/count review for audit runs
- context review facts for receiving and inventory history from places,
  resources, and responsibilities
- record run
- upload evidence
- record approval
- browse places and resources
- search

Current browser regression coverage now also includes a headless browser smoke
test against a live test server, so the shipped UI is checked as a rendered
page instead of only as embedded asset text.

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

Current Neovim phase 1 surface:

- open a live draft by item ID
- refresh from the runtime
- push the current working body
- inspect current live participants and version state
- inspect projected item detail, revisions, approvals, and related runs

This is a live-draft embodiment, not yet a full workflow embodiment. It does
not currently expose approvals, run creation, or typed-link browsing directly
inside Neovim.

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

## Current direction

The current product direction is now:

- do not port the full `ex3` websocket collaboration model into `ex5`
- treat collaborative editing as optional rather than core
- keep a richer future Neovim embodiment as a desirable follow-on because many real teams operate there
