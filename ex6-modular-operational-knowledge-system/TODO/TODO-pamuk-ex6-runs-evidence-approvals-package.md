# TODO pamuk - ex6 runs, evidence, and approvals package

## Decision Intent Log

ID: DI-pamuk
Date: 2026-07-28 14:00:00
Status: active
Decision: Implement the first runs package as a first-party built-in egg with run, evidence, and approval families while the installed-package mutation contract remains shallow.
Intent: Get the durable performed-work surface into ex6 code now so the basket has real operational records beyond context and document state.
Constraints: Keep the package under `packages/runs/`; preserve durable history writes; use CAS for evidence body payloads when text evidence bodies are supplied.
Affects: `packages/runs`, `cmd/moks`, runtime tests, and ex6 current-state docs.

## Goal

Build the first-party `runs` package.

## Scope

- run recording
- evidence attachment model
- evidence facts
- item and run approvals
- durable family declarations

## Why

This package carries the operational record of work actually performed.
