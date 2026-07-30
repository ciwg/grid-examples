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
        workflow.yaml
        procedure.md
        protocols/
          run-request.md
          run-evidence.md
          approval.md
        policies/
          evidence-rules.yaml
          approval-rules.yaml
        examples/
          successful-run.json
          rejected-run.json
        tests/
          workflow_test.go

- **README.md** gives an operator-facing summary.
- **workflow.yaml** identifies the workflow version and declares inputs,
  outputs, and package/protocol dependencies.
- **procedure.md** is the human-readable operating procedure.
- **protocols/** contains the pCID-defined message contracts the workflow uses.
- **policies/** states local evidence and approval requirements.
- **examples/** supplies realistic worked cases.
- **tests/** verifies the workflow's declared behavior.

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

Receiving inspection, maintenance closure, training qualification, and
knowledge review can follow the same artifact shape while applying different
domain rules.

## Current status

This directory layout and deterministic archive format are proposed. The
runtime currently supports CAS-backed lifecycle events for imported artifacts,
but does not yet capture a local workflow directory into CAS or dispatch a
workflow through workers.
