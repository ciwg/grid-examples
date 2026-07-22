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

- Browser: broadest operational surface, shared live draft editing, record
  inspection, contextual drilldowns, and grouped review panels.
- CLI: shell-first creation, inspection, search, evidence upload, approvals,
  and review queues.
- Neovim: live draft editing, structured search, grouped problem review,
  pending review, item/run approval, item supersede, and linked inspection
  inside one editor session.

These are not separate backends. They all read and write the same projected
state through the same local HTTP runtime. Source: `DI-fudok`; `DI-givot`;
`DI-lorav`; `DI-vamor`; `DI-bafor`; `DI-pudor`; `DI-ravum`.

## Current Boundaries

`ex5` is intentionally not trying to be:

- a websocket collaboration system yet
- a peer-to-peer relay
- a signed wire protocol implementation
- an ERP or MRP quantity-planning system

The current focus is durable local operational memory with strong browser,
terminal, and Neovim workflows. Source: `DI-tabiv`; `DI-ranor`; `DI-fudok`.

## Where To Read Next

- [User Guide](./user-guide.md)
- [Terminal Capability Matrix](./terminal-capability-matrix.md)
- [Features Guide](./features-guide.md)
- [HTTP API Guide](./http-api-guide.md)
- [Architecture](./architecture.md)
