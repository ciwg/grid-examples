# Workflow Artifacts

Workflows are versioned operating playbooks that compose package capabilities.
They are not a new top-level PromiseGrid action kind. A workflow artifact tells
people and participating packages how to perform work; lifecycle events record
the local runtime's decision to retain or make that artifact available.

Source: DI-lovek; DI-bavuk.

## Proposed layout

Each workflow lives in its own self-contained directory:

    workflows/
      procedure-execution/
        README.md
        workflow.json
        procedure.md
        policies/
          evidence-rules.json
          approval-rules.json
        examples/
          successful-run.json
      inventory-receipt/
        README.md
        workflow.json
        receipt.md
        policies/
          evidence-rules.json
      maintenance-round/
        README.md
        workflow.json
        round.md
        policies/
          evidence-rules.json
      receiving-check/
        README.md
        workflow.json
        check.md
        policies/
          evidence-rules.json
      training-qualification/
        README.md
        workflow.json
        qualification.md
        policies/
          evidence-rules.json
      inventory-discrepancy-review/
        README.md
        workflow.json
        review.md
        policies/
          evidence-rules.json
      knowledge-review/
        README.md
        workflow.json
        review.md
        policies/
          review-rules.json

- **README.md** gives an operator-facing summary.
- **workflow.json** identifies the workflow version and declares its summary
  and package/protocol dependencies.
- **procedure.md** is the human-readable operating procedure.
- **policies/** states local evidence and approval requirements.
- **examples/** supplies realistic worked cases.

## Lifecycle and storage

The runtime will deterministically archive a workflow directory into CAS. The
resulting artifact CID identifies the exact workflow version. The runtime's
existing lifecycle protocol records whether that CID is locally imported,
active, deactivated, or revoked. It retains the artifact and event history
even when local availability is withdrawn.

The lifecycle registry does not execute the workflow. Import also does not
grant route or worker-execution authority. Future route/worker registration
must make that authority explicit.

## Composition

Packages provide reusable capabilities; workflows compose them for an
operational outcome.

For example, a procedure-execution workflow can combine procedures for the
definition, runs for performed-work records/evidence/approvals, context for
people/places/resources, knowledge for approved reference material, and links
for durable relationships between the resulting records.

The repository also includes `inventory-receipt`, which composes context,
receiving, inventory, and runs to document receipt, disposition, counting, and
reconciliation, `maintenance-round`, which composes context, maintenance, and
runs for resource inspection, service, and findings, and `receiving-check`,
which keeps inbound inspection and disposition separate from inventory, and
`training-qualification`, which keeps training sessions separate from explicit
certification decisions. Together, the seven artifacts demonstrate that the
loader is not specific to procedure execution. `inventory-discrepancy-review`
adds explicit count reconciliation for adjust, investigate, or reject decisions,
and `knowledge-review` retains revision, approval, and supersedence review.

Source: DI-voruk; DI-favuk; DI-zovuk; DI-yavuk; DI-pavuk; DI-dovuk.

Receiving inspection, maintenance closure, training qualification, and
knowledge review can follow the same artifact shape while applying different
domain rules.

## Current status

The runtime captures a valid workflow directory as a deterministic tar archive
in CAS, imports it under a local alias, and requires an explicit activation
step. It does not yet dispatch a workflow through workers.

The complete generic loader contract is documented in
[Workflow Loader: The Basket](./workflow-loader.md).
