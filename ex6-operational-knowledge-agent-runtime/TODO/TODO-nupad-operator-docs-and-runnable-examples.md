# TODO nupad - operator docs and runnable examples

## Decision Intent Log

ID: DI-voruk
Date: 2026-07-30 21:05:00
Status: active
Decision: Add `inventory-receipt` as the second loadable workflow example, composing existing context, receiving, inventory, and runs packages.
Intent: Demonstrate that workflow artifacts are reusable package compositions rather than a procedure-execution-specific feature.
Constraints: Keep the artifact non-executing; add no new top-level action, runtime behavior, or package command.
Affects: `workflows/inventory-receipt/`, loader demonstration coverage, operator documentation, and `docs/thought-experiments/TE-voruk-inventory-receipt-workflow.md`.

ID: DI-nupad
Date: 2026-07-28 16:10:00
Status: active
Decision: Add runnable operator-facing examples and make the docs point at the actual built-in and installed-package flows that work today.
Intent: Keep ex6 from presenting itself as a finished product without giving readers an obvious path to prove the current package stack and installed-package contract.
Constraints: Keep all example artifacts inside `ex6-operational-knowledge-agent-runtime/`; use the real `moks` CLI paths; avoid inventing a fake finished UX.
Affects: `README.md`, `docs/current-state.md`, `docs/package-author-guide.md`, `docs/runnable-examples.md`, `examples/`, and CLI coverage.

## Goal

Make it easy to run the current ex6 system without reverse-engineering the
tests or package source.
