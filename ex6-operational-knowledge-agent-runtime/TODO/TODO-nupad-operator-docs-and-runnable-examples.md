# TODO nupad - operator docs and runnable examples

## Decision Intent Log

ID: DI-yavuk
Date: 2026-07-30 22:00:00
Status: active
Decision: Add `training-qualification` as a distinct loadable workflow artifact using existing training and runs packages.
Intent: Preserve the ex5 distinction between performing a training session and explicitly recording a qualification decision.
Constraints: Keep the artifact non-executing; add no new top-level action, runtime behavior, or package command.
Affects: `workflows/training-qualification/`, loader demonstration coverage, operator documentation, and `docs/thought-experiments/TE-yavuk-training-qualification-workflow.md`.

ID: DI-zovuk
Date: 2026-07-30 21:45:00
Status: active
Decision: Add `receiving-check` as a distinct loadable workflow artifact using existing context, receiving, and runs packages.
Intent: Preserve the ex5 receiving-check boundary between inbound inspection/disposition and later inventory counting/reconciliation.
Constraints: Keep the artifact non-executing; add no new top-level action, runtime behavior, or package command.
Affects: `workflows/receiving-check/`, loader demonstration coverage, operator documentation, and `docs/thought-experiments/TE-zovuk-receiving-check-workflow.md`.

ID: DI-favuk
Date: 2026-07-30 21:25:00
Status: active
Decision: Add `maintenance-round` as the third loadable workflow example, composing existing context, maintenance, and runs packages.
Intent: Demonstrate a resource-inspection and finding workflow that is distinct from procedure execution and inventory receipt.
Constraints: Keep the artifact non-executing; add no new top-level action, runtime behavior, or package command.
Affects: `workflows/maintenance-round/`, loader demonstration coverage, operator documentation, and `docs/thought-experiments/TE-favuk-maintenance-round-workflow.md`.

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
