# TODO nuvaz - ex5 CLI problem review

## Decision Intent Log

ID: DI-nuvaz
Date: 2026-07-21 12:32:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add a CLI `problem-review` command that reuses the existing `/api/problem-review` grouped hotspot surface.
Intent: Close the next terminal review gap by letting shell-first operators see repeated receiving and count problems by place and resource without opening the browser.
Constraints: Reuse the existing grouped problem-review endpoint and response shape, keep the CLI command read-only, and avoid inventing a second terminal-only review API.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-nuvaz-ex5-cli-problem-review.md`, `ex5-operational-knowledge-system/cmd/oks-cli/main.go`, `ex5-operational-knowledge-system/cmd/oks-cli/main_test.go`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`

## Goal

Add a shell-first grouped problem review command for hotspot triage.

## Tasks

- [x] nuvaz.1 Define the CLI problem-review command.
- [x] nuvaz.2 Reuse the existing grouped hotspot endpoint for terminal review.
- [x] nuvaz.3 Add CLI regression coverage for the route.
- [x] nuvaz.4 Update the ex5 docs to describe the new terminal review behavior honestly.
