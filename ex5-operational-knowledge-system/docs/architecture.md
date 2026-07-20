# operational knowledge system architecture

`ex5-operational-knowledge-system` keeps durable operational history in one
local Go runtime and presents equal browser and CLI embodiments over that
shared state.

## Topology

```text
                 protocol docs / pCID-selected meaning
                                  │
                       append-only operational events
                                  │
                    local Go runtime and projections
                       /                        \
                      /                          \
                 local HTTP API                 CLI
                      │
          browser UI + live draft / presence surface
```

The runtime owns:

- append-only event persistence
- shared working-draft persistence
- evidence attachment storage
- current-state projections for places, resources, responsibilities, items,
  runs, approvals, and links
- local HTTP routes
- CLI-visible state through the same service model

The browser owns:

- dashboard presentation
- local form state
- shared live drafting for knowledge-item bodies
- participant presence for the current draft
- revision snapshot actions
- search display

The CLI owns:

- operator-oriented creation and inspection commands
- shell-friendly JSON/text output

Source: `DI-radok`; `DI-zuvob`.

## Shared Model

The current implementation has five central durable record families:

### Places

Places are generic hierarchical operational locations. They can represent a
site, room, area, bench, rack, or bin without splitting the model into
separate domain-specific hierarchies.

### Resources

Resources are the operational things that belong to or move through those
places: machines, tools, containers, bins, parts, or fixtures.

### Responsibilities

Responsibilities are first-class records with their own IDs, summaries, role
keys, and timelines. They are not embedded only inside procedure metadata.

### Knowledge items

Knowledge items are versioned records for:

- procedures
- training content
- maintenance content
- inventory-audit content

Each knowledge item keeps:

- current title and summary
- compact status (`draft`, `approved`, `superseded`)
- responsibility links
- revision history with body text
- current shared working draft body
- approvals
- typed links
- append-only timeline

### Runs

Runs are the durable anchor for completed work. A run points to:

- exact knowledge item
- exact revision number
- actor
- outcome
- notes
- optional machine/location context
- optional place/resource context
- linked responsibilities
- evidence
- approvals

## Evidence and attachments

Evidence is stored as structured summary plus facts, with optional copied file
attachments under the runtime root. The attachment path is runtime-local; the
durable story is the event and its copied artifact, not the original source
path on the operator machine.

## Live draft subsystem

The browser-facing live draft subsystem is separate from durable revision
history.

It keeps:

- the current shared working body per item
- a monotonically increasing live version
- current participant presence for the item

Creating a new durable revision snapshots the current working body into the
append-only knowledge-item history. This keeps live drafting and auditable
revisions related but distinct.

## Current implementation note

The code currently implements:

- live shared drafting through the local HTTP runtime
- participant presence on the current draft
- durable versioned document bodies inside knowledge-item revisions

It still does **not** yet implement:

- websocket-based collaboration transport
- relay-visible peer exchange
- signed grid envelopes on the wire
