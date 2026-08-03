# EX6 operator deployment checklist

Use this checklist to run EX6 safely without assuming every workflow or image
is automatically trusted or executable. Workflow artifact receipt, local
activation, registry permission, image acquisition, and execution are separate
operator decisions. Source: `DI-lovek`; `DI-novuk`; `DI-harib`; `DI-zivut`.

## 1. Local workflow use

- [ ] Start from the EX6 root in a clean working directory.
- [ ] Begin with the read-only team briefing. It reports readiness, inbox
  attention, current runs, and one safe `NEXT:` command. Source: `DI-sotad`.

  ```bash
  go run ./cmd/moks workflow overview
  ```
- [ ] Capture or import a workflow artifact.

  ```bash
  go run ./cmd/moks workflow capture workflows/procedure-execution procedure-execution
  ```

- [ ] Inspect readiness. This does not change lifecycle state or pull images.

  ```bash
  go run ./cmd/moks workflow verify procedure-execution
  ```

- [ ] Activate only the workflow artifact you intend to make locally runnable.

  ```bash
  go run ./cmd/moks workflow activate procedure-execution
  ```

- [ ] Start work with explicit input values. Loading or activation alone never
  starts a workflow.

## 2. Optional two-node artifact relay

- [ ] Run a relay endpoint on the receiving node:
  `go run ./cmd/moks relay serve <addr>`.
- [ ] Exchange peer cards and add explicit local peer permissions on both
  nodes.
- [ ] Send the exact workflow artifact and sender lifecycle evidence:

  ```bash
  go run ./cmd/moks workflow relay push <alias> <peer-id>
  ```

- [ ] Receipt retains exact artifact and lifecycle-evidence bytes, plus local
  authenticated sender identities in `state/workflow-receipts.json`; it does not
  add a local workflow listing, activate, or execute anything. Source:
  `DI-jifuk`; `DI-rufir`.
- [ ] On the receiver, scan the CAS-derived inbox and inspect the artifact's
  evidence before making a local lifecycle decision:

  ```bash
  go run ./cmd/moks workflow inbox list
  go run ./cmd/moks workflow inbox inspect <artifact-cid>
  ```

- [ ] Import only an inbox entry whose JSON says `ready_to_import: true`. This
  creates a local alias; inspection and receipt never do:

  ```bash
  go run ./cmd/moks workflow inbox import <artifact-cid> <alias>
  ```

- [ ] Then verify and activate that local alias as separate actions.

- [ ] Use `moks workflow overview` for a human-readable current status view.
  Its Recent activity section is newest-first by durable local run-event time;
  retained v1 heads explicitly show when no event time is available. Source:
  `DI-gihor`.

## 3. Optional portable Docker adapter image

This is needed only when a second node must execute an adapter image that it
does not already have locally.

- [ ] Use an operator-approved registry host and an immutable OCI digest.
- [ ] Allow that exact registry host locally:

  ```bash
  go run ./cmd/moks registry allow registry.example.com
  ```

- [ ] Install a package whose adapter declaration is pinned to that registry
  digest and exactly matches the workflow adapter name and pCIDs.
- [ ] Use `moks workflow verify <alias>` to confirm `registry_allowed` and
  `image_available`.
- [ ] Pull only through the explicit command:

  ```bash
  go run ./cmd/moks workflow image pull <alias>
  ```

- [ ] Verify again, then activate and run only if the artifact is locally
  approved.

## 4. Things EX6 never does automatically

- It does not activate a received workflow.
- It does not execute a workflow merely because it is imported or active.
- It does not allow a package to select an arbitrary registry host.
- It does not pull an image during verification.
- It does not give Docker workers runtime files, secrets, network access, or
  direct durable-write authority.

## 5. Before an operational run

- [ ] Confirm the workflow alias and artifact CID.
- [ ] Confirm `workflow verify` reports eligible execution.
- [ ] Confirm any required procedure, resource, or local context records are
  present.
- [ ] Retain the resulting run, evidence, and approvals through the normal
  runtime commands.
