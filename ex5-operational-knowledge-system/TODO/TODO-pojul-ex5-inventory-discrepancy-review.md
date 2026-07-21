# TODO pojul - ex5 inventory discrepancy review

## Decision Intent Log

ID: DI-pojul
Date: 2026-07-20 17:14:54 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add a dedicated inventory-audit review layer to the browser inspector so discrepancy/count facts and inventory-audit history are easier to read than generic evidence blobs.
Intent: Make inventory-shaped operator work feel like a first-class ex5 workflow without expanding into full ERP/MRP quantity logic.
Constraints: Stay within the current operational-memory model; use the existing evidence facts and related-run history; update docs and tests in the same pass.
Affects: `ex5-operational-knowledge-system/web/**`, `ex5-operational-knowledge-system/service/**`, `ex5-operational-knowledge-system/docs/**`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Improve inventory audit review by surfacing discrepancy/count facts and
inventory-audit history more clearly in the browser inspector and API-backed
tests.

## Tasks

- [x] pojul.1 Show inventory discrepancy/count facts clearly for `inventory_audit` runs in the browser inspector.
- [x] pojul.2 Highlight inventory-audit history from item, place, and resource context.
- [x] pojul.3 Add tests and docs for the inventory discrepancy review flow.
