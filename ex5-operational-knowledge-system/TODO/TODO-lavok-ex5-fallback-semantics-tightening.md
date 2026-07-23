# TODO lavok - ex5 fallback semantics tightening

## Decision Intent Log

ID: DI-lavok
Date: 2026-07-22 18:12:55 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track a future-scope pass to tighten remaining compatibility fallback semantics after the CLI fail-closed change.
Intent: Push more embodiment behavior toward explicit primary contracts instead of silent or overly broad compatibility fallback.
Constraints: Treat CLI fail-closed as already settled; focus on remaining browser and Neovim fallback choices.
Affects: `ex5-operational-knowledge-system/nvim/lua/oks/init.lua`, `ex5-operational-knowledge-system/web/app.js`, `ex5-operational-knowledge-system/service/live_socket.go`, `ex5-operational-knowledge-system/service/local_socket.go`, `ex5-operational-knowledge-system/docs/*`, `ex5-operational-knowledge-system/TODO/TODO.md`

ID: DI-fonuv
Date: 2026-07-22 20:45:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Tighten only cross-adapter fallback, keeping browser websocket-to-HTTP fallback implicit inside `local_http` while making Neovim compatibility transport explicit through `oks-nvim --socket=off`.
Intent: Preserve resilience inside one embodiment adapter while stopping silent demotion from the direct local socket contract into browser-adapter compatibility transport.
Constraints: Do not change browser live-draft behavior in this pass; keep CLI fail-closed as already locked; expose the Neovim compatibility choice in operator-facing docs and metadata.
Affects: `ex5-operational-knowledge-system/nvim/lua/oks/init.lua`, `ex5-operational-knowledge-system/scripts/oks-nvim`, `ex5-operational-knowledge-system/service/app.go`, `ex5-operational-knowledge-system/service/server_test.go`, `ex5-operational-knowledge-system/nvim/snapshot_test.go`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Decide where remaining compatibility transport should stay explicit fallback,
where it should tighten further, and how operators should see those states.

## Tasks

- [x] lavok.1 Review the remaining Neovim and browser fallback lanes and classify which are still necessary. See `../../docs/thought-experiments/TE-zunav-ex5-fallback-semantics-tightening.md`.
- [x] lavok.2 Lock the explicit fallback policy per embodiment. Browser keeps implicit fallback inside `local_http`; Neovim compatibility now requires `oks-nvim --socket=off`.
- [x] lavok.3 Implement the chosen tightening and align tests/docs.

## Status

- closed
- PromiseGrid refinement completed
- TE `TE-zunav` completed and locked through `DI-fonuv`.
