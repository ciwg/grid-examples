# TODO lurek - ex5 socket contract root discovery

## Decision Intent Log

ID: DI-sorek
Date: 2026-07-22 17:15:03 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Make the runtime authoritative for terminal socket discovery by having CLI and Neovim query `/api/meta` first for the canonical socket path, then use the direct Unix-socket contract from there.
Intent: Eliminate custom-root ambiguity without demoting the terminal embodiments back to HTTP as their real transport.
Constraints: Keep the Unix socket as the preferred terminal contract; HTTP remains discovery and compatibility, not the durable terminal transport itself.
Affects: `ex5-operational-knowledge-system/cmd/oks-cli/main.go`, `ex5-operational-knowledge-system/nvim/lua/oks/init.lua`, `ex5-operational-knowledge-system/scripts/oks-nvim`, `ex5-operational-knowledge-system/service/app.go`, `ex5-operational-knowledge-system/service/types.go`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/TODO/TODO.md`, `docs/thought-experiments/TE-zurek-ex5-socket-contract-root-discovery.md`
Supersedes: DI-movik

ID: DI-movik
Date: 2026-07-22 19:12:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track non-default runtime-root discovery for the direct terminal socket contract as a dedicated ex5 follow-on.
Intent: Eliminate the remaining gap where CLI and Neovim cannot reliably discover the direct socket when `operational-knowledge` runs with a custom `-data-root`.
Constraints: Keep the current socket-first terminal embodiment slice; focus on making the contract self-locating rather than reopening transport selection.
Affects: `ex5-operational-knowledge-system/cmd/operational-knowledge/main.go`, `ex5-operational-knowledge-system/cmd/oks-cli/*.go`, `ex5-operational-knowledge-system/nvim/lua/oks/init.lua`, `ex5-operational-knowledge-system/scripts/oks-nvim`, `ex5-operational-knowledge-system/service/types.go`, `ex5-operational-knowledge-system/docs/*`, `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-lurek-ex5-socket-contract-root-discovery.md`

## Goal

Make the direct local Unix-socket embodiment contract self-locating even when
the runtime uses a non-default `-data-root`.

## Tasks

- [x] lurek.1 Define how terminal embodiments discover the intended socket path when the runtime root is not `.operational-knowledge-system/`.
- [x] lurek.2 Implement that discovery or advertisement path across CLI, Neovim, and the launcher.
- [x] lurek.3 Add regression coverage for non-default runtime-root terminal attachment.

## Status

- closed
- resolved by runtime-first `/api/meta` socket discovery for CLI and Neovim
