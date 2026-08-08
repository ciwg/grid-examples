# TODO ronos — Quarantine package and receiving-exception workflow

## Decision Intent Log

ID: DI-luval
Date: 2026-08-08 11:06:28
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Model quarantine as a reusable Ex6 domain package, with a receiving-exception workflow that opens a quarantine case but does not resolve it.
Intent: Preserve a durable, reusable quarantine boundary for receiving, inventory, and later corrective-action work while keeping the generic runtime free of receiving-domain policy.
Constraints: Use append-only quarantine events with a rebuildable local case view; the first terminal transitions are release and reject only; do not include corrective-action transfer in this slice; the receiving-exception workflow only opens cases; avoid the retired analogy labels in newly authored documentation.
Affects: Future quarantine package source, receiving-exception workflow artifact, Ex6 workflow registry, tests, documentation, and runtime-local evidence projections.

ID: DI-giriz
Date: 2026-08-08 11:08:47
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Use one `moks.quarantine.event.v1` append-only event family, expose `quarantine open|release|reject|list|inspect`, and add the `receiving-exception` workflow under `packages/quarantine/` and `workflows/receiving-exception/`.
Intent: Give every quarantine transition one inspectable pCID-defined evidence shape, retain a compact public interface consistent with existing Ex6 packages, and make receiving exception handling independently reusable.
Constraints: The event family represents only open, release, and reject in this slice; its local case projection is rebuildable; the workflow opens a case only; code and documentation use `quarantine` and `receiving-exception` names; avoid the retired analogy labels in newly authored documentation.
Affects: `packages/quarantine/`, `workflows/receiving-exception/`, built-in workflow dispatch, canonical adapter schemas, kernel workflow tests, CLI workflow tests, and Ex6 documentation.

ID: DI-hogid
Date: 2026-08-08 11:09:23
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Use caller-supplied distinct quarantine case IDs; store append-only typed events containing case ID, transition, actor, evidence reference, and optional notes; require receiving ID and failed receipt-run reference on open; and expose the typed receiving-exception workflow contract under the approved standard source and test paths.
Intent: Make the evidence boundary and reusable case identity explicit, permit multiple independent cases for one receiving item, and make a failed receiving inspection mechanically distinguishable from a completed quarantine hold.
Constraints: The only transitions are open, release, and reject; open is the receiving-exception workflow's only effect; the workflow input fields are `receiving_id`, `receipt_run_id`, `case_id`, `actor`, `evidence_id`, `exception`, and `notes`; output returns case and opening-event IDs; use `openQuarantine`, `releaseQuarantine`, `rejectQuarantine`, `listQuarantines`, and `inspectQuarantine`; use `packages/quarantine/package.go`, package-adjacent tests, `workflows/receiving-exception/`, extended existing kernel/CLI workflow tests, and automatically cleaned Go test temporary directories.
Affects: `packages/quarantine/package.go`, `packages/quarantine/package_test.go`, `workflows/receiving-exception/`, `builtin/workflow_operations.go`, canonical workflow-adapter schemas, `kernel/workflow_schema_test.go`, `kernel/workflow_runs_test.go`, and `cmd/moks/main_test.go`.

## Goal

Add a reusable quarantine package and a receiving-exception workflow that opens
a durable quarantine case from an inbound receiving exception.

## Scope

- Define the quarantine package's append-only record and rebuildable local view.
- Record explicit open, release, and reject transitions.
- Add a receiving-exception workflow that records the failed inspection and
  opens, but never resolves, a quarantine case.
- Preserve the runtime as a generic persistence, routing, and workflow host.
- Add deterministic coverage for normal, incomplete, concurrent, and retained
  workflow-artifact scenarios identified in TE-maluz.

## Out of scope

- Corrective-action ownership transfer.
- Runtime-enforced domain quarantine policy.
- Retroactively rewriting the completed receiving-package scope in TODO-zibek.

## Source analysis

- `docs/thought-experiments/TE-maluz-receiving-quarantine-package.md`
- `ex6-operational-knowledge-agent-runtime/TODO/TODO-zibek-ex6-receiving-package.md`

## Status

- [x] Complete decision framing for record names, source paths, schemas,
  adapter contract, test paths, and runtime-path patterns.
- [x] Implement the locked package and workflow slice.
