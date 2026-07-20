# TODO foluk - operational knowledge system live workflow

## Decision Intent Log

ID: DI-foluk
Date: 2026-07-20 11:08:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Use one shared workflow model across procedures, training, maintenance, and inventory-audit work, with first-class generic places and resources instead of domain-specific warehouse-only entities.
Intent: Keep `ex5` broad enough for buildings, rooms, benches, storage areas, and tracked resources without hardcoding one physical layout vocabulary or turning the app into a specialized ERP system.
Constraints: Inventory remains an operational-memory slice, not a quantity-ledger; places must support hierarchy; resources must remain generic and linkable.
Affects: `service/types.go`, `service/app.go`, `service/server.go`, `cmd/oks-cli/main.go`, `web/index.html`, `web/app.js`, `docs/**`, `README.md`

ID: DI-lusov
Date: 2026-07-20 11:09:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add browser-only live collaboration for all knowledge-item document bodies while keeping CLI equal as an operational embodiment over the same durable runtime.
Intent: Let `ex5` support real shared drafting and operational editing in the browser without making Neovim or `ex3` a prerequisite for the broader operational-memory tool.
Constraints: Keep `ex5` standalone; do not require runtime imports from `ex3`; durable item revisions remain distinct from ephemeral live draft and awareness state.
Affects: `service/**`, `web/**`, `README.md`, `docs/**`

ID: DI-zoruk
Date: 2026-07-20 11:10:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Use a compact knowledge-item lifecycle of `draft`, `approved`, and `superseded`, with revision snapshots, approvals, and supersedence recorded in the shared workflow model.
Intent: Make review and historical replacement explicit without overcomplicating the first full `ex5` workflow implementation.
Constraints: Lifecycle applies across all knowledge-item kinds; approvals remain append-only records; supersedence must not rewrite history.
Affects: `service/types.go`, `service/app.go`, `service/server.go`, `cmd/oks-cli/main.go`, `web/index.html`, `web/app.js`, `service/*_test.go`, `README.md`, `docs/**`

## Goal

Implement the broader `ex5` operational-memory slice with generic places and
resources, lifecycle-aware knowledge items, and browser live collaboration for
knowledge-item documents.

## Tasks

- [x] foluk.1 Add place/resource entities, lifecycle state, and related search/link projections.
- [x] foluk.2 Add live draft and awareness state for browser collaboration on knowledge-item bodies.
- [x] foluk.3 Update CLI and browser surfaces to create, browse, and connect the broader operational model.
- [x] foluk.4 Extend tests and docs for the broader operational-memory workflow.
