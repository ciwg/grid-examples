# Ex5 Product Overview

`ex5-operational-knowledge-system` is a local operational memory system. It is
meant for teams that need to keep procedures, training notes, maintenance
history, receiving checks, inventory audits, approvals, evidence, and typed
context links attached to the actual work that happened. Source: `DI-radok`;
`DI-foluk`; `DI-vemok`; `DI-fudok`.

## What The Product Does

The product keeps operational records as append-only events and then projects
them into views that operators can search, inspect, review, and act on. That
lets a team answer practical questions later: which revision existed, which run
used it, what evidence was captured, who approved it, and what place, resource,
or responsibility it belonged to. Source: `DI-radok`; `DI-zuvob`; `DI-luzaf`;
`DI-farun`.

Today the system stores and projects:

- responsibilities
- places and place hierarchy
- resources
- versioned knowledge items
- performed runs
- approvals
- evidence
- typed links
- live working drafts

Those records are all part of one local runtime instead of separate tools for
documents, review, and operational context. Source: `DI-radok`; `DI-foluk`;
`DI-zoruk`; `DI-fudok`.

## Main Workflows

The main workflow shape is:

1. Define context with responsibilities, places, and resources.
2. Create or revise a knowledge item such as a procedure, receiving check, or
   inventory audit.
3. Record a run against an exact revision.
4. Add evidence and approvals.
5. Search, review, and drill back into related context later.

That same workflow works for procedure history, receiving review, count
discrepancies, training artifacts, and maintenance logs because the durable
model is operational memory, not just a document editor. Source: `DI-vemok`;
`DI-zemok`; `DI-pogul`; `DI-fudok`.

## Embodiments

The current product has three practical embodiments over the same local Go
runtime:

- Browser: broadest operational surface, a review-first home with one active
  review queue lane at a time, shared live draft editing, staged operate
  actions, contextual drilldowns, and grouped review panels.
- CLI: shell-first creation, inspection, search, evidence upload, approvals,
  and review queues.
- Neovim: live draft editing, durable revision snapshots, structured search,
  grouped problem review, pending review, item/run approval, item supersede,
  and linked inspection inside one editor session.

These are not separate backends. They all read and write the same projected
state through one shared local runtime, with the browser staying on the HTTP
adapter while CLI and Neovim now prefer the direct local Unix-socket contract.
Source: `DI-fudok`; `DI-givot`; `DI-lorav`; `DI-vamor`; `DI-bafor`;
`DI-pudor`; `DI-ravum`; `DI-favel`.

## Current Boundaries

`ex5` is intentionally not trying to be:

- a peer-to-peer relay
- a signed wire protocol implementation
- an ERP or MRP quantity-planning system

The current focus is durable local operational memory with strong browser,
terminal, and Neovim workflows. Source: `DI-tabiv`; `DI-ranor`; `DI-fudok`.

## PromiseGrid Boundary

`ex5` ships with the PromiseGrid examples and follows that model, but the
shipped product currently implements the local runtime and local embodiment
layer rather than the full signed-envelope / relay layer. Source: `DI-sobek`.

What that means in practice:

- the current local HTTP API is the shipped embodiment adapter
- append-only operational history is real and durable today
- protocol-family and `pCID` language are part of the shipped PromiseGrid
  framing for ex5
- signed grid envelopes, frozen `pCID`-selected runtime behavior, origin-aware
  peer exchange, incremental relay feed plus CID-addressed blob transfer, and
  CAS-backed draft-body reload are now shipped runtime behavior

What is still not shipped in that layer:

- a browser-side non-HTTP embodiment contract

Those are now future-scope choices rather than missing pieces inside the
current shipped `ex5` PromiseGrid slice. Source: `DI-lavek`.

For the technical claims list, see
[PromiseGrid Implementation Claims](./promisegrid-implementation-claims.md).

Terminal authoring is now split more clearly too. Neovim can carry an item
from live draft editing into a durable revision snapshot without leaving the
editor, but broader create/run/evidence workflows still live in the browser or
CLI. Source: `DI-jabup`; `DI-vogar`.

## Where To Read Next

- [Browser UI Guide](./browser-ui-guide.md)
- [User Guide](./user-guide.md)
- [Terminal Capability Matrix](./terminal-capability-matrix.md)
- [Features Guide](./features-guide.md)
- [HTTP API Guide](./http-api-guide.md)
- [Architecture](./architecture.md)
