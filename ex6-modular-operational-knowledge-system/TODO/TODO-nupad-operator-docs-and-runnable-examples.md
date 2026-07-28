# TODO nupad - operator docs and runnable examples

## Decision Intent Log

ID: DI-nupad
Date: 2026-07-28 16:10:00
Status: active
Decision: Add runnable operator-facing examples and make the docs point at the actual built-in and installed-package flows that work today.
Intent: Keep ex6 from presenting itself as a finished product without giving readers an obvious path to prove the current package stack and installed-package contract.
Constraints: Keep all example artifacts inside `ex6-modular-operational-knowledge-system/`; use the real `moks` CLI paths; avoid inventing a fake finished UX.
Affects: `README.md`, `docs/current-state.md`, `docs/package-author-guide.md`, `docs/runnable-examples.md`, `examples/`, and CLI coverage.

## Goal

Make it easy to run the current ex6 system without reverse-engineering the
tests or package source.
