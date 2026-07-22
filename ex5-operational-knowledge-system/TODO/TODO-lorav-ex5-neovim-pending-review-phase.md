# TODO lorav - ex5 Neovim pending review phase

## Decision Intent Log

ID: DI-lorav
Date: 2026-07-21 10:05:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add a read-only Neovim pending-review view that groups draft items, unreviewed runs, and problem runs by reusing the existing search projections.
Intent: Give terminal-heavy operators and reviewers an immediate “what should I inspect next” surface inside Neovim without jumping to write-side approval actions yet.
Constraints: Stay on the current local HTTP runtime, keep the feature read-only, reuse `/api/search` instead of inventing a new endpoint, and route deeper inspection through the existing `:OksInspect`, `:OksInspectRun`, and `:OksInspectEntity` commands.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md`, `ex5-operational-knowledge-system/TODO/TODO-lorav-ex5-neovim-pending-review-phase.md`, `ex5-operational-knowledge-system/nvim/lua/oks/init.lua`, `ex5-operational-knowledge-system/nvim/assets_test.go`, `ex5-operational-knowledge-system/nvim/pending_test.go`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/architecture.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`

## Goal

Add a read-only Neovim pending-review buffer so terminal users can see draft
items and runs that likely need attention before choosing a deeper inspector.

## Tasks

- [x] lorav.1 Define the next Neovim follow-on after search/browse as a read-only pending-review view.
- [x] lorav.2 Add a Neovim command that groups draft items, unreviewed runs, and problem runs from the existing search projections.
- [x] lorav.3 Show direct inspect hints in the pending-review buffer so terminal users can jump into the existing inspectors.
- [x] lorav.4 Add Neovim regression coverage for the pending-review buffer and command markers.
- [x] lorav.5 Update the ex5 docs to describe the new terminal-first pending-review phase honestly.
