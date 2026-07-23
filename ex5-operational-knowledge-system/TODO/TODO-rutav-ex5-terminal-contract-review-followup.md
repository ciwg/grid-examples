# TODO rutav - ex5 terminal contract review followup

## Decision Intent Log

ID: DI-lurav
Date: 2026-07-22 18:03:26 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Cover the user-facing `-socket=off` CLI contract through a small pure startup-resolution helper and direct unit tests, rather than a subprocess harness.
Intent: Keep the embodiment boundary explicit by testing startup transport selection directly while leaving the existing socket and HTTP transport tests in place.
Constraints: Preserve the already-locked fail-closed default; keep the helper narrow and avoid broad CLI startup refactors.
Affects: `ex5-operational-knowledge-system/cmd/oks-cli/main.go`, `ex5-operational-knowledge-system/cmd/oks-cli/main_test.go`, `docs/thought-experiments/TE-lurav-ex5-cli-http-opt-in-test-shape.md`

ID: DI-rutav
Date: 2026-07-22 18:03:26 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track the remaining ex5 review findings around terminal transport documentation, CLI `-socket=off` coverage, and stale browser/CLI-only wording.
Intent: Close the last small PromiseGrid-alignment and review gaps without reopening the larger transport design that was already locked in TODO `127`.
Constraints: Keep the current fail-closed CLI transport rule; focus on test coverage for the explicit HTTP opt-in plus doc/comment honesty.
Affects: `ex5-operational-knowledge-system/cmd/oks-cli/main.go`, `ex5-operational-knowledge-system/cmd/oks-cli/main_test.go`, `ex5-operational-knowledge-system/cmd/operational-knowledge/main.go`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/terminal-capability-matrix.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/TODO/TODO.md`, `docs/thought-experiments/TE-lurav-ex5-cli-http-opt-in-test-shape.md`

## Goal

Close the remaining ex5 review findings by:

- documenting the CLI's fail-closed local-socket rule and explicit
  `-socket=off` compatibility mode across the remaining operator docs,
- adding direct test coverage for the user-facing `-socket=off` contract,
- and removing the last browser/CLI-only wording drift.

## Tasks

- [x] rutav.1 Run a TE on the right test shape for the `-socket=off` contract and lock the coverage approach.
- [x] rutav.2 Add the chosen CLI coverage and any small supporting refactor needed for it.
- [x] rutav.3 Align the remaining HTTP-guide, terminal-matrix, feature-guide, and server-entrypoint wording to the shipped three-embodiment state.

## Status

- closed
- resolved by helper-level `-socket=off` coverage plus remaining transport wording cleanup
