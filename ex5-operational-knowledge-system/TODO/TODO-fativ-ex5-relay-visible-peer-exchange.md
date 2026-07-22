# TODO fativ - decide and stage relay-visible ex5 PromiseGrid peer exchange

## Decision Intent Log

ID: DI-guzab
Date: 2026-07-22 10:22:17 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Stage the first relay-visible `ex5` peer-exchange slice around the four attachment-free families (`knowledge-item`, `knowledge-approval`, `knowledge-link`, and `knowledge-responsibility`) and defer peer-visible `knowledge-evidence` until the later CAS/blob-carriage decision.
Intent: Start real peer-visible PromiseGrid exchange on artifacts the runtime already signs and can represent portably, without pretending local evidence attachment paths are meaningful off-host.
Constraints: Do not bundle CAS-backed storage into this decision slice; do not tighten the embodiment contract in the same pass; document clearly that evidence stays local-only at peer-exchange time until a later storage/carriage decision lands.
Affects: `ex5-operational-knowledge-system/TODO/TODO-fativ-ex5-relay-visible-peer-exchange.md`, `docs/thought-experiments/TE-tavok-ex5-first-relay-exchange-scope.md`, `ex5-operational-knowledge-system/docs/promisegrid-peer-exchange-staging.md`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `ex5-operational-knowledge-system/docs/architecture.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Define the first relay-visible peer-exchange slice for ex5 now that the core
local durable families are frozen.

## Why this exists

The current shipped ex5 runtime is PromiseGrid-native only at the local durable
family layer. Relay-visible peer exchange is still explicitly not implemented,
and evidence attachment references are not yet portable across hosts.

## Tasks

- [x] fativ.1 Run the required TE for the first relay-visible exchange scope.
- [x] fativ.2 Lock what is exchanged first, by whom, and under what trust
  assumptions.
- [x] fativ.3 Define the first staged runtime and storage changes.
- [x] fativ.4 Add tracking docs for what becomes peer-visible.

## Status

- done
- first staged relay-visible exchange is the four attachment-free families
- peer-visible `knowledge-evidence` stays deferred until the later
  CAS/blob-carriage decision
