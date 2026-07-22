# TODO faruv - make ex5 knowledge-evidence peer-visible and portable

## Decision Intent Log

ID: DI-zuvem
Date: 2026-07-22 11:01:56 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Do not make peer-visible `knowledge-evidence` depend on carried-along unfrozen `run_recorded` compatibility events. Freeze a PromiseGrid-native run family first, then implement peer-visible evidence exchange on top of that stronger run-context contract.
Intent: Keep the peer-visible PromiseGrid surface centered on frozen families and durable contracts rather than on an ad-hoc compatibility exception for evidence replay.
Constraints: The self-contained CID-keyed evidence-plus-blob carriage decision remains valid, but its implementation is blocked until the run-context family exists.
Affects: `ex5-operational-knowledge-system/TODO/TODO-faruv-ex5-peer-visible-evidence-exchange.md`, `docs/thought-experiments/TE-fubok-ex5-peer-visible-evidence-blob-carriage.md`, `docs/thought-experiments/TE-zuvem-ex5-peer-visible-evidence-run-context.md`, `ex5-operational-knowledge-system/TODO/TODO-zaruv-ex5-sixth-frozen-protocol-family-run.md`, `ex5-operational-knowledge-system/TODO/TODO.md`

ID: DI-faruv
Date: 2026-07-22 12:18:04 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Extend bootstrap peer exchange to carry signed `knowledge-evidence` records plus inline CID-keyed CAS blobs, and materialize imported attachment blobs back into a local compatibility attachment path during replay/import.
Intent: Make evidence actually peer-visible and usable on another host now, not merely signed locally or portable only on paper.
Constraints: The importer remains bootstrap-only. Blob carriage is self-contained inside the bundle for the first slice rather than a separate fetch protocol. Run context now comes from the frozen `operational-run` family instead of carried-along unfrozen compatibility events.
Affects: `ex5-operational-knowledge-system/service/peer_exchange.go`, `ex5-operational-knowledge-system/service/persistence.go`, `ex5-operational-knowledge-system/service/app.go`, `ex5-operational-knowledge-system/service/types.go`, `docs/thought-experiments/TE-fubok-ex5-peer-visible-evidence-blob-carriage.md`, `ex5-operational-knowledge-system/docs/promisegrid-peer-exchange-staging.md`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `ex5-operational-knowledge-system/TODO/TODO-zaruv-ex5-sixth-frozen-protocol-family-run.md`, `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Add peer-visible `knowledge-evidence` exchange for ex5, including portable
blob carriage that another host can actually resolve.

## Why this exists

The evidence family is frozen and signed locally, but it is still excluded from
peer exchange because attachment bytes are not yet carried in a peer-portable
way.

## Tasks

- [x] faruv.1 Run the required TE for evidence blob carriage and remote
  resolvability.
- [x] faruv.2 Lock how evidence metadata, attachment CIDs, and blob bytes move
  between peers.
- [x] faruv.3 Extend peer-exchange bundle rules to cover `knowledge-evidence`.
- [x] faruv.4 Implement evidence export/import over the settled portable blob
  carriage path.
- [x] faruv.5 Add round-trip, missing-blob, and tamper coverage for peer
  evidence exchange.

## Status

- closed
- `knowledge-evidence` is now peer-visible through bootstrap exchange with
  inline CID-keyed blob carriage
- `TE-fubok` completed; first blob carriage is locked to self-contained CID-keyed
  evidence plus blobs
- `TE-zuvem` completed; locked to freeze a run family first instead of carrying
  along unfrozen `run_recorded` compatibility events
- TODO `108` completed; imported evidence now anchors to the frozen
  `operational-run` family
