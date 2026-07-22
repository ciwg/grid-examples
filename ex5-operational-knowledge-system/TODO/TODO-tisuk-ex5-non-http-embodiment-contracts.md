# TODO tisuk - ex5 non-HTTP embodiment contracts

## Decision Intent Log

ID: DI-suvet
Date: 2026-07-22 15:41:15 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Log a fresh follow-on TODO for direct non-HTTP embodiment contracts after the current adapter-scoped PromiseGrid slice is treated as complete.
Intent: Preserve the current “complete for shipped scope” status while making a future non-HTTP embodiment expansion explicit and trackable.
Constraints: Do not treat this TODO as current migration debt; keep the present browser, CLI, and Neovim HTTP-adapter contract unchanged until a TE narrows a specific replacement or parallel embodiment path.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-tisuk-ex5-non-http-embodiment-contracts.md`

ID: DI-ronut
Date: 2026-07-22 16:23:58 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Implement `117A` as a terminal-first non-HTTP embodiment slice: CLI and Neovim move first, while the browser remains on the current local HTTP adapter as compatibility.
Intent: Remove HTTP-route dependence from the two terminal-facing embodiments without reopening browser transport, browser live-draft wiring, or relay transport scope in the same wave.
Constraints: Keep the browser usable through the current HTTP routes; do not claim a shared non-HTTP contract for all embodiments yet; keep the non-HTTP slice local-only.
Affects: `ex5-operational-knowledge-system/cmd/operational-knowledge/main.go`, `ex5-operational-knowledge-system/cmd/oks-cli/*`, `ex5-operational-knowledge-system/nvim/lua/oks/init.lua`, `ex5-operational-knowledge-system/service/*`, `ex5-operational-knowledge-system/docs/*`

ID: DI-favel
Date: 2026-07-22 16:23:58 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Use a Unix domain socket as the primary direct local contract for CLI and Neovim, carried as a JSON message stream over `.operational-knowledge-system/embodiment.sock`, while retaining HTTP as explicit browser and compatibility fallback.
Intent: Keep the first direct embodiment contract non-HTTP on the wire, long-lived, and local-runtime-scoped without inventing a second durable semantics model.
Constraints: Reuse the existing app semantics; support request/response operations and live-draft streaming over the same socket; do not route the browser through the new socket in this wave.
Affects: `ex5-operational-knowledge-system/cmd/operational-knowledge/main.go`, `ex5-operational-knowledge-system/cmd/oks-cli/*`, `ex5-operational-knowledge-system/nvim/lua/oks/init.lua`, `ex5-operational-knowledge-system/service/*`, `ex5-operational-knowledge-system/docs/*`

## Goal

Define direct non-HTTP embodiment contracts for `ex5` as a future-scope
expansion beyond the current local HTTP adapter.

## Tasks

- [x] tisuk.1 Run the required TE for the first non-HTTP embodiment slice.
- [x] tisuk.2 Lock the surviving embodiment-contract scope.
- [x] tisuk.3 Implement the chosen non-HTTP embodiment slice with matching tests and docs.

## Status

- completed
- `TE-noruk` completed
- `117A` locked
- Unix socket locked as the direct local contract for CLI and Neovim
- CLI and Neovim now prefer the local Unix-socket contract while browser remains on HTTP compatibility
