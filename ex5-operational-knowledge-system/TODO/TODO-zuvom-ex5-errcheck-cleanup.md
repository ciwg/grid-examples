# TODO zuvom - ex5 errcheck cleanup

## Decision Intent Log

ID: DI-zuvom
Date: 2026-07-21 15:05:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track the current ex5 `errcheck` failures as one focused cleanup slice after the review backlog is complete.
Intent: Restore the ex5 quality gate so `errcheck ./...` passes again without mixing this debt into unrelated feature work.
Constraints: Stay inside `ex5`; fix the currently reported ignored-error sites in the CLI, persistence layer, service helpers, and multipart cleanup path; add tests only where behavior meaningfully changes.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-zuvom-ex5-errcheck-cleanup.md`, `ex5-operational-knowledge-system/cmd/oks-cli/main.go`, `ex5-operational-knowledge-system/service/app.go`, `ex5-operational-knowledge-system/service/persistence.go`, `ex5-operational-knowledge-system/service/server.go`

## Goal

Make `errcheck ./...` pass for `ex5` by handling the currently ignored errors
explicitly instead of relying on unchecked defers or discarded parse results.

## Tasks

- [x] zuvom.1 Handle response body closes explicitly in the ex5 CLI helpers.
- [x] zuvom.2 Handle file close/seek cleanup explicitly in the ex5 persistence layer.
- [x] zuvom.3 Handle multipart cleanup and ID parsing results explicitly in the ex5 service layer.

## Status

- done
- derived from the 2026-07-21 ex5 verification follow-up
