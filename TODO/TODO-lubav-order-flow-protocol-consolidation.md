# TODO lubav - Order flow protocol consolidation

## Decision Intent Log

ID: DI-rafud
Date: 2026-07-10 14:51:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Revise the order-flow design to remove `_v1` suffixes and remove payload protocol-name fields.
Intent: Make protocol identity derive from the selected pCID and the payload shape rather than a duplicated wire-field name.
Constraints: Rename kernel registration references too; keep human-readable protocol names in docs, code symbols, and CLI output only; payloads must not contain protocol names.
Affects: ex1-order-flow/docs/design.md

ID: DI-movab
Date: 2026-07-10 14:51:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Consolidate app-level pCIDs to one per business domain and remove the separate post-order recording flow.
Intent: Minimize pCID count while preserving one inspectable protocol family per materially distinct business domain.
Constraints: Keep distinct domain pCIDs only for order, pick_pack, shipment, and accounting; do not allocate separate pCIDs for request/result direction; do not add a separate post-order recording flow.
Affects: ex1-order-flow/docs/design.md

ID: DI-lisut
Date: 2026-07-10 14:51:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: MVP requires actual signatures and cryptographic capability tokens on every message path, including kernel registration.
Intent: Make the MVP demonstrate authentic signed traffic and promise-carrying authority rather than placeholders.
Constraints: Capability tokens are promises; every message path is signed; the doc must define refusal behavior for missing or invalid signatures or capability tokens.
Affects: ex1-order-flow/docs/design.md

ID: DI-pagud
Date: 2026-07-10 14:51:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Recommended Go package paths in this repo must not use `internal/`.
Intent: Keep the design aligned with repo package-layout policy.
Constraints: Use only top-level or other approved purpose-named paths in examples.
Affects: ex1-order-flow/docs/design.md

