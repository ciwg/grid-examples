# TODO zalor - ex5 Neovim typed-link browsing phase

## Decision Intent Log

ID: DI-zalor
Date: 2026-07-21 08:55:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add read-only typed-link browsing in Neovim by exposing link sections in inspectors and a generic entity inspection command over existing detail APIs.
Intent: Let Neovim users follow operational context across linked items, runs, places, resources, and responsibilities without leaving the editor, while keeping the embodiment aligned with the current staged read-only approach.
Constraints: Reuse existing `GET /api/items/{id}`, `GET /api/runs/{id}`, `GET /api/places/{id}`, `GET /api/resources/{id}`, and `GET /api/responsibilities/{id}`; keep the feature read-only; update docs and tests in the same slice.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md`, `ex5-operational-knowledge-system/TODO/TODO-zalor-ex5-neovim-typed-link-browsing-phase.md`, `ex5-operational-knowledge-system/nvim/lua/oks/init.lua`, `ex5-operational-knowledge-system/nvim/assets_test.go`, `ex5-operational-knowledge-system/service/server_test.go`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/architecture.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`

## Goal

Add read-only typed-link browsing for Neovim inspectors so linked operational
entities are visible and inspectable from the editor.

## Intended repo paths

- `ex5-operational-knowledge-system/TODO/TODO.md`
- `ex5-operational-knowledge-system/TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md`
- `ex5-operational-knowledge-system/TODO/TODO-zalor-ex5-neovim-typed-link-browsing-phase.md`
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
  - purpose: fetch item detail for link-aware item inspection
  - lifecycle: runtime API call; no local file artifact

- `http://127.0.0.1:7045/api/runs/<run-id>`
  - class: `prod-data`
  - actions: `read`
  - purpose: fetch run detail for link-aware run inspection
  - lifecycle: runtime API call; no local file artifact

- `http://127.0.0.1:7045/api/places/<place-id>`
  - class: `prod-data`
  - actions: `read`
  - purpose: fetch place detail for generic entity inspection
  - lifecycle: runtime API call; no local file artifact

- `http://127.0.0.1:7045/api/resources/<resource-id>`
  - class: `prod-data`
  - actions: `read`
  - purpose: fetch resource detail for generic entity inspection
  - lifecycle: runtime API call; no local file artifact

- `http://127.0.0.1:7045/api/responsibilities/<responsibility-id>`
  - class: `prod-data`
  - actions: `read`
  - purpose: fetch responsibility detail for generic entity inspection
  - lifecycle: runtime API call; no local file artifact

- `t.TempDir()/**`
  - class: `test`
  - actions: `read/write`
  - purpose: server tests that validate link-bearing detail responses for the inspector
  - lifecycle: test-only; auto-cleaned by the Go test harness

## Tasks

- [x] zalor.1 Add link sections to Neovim inspectors.
- [x] zalor.2 Add a generic read-only Neovim entity inspector command for linked records.
- [x] zalor.3 Update docs and tests for typed-link browsing.
