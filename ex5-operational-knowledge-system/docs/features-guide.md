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
append-only operational history and projecting it into shared browser and CLI
views over the same runtime. The embodiments are intentionally uneven today:
the browser is deeper, the CLI is thinner, and both read/write the same
durable model. Source: `DI-radok`; `DI-kovup`; `DI-zuvob`.

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
- record approvals from the CLI with the real approver identity instead of a
  hardcoded placeholder. Source: `DI-tarok`.

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

The approval path is now revision-aware. Approving a knowledge item only moves
it to `approved` when the approval targets the current revision, so an old
review cannot silently bless a newer draft. Source: `DI-dazim`.

The current durability pass also keeps larger stored knowledge bodies
replayable after restart, so the size of a legitimate procedure or audit body
does not quietly make the runtime unreadable later. Source: `DI-busor`.

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

That live endpoint now also lets clients intentionally clear a shared draft to
empty text while still keeping Neovim presence heartbeats body-neutral. Source:
`DI-dazim`.

### Neovim live draft phase 1

`ex5` now also has a first-phase Neovim embodiment for knowledge-item live
drafts.

It is intentionally narrower than the browser surface:

- open one item into an `acwrite` buffer
- refresh the shared live draft from the runtime
- push the current buffer body with `:write` or `:OksPush`
- inspect the current item/version/participant state with `:OksInfo`
- inspect projected item metadata, revisions, approvals, and related runs with `:OksInspect`
- inspect projected run evidence and approvals directly with `:OksInspectRun`
- inspect linked entities directly with `:OksInspectEntity TYPE ID`
- search grouped projected records with `:OksSearch QUERY`
- open a grouped pending-review buffer with `:OksPending`
- approve the current or specified item with `:OksApproveItem [ITEM_ID] ROLE DECISION [NOTES...]`
- approve the current or specified run with `:OksApproveRun [RUN_ID] ROLE DECISION [NOTES...]`
- supersede the current or specified item with `:OksSupersedeItem [ITEM_ID] [NOTES...]`
- publish presence and typing heartbeats over the same local HTTP live endpoint

Cursor and presence reporting now stay anchored to the live-draft window
instead of whichever split is currently focused, so opening inspectors does not
distort shared cursor offsets. Source: `DI-pazud`.

The close path is now explicit too. `:OksClose` tears down the live session by
wiping the live-draft buffer and any open read-only inspector buffer, so the
editor does not leave a detached `acwrite` buffer behind after the session
hooks are gone. Source: `DI-mabek`.

This phase deliberately reuses the existing live-draft API rather than adding a
separate websocket sidecar or remote-cursor renderer. The point is to give
Neovim-heavy teams a real operational embodiment now without reopening the
larger transport decision. Source: `DI-tabiv`; `DI-fudok`.

The first richer follow-on stays read-only on purpose. The inspector reads the
same item detail projection as the browser and CLI, so Neovim users can review
revision history, approvals, and related runs without pretending the editor
already supports every workflow action. Source: `DI-lonuk`.

The next richer follow-on keeps the same rule. Direct run inspection reads the
same run detail projection as the browser and CLI, so Neovim users can review
evidence facts and run approvals without adding write-side workflow actions to
the editor yet. Source: `DI-ravok`.

The next follow-on after that keeps the same read-only posture. Typed-link
browsing exposes link sections inside inspectors and lets Neovim jump to the
existing detail projections for places, resources, responsibilities, items, and
runs without inventing a second navigation model. Source: `DI-zalor`.

The next follow-on after that keeps the same rule too. Neovim search/browse
reads the existing `/api/search` projection and renders grouped places,
resources, responsibilities, items, and runs in a read-only search buffer,
with explicit hints for the inspect commands that already exist. It improves
discovery inside the editor without adding write-side workflow actions. Source:
`DI-givot`.

The next follow-on after that stays terminal-first and read-only. Neovim
pending review reuses the same search projections to group draft items,
unreviewed runs, and problem runs into one “what should I inspect next”
buffer, with direct hints for the existing inspectors. It improves reviewer
flow without adding write-side approval actions yet. Source: `DI-lorav`.

The next follow-on after that adds the first small write-side review action.
Neovim item approval resolves the current revision from the existing item
detail API, posts the approval through the existing item approval endpoint, and
then refreshes the relevant live, inspector, or pending-review context. That
keeps the action small and revision-safe instead of inventing a broader editor
workflow all at once. Source: `DI-vamor`.

