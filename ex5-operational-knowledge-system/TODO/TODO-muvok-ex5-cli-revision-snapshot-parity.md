# TODO muvok - ex5 CLI revision snapshot parity

## Decision Intent Log

ID: DI-muvok
Date: 2026-07-22 02:05:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track shell-only durable revision creation as its own `016` child TODO.
Intent: Close the remaining terminal authoring dead-end by letting a CLI-only operator cut a durable item revision without switching to the browser or Neovim.
Constraints: Reuse the existing item revision route and keep this scoped to CLI parity for durable snapshots, not to broader browser-parity workflow expansion.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-muvok-ex5-cli-revision-snapshot-parity.md`, `ex5-operational-knowledge-system/TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md`, `ex5-operational-knowledge-system/cmd/oks-cli/main.go`, `ex5-operational-knowledge-system/cmd/oks-cli/main_test.go`, `ex5-operational-knowledge-system/docs/terminal-capability-matrix.md`, `ex5-operational-knowledge-system/docs/user-guide.md`

## Goal

Add durable revision snapshot creation to the CLI embodiment.

## Links

- Parent follow-on: `TODO/TODO-nuvok-ex5-neovim-embodiment-followup.md` (`016`)

## Tasks

- [x] muvok.1 Define the smallest CLI command shape for durable revision creation.
- [x] muvok.2 Reuse the existing item revision route from the CLI.
- [x] muvok.3 Add CLI regression coverage and terminal docs.
