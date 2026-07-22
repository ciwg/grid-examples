# TODO navud - introduce peer-stable identity and ordering for ex5 peer exchange

## Decision Intent Log

ID: DI-ruzok
Date: 2026-07-22 11:27:18 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Lock the first peer-stable identity and ordering layer to explicit `origin_peer_id` plus `origin_sequence`. Treat that tuple as the durable replay and dedupe identity for multi-origin exchange, while keeping envelope CID as content identity and local `Sequence` as a compatibility projection only.
Intent: Move `ex5` toward honest multi-origin PromiseGrid behavior without relying on timestamp heuristics or one fake shared local sequence across peers.
Constraints: The runtime must preserve existing signed family payload semantics for already-written local histories, must keep the current bootstrap exchange working during migration, and must reopen TODO `103` on top of the settled origin-aware model instead of inventing a weaker stepping stone.
Affects: `docs/thought-experiments/TE-ravum-ex5-peer-stable-identity-and-ordering.md`, `ex5-operational-knowledge-system/service/types.go`, `ex5-operational-knowledge-system/service/app.go`, `ex5-operational-knowledge-system/service/persistence.go`, `ex5-operational-knowledge-system/service/peer_exchange.go`, `ex5-operational-knowledge-system/service/server_test.go`, `ex5-operational-knowledge-system/service/app_test.go`, `ex5-operational-knowledge-system/TODO/TODO-rumek-ex5-peer-exchange-beyond-bootstrap.md`, `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Define and implement peer-stable durable identity and ordering semantics so
`ex5` can accept non-bootstrap peer exchange into already-populated runtimes
honestly.

## Why this exists

The current compatibility event model still uses runtime-local event sequences
and runtime-local entity IDs, so arbitrary non-bootstrap import across peers
cannot be implemented honestly yet.

## Tasks

- [x] navud.1 Run the required TE for peer-stable identity and ordering across
  multi-origin runtimes.
- [x] navud.2 Lock the first durable identity layer for imported artifacts and
  compatibility replay.
- [x] navud.3 Lock how duplicate delivery, origin tracking, and ordering work
  once multiple peers contribute history.
- [x] navud.4 Implement the first peer-stable identity/order slice.
- [x] navud.5 Re-open TODO `103` implementation on top of that settled model.

## Status

- closed
- created because TODO `103` is now locked to solve peer-stable identity and
  ordering before non-bootstrap import
- `TE-ravum` completed; Alternative B is locked by `DI-ruzok`
- local events and signed family records now carry origin metadata; TODO `103`
  is implemented on top of that model
