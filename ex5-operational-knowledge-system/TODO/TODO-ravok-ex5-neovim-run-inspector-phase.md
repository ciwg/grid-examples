# TODO ravok - ex5 Neovim run inspector phase

## Decision Intent Log

ID: DI-ravok
Date: 2026-07-21 08:40:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add a direct read-only Neovim run inspector that renders run detail, evidence, and approvals from the existing run-detail HTTP API before attempting write-side workflow actions in the editor.
Intent: Give Neovim users direct review of operational execution records, not just item-linked summaries, while keeping the embodiment aligned with the current `ex5` runtime and staged read-mostly approach.
Constraints: Reuse `GET /api/runs/{id}`; keep the inspector read-only; update docs and tests in the same slice.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-ravok-ex5-neovim-run-inspector-phase.md`, `ex5-operational-knowledge-system/TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md`, `ex5-operational-knowledge-system/nvim/lua/oks/init.lua`, `ex5-operational-knowledge-system/nvim/assets_test.go`, `ex5-operational-knowledge-system/service/server_test.go`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/architecture.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`

## Goal

Add a direct Neovim run-inspection surface so users can review evidence and
approvals for a specific run without leaving the editor.

## Intended repo paths

- `ex5-operational-knowledge-system/TODO/TODO.md`
- `ex5-operational-knowledge-system/TODO/TODO-ravok-ex5-neovim-run-inspector-phase.md`
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

- `http://127.0.0.1:7045/api/runs/<run-id>`
  - class: `prod-data`
  - actions: `read`
  - purpose: fetch projected run detail for the Neovim run inspector
  - lifecycle: runtime API call; no local file artifact

- `t.TempDir()/**`
  - class: `test`
  - actions: `read/write`
  - purpose: server tests that validate the run-detail shape the inspector depends on
  - lifecycle: test-only; auto-cleaned by the Go test harness

## Tasks

- [x] ravok.1 Add a direct read-only Neovim run inspector command.
- [x] ravok.2 Cover run evidence and approval detail in tests.
- [x] ravok.3 Update docs to describe the richer Neovim run inspection phase.
