# TODO fudok - ex5 Neovim live draft phase 1

## Decision Intent Log

ID: DI-fudok
Date: 2026-07-20 22:05:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Implement the first `ex5` Neovim embodiment as a thin plugin and launcher over the existing local HTTP live-draft API, without adding a websocket sidecar or porting the full `ex3` collaboration stack.
Intent: Give Neovim-heavy teams a real first operational embodiment for `ex5` that can open, refresh, inspect, and push live knowledge-item drafts while staying aligned with the current `ex5` collaboration decision.
Constraints: Use the existing `GET/POST /api/items/{id}/live` surface; keep Neovim phase 1 focused on item live drafts and presence, not remote cursor rendering; update docs and tests in the same slice.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-fudok-ex5-neovim-live-draft-phase1.md`, `ex5-operational-knowledge-system/TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md`, `ex5-operational-knowledge-system/nvim/plugin/oks.lua`, `ex5-operational-knowledge-system/nvim/lua/oks/init.lua`, `ex5-operational-knowledge-system/nvim/assets_test.go`, `ex5-operational-knowledge-system/scripts/oks-nvim`, `ex5-operational-knowledge-system/service/server_test.go`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`

## Goal

Add a usable first-phase Neovim embodiment for `ex5` that opens a knowledge
item's live draft, refreshes from the local runtime, pushes the current body
with `:write`, and exposes basic status/participant info.

## Intended repo paths

- `ex5-operational-knowledge-system/TODO/TODO.md`
- `ex5-operational-knowledge-system/TODO/TODO-fudok-ex5-neovim-live-draft-phase1.md`
- `ex5-operational-knowledge-system/TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md`
- `ex5-operational-knowledge-system/nvim/plugin/oks.lua`
- `ex5-operational-knowledge-system/nvim/lua/oks/init.lua`
- `ex5-operational-knowledge-system/nvim/assets_test.go`
- `ex5-operational-knowledge-system/scripts/oks-nvim`
- `ex5-operational-knowledge-system/service/server_test.go`
- `ex5-operational-knowledge-system/README.md`
- `ex5-operational-knowledge-system/docs/features-guide.md`
- `ex5-operational-knowledge-system/docs/http-api-guide.md`
- `ex5-operational-knowledge-system/docs/practical-implementation.md`

## Intended runtime path patterns

- `http://127.0.0.1:7045/api/items/<item-id>/live`
  - class: `prod-data`
  - actions: `read/write`
  - purpose: open, refresh, presence heartbeat, and push the Neovim live-draft body
  - lifecycle: runtime API call; no local file artifact

- `t.TempDir()/**`
  - class: `test`
  - actions: `read/write`
  - purpose: existing Go HTTP tests that validate the live endpoint behavior used by the plugin
  - lifecycle: test-only; auto-cleaned by the Go test harness

## Tasks

- [x] fudok.1 Add a repo-local Neovim launcher and plugin commands for opening and closing `ex5` live drafts.
- [x] fudok.2 Support refresh, write/push, and participant/status inspection over the existing live-draft HTTP API.
- [x] fudok.3 Add tests and docs for the first Neovim embodiment slice.
