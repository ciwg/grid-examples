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

ID: DI-rovik
Date: 2026-07-22 16:06:56 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Implement TODO `116` as a dedicated neutral relay service rather than a remote proxy of the current local adapter. The first slice lands as a separate binary at `cmd/operational-relay/main.go`, stores relay state under `.operational-relay/`, and exposes versioned remote routes under `/relay/v1`.
Intent: Keep runtime, embodiment adapter, and remote relay deployment cleanly separated so the first remote deployment slice stays aligned with the PromiseGrid role split.
Constraints: Do not merge this service back into `cmd/operational-knowledge`; do not reuse browser/CLI embodiment routes as the remote relay contract; keep the local adapter and remote relay as distinct operational surfaces.
Affects: `ex5-operational-knowledge-system/cmd/**`, `ex5-operational-knowledge-system/service/**`, `ex5-operational-knowledge-system/docs/**`, `ex5-operational-knowledge-system/TODO/TODO-ravel-ex5-remote-relay-deployment.md`

ID: DI-tasov
Date: 2026-07-22 16:06:56 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: The first relay service stores normalized per-origin append-only relay records and uses explicit per-origin cursor pulls that match the shipped local relay-feed request shape.
Intent: Keep the relay anchored to origin-aware signed artifact semantics instead of introducing relay-assigned sequencing as the main source of truth.
Constraints: The relay must not redefine durable identity or ordering; it must serve cursor-based pulls against stored origin-aware artifacts and remain compatible with the current `origin_peer_id + origin_sequence` model.
Affects: `ex5-operational-knowledge-system/service/**`, `ex5-operational-knowledge-system/docs/**`, `ex5-operational-knowledge-system/TODO/TODO-ravel-ex5-remote-relay-deployment.md`

ID: DI-nulav
Date: 2026-07-22 16:06:56 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: The first relay service is durable store-and-forward, and evidence-bearing feed publishes are accepted only after all referenced blobs are already staged into relay CAS.
Intent: Avoid half-valid relay state where signed evidence history is published before the referenced content-addressed blobs are durably present.
Constraints: No feed-first placeholder acceptance for missing blobs; blob staging must succeed before evidence-bearing relay publish succeeds.
Affects: `ex5-operational-knowledge-system/service/**`, `ex5-operational-knowledge-system/docs/**`, `ex5-operational-knowledge-system/TODO/TODO-ravel-ex5-remote-relay-deployment.md`

## Goal

Define a dedicated remote relay deployment shape for `ex5` beyond the current
local-adapter exchange layer.

## Tasks

- [x] ravel.1 Run the required TE for the first remote relay deployment slice.
- [x] ravel.2 Lock the surviving remote relay deployment scope.
- [x] ravel.3 Implement the chosen remote relay deployment slice with matching tests and docs.

## Status

- completed
- dedicated neutral remote relay service now ships as a separate binary
- `TE-vurek` completed
- `116B.1` locked and implemented
