# TODO zaruv - freeze and claim the sixth ex5 PromiseGrid protocol family for runs

## Decision Intent Log

ID: DI-vamok
Date: 2026-07-22 11:44:12 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Freeze the sixth PromiseGrid family as `operational-run` and lock the first family scope to the richer run boundary. The family carries current `place_id`, `resource_ids`, and `responsibility_ids` exactly as they exist today, along with run identity, item identity, revision, actor, timestamp, outcome, notes, machine, and location.
Intent: Get `ex5` as close to PromiseGrid-ready as possible now by freezing a real operational execution record instead of a skeletal anchor that would need a second expansion later.
Constraints: Evidence, approvals, and links remain separate frozen families. Cross-peer tightening of place/resource/responsibility identity is deferred to later peer-stable identity work under TODO `107`.
Affects: `docs/thought-experiments/TE-vamok-ex5-run-family-boundary.md`, `ex5-operational-knowledge-system/protocols/operational-run.md`, `ex5-operational-knowledge-system/protocols/profiles.go`, `ex5-operational-knowledge-system/service/app.go`, `ex5-operational-knowledge-system/service/persistence.go`, `ex5-operational-knowledge-system/service/types.go`, `ex5-operational-knowledge-system/service/peer_exchange.go`, `ex5-operational-knowledge-system/service/server_test.go`, `ex5-operational-knowledge-system/service/app_test.go`, `ex5-operational-knowledge-system/TODO/TODO-faruv-ex5-peer-visible-evidence-exchange.md`, `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Freeze and implement a PromiseGrid-native run family so peer-visible evidence
can attach to a durable run-context contract instead of relying on unfrozen
compatibility events.

## Why this exists

`knowledge-evidence` already has a settled self-contained blob-carriage shape,
but imported evidence still depends on `run_recorded` context that is outside
the current peer-visible family set.

## Tasks

- [x] zaruv.1 Run the required TE for the run-family boundary and scope.
- [x] zaruv.2 Lock what the first frozen run family includes and excludes.
- [x] zaruv.3 Publish the implementation claim for the run family.
- [x] zaruv.4 Add the signed-envelope runtime slice for run records.
- [x] zaruv.5 Add replay/tamper coverage and update the PromiseGrid docs.
- [x] zaruv.6 Re-open TODO `105` implementation on top of the frozen run family.

## Status

- closed
- created because TODO `105` is locked to run-family-first rather than
  compatibility-event carry-along
- `TE-vamok` completed; `operational-run` and the richer first-family scope are
  locked by `DI-vamok`
