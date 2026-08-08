# TODO nosij — Quarantine resolution workflow

## Decision Intent Log

ID: DI-nufav
Date: 2026-08-08 11:16:28
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add one `quarantine-resolution` workflow with typed `release` or `reject` decision input; validate it in `quarantineResolutionWorkflowOperation` and dispatch only to the existing explicit package command.
Intent: Retain one pCID-defined workflow contract while preserving explicit, append-only terminal transitions and avoiding a new generic package command.
Constraints: Input is `case_id`, `event_id`, `actor`, `evidence_id`, `decision`, and `notes`; output preserves those fields plus `stage`; case ID identifies its opening event; one evidence reference is required; paths are `workflows/quarantine-resolution/`, canonical schemas, existing adapter registry, and existing kernel/CLI tests.
Affects: `builtin/workflow_operations.go`, `workflows/quarantine-resolution/`, canonical schemas, kernel workflow tests, CLI workflow tests, and package tests.

## Goal

Add the separately authorized terminal-resolution workflow for open quarantine
cases.

## Status

- [x] Implement and verify the locked workflow slice.
