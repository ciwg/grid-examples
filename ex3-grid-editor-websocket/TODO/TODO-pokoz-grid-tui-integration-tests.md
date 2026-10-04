# TODO pokoz - Grid TUI integration tests

## Decision Intent Log

ID: DI-gufuj
Date: 2026-10-04 00:00:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add isolated end-to-end coverage for browser↔Grid TUI and Grid TUI↔Grid TUI fresh-document collaboration.
Intent: Prove that Ex3's terminal embodiment joins the same relay-backed
document and awareness contracts as the browser and a second terminal client,
rather than relying only on model-level rendering tests or browser↔sidecar
coverage.
Constraints: Exercise the existing relay and Automerge sidecar processes; use
temporary relay roots; assert shared text plus peer identity, color, and cursor
awareness; do not alter pCID-selected protocol meanings or promote local test
evidence into a PromiseGrid interoperability claim.
Affects: `tui/integration_test.go`, `docs/testing.md`, this TODO, and the root
TODO index.

## Goal

Close the documented fresh-document integration gap for Grid TUI without
changing Ex3's relay, Automerge, WebSocket, or awareness contracts.

## Tasks

- [x] pokoz.1 Add browser↔Grid TUI convergence coverage through an isolated relay.
- [x] pokoz.2 Add Grid TUI↔Grid TUI convergence and awareness coverage through an isolated relay.
- [x] pokoz.3 Document the exact proof boundary and verification commands.
