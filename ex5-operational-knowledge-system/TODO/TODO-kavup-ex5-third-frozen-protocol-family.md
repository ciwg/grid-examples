# TODO kavup - freeze and claim the third ex5 PromiseGrid protocol family

## Decision Intent Log

ID: DI-kavup
Date: 2026-07-22 09:45:15 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Use `knowledge-evidence` as the third ex5 frozen PromiseGrid family, freeze it around durable run evidence artifacts, publish the next implementation claim against it, and add the third local signed-envelope runtime slice while keeping browser/CLI/Neovim on the current HTTP adapter.
Intent: Continue the ex5 PromiseGrid migration with the next durable trust-bearing family that already hangs off run history and review, without reopening transport or broad runtime rewrites.
Constraints: Keep the migration staged and additive; do not change the local embodiment adapter contract; freeze only the durable evidence artifact contract in this slice; keep attachment bytes on the current copied-file storage path while the evidence family claims the durable evidence metadata and attachment reference.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-kavup-ex5-third-frozen-protocol-family.md`, `docs/thought-experiments/TE-ribof-ex5-knowledge-evidence-family-boundary.md`, `ex5-operational-knowledge-system/protocols/knowledge-evidence.md`, `ex5-operational-knowledge-system/protocols/**`, `ex5-operational-knowledge-system/service/**`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `ex5-operational-knowledge-system/docs/architecture.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/CHANGELOG.md`

ID: DI-ribof
Date: 2026-07-22 09:45:15 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Lock the first frozen `knowledge-evidence` family to evidence metadata plus attachment references, not raw attachment bytes, and add a stable durable evidence ID to the underlying event/payload before signing.
Intent: Make the third frozen family trustworthy and replay-stable without overreaching into attachment-byte transport or storage redesign.
Constraints: Attachment name, path/reference, and size are in scope; copied attachment bytes stay on the current runtime storage path; old evidence events without a stored evidence ID must still replay correctly.
Affects: `ex5-operational-knowledge-system/TODO/TODO-kavup-ex5-third-frozen-protocol-family.md`, `docs/thought-experiments/TE-ribof-ex5-knowledge-evidence-family-boundary.md`, `ex5-operational-knowledge-system/protocols/knowledge-evidence.md`, `ex5-operational-knowledge-system/service/**`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `ex5-operational-knowledge-system/CHANGELOG.md`

## Goal

Freeze `knowledge-evidence` as the third ex5 PromiseGrid family and add the
next local signed-envelope runtime slice over the current evidence workflow.

## Tasks

- [x] kavup.1 Freeze the `knowledge-evidence` protocol boundary and publish the next implementation claim.
- [x] kavup.2 Add the third local signed-envelope runtime slice for durable evidence artifacts.
- [x] kavup.3 Keep browser/CLI/Neovim on the current HTTP adapter while replay and verification cover the new family.

## Status

- done
- third frozen family and signed evidence runtime slice implemented
