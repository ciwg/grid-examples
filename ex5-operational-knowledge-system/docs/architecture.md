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
                  browser UI
```

The runtime owns:

- append-only event persistence
- evidence attachment storage
- current-state projections for responsibilities, items, runs, approvals, and
  links
- local HTTP routes
- CLI-visible state through the same service model

The browser owns:

- dashboard presentation
- local form state
- search display

The CLI owns:

- operator-oriented creation and inspection commands
- shell-friendly JSON/text output

Source: `DI-radok`; `DI-zuvob`.

## Shared Model

The current implementation has three central durable record families:

### Responsibilities

Responsibilities are first-class records with their own IDs, summaries, role
keys, and timelines. They are not embedded only inside procedure metadata.

### Knowledge items

Knowledge items are versioned records for:

- procedures
- training content
- maintenance content

Each knowledge item keeps:

- current title and summary
- responsibility links
- revision history with body text
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
- linked responsibilities
- evidence
- approvals

## Evidence and attachments

Evidence is stored as structured summary plus facts, with optional copied file
attachments under the runtime root. The attachment path is runtime-local; the
durable story is the event and its copied artifact, not the original source
path on the operator machine.

## Current implementation note

The code currently implements durable versioned document bodies inside
knowledge-item revisions, but it does **not** yet implement live collaborative
editing or awareness transport. This doc describes the actual foundation in the
repo today.
