# TODO tazum - ex5 runtime contract above adapters

## Decision Intent Log

ID: DI-tazum
Date: 2026-07-22 18:12:55 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track a future-scope pass to move more behavior out of adapter-shaped routes and into clearer shared runtime contract seams.
Intent: Reduce HTTP-route-shaped assumptions and make the runtime contract more directly PromiseGrid-aligned without immediately replacing every adapter.
Constraints: Preserve current shipped behavior during the investigation; this is a contract-shape refinement, not a broad rewrite yet.
Affects: `ex5-operational-knowledge-system/service/*`, `ex5-operational-knowledge-system/cmd/oks-cli/*`, `ex5-operational-knowledge-system/nvim/lua/oks/init.lua`, `ex5-operational-knowledge-system/web/app.js`, `ex5-operational-knowledge-system/docs/*`, `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Identify which current embodiment interactions still lean too heavily on
adapter-shaped route semantics and define the next shared runtime contract
surface above them.

## Tasks

- [ ] tazum.1 Map the current adapter-shaped seams that still carry too much runtime meaning.
- [ ] tazum.2 Define the next runtime-contract shape for those seams.
- [ ] tazum.3 Stage implementation and embodiment adoption order.

## Status

- open
- future-scope PromiseGrid refinement
