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
