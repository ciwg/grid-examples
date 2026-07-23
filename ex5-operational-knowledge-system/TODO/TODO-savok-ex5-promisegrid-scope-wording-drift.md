# TODO savok - ex5 PromiseGrid scope wording drift

## Decision Intent Log

ID: DI-murev
Date: 2026-07-22 17:15:03 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Rewrite the top-level ex5 PromiseGrid summary docs so they directly describe the shipped signed-envelope and relay layers instead of retaining the older “local-runtime layer only” framing.
Intent: Make the repo's top-level PromiseGrid wording match the runtime that already ships frozen families, signed envelopes, local relay-feed exchange, and the dedicated remote relay binary.
Constraints: This is a wording honesty pass, not a new transport or protocol expansion.
Affects: `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `ex5-operational-knowledge-system/docs/product-overview.md`, `ex5-operational-knowledge-system/TODO/TODO.md`, `docs/thought-experiments/TE-pavok-ex5-promisegrid-scope-wording-drift.md`
Supersedes: DI-rubek

ID: DI-rubek
Date: 2026-07-22 19:12:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track stale PromiseGrid scope wording that still understates the shipped signed-envelope and relay layers.
Intent: Bring top-level ex5 PromiseGrid prose back into line with the runtime that now ships frozen families, signed envelopes, local relay-feed exchange, and the dedicated remote relay binary.
Constraints: This is an honesty/scope correction pass, not a new transport implementation wave.
Affects: `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `ex5-operational-knowledge-system/docs/product-overview.md`, `ex5-operational-knowledge-system/docs/*`, `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-savok-ex5-promisegrid-scope-wording-drift.md`

## Goal

Correct the remaining ex5 PromiseGrid wording that still describes the runtime
as pre-signed-envelope or pre-relay in places where those layers are already
shipped.

## Tasks

- [x] savok.1 Find and classify the remaining stale PromiseGrid scope statements across the top-level ex5 docs.
- [x] savok.2 Rewrite those statements so the shipped signed-envelope and relay scope is described consistently.
- [x] savok.3 Re-run a doc-alignment pass to confirm the PromiseGrid summary surfaces agree.

## Status

- closed
- resolved by rewriting the top-level PromiseGrid summary surfaces to describe the shipped signed-envelope and relay layers directly
