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