The next follow-on after that adds the matching run-side action. Neovim run
approval posts directly through the existing run approval endpoint, uses the
current inspected run when possible, and then refreshes the run inspector or
pending-review queue. That gives terminal reviewers a direct next step after
finding run work in `:OksPending` without opening a broader editor workflow.
Source: `DI-bafor`.

The next follow-on after that adds the matching item lifecycle action. Neovim
item supersede posts through the existing item supersede endpoint, uses the
current item context when possible, and then refreshes the live draft,
inspector, or pending-review queue. That gives terminal reviewers the next
obvious item-state change without widening the editor into a broad mutation
surface. Source: `DI-pudor`.

### Terminal-first behavior

The current terminal behavior is intentionally split between CLI and Neovim
instead of forcing one tool to do everything badly.

CLI behavior today:

- fast shell-oriented create/list/show commands
- explicit record-run and approval actions
- explicit run evidence upload with optional facts JSON and optional attachment
- direct free-text search

Neovim behavior today:

- live draft editing for one knowledge item
- read-only item/run/entity inspection
- grouped search and browse over the operational graph
- grouped pending-review browsing for draft items and review-worthy runs
- limited item approval for the current or specified item
- limited run approval for the current or specified run
- limited item supersede for the current or specified item

That means terminal-first work in `ex5` currently follows a practical pattern:

- create or mutate directly from the CLI when a one-shot command is clearer
- stay in Neovim when you want to read, compare, inspect, browse, or queue
  review work inside one editor session

This is a deliberate staged embodiment strategy, not an accident. Source:
`DI-fudok`; `DI-givot`; `DI-lorav`.

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

Those problem drilldowns now use the same receiving/inventory problem logic as
the grouped hotspot review, rather than filtering only one receiving outcome.
Source: `DI-vemur`.

Run search now also reaches evidence summaries, evidence facts, and approval
notes. That means later operators can find runs by details like supplier names,
packing-slip identifiers, discrepancy facts, or recorded review notes instead
of only outcome and freeform run notes. Source: `DI-farun`.

The browser startup path is also hardened for restrictive/private environments.
If `localStorage` access is blocked or `crypto.randomUUID()` is unavailable,
the UI falls back to an in-memory participant identity instead of failing to
boot the live-draft surface. Source: `DI-mitob`.

The browser request path now uses the same shared error handler for create,
approval, evidence, refresh, and search flows. That keeps routine validation
failures inside the UI instead of depending on unhandled async behavior.
Source: `DI-ruvot`.

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

Evidence attachments are also stored immutably per upload, even if two entries
reuse the same filename. That preserves historical review instead of letting a
later upload replace the bytes behind an older record. Source: `DI-busor`.

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
- CLI upload flow for summary-only, summary-plus-facts, and
  summary-plus-facts-plus-attachment evidence entry
- durable attachment path under the local runtime root

Attachment uploads are enforced at the HTTP boundary: files larger than 8 MiB
are rejected instead of being partially accepted as evidence. Source:
`DI-navos`.

That same evidence route is now practical from a terminal-first workflow too.
Shell users can attach evidence with just a summary, with summary plus facts
JSON, or with summary plus facts JSON plus a copied attachment file without
opening the browser. Source: `DI-zanub`.

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

Typed links connect responsibilities, knowledge items, runs, places, and
resources.

This is how the app begins to behave like one operational knowledge system
instead of separate record silos. Links let users move through related context
instead of reconstructing the story manually from separate tools.

The link write path is now stricter too. A typed link must name a supported
endpoint type and an existing record ID, and responsibility detail now projects
its own `links` array the same way place, resource, item, and run detail
already do. Source: `DI-luzaf`.

### Equal embodiments

The browser and CLI are first-class embodiments over the same runtime, but they
are not equal in depth today.

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
layer that loads the shipped UI as a rendered page, but most of those tests use
stubbed `/api/*` responses rather than the full service stack. That is useful
coverage, but it is not a full browser-to-service integration suite yet.

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
- inspect projected run detail, evidence, and approvals
- inspect linked entities across existing detail APIs

This is a live-draft embodiment, not yet a full workflow embodiment. It does
not currently expose run creation or write-side approval actions directly
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
