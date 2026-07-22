# TODO zanub - ex5 CLI evidence upload

## Decision Intent Log

ID: DI-zanub
Date: 2026-07-21 11:20:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add a CLI evidence upload command that reuses the existing run evidence multipart API for summary-only, fact-bearing, and attachment-bearing writes.
Intent: Close one of the largest remaining terminal gaps by letting shell-first operators attach real evidence to runs without opening the browser.
Constraints: Stay on the current local HTTP runtime, reuse the existing `/api/runs/{id}/evidence` contract, keep the CLI command narrow and explicit, and support the same three practical cases the browser already supports: summary only, summary plus facts, and summary plus facts plus attachment.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-zanub-ex5-cli-evidence-upload.md`, `ex5-operational-knowledge-system/cmd/oks-cli/main.go`, `ex5-operational-knowledge-system/cmd/oks-cli/main_test.go`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`

## Goal

Add a shell-first CLI evidence upload command for runs, including optional
facts JSON and optional file attachment.

## Tasks

- [x] zanub.1 Define the CLI evidence upload command and argument shape.
- [x] zanub.2 Add a CLI multipart upload path for run evidence.
- [x] zanub.3 Support summary-only, summary-plus-facts, and summary-plus-facts-plus-attachment cases.
- [x] zanub.4 Add CLI regression coverage for the multipart request and optional file handling.
- [x] zanub.5 Update the ex5 docs to describe the new terminal-side evidence behavior honestly.
