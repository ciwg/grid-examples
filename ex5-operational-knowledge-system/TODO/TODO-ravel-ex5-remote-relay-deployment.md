# TODO ravel - ex5 remote relay deployment

## Decision Intent Log

ID: DI-lavek
Date: 2026-07-22 15:41:15 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Treat the current `ex5` PromiseGrid slice as complete for its shipped scope and open dedicated remote relay deployment as a new future-scope wave instead of as hidden remaining debt inside the current scope.
Intent: Separate “complete for current shipped scope” from “possible next expansion” so the repo can stay honest about what is done now versus what is a new product/runtime choice.
Constraints: Do not reopen the just-finished local-adapter relay-feed slice under TODO `115`; treat this TODO as a new expansion lane beyond the current adapter-scoped exchange model.
Affects: `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `ex5-operational-knowledge-system/docs/product-overview.md`, `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-ravel-ex5-remote-relay-deployment.md`

ID: DI-ponur
Date: 2026-07-22 15:41:15 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Log a fresh follow-on TODO for dedicated remote relay deployment beyond the current local-adapter exchange layer.
Intent: Keep remote relay work explicit as a new scope choice instead of implying that the current shipped ex5 layer is incomplete.
Constraints: This TODO is planning/decision scope only until a TE narrows the first remote relay slice.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-ravel-ex5-remote-relay-deployment.md`

## Goal

Define a dedicated remote relay deployment shape for `ex5` beyond the current
local-adapter exchange layer.

## Tasks

- [ ] ravel.1 Run the required TE for the first remote relay deployment slice.
- [ ] ravel.2 Lock the surviving remote relay deployment scope.
- [ ] ravel.3 Implement the chosen remote relay deployment slice with matching tests and docs.

## Status

- open
- future-scope follow-on beyond the current shipped PromiseGrid slice
- no TE filed yet
