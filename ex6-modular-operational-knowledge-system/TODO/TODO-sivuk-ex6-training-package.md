# TODO sivuk - ex6 training package

## Decision Intent Log

ID: DI-sivuk
Date: 2026-07-28 17:00:00
Status: active
Decision: Implement the first training package as a first-party built-in egg with explicit training item, session, and completion families over shared knowledge and runs records.
Intent: Add the next real domain egg without collapsing training behavior back into the basket or hiding it inside generic run payloads alone.
Constraints: Keep the package under `packages/training/`; preserve basket-mediated durable history; compose with `knowledge` and `runs` instead of duplicating their cores.
Affects: `packages/training`, `cmd/moks`, runtime tests, and ex6 status docs.

## Goal

Build the first-party `training` package.

## Scope

- training item kind
- training session recording
- training completion recording
- training inspection and listing
- durable family declarations

## Why

Training is one of the next real ex5-derived domain surfaces and should exist
as its own egg instead of waiting behind later embodiments.
