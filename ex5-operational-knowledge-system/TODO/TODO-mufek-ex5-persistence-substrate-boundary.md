# TODO mufek - ex5 persistence substrate boundary

## Decision Intent Log

ID: DI-mufek
Date: 2026-07-22 21:24:27 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track a future PromiseGrid alignment wave to decide whether the next proven reusable substrate slice after `promisegrid/records` and `promisegrid/transport` should be store and CAS persistence wiring.
Intent: Close the remaining gap where persistence mechanics still live entirely under `service/*`, while keeping the substrate boundary evidence-based instead of extracting app-specific projections prematurely.
Constraints: Preserve ex5 ownership of operational projections and workflows unless a narrower reusable persistence seam is clearly proven first.
Affects: `ex5-operational-knowledge-system/service/*`, `ex5-operational-knowledge-system/promisegrid/*`, `docs/thought-experiments/*`, `ex5-operational-knowledge-system/TODO/TODO.md`

ID: DI-lemor
Date: 2026-07-22 21:33:12 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Extract the next reusable persistence substrate as `promisegrid/store/`, carrying append-only JSONL log helpers, CAS object storage, and authoritative frozen-envelope hydration while leaving draft manifests, attachment rematerialization, and ex5 file-layout policy in `service/`.
Intent: Move the proven shared durability mechanics out of the ex5 app without accidentally freezing example-local storage policy into PromiseGrid substrate.
Constraints: Preserve current runtime and relay behavior; do not extract drafts, attachments, projections, or operator workflow storage semantics in this pass.
Affects: `docs/thought-experiments/TE-rakem-ex5-persistence-substrate-boundary.md`, `ex5-operational-knowledge-system/promisegrid/store/*`, `ex5-operational-knowledge-system/service/persistence.go`, `ex5-operational-knowledge-system/service/peer_exchange.go`, `ex5-operational-knowledge-system/service/relay_feed.go`, `ex5-operational-knowledge-system/service/relay_service.go`, `ex5-operational-knowledge-system/service/app_test.go`, `ex5-operational-knowledge-system/service/relay_server_test.go`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Decide whether store, CAS, and replay wiring should become the next reusable
PromiseGrid substrate slice beyond the current record and transport cores.

## Tasks

- [x] mufek.1 Identify which persistence behaviors are protocol-agnostic substrate versus ex5-specific storage policy. See `../../docs/thought-experiments/TE-rakem-ex5-persistence-substrate-boundary.md`.
- [x] mufek.2 Define the smallest reusable persistence boundary that does not drag projections or embodiment concerns into substrate. The locked boundary is `promisegrid/store/` for append-only logs, CAS objects, and authoritative frozen-envelope hydration only.
- [x] mufek.3 Align docs and claims after the reusable persistence slice leaves `service/`.

## Status

- completed
- created from the remaining “store wiring is still app-owned” PromiseGrid alignment gap after the `141`–`143` wave
- `promisegrid/store/` now owns append-only JSONL log helpers, CAS object storage, and authoritative frozen-envelope hydration, while drafts and attachment policy remain ex5-local
