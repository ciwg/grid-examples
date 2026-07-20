# Operational Knowledge System

## A PromiseGrid Reference Application

**Status:** Draft 0.1

**Primary Language:** Go

**Interfaces**

- Go CLI (Unix pipes & filters)
- REST/HTTP API
- Web UI

## Vision

> **Every organization should become smarter every time someone does
> work.**

Every task completed, every repair performed, every product
manufactured, every improvement suggested, every lesson learned, and
every customer served is an opportunity for an organization to become
smarter.

Today's organizations lose knowledge every day. It lives in
conversations, emails, spreadsheets, disconnected software, and people's
memories.

> **Organizations should improve without forgetting.**

The Operational Knowledge System preserves organizational memory.

## Design Principles

1.  Every organization should become smarter every time someone does
    work.
2.  Organizations should improve without forgetting.
3.  Knowledge accumulates.
4.  Knowledge lives where the work happens.
5.  Preserve history instead of overwriting it.
6.  Every employee can contribute to the organization's memory.
7.  Support both search and contextual knowledge delivery.
8.  The interface is replaceable; the knowledge is not.

## Purpose

This reference application demonstrates how PromiseGrid can support a
real organization while showcasing immutable history, approval
workflows, organizational learning, and long-term operational memory.

It should not describe "grid technology" abstractly. In this repo, the
working model is already visible in the earlier examples:

- peer-visible meaning is selected by protocol CID (`pCID`)
- peer-visible traffic is carried as signed grid envelopes
- a local Go runtime owns signing, verification, append-only logging, and
  durable storage
- CLI, browser, and other embodiments are replaceable local views over the
  same shared contract
- durable workflow state is kept separate from ephemeral presence or local UI
  convenience state

`ex5` should build on that actual shape instead of inventing a separate one.

## Core Memories

### Physical Memory

-   Inventory
-   Tools
-   Machines
-   Locations

### Operational Memory

-   Procedures
-   Training
-   Maintenance
-   Quality

### Organizational Memory

-   Improvement proposals
-   Design reviews
-   Evidence
-   Lessons learned
-   Approvals

## Historical Preservation

Products are permanently linked to the exact procedure revision used to
build them. Organizations improve without losing the ability to
reproduce historical products years later.

## PromiseGrid Shape

This application should follow the same concrete PromiseGrid model already
proven in `ex1` through `ex4`.

### 1. Public protocol meaning is explicit

The peer-visible contract should not be "the HTTP API" or "the database
schema." It should be a small set of spec documents whose exact bytes select
meaning through `pCID`.

That means the shared operational record is not defined by:

- route names
- UI forms
- CLI flags
- internal Go structs by themselves

It is defined by named protocol docs plus signed envelopes that cite those
docs.

### 2. Messages are signed envelopes

Peer-visible records should use the same envelope discipline as the current
grid examples:

`grid([42(pCID), payload, proof])`

In practice this means:

- slot `0` selects the exact protocol spec
- slot `1` carries the payload for that spec
- slot `2` carries the signing proof

This keeps the contract auditable and stable across browser, CLI, and future
integrations.

### 3. The Go runtime is the relay/service boundary

The Go runtime should own:

- local identity and signing keys
- pCID discovery from exact local spec bytes
- signed envelope creation and verification
- append-only message logging
- current-state projection where needed for fast reads
- CAS-backed durable object storage
- local HTTP and CLI surfaces
- optional peer relay exchange

The CLI and Web UI should be embodiments over that runtime, not separate
protocol centers.

### 4. Durable history and ephemeral state are different

`ex5` is mostly about durable organizational history. It should therefore keep
its durable operational records separate from ephemeral collaboration state.

Durable examples:

- procedure revisions
- completion evidence
- approvals
- training records
- maintenance events
- improvement proposals
- inventory and machine state changes
- versioned links between products, procedures, and evidence

Ephemeral or embodiment-local examples:

- current browser form state
- currently selected filters
- local editor draft text
- viewport/layout preferences
- live presence or cursor awareness if a collaborative editor is later added

### 5. HTTP is an embodiment surface, not the peer contract

The local HTTP API can still exist and will likely be useful, but it should be
presented as a local adapter over the signed relay/runtime state, not as the
primary PromiseGrid-facing contract.

That is the same pattern already used in `ex3`.

## Candidate Protocol Families

This spec should talk in terms of protocol families instead of one giant
"knowledge graph" blob.

Reasonable first protocol families for `ex5` are:

- `knowledge-item`
  - durable records for procedures, lessons, training material, machine notes,
    and similar versioned knowledge artifacts
- `knowledge-approval`
  - approvals, rejections, supersedence, and review outcomes for knowledge
    items
- `knowledge-evidence`
  - durable evidence records tied to completed work, maintenance, inspections,
    training runs, or production events
- `knowledge-link`
  - typed links between products, procedure revisions, evidence, approvals,
    locations, tools, and machines
- `knowledge-search-metadata`
  - title, summary, tags, collections, archival state, and other latest-state
    metadata used for discovery

These names are placeholders, but the shape matters:

- small protocol families
- explicit pCID ownership
- append-only history first
- latest-state projection only where useful for operators

## Embodiment Topology

```text
                           PromiseGrid protocols
                                   │
                         signed grid envelopes
                                   │
                     local Go relay / business runtime
                    /              |                  \
                   /               |                   \
              Go CLI          local HTTP API          Web UI
                   \               |                   /
                    \              |                  /
                     operator workflows and local views
```

The CLI and Web UI are different views over the same operational state, but
the durable shared truth lives in the signed protocol artifacts and their
projections.

## Storage Model

The storage model should also follow the actual repo direction.

### Durable runtime storage

The local runtime should persist:

- signing identity
- append-only signed message log
- CAS-backed envelopes
- CAS-backed large attached objects when needed
- projected indexes for read/query convenience

Examples of large attached objects:

- procedure exports
- drawings
- training attachments
- machine manuals
- photos or inspection artifacts

### Embodiment-local state

The browser and CLI can keep local convenience state such as:

- last-opened queries
- recent filters
- local formatting preferences
- local cached exports

Those should not be confused with the PromiseGrid-facing organizational
record.

## What This Example Should Prove

If `ex5` is the next flagship reference app, it should prove the following
with real behavior rather than slogans:

1. A real organization can keep durable operational memory in signed,
   append-only, protocol-addressed records.
2. Procedures, evidence, approvals, and historical reproduction can all be
   linked without collapsing into one mutable row-per-object database story.
3. A Go CLI and a browser UI can both operate on the same shared operational
   state without becoming separate protocol systems.
4. Historical reproduction remains possible because products, events, and
   procedure revisions stay linked by durable artifact identity.

## Non-Goals For V0.1

To stay aligned with the actual technology already in the repo, this draft
should explicitly avoid overclaiming.

V0.1 should not assume:

- one giant globally synchronized "knowledge graph"
- direct browser-to-browser contract ownership
- that every visible UI feature is already a public protocol
- broad auth/federation/user-management solved up front
- that live collaboration and durable operational history are the same kind of
  record

## Practical Direction

The most realistic near-term shape for `ex5` is:

- one local Go runtime
- explicit protocol docs
- CLI + browser embodiments
- append-only durable operational events and artifacts
- projected query/read models for operators

That dovetails with the actual grid technology already demonstrated in this
repo instead of describing a separate imagined platform.

## Long-Term Vision

This project is intended to become the flagship PromiseGrid reference
application.
