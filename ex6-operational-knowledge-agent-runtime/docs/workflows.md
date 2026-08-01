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

Import does not grant route or worker-execution authority. Activation and the
manifest-selected built-in adapter make that authority explicit.

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

## Shared-runtime scenario coverage

The CLI scenario test loads and activates all seven artifacts in one runtime,
then drives receiving, inventory, maintenance, training, and knowledge commands
through the main `moks` command dispatcher. This proves shared operational
interaction without claiming that workflow artifacts autonomously execute.

Source: DI-sovuk.

## Basket lifecycle and executable orchestration

The runtime is the basket; workflow artifacts are removable eggs. Capturing or
importing an artifact only puts exact bytes in CAS. Activating it makes that
specific artifact eligible to run. Deactivating or revoking it blocks both new
runs and handoffs, even when its built-in package capability remains installed.
Source: DI-lumek.

1. **Capture/import.** Retain the deterministic tar artifact in CAS. It is not
   executable yet.
2. **Verify.** Check its manifest and required local package/protocol
   dependencies.
3. **Activate.** The artifact must be active before the runtime accepts a run.
4. **Start.** `moks workflow run start <alias> <key> <value> ...` stores a
   canonical pCID-selected CBOR input envelope and CAS lifecycle event. Each
   active artifact validates its own required fields before calling its owning
   built-in package command.
5. **Execute.** The manifest-selected trusted built-in adapter validates the
   input and emits a canonical pCID-selected CBOR output. Docker execution is a
   later backend under TE-dovek.
6. **Inspect.** `moks workflow run status <run-id>` shows durable state,
   input/output CIDs, and any failure reason.
7. **Handoff.** `moks workflow run handoff <run-id> <target-alias>` passes the
   exact completed CBOR output to an active target. A matching target input
   pCID executes immediately; a different target schema records a durable
   `waiting-for-input` run with the source envelope retained as evidence.
8. **Policy handoff.** `moks workflow policy handoff set <source>
   <output-pcid> <target> <input-pcid>` records one explicit local next step.
   A completed matching source run starts that active target automatically;
   differing schemas wait for explicit target-schema input instead of silently
   coercing data. An unavailable target records a durable failed source state
   instead of retrying.
9. **Wait, supply, retry.** An adapter that needs a required field records
   `waiting-for-input`; `moks workflow run supply <run-id> <key> <value> ...`
   resumes it. A failed run may be retried with `moks workflow run retry
   <run-id>` using the exact retained input. If a physical persistence failure
   leaves a run in `running`, the same explicit retry command is the manual
   recovery path. There is no automatic retry.
10. **Deactivate/revoke.** The egg stays retained for audit/extraction, but the
   runtime refuses to start it again.

The current seven artifacts each declare an adapter plus distinct input and
output pCIDs. The common outer CBOR envelope carries sorted string fields, but
the pCID identifies the individual adapter contract; the current adapters call
these real package commands:

- `procedure-execution`: `procedures record-use`
- `inventory-receipt`: `inventory record-count`
- `inventory-discrepancy-review`: `inventory record-reconcile`
- `receiving-check`: `receiving record-receipt`
- `maintenance-round`: `maintenance record-service`
- `training-qualification`: `training record-session`
- `knowledge-review`: `knowledge item approve`

The frozen canonical JSON specifications for those contracts are published in
[`docs/protocols/workflow-adapter-schemas/`](./protocols/workflow-adapter-schemas/).
Each declared pCID is the CIDv1 raw-SHA-256 identifier of its exact schema
bytes; the schema test recomputes every shipped mapping. Each workflow artifact
also carries its input and output schema under `schemas/`, and verification
retains those exact bytes in local CAS. Source: DI-lumek.

### Retained v1 artifacts

Artifacts captured before the canonical schema publication can remain active.
Their old input/output pCID pairs are supported only by the corresponding
shipped trusted adapter, which emits the old declared output pCID when it
receives the old declared input pCID. The runtime does not translate an old
envelope into a new one, and new capture rejects an old pCID.
This preserves the exact contract of a retained egg while newer artifacts use
the canonical embedded schemas. Source: DI-lumek.

`moks workflow verify <alias-or-cid>` prints the manifest plus `contract`,
`adapter_available`, `schema_cas_ready`, and `eligible_to_execute`. A retained
v1 artifact can be eligible through its supported adapter even though it has no
embedded canonical schema bytes; canonical artifacts report schema/CAS readiness
only after their embedded schemas have verified and entered CAS, and cannot be
eligible without that readiness. A structurally valid artifact with an
unavailable package or protocol still prints this report with
`eligible_to_execute: false` and a dependency `reason`. Source: DI-lumek.

Missing adapter fields and schema-changing handoffs produce
`waiting-for-input`; malformed or rejected commands produce durable `failed`
state. The run cache is disposable: runtime open rebuilds it from a full CAS
scan. Source: DI-lumek.

Receiving inspection, maintenance closure, training qualification, and
knowledge review can follow the same artifact shape while applying different
domain rules.

## Current status

The runtime captures a valid workflow directory as a deterministic tar archive
in CAS, imports it under a local alias, requires explicit activation, and
dispatches through manifest-selected trusted built-in adapters. Worker backends
remain future work.

The complete generic loader contract is documented in
[Workflow Loader: The Basket](./workflow-loader.md).
