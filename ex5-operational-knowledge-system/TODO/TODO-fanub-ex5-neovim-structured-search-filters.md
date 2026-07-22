# TODO fanub - ex5 Neovim structured search filters

## Decision Intent Log

ID: DI-fanub
Date: 2026-07-21 23:58:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track Neovim-side structured search filters as the next terminal follow-on after grouped problem review.
Intent: Keep the remaining search-parity gap explicit now that CLI and browser already expose shared `/api/search` filters more broadly than `:OksSearch QUERY`.
Constraints: This TODO stays linked to deferred TODO `016`; it does not imply a new search API or browser-first work.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-fanub-ex5-neovim-structured-search-filters.md`, `ex5-operational-knowledge-system/TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md`, `ex5-operational-knowledge-system/nvim/lua/oks/init.lua`, `ex5-operational-knowledge-system/nvim/search_test.go`, `ex5-operational-knowledge-system/nvim/assets_test.go`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`, `ex5-operational-knowledge-system/docs/terminal-capability-matrix.md`, `ex5-operational-knowledge-system/docs/user-guide.md`

## Goal

Track the next Neovim search-parity gap after grouped problem review.

## Links

- Parent follow-on: `TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md` (`016`)

## Tasks

- [x] fanub.1 Define the desired Neovim structured-filter entrypoint over `/api/search`.
- [x] fanub.2 Decide how filter state should be expressed without bloating the editor command surface.
- [x] fanub.3 Add the implementation, tests, and docs in a later slice.
