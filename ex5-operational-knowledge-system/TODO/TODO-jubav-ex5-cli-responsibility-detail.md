# TODO jubav - ex5 CLI responsibility detail inspection

## Decision Intent Log

ID: DI-jubav
Date: 2026-07-21 12:42:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add a CLI `show-responsibility` command that reuses the existing projected responsibility detail route.
Intent: Close a simple but real terminal inspection gap by letting shell-first operators inspect responsibility records, linked items, linked runs, and typed links directly from the CLI.
Constraints: Reuse the existing `/api/responsibilities/{id}` route, keep the command read-only, and avoid inventing a second responsibility-detail format.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-jubav-ex5-cli-responsibility-detail.md`, `ex5-operational-knowledge-system/cmd/oks-cli/main.go`, `ex5-operational-knowledge-system/cmd/oks-cli/main_test.go`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`

## Goal

Add a shell-first responsibility detail command over the existing projected
detail route.

## Tasks

- [x] jubav.1 Define the CLI responsibility detail command.
- [x] jubav.2 Reuse the existing projected responsibility detail route.
- [x] jubav.3 Add CLI regression coverage for the route.
- [x] jubav.4 Update the ex5 docs to describe the terminal-side responsibility detail behavior honestly.
