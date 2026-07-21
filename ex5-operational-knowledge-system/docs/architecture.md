# operational knowledge system architecture

`ex5-operational-knowledge-system` keeps durable operational history in one
local Go runtime and presents browser, CLI, and a first-phase Neovim
live-draft embodiment over that shared state. Source: `DI-radok`; `DI-fudok`.

## Topology

```text
                 protocol docs / pCID-selected meaning
                                  │
                       append-only operational events
                                  │
                    local Go runtime and projections
                       /                        \
                      /                          \
                 local HTTP API            CLI
                      │
          browser UI + live draft / presence surface
                      │
                 Neovim live-draft client
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

The first-phase Neovim embodiment owns:

- opening one knowledge-item live draft into an editor buffer
- refreshing the shared draft state from the runtime
- pushing the current body back through the live-draft API
- reporting participant presence and typing state through the same local
  endpoint
- opening a read-only item inspector from the projected item-detail API
- opening a read-only run inspector from the projected run-detail API
- opening a generic read-only entity inspector for linked places, resources,
  responsibilities, items, and runs

It deliberately does not own a separate transport, remote cursor rendering, or
full workflow review surface in this phase. Source: `DI-fudok`.

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

The live draft subsystem used by the browser and first-phase Neovim embodiment
is separate from durable revision history.

It keeps:

- the current shared working body per item
- a monotonically increasing live version
- current participant presence for the item

Creating a new durable revision snapshots the current working body into the
append-only knowledge-item history. This keeps live drafting and auditable
revisions related but distinct.

In the current Neovim phase, the editor participates by polling and posting to
`GET/POST /api/items/{id}/live`, and it reads projected detail from
`GET /api/items/{id}` plus `GET /api/runs/{id}` for inspection. That keeps the
embodiment aligned with the same runtime truth the browser uses instead of
creating a second collaboration channel. Linked-entity browsing extends that
same rule to `GET /api/places/{id}`, `GET /api/resources/{id}`, and
`GET /api/responsibilities/{id}`. Source: `DI-fudok`; `DI-lonuk`;
`DI-ravok`; `DI-zalor`.

## Current implementation note

The code currently implements:

- live shared drafting through the local HTTP runtime
- participant presence on the current draft
- durable versioned document bodies inside knowledge-item revisions

It still does **not** yet implement:

- websocket-based collaboration transport
- relay-visible peer exchange
- signed grid envelopes on the wire
