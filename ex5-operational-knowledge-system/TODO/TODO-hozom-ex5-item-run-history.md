# TODO hozom - ex5 item run history

## Decision Intent Log

ID: DI-hozom
Date: 2026-07-20 16:43:46 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add related-run history to the ex5 item detail view so operators can move from a procedure, training item, maintenance doc, or inventory audit directly into the runs that used that exact item.
Intent: Make the operational-memory story easier to follow by showing not only revisions and approvals but also the performed history attached to an item.
Constraints: Keep the current local HTTP runtime; do not change the optional collaboration direction; update docs and tests in the same pass.
Affects: `ex5-operational-knowledge-system/service/**`, `ex5-operational-knowledge-system/web/**`, `ex5-operational-knowledge-system/docs/**`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Expose related run history from a knowledge item so the item detail view can
show the performed operational history tied to that item.

## Tasks

- [x] hozom.1 Add related run history to the knowledge item projection/API.
- [x] hozom.2 Show related runs in the browser record inspector with drilldown links.
- [x] hozom.3 Add tests and docs for the related-run history flow.