ID: DI-kozod
Date: 2026-07-10 15:26:02 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: The implementation lives in a nested Go module rooted at `ex1-order-flow/`, with module path `github.com/computerscienceiscool/grid-examples/ex1-order-flow`, direct shared packages under `ex1-order-flow/`, role entrypoints under `ex1-order-flow/cmd/`, and Docker assets under `ex1-order-flow/docker/`.
Intent: Keep every artifact for this example under the example directory while still producing normal Go imports and one binary per role.
Constraints: Do not use `internal/`; do not use a `vendor/` directory; keep all example-specific files under `ex1-order-flow/`.
Affects: ex1-order-flow/go.mod, ex1-order-flow/cmd, ex1-order-flow/docker, ex1-order-flow/*

ID: DI-nonad
Date: 2026-07-10 15:26:02 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Use the signed outer envelope shape `grid([42(pCID), payload, proof])`, and make capability tokens pCID-owned payload fields rather than a universal envelope slot.
Intent: Stay aligned with the upstream signed-message profile and keep token placement owned by protocol semantics instead of freezing a universal wire slot too early.
Constraints: Payloads still must not contain protocol names; every message path remains signed; capability tokens remain required where the protocol profile says they are required.
Affects: ex1-order-flow/docs/design.md, ex1-order-flow/protocol, ex1-order-flow/token, ex1-order-flow/*
Supersedes: DI-lisut

ID: DI-rokol
Date: 2026-07-10 15:26:02 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Clean the host run-data tree `/tmp/grid-examples-ex1-data/<role>/...` before each run and preserve the resulting artifacts after the run for manual review.
Intent: Ensure runs start from a deterministic empty state without deleting the completed artifacts operators may want to inspect.
Constraints: The approved temp pattern is `/tmp/grid-examples-ex1-data/<role>/...`; the run wrapper removes it before a run and does not remove it after a run.
Affects: ex1-order-flow/docker, ex1-order-flow/*run*, /tmp/grid-examples-ex1-data/<role>/...

ID: DI-sabol
Date: 2026-07-10 15:26:02 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: The collector runs as its own long-running container service, and the analyzer runs as its own separate post-run container invocation from the built image.
Intent: Match the POC harness shape more closely while keeping the analyzer outside the live routing topology.
Constraints: The analyzer is not a resident routing participant; it runs after the scenario against the preserved run root.
Affects: ex1-order-flow/docker, ex1-order-flow/collector, ex1-order-flow/analyzer

ID: DI-lihit
Date: 2026-07-10 15:26:02 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Reuse upstream `wire-lab` work by copying and adapting selected implementation pieces into this repo, using external CBOR and COSE libraries plus either a credible external CWT library or the adapted upstream local CWT layer when needed.
Intent: Ship a runnable example quickly without coupling this repo directly to sibling-repo imports.
Constraints: Do not import runtime packages directly from `~/lab/wire-lab`; copy and adapt only the pieces needed here; keep external dependencies standards-based for CBOR and COSE support.
Affects: ex1-order-flow/go.mod, ex1-order-flow/token, ex1-order-flow/collector, ex1-order-flow/analyzer, ex1-order-flow/*

ID: DI-vurad
Date: 2026-07-10 16:00:10 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: In the signed outer envelope, slot `1` must carry the payload as a CBOR item directly, not as a CBOR byte string that wraps encoded payload bytes.
Intent: Keep the wire shape aligned with `grid([42(pCID), payload, proof])` and avoid a wrapper that changes the meaning of slot `1`.
Constraints: Preserve exact payload bytes for signing and downstream decoding; keep `proof` in slot `2`; do not change capability-token placement or the business payload structs in this correction.
Affects: ex1-order-flow/protocol, ex1-order-flow/docs/design.md, ex1-order-flow/*tests*

### DI-garis

- ID: DI-garis
- Date: 2026-08-10 08:39:01 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Expand the five existing Ex1 local draft profile specs in place with a common contract template; document the current signed bearer capability tokens as reusable until expiry; and specify local rejection for malformed or unauthorized input, with signed results or finals only for defined valid-profile response paths.
- Intent: Make every selected local pCID auditable without pretending it is a frozen upstream spec, adding a generic wire protocol, or silently changing current replay and failure behavior.
- Constraints: Keep capability-token meaning pCID-owned payload semantics; preserve reusable-until-expiry behavior and the absence of a redemption ledger; do not require signed responses to malformed or unauthorized input; retain the existing five spec paths; mark the profiles as local draft/example contracts; update derived pCID inventory values in the same change; defer durable observer timeout evidence to lubav.4.
- Affects: `ex1-order-flow/specdocs/*.md`, `ex1-order-flow/docs/design.md`, `ex1-order-flow/README.md`, `TODO/TODO-lubav-order-flow-protocol-consolidation.md`, `TODO/handle-namespace.tsv`

### DI-josir

- ID: DI-josir
- Date: 2026-08-10 08:42:07 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Publish an Ex1 local-draft implementation-scope declaration in `ex1-order-flow/CHANGELOG.md` rather than a formal frozen-spec conformance claim.
- Intent: Make the currently implemented local profile pCIDs and component scope auditable while preserving the distinction between an example profile and a frozen upstream PromiseGrid contract.
- Constraints: Name the five current local pCIDs and implemented components; state that no upstream interoperability or frozen-spec conformance is claimed; list host-local and production-grade behavior that remains unclaimed; reserve normal implementation-promise claim fields for a future frozen-spec adoption.
- Affects: `ex1-order-flow/CHANGELOG.md`, `ex1-order-flow/docs/thought-experiments/TE-birut-ex1-local-draft-implementation-publication.md`, `TODO/TODO-lubav-order-flow-protocol-consolidation.md`, `TODO/handle-namespace.tsv`

### DI-vihoz

- ID: DI-vihoz
- Date: 2026-08-10 08:42:07 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Retain Ex1 exceptional-input evidence as append-only local `ObservationRecord` entries in `observations.jsonl`, with raw-byte retention where bytes exist; have both agents and the kernel retain their own observations; and append `refusal_observed` records that point to verified signed refusal artifacts.
- Intent: Preserve what each component actually observed without converting a timeout, malformed input, or local rejection into a claim about another agent's intent or adding a new PromiseGrid wire action.
- Constraints: Keep `messages.jsonl` for valid sent/received envelope records; retain raw bytes before parse/validation when observed; use local evidence only; keep kernel records limited to ingress/dispatch observations; require no signed response to malformed or unauthorized input; timeout records reference the expected request and local deadline; do not add a redemption ledger or replay guarantee.
- Affects: `ex1-order-flow/artifact`, `ex1-order-flow/agent`, `ex1-order-flow/kernel`, `ex1-order-flow/intake`, `ex1-order-flow/docker`, `ex1-order-flow/e2e`, `ex1-order-flow/specdocs/*.md`, `ex1-order-flow/docs/design.md`, `TODO/TODO-lubav-order-flow-protocol-consolidation.md`, `TODO/handle-namespace.tsv`

### DI-riguz

- ID: DI-riguz
- Date: 2026-08-10 08:46:11 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Give the kernel its own `PG_DATA_DIR` rooted at `/tmp/grid-examples-ex1-data/kernel/...`; use `ObservationRecord`, `Store.SaveRawBytes`, and `Store.AppendObservation`; persist raw bytes at ingress; and create observations at the receiving boundary appropriate to parse/pCID/proof, token/refusal/timeout, or kernel ingress/dispatch conditions.
- Intent: Keep durable evidence local to its observer and make every stored exceptional-input record refer to the bytes or request the observer actually assessed.
- Constraints: Preserve the existing agent store layout beside `messages.jsonl`; give Docker and tests a kernel directory; do not centralize observation interpretation in the collector; do not make the kernel assess business promises; do not invent a new wire action; use append-only `observations.jsonl`.
- Affects: `ex1-order-flow/artifact`, `ex1-order-flow/agent`, `ex1-order-flow/kernel`, `ex1-order-flow/intake`, `ex1-order-flow/{seller,warehouse,accounting,carrier}`, `ex1-order-flow/cmd/pg-order-kernel`, `ex1-order-flow/docker`, `ex1-order-flow/e2e`, `TODO/TODO-lubav-order-flow-protocol-consolidation.md`, `TODO/handle-namespace.tsv`

### DI-purum

- ID: DI-purum
- Date: 2026-08-10 08:47:27 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Expose `Client.RecordObservation` to role handlers while retaining `Store.AppendObservation` as the append-only storage operation.
- Intent: Let handlers state their local observation without exposing the store or conflating domain evidence with file-append mechanics.
- Constraints: `Client.RecordObservation` records only local evidence; it must not send a wire message or assert another agent's intent.
- Affects: `ex1-order-flow/agent`, `ex1-order-flow/{seller,warehouse,accounting,carrier,intake}`, `TODO/TODO-lubav-order-flow-protocol-consolidation.md`, `TODO/handle-namespace.tsv`

### DI-zosiz

- ID: DI-zosiz
- Date: 2026-08-10 08:47:27 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Record `no_registered_recipient` when the kernel receives a valid envelope whose selected pCID has no current local recipient, retain the exact bytes, and decline local dispatch without claiming global pCID validity.
- Intent: Keep the kernel a local pCID dispatcher rather than a protocol registry or semantic authority.
- Constraints: No signed reply; raw bytes remain local evidence; an agent's `unexpected_pcid` observation remains separate; do not maintain a kernel pCID allow-list.
- Affects: `ex1-order-flow/kernel/server.go`, `ex1-order-flow/docs/design.md`, `ex1-order-flow/docs/thought-experiments/TE-vibab-ex1-kernel-unroutable-pcid-policy.md`, `TODO/TODO-lubav-order-flow-protocol-consolidation.md`, `TODO/handle-namespace.tsv`

### DI-potoj

- ID: DI-potoj
- Date: 2026-08-10 09:06:44 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Cover Ex1's published local-profile claims with layered regression tests: derive profile pCIDs from the five source specs and compare them with the design inventory and implementation-scope declaration; add focused package tests for rejection, raw retention, and local observations; and retain one E2E refusal path plus one E2E timeout path with evidence assertions.
- Intent: Verify the actual pCID-defined local contracts and each observer's retained evidence without converting a test harness into semantic authority or inferring another agent's intent from a timeout or refusal.
- Constraints: Put focused tests beside their respective packages; extend the existing E2E package only for cross-role behavior; use Go-managed temporary roots; do not add a kernel pCID allow-list, a new top-level wire action, or E2E duplicates for every rejection class.
- Affects: `ex1-order-flow/{artifact,agent,protocol,token,kernel}/*_test.go`, `ex1-order-flow/e2e/e2e_test.go`, `ex1-order-flow/docs/design.md`, `ex1-order-flow/CHANGELOG.md`, `ex1-order-flow/docs/thought-experiments/TE-pokis-ex1-regression-coverage.md`, `TODO/TODO-lubav-order-flow-protocol-consolidation.md`, `TODO/handle-namespace.tsv`

### DI-motiv

- ID: DI-motiv
- Date: 2026-08-10 09:41:27 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Complete Ex1's guide-facing documentation by extending its README with a concise walkthrough for the happy path, warehouse refusal, and carrier timeout fixtures; name the collector DAG, per-role raw artifacts, and local observations with their evidence boundaries; and state that a failed run is not a completed demo claim and must be rerun from a fresh selected root.
- Intent: Make the local-draft implementation scope, repeatable demo, and observer-local evidence legible from the entry document without creating a second guide that could drift from the existing wrapper, profile, and testing documentation.
- Constraints: Keep detailed contracts in `docs/design.md` and `specdocs/`; keep verification detail in `docs/testing.md`; link the implementation-scope declaration; do not imply frozen-spec conformance, peer interoperability, shared evidence, or an inference about another agent's intent.
- Affects: `ex1-order-flow/README.md`, `ex1-order-flow/docs/thought-experiments/TE-tivot-ex1-guide-facing-documentation.md`, `TODO/TODO-lubav-order-flow-protocol-consolidation.md`, `TODO/handle-namespace.tsv`

## PromiseGrid Alignment Plan

Source review: `~/lab/cswg/promisegrid-dev-guide/README.md`, App Devs
guidance, reviewed 2026-08-10. The current Ex1 mechanics already demonstrate
pCID-selected envelopes, app-local promise interpretation, signed traffic,
capability-token fields, pCID-only kernel routing, and raw-CBOR artifact
retention. The remaining work makes its protocol contract and implementation
claims legible to readers of the development guide.

The tasks below are deliberately ordered so that the protocol's published
meaning is made explicit before any behavior, persistence, or test changes.
Every task that changes wire semantics, capability semantics, durable evidence,
or names must first complete its own TE and Decision Framing; this plan does
not lock those later decisions.

- [ ] lubav.1 Publish an Ex1 protocol inventory that lists every local profile,
  its derived pCID, its owning spec document, its handler roles, and its status
  as a local draft/example profile rather than a claimed frozen upstream spec.
  Update the README to link that inventory and to identify the retained raw-CBOR
  evidence and analyzer output.
- [x] lubav.2 Expand each local profile spec into an auditable contract:
  payload fields and valid kinds; signing/proof coverage; capability issuer,
  audience, expiry, redemption/reuse, and refusal behavior; parent/evidence
  links; validation failures; and any durable-versus-transient distinction.
- [x] lubav.3 Add an implementation promise publication record that names the
  exact spec doc-CIDs/pCIDs, claim scope, partial status, and explicitly
  unclaimed host-local behavior. Do not represent local draft profiles as
  frozen upstream PromiseGrid specifications.
- [x] lubav.4 Run a TE and DF for durable observer evidence of timeout,
  malformed input, refusal, and unknown-pCID receipt. Implement only the
  selected local evidence record and retention policy; preserve the rule that a
  timeout is the observer's record, not proof of another agent's intent.
- [x] lubav.5 Document and test the kernel's unknown-pCID policy, registration
  boundary, and exact-byte forwarding guarantee. The kernel must remain a
  pCID dispatcher, not a business workflow authority or generic RPC router.
- [x] lubav.6 Add regression coverage that proves the published claims:
  profile/pCID inventory consistency, signature and capability rejection,
  refusal versus promise assessment, timeout/observer-evidence behavior after
  its decision is locked, and replay/audit of retained raw artifacts.
- [x] lubav.7 Perform a final guide-facing documentation pass: a reader can
  identify what Ex1 claims to implement, what remains local/provisional, how
  to reproduce the demo, and how to inspect evidence without mistaking the
  Docker topology or current source path for the protocol contract.
