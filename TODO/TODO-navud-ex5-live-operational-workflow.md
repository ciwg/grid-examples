# TODO navud - ex5 live operational workflow

## Decision Intent Log

ID: DI-navud
Date: 2026-07-20 11:07:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Expand `ex5` from the durable workflow foundation into a broader operational-memory app with browser live collaboration for all knowledge-item documents, one shared workflow model across procedures/training/maintenance/inventory-audit work, and first-class generic places and resources.
Intent: Finish the next meaningful `ex5` phase without making `ex2` or `ex3` prerequisites, while keeping the model generic enough for places, benches, rooms, storage areas, and tracked resources of many kinds.
Constraints: Keep `ex5` standalone; keep quantity-ledger/ERP/MRP logic out of scope; keep CLI equal as an operational embodiment even if live editing remains browser-only; keep browser collaboration subordinate to durable workflow history.
Affects: `ex5-operational-knowledge-system/**`, `TODO/TODO.md`, `TODO/TODO-navud-ex5-live-operational-workflow.md`

## Goal

Add the next complete `ex5` slice: generic place/resource modeling, compact
knowledge-item lifecycle, browser live collaboration for all knowledge-item
documents, and broader operational navigation/search across all domains.

## Tasks

- [x] navud.1 Add generic places and resources plus compact item lifecycle to the shared service model and HTTP/CLI surfaces.
- [x] navud.2 Add browser live collaboration and awareness for knowledge-item document bodies inside `ex5`.
- [x] navud.3 Update the browser UI to operate on places, resources, lifecycle, and collaborative knowledge items in one operational workspace.
- [x] navud.4 Expand tests and docs to cover the broader operational-memory slice.
