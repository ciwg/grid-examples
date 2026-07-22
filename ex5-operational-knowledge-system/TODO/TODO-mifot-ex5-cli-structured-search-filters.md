# TODO mifot - ex5 CLI structured and problem search filters

## Decision Intent Log

ID: DI-mifot
Date: 2026-07-21 12:18:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Extend the existing CLI `search` command with optional `KEY=VALUE` structured filters that map directly onto `/api/search`.
Intent: Close the next terminal review gap by letting shell-first operators use the same structured and `problem=true` search surface that already powers browser drilldowns and Neovim review views.
Constraints: Reuse the existing `/api/search` contract, keep the CLI filter syntax explicit and narrow, reject malformed or unsupported filter keys locally, and avoid creating a second CLI-only search endpoint.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-mifot-ex5-cli-structured-search-filters.md`, `ex5-operational-knowledge-system/cmd/oks-cli/main.go`, `ex5-operational-knowledge-system/cmd/oks-cli/main_test.go`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`

## Goal

Add terminal-side structured and problem-only search filters without inventing
a separate CLI search backend.

## Tasks

- [x] mifot.1 Define the CLI search filter argument shape.
- [x] mifot.2 Route CLI structured search through the existing `/api/search` query params.
- [x] mifot.3 Reject malformed or unsupported filter arguments locally.
- [x] mifot.4 Add CLI regression coverage for encoded queries, structured filters, and invalid filter keys.
- [x] mifot.5 Update the ex5 docs to describe the richer terminal search behavior honestly.
