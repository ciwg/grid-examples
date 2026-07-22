# TODO rumek - extend ex5 PromiseGrid peer exchange beyond bootstrap import

## Decision Intent Log

ID: DI-suvem
Date: 2026-07-22 10:53:34 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Do not implement lineage-limited continuation import as the next non-bootstrap step. Instead, solve peer-stable identity and ordering first, then return to non-bootstrap exchange on that stronger model.
Intent: Avoid a stepping-stone import mode that would still be weaker than the PromiseGrid target and would have to be caveated or superseded once true multi-origin semantics land.
Constraints: Keep the current bootstrap exchange intact; create a focused follow-on for peer-stable identity and ordering; do not claim non-bootstrap peer exchange before the new identity/order layer exists.
Affects: `ex5-operational-knowledge-system/TODO/TODO-rumek-ex5-peer-exchange-beyond-bootstrap.md`, `docs/thought-experiments/TE-vunok-ex5-non-bootstrap-peer-exchange-semantics.md`, `ex5-operational-knowledge-system/TODO/TODO-navud-ex5-peer-stable-identity-and-ordering.md`, `ex5-operational-knowledge-system/TODO/TODO.md`

ID: DI-rumek
Date: 2026-07-22 11:27:18 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Reopen non-bootstrap import on top of the settled origin-aware model. Accept unseen peer events into non-empty runtimes by deduping on `(origin_peer_id, origin_sequence)`, assign fresh local compatibility `Sequence` on acceptance, and reject imported create-event IDs that collide with an already-populated local namespace.
Intent: Ship the first honest ongoing peer-exchange slice now without pretending the repo has already solved cross-peer entity-ID reconciliation.
Constraints: Keep bootstrap import/export working. Preserve signed family content as authored by the exporting peer. Treat local `Sequence` as a projection order only. Leave true cross-peer entity namespace convergence to a later follow-on instead of silently rewriting or merging colliding IDs.
Affects: `docs/thought-experiments/TE-ravum-ex5-peer-stable-identity-and-ordering.md`, `ex5-operational-knowledge-system/service/peer_exchange.go`, `ex5-operational-knowledge-system/service/app.go`, `ex5-operational-knowledge-system/service/types.go`, `ex5-operational-knowledge-system/service/app_test.go`, `ex5-operational-knowledge-system/service/server_test.go`, `ex5-operational-knowledge-system/TODO/TODO-navud-ex5-peer-stable-identity-and-ordering.md`, `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Move ex5 from bootstrap-only peer exchange toward an ongoing peer-visible
exchange model that can accept valid signed family artifacts into a non-empty
runtime.

## Why this exists

The shipped peer exchange is still bootstrap-only and requires an empty
runtime, which means ex5 can clone a peer state but cannot yet behave like an
ongoing multi-peer grid node.

## Tasks

- [x] rumek.1 Run the required TE for non-bootstrap peer import semantics.
- [x] rumek.2 Lock how ex5 detects duplicates, conflicts, and replay across
  already-populated runtimes.
- [x] rumek.3 Define the first non-bootstrap import/export contract.
- [x] rumek.4 Implement import into non-empty runtimes for the currently
  portable families.
- [x] rumek.5 Add coverage for duplicate delivery, replay, and mixed local plus
  imported artifact histories.

## Status

- closed
- non-bootstrap peer exchange now accepts origin-aware unseen histories into
  non-empty runtimes when local entity namespaces do not collide
- `TE-vunok` completed; locked to solve peer-stable identity and ordering
  before non-bootstrap import implementation
- TODO `107` completed; `DI-rumek` locks the first non-bootstrap contract
