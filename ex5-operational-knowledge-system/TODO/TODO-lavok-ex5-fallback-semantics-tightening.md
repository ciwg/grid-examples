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

## Goal

Decide where remaining compatibility transport should stay explicit fallback,
where it should tighten further, and how operators should see those states.

## Tasks

- [ ] lavok.1 Review the remaining Neovim and browser fallback lanes and classify which are still necessary.
- [ ] lavok.2 Lock the explicit fallback policy per embodiment.
- [ ] lavok.3 Implement the chosen tightening and align tests/docs.

## Status

- open
- future-scope PromiseGrid refinement
