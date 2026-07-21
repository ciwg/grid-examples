# TODO farun - ex5 search evidence and approval history

## Decision Intent Log

ID: DI-farun
Date: 2026-07-21 13:25:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track evidence/approval search coverage as its own ex5 review bug TODO.
Intent: Align the implemented search behavior with the ex5 promise that later operators can find evidence and follow-up history again.
Constraints: Keep the current search surface and mixed-result model; expand the indexed search text and regression coverage without inventing a separate search subsystem.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-farun-ex5-search-evidence-and-approval-history.md`, `ex5-operational-knowledge-system/service/app.go`, `ex5-operational-knowledge-system/service/app_test.go`, `ex5-operational-knowledge-system/service/server_test.go`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/README.md`

## Goal

Make ex5 free-text search actually reach the evidence facts and approval notes
that operators need when they come back later to reconstruct what happened.

## Review finding

`SearchWithOptions()` builds run search text from outcome, notes, machine,
location, place, and resource names, but it does not include evidence
summaries/facts or approval notes. That undercuts the documented ex5 story that
later operators can search operational history to find what evidence was
captured and what follow-up happened.

## Tasks

- [ ] farun.1 Expand run search indexing so evidence summaries/facts and approval notes participate in free-text search.
- [ ] farun.2 Add service/API regression coverage and update docs to describe the searchable history surface accurately.

## Status

- open
- derived from the 2026-07-21 deep ex5 review
