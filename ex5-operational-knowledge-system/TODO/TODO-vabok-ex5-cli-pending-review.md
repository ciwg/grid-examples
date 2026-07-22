# TODO vabok - ex5 CLI pending-review

## Decision Intent Log

ID: DI-vabok
Date: 2026-07-21 14:35:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add a CLI pending-review command that reuses the same `/api/search` draft-item, all-run, and problem-run projections as `:OksPending`.
Intent: Let shell-first reviewers see the same pending work queue from the terminal without inventing a separate review endpoint or forcing them into Neovim for simple triage.
Constraints: Stay on the current local HTTP runtime, reuse only the existing search projections, keep the command read-only, and make the CLI output honest about mirroring the current pending-review route mix.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-vabok-ex5-cli-pending-review.md`, `ex5-operational-knowledge-system/cmd/oks-cli/main.go`, `ex5-operational-knowledge-system/cmd/oks-cli/main_test.go`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`, `ex5-operational-knowledge-system/docs/architecture.md`

ID: DI-vabop
Date: 2026-07-21 15:10:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Harden CLI pending-review so malformed `/api/search` payloads fail loudly instead of silently collapsing into incomplete queues.
Intent: Keep terminal review trustworthy by surfacing projection drift immediately, and keep CLI pending-review semantics aligned with Neovim's actual unreviewed-run filtering.
Constraints: Stay on the current shared search routes, reject malformed run and approval shapes with explicit errors, and add regression coverage for malformed payloads and route failures.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-vabok-ex5-cli-pending-review.md`, `ex5-operational-knowledge-system/cmd/oks-cli/main.go`, `ex5-operational-knowledge-system/cmd/oks-cli/main_test.go`

## Goal

Add a terminal-friendly CLI pending-review queue over the same projected search
surfaces the browser and Neovim already use.

## Tasks

- [x] vabok.1 Define the CLI pending-review command and keep it read-only.
- [x] vabok.2 Reuse the existing draft-item, all-run, and problem-run search projections without adding a new endpoint.
- [x] vabok.3 Add CLI regression coverage for the three-route aggregation behavior.
- [x] vabok.4 Update the ex5 docs to describe terminal pending-review behavior honestly.
- [x] vabok.5 Reject malformed shared-search payloads instead of silently hiding projection drift.
- [x] vabok.6 Add regression coverage for malformed payloads and failing pending-review route reads.
