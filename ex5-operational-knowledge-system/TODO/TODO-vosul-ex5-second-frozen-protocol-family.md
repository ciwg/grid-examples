# TODO vosul - freeze and claim the second ex5 PromiseGrid protocol family

## Decision Intent Log

ID: DI-vosul
Date: 2026-07-22 00:07:49 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Use `knowledge-approval` as the second ex5 frozen PromiseGrid family, freeze it as one family covering both knowledge-item and run approvals, publish the next implementation claim against it, and add the second local signed-envelope runtime slice for approval artifacts while keeping browser/CLI/Neovim on the current HTTP adapter.
Intent: Continue the ex5 PromiseGrid migration with the next trust-bearing durable family that already hangs directly off the current item/run workflow, without reopening transport or broad runtime rewrites.
Constraints: Keep the migration staged and additive; do not change the local embodiment adapter contract; keep lifecycle status changes in the `knowledge-item` family; freeze only the durable approval artifact contract in this slice.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-vosul-ex5-second-frozen-protocol-family.md`, `docs/thought-experiments/TE-tipav-ex5-knowledge-approval-family-boundary.md`, `ex5-operational-knowledge-system/protocols/knowledge-approval.md`, `ex5-operational-knowledge-system/protocols/**`, `ex5-operational-knowledge-system/service/**`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `ex5-operational-knowledge-system/docs/architecture.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/CHANGELOG.md`

## Goal

Freeze `knowledge-approval` as the second ex5 PromiseGrid family and add the
next local signed-envelope runtime slice over the current approval workflow.

## Tasks

- [x] vosul.1 Freeze the `knowledge-approval` protocol boundary and publish the next implementation claim.
- [x] vosul.2 Add the second local signed-envelope runtime slice for durable approval artifacts.
- [x] vosul.3 Keep browser/CLI/Neovim on the current HTTP adapter while replay and verification cover the new family.

## Status

- done
- second frozen family and signed approval runtime slice implemented
