# operational knowledge system architecture

`ex5-operational-knowledge-system` keeps durable operational history in one
local Go runtime and presents browser, CLI, and a first-phase Neovim
live-draft embodiment over that shared state. Source: `DI-radok`; `DI-fudok`.

## Topology

```text
             protocol docs / PromiseGrid public meaning
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

That top line is part of the PromiseGrid framing that ships with ex5. The
shipped runtime already uses append-only events and shared embodiments, and it
now implements five runtime-selected frozen `pCID` documents with signed grid
envelopes for `knowledge-item`, `knowledge-approval`, `knowledge-evidence`,
`knowledge-link`, and `knowledge-responsibility`. It still does not yet
implement the relay-visible peer-exchange and CAS-backed storage layers, and
search metadata remains derived projection state instead of a sixth durable
family. Source: `DI-sobek`; `DI-mibor`; `DI-vosul`; `DI-kavup`; `DI-votek`;
`DI-sarib`; `DI-fusok`.

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

- operator-oriented creation and mutation commands
- shell-only durable revision snapshots over the shared item revision route
- run evidence upload over the shared multipart evidence route
- typed-link creation over the shared validated graph route
- structured search and `problem=true` review over the shared search route
- pending-review and problem-review terminal summaries
- terminal drilldowns for place, resource, responsibility, item, and run detail

The current CLI surface is broad, but it is still not browser parity. It is a
terminal-first operational surface built from shared projected routes rather
than a second backend. Source: `DI-ravum`; `DI-salup`; `DI-zanub`; `DI-vuteg`.

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
- opening a read-only grouped search buffer from the projected `/api/search`
  API, then handing off deeper browsing to the existing inspectors
- extending that same grouped search buffer with the shared structured filter
  vocabulary the CLI already uses
- opening a read-only pending-review buffer by composing draft-item and
  run-review slices from the same search projection family
- opening a read-only grouped problem-review buffer over the shared hotspot
  projection
- posting durable item revision snapshots by flushing the live-draft body and
  then using the existing item revision route
- posting limited item approvals by reading the current item projection first
  and then using the existing item approval route
- posting limited run approvals through the existing run approval route
- posting limited item supersede actions through the existing item supersede route

Together, those two terminal embodiments now form a staged terminal-first
surface:

- CLI handles direct command-style mutation and shell inspection
- CLI now also handles shell-only durable revision snapshots
- CLI also assembles shell-facing pending and problem review summaries from the
  same shared projections Neovim already uses
- Neovim handles text editing, revision snapshots, review, browse,
  pending-work triage, limited approval actions, and limited supersede action
- both surfaces now provide deeper contextual drilldown from run and queue work
  into related item, place, resource, and responsibility records

They still share one runtime and one projection layer. Source: `DI-fudok`;
`DI-givot`; `DI-lorav`; `DI-vabok`; `DI-vunep`.

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
keys, timelines, and typed-link projections. They are not embedded only inside
procedure metadata, and they now expose the same `links` surface as the other
projected entity types. Source: `DI-luzaf`.

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
`GET /api/responsibilities/{id}`. Search/browse extends it to
`GET /api/search?q=...` for grouped discovery over the same projections.
Pending review extends it to filtered and unfiltered `/api/search` reads for
draft items plus run-review slices. Item approval extends it to one small
write-side action that still depends on the same `GET /api/items/{id}` truth
and the existing `POST /api/items/{id}/approvals` route. Source: `DI-fudok`;
`DI-lonuk`; `DI-ravok`; `DI-zalor`; `DI-givot`; `DI-lorav`; `DI-vamor`;
`DI-bafor`; `DI-pudor`.

## Current implementation note

The code currently implements:

- live shared drafting through the local HTTP runtime
- participant presence on the current draft
- durable versioned document bodies inside knowledge-item revisions
- one shared local HTTP embodiment adapter for browser, CLI, and Neovim
- one signed-envelope runtime slice for durable knowledge-item create/revision/lifecycle events
- one frozen `knowledge-item` `pCID` selected from the shipped protocol bytes

It still does **not** yet implement:

- websocket-based collaboration transport
- relay-visible peer exchange
- signed grid envelopes on the wire for the remaining ex5 families
- frozen runtime behavior selected by a shipped `pCID` for the remaining ex5 families

So in current ex5, protocol-family and `pCID` language are part of the shipped
PromiseGrid framing, and the `knowledge-item` family now also has a real
runtime/wire implementation. The remaining families are still on the staged
bridge layer. For the explicit current
claims list, see [PromiseGrid Implementation Claims](./promisegrid-implementation-claims.md).
