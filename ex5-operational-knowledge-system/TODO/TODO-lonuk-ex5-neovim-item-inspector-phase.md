# TODO lonuk - ex5 Neovim item inspector phase

## Decision Intent Log

ID: DI-lonuk
Date: 2026-07-20 23:05:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add a read-only Neovim inspector command that renders knowledge-item detail, revisions, approvals, and related runs from the existing item-detail HTTP API before attempting richer in-editor workflow actions.
Intent: Give Neovim users better operational inspection value immediately, while keeping the embodiment aligned with the current `ex5` HTTP runtime and avoiding a premature expansion into approval or run-creation actions.
Constraints: Reuse `GET /api/items/{id}`; keep the inspector read-only; update docs and tests in the same slice.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-lonuk-ex5-neovim-item-inspector-phase.md`, `ex5-operational-knowledge-system/TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md`, `ex5-operational-knowledge-system/nvim/lua/oks/init.lua`, `ex5-operational-knowledge-system/nvim/assets_test.go`, `ex5-operational-knowledge-system/service/server_test.go`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/architecture.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`

## Goal

Add a Neovim inspection surface for the current live-draft item so users can
review the surrounding operational record without leaving the editor.

## Intended repo paths

- `ex5-operational-knowledge-system/TODO/TODO.md`
- `ex5-operational-knowledge-system/TODO/TODO-lonuk-ex5-neovim-item-inspector-phase.md`
- `ex5-operational-knowledge-system/TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md`
- `ex5-operational-knowledge-system/nvim/lua/oks/init.lua`
- `ex5-operational-knowledge-system/nvim/assets_test.go`
- `ex5-operational-knowledge-system/service/server_test.go`
- `ex5-operational-knowledge-system/README.md`
- `ex5-operational-knowledge-system/docs/architecture.md`
- `ex5-operational-knowledge-system/docs/features-guide.md`
- `ex5-operational-knowledge-system/docs/http-api-guide.md`
- `ex5-operational-knowledge-system/docs/practical-implementation.md`

## Intended runtime path patterns

- `http://127.0.0.1:7045/api/items/<item-id>`
  - class: `prod-data`
  - actions: `read`
  - purpose: fetch the current projected item detail for the Neovim inspector
  - lifecycle: runtime API call; no local file artifact

- `t.TempDir()/**`
  - class: `test`
  - actions: `read/write`
  - purpose: server tests that validate the item detail shape the inspector depends on
  - lifecycle: test-only; auto-cleaned by the Go test harness

## Tasks

- [x] lonuk.1 Add a read-only Neovim inspector command for the current item.
- [x] lonuk.2 Cover the item-detail review shape in tests.
- [x] lonuk.3 Update docs to describe the richer Neovim inspection phase.
