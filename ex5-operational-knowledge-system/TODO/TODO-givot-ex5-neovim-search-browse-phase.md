# TODO givot - ex5 Neovim search and browse phase

## Decision Intent Log

ID: DI-givot
Date: 2026-07-21 09:25:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add a read-only Neovim search and browse phase over the existing `/api/search` projection before any write-side review or approval actions are added to the editor.
Intent: Give Neovim-heavy teams a practical way to discover relevant places, resources, responsibilities, items, and runs from inside the editor, while reusing the same projection and browse commands already established in the browser, CLI, and earlier Neovim phases.
Constraints: Stay on the current local HTTP runtime, do not reopen the websocket decision, keep the result surface read-only, and route deeper browsing through the existing `:OksInspect`, `:OksInspectRun`, and `:OksInspectEntity` commands.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md`, `ex5-operational-knowledge-system/TODO/TODO-givot-ex5-neovim-search-browse-phase.md`, `ex5-operational-knowledge-system/nvim/lua/oks/init.lua`, `ex5-operational-knowledge-system/nvim/assets_test.go`, `ex5-operational-knowledge-system/nvim/search_test.go`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/architecture.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`

## Goal

Add a read-only Neovim search result buffer over the existing `ex5` search
projection so users can discover and then inspect operational records without
leaving the editor.

## Tasks

- [x] givot.1 Define the next Neovim follow-on after typed-link browsing as read-only search and browse over `/api/search`.
- [x] givot.2 Add a Neovim search command that renders grouped result sections for places, resources, responsibilities, items, and runs.
- [x] givot.3 Show direct inspect hints in the search result buffer so users can jump into existing inspector commands.
- [x] givot.4 Add Neovim regression coverage for the search buffer and command markers.
- [x] givot.5 Update the ex5 docs to describe the new Neovim search/browse phase honestly.
