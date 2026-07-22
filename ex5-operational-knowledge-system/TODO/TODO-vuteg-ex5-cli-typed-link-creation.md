# TODO vuteg - ex5 CLI typed-link creation

## Decision Intent Log

ID: DI-vuteg
Date: 2026-07-21 12:05:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add a CLI typed-link creation command that reuses the existing `/api/links` contract for shell-first graph mutation.
Intent: Close the next terminal gap after evidence upload by letting shell-first operators create validated operational links without opening the browser.
Constraints: Stay on the current local HTTP runtime, reuse the existing typed-link endpoint and payload shape, keep the command narrow and explicit, and make CLI behavior honest when the server rejects unsupported or missing endpoints.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-vuteg-ex5-cli-typed-link-creation.md`, `ex5-operational-knowledge-system/cmd/oks-cli/main.go`, `ex5-operational-knowledge-system/cmd/oks-cli/main_test.go`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`

## Goal

Add a shell-first CLI typed-link creation command for validated operational
graph edges.

## Tasks

- [x] vuteg.1 Define the CLI typed-link command and argument shape.
- [x] vuteg.2 Add a CLI write path for the existing typed-link HTTP contract.
- [x] vuteg.3 Surface both successful writes and server-side endpoint rejection honestly.
- [x] vuteg.4 Add CLI regression coverage for valid and invalid typed-link flows.
- [x] vuteg.5 Update the ex5 docs to describe the terminal-side typed-link behavior honestly.
