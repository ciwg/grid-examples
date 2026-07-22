# TODO movar - ex5 user documentation pass

## Decision Intent Log

ID: DI-movar
Date: 2026-07-21 22:10:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add a user-facing documentation pass for ex5 that includes a user guide, a product overview, and a terminal/Neovim capability matrix.
Intent: Make the current ex5 surface understandable to operators and reviewers without forcing them to reconstruct the product from the README, feature log, or source tree.
Constraints: Stay inside `ex5-operational-knowledge-system`, describe current shipped behavior rather than aspiration, and make the new docs easy to discover from the README.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-movar-ex5-user-documentation-pass.md`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/user-guide.md`, `ex5-operational-knowledge-system/docs/product-overview.md`, `ex5-operational-knowledge-system/docs/terminal-capability-matrix.md`

## Goal

Add a user-facing doc set that explains what ex5 does, how to use it, and what
the terminal surfaces can do today.

## Tasks

- [x] movar.1 Add an ex5 user guide for common operator workflows.
- [x] movar.2 Add an ex5 product overview that explains the system in plain language.
- [x] movar.3 Add an ex5 terminal/Neovim capability matrix.
- [x] movar.4 Link the new docs from the README and the ex5 TODO index.
