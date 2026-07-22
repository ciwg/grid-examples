# TODO sobek - ex5 PromiseGrid doc honesty and implementation claims

## Decision Intent Log

ID: DI-sobek
Date: 2026-07-21 00:00:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track the ex5 PromiseGrid doc-honesty pass and implementation-claims doc as one focused documentation slice.
Intent: Separate shipped local ex5 behavior from planned PromiseGrid wire behavior, then state the current implementation promises in one explicit technical document.
Constraints: Stay docs-only in this slice; do not claim frozen `pCID` handling, signed envelopes, or relay-visible peer exchange that the runtime does not yet implement.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-sobek-ex5-promisegrid-doc-honesty-and-claims.md`, `ex5-operational-knowledge-system/TODO/TODO-ragup-ex5-promisegrid-wire-slice-decision.md`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/product-overview.md`, `ex5-operational-knowledge-system/docs/architecture.md`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `ex5-operational-knowledge-system/operational-knowledge-system-spec-v0.1.md`

## Goal

Make the ex5 technical docs follow the PromiseGrid dev-guide boundary more
honestly and add one explicit statement of what the shipped implementation
actually promises today.

## Tasks

- [x] sobek.1 Audit the main ex5 technical docs for PromiseGrid overstatement or ambiguity.
- [x] sobek.2 Rewrite those docs to separate shipped local HTTP/runtime behavior from planned PromiseGrid wire behavior.
- [x] sobek.3 Add a dedicated implementation-claims doc for the current shipped ex5 contract.
- [x] sobek.4 Leave the real pCID/signed-envelope implementation decision as a separate later TODO.

## Status

- done
- derived from the 2026-07-21 PromiseGrid alignment pass
