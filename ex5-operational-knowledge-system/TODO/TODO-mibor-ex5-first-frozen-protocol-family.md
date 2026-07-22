# TODO mibor - choose and freeze the first ex5 PromiseGrid protocol family

## Decision Intent Log

ID: DI-mibor
Date: 2026-07-21 00:00:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Use `knowledge-item` as the first ex5 frozen PromiseGrid family, publish the first implementation claim against it, and add the first local signed-envelope runtime slice for item create/revision/lifecycle events.
Intent: Start the real ex5 PromiseGrid-native runtime work with the central durable artifact family that already has the clearest semantics and the least avoidable cross-family ambiguity.
Constraints: Keep the first runtime slice additive and staged; preserve browser/CLI/Neovim behavior through current projections; leave approval, evidence, link, responsibility, and search-metadata families on the bridge layer for now.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-mibor-ex5-first-frozen-protocol-family.md`, `ex5-operational-knowledge-system/protocols/knowledge-item.md`, `ex5-operational-knowledge-system/CHANGELOG.md`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `ex5-operational-knowledge-system/service/**`, `docs/thought-experiments/TE-lafiz-ex5-promisegrid-wire-slice-decision.md`

## Goal

Choose the first narrow ex5 protocol family that should become a real frozen
PromiseGrid contract, then define the freeze/claim work needed before any
signed-envelope runtime slice opens.

## Tasks

- [x] mibor.1 Identify the best first ex5 protocol family candidate.
- [x] mibor.2 Define the exact contract boundary and what durable artifacts it owns.
- [x] mibor.3 Define the first implementation promise claim ex5 should publish against that frozen spec.

## Status

- done
- first frozen family, first claim, and first runtime slice implemented
