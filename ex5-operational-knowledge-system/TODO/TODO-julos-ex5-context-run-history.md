# TODO julos - ex5 context run history

## Decision Intent Log

ID: DI-julos
Date: 2026-07-20 16:55:43 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add related run history to place, resource, and responsibility detail views so operators can inspect operational history from the surrounding context, not only from the item itself.
Intent: Make receiving areas, containers, machines, and responsibilities act more like operational anchors by showing the runs connected to them directly in the browser and API.
Constraints: Stay within the current local HTTP runtime and optional-collaboration product direction; update docs and tests in the same pass.
Affects: `ex5-operational-knowledge-system/service/**`, `ex5-operational-knowledge-system/web/**`, `ex5-operational-knowledge-system/docs/**`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Expose related run history from place, resource, and responsibility detail
views so operators can review operational history by context.

## Tasks

- [x] julos.1 Add related run history to place/resource/responsibility projections and API responses.
- [x] julos.2 Show related run history in the browser record inspector with drilldown actions.
- [x] julos.3 Add tests and docs for the context run-history flow.
