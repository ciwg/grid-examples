# TODO sakis — Corrective-action review

## Decision Intent Log

ID: DI-hiboj
Date: 2026-08-08 11:27:26
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add a `correctiveaction` package with append-only opening events and a rebuildable action view, plus a `corrective-action-review` workflow that opens one action linked to a rejected quarantine case.
Intent: Preserve accountable, reusable follow-up for rejected material without treating a quarantine rejection as corrective-action completion or coupling non-resource corrections to maintenance.
Constraints: Require `quarantine_case_id`, `action_id`, `actor`, `evidence_id`, `summary`, and `notes`; use `correctiveaction open <action-id> <quarantine-case-id> <actor> <evidence-id> <summary> [notes...]`; no closure behavior; use `packages/correctiveaction/`, `workflows/corrective-action-review/`, canonical schemas, existing adapter registry, kernel/CLI tests, and automatically cleaned Go test directories.
Affects: Corrective-action package, review workflow artifact, adapter registry, schemas, tests, and Ex6 documentation.

## Goal

Add the first explicit corrective-action review workflow for a rejected
quarantine case.

## Status

- [ ] Implement and verify the locked opening-review slice.
