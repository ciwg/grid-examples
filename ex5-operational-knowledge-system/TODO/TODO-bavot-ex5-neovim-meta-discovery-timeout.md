# TODO bavot - ex5 Neovim meta discovery timeout

## Decision Intent Log

ID: DI-batov
Date: 2026-07-22 17:35:03 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Bound Neovim runtime-first `/api/meta` socket discovery with a short timeout and immediate repo-root socket fallback.
Intent: Keep the runtime as the preferred discovery truth source without letting a dead `OKS_BASE_URL` stall editor startup.
Constraints: Preserve runtime-first discovery and the direct local Unix-socket preference; avoid asynchronous transport switching during startup.
Affects: `ex5-operational-knowledge-system/nvim/lua/oks/init.lua`, `ex5-operational-knowledge-system/nvim/snapshot_test.go`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/TODO/TODO.md`, `docs/thought-experiments/TE-borav-ex5-neovim-meta-discovery-timeout.md`
Supersedes: DI-bavot

ID: DI-bavot
Date: 2026-07-22 17:23:59 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track the Neovim startup stall risk introduced by unbounded runtime-first `/api/meta` socket discovery.
Intent: Keep runtime-first discovery PromiseGrid-aligned while making Neovim startup fail fast and fall back predictably when the configured HTTP base URL is dead or unreachable.
Constraints: Preserve runtime-first discovery and the direct local Unix-socket preference; focus on discovery timeout and fallback behavior rather than transport redesign.
Affects: `ex5-operational-knowledge-system/nvim/lua/oks/init.lua`, `ex5-operational-knowledge-system/scripts/oks-nvim`, `ex5-operational-knowledge-system/nvim/snapshot_test.go`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-bavot-ex5-neovim-meta-discovery-timeout.md`

## Goal

Prevent `oks` Neovim startup from hanging on a dead or blackholed
`OKS_BASE_URL` while it is trying to discover the canonical socket path from
`/api/meta`.

## Tasks

- [x] bavot.1 Define the discovery timeout and fallback behavior for Neovim runtime-first socket resolution.
- [x] bavot.2 Implement bounded `/api/meta` discovery so startup can fall back locally when the HTTP path is unavailable.
- [x] bavot.3 Add regression coverage for dead-HTTP discovery with successful local fallback.

## Status

- closed
- resolved by short-timeout runtime discovery with immediate local fallback
