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

ID: DI-monuv
Date: 2026-07-22 18:43:28 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Land a first typed local runtime contract on the direct Unix-socket embodiment path for the inspect/read slice only, using `type:"operation"` plus the named operations `inspect_item`, `inspect_run`, `inspect_entity`, `search`, `pending_review`, and `problem_review`.
Intent: Move selected terminal read workflows above route-shaped socket forwarding without forcing the browser off HTTP or rewriting the full adapter stack in one wave.
Constraints: Keep browser HTTP behavior unchanged; keep CLI/Neovim write and live-draft semantics unchanged; preserve existing payload shapes so terminal renderers keep working while the local socket contract becomes more runtime-native.
Affects: `ex5-operational-knowledge-system/service/local_socket*.go`, `ex5-operational-knowledge-system/service/app.go`, `ex5-operational-knowledge-system/service/types.go`, `ex5-operational-knowledge-system/cmd/oks-cli/*`, `ex5-operational-knowledge-system/nvim/lua/oks/init.lua`, `ex5-operational-knowledge-system/nvim/*_test.go`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `ex5-operational-knowledge-system/TODO/TODO.md`
Supersedes: DI-tazum

## Goal

Identify which current embodiment interactions still lean too heavily on
adapter-shaped route semantics and define the next shared runtime contract
surface above them.

## Tasks

- [x] tazum.1 Map the current adapter-shaped seams that still carry too much runtime meaning. See `../../docs/thought-experiments/TE-zoruk-ex5-runtime-contract-above-adapters.md`.
- [x] tazum.2 Define the next runtime-contract shape for those seams.
- [x] tazum.3 Stage implementation and embodiment adoption order.

## Status

- completed
- typed local runtime operations now back the selected terminal inspect/read slice
- TE `TE-zoruk` completed and implemented
