# TODO musav - ex5 terminal transport alignment pass

## Decision Intent Log

ID: DI-zorav
Date: 2026-07-22 17:45:47 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Make the CLI fail closed when the preferred local Unix socket is unavailable, and require explicit HTTP opt-in through `-socket=off`.
Intent: Keep the shipped direct terminal embodiment contract operationally honest instead of silently demoting commands onto the HTTP compatibility surface.
Constraints: Preserve HTTP compatibility transport for explicit operator choice; do not change browser or Neovim transport policy in this pass.
Affects: `ex5-operational-knowledge-system/cmd/oks-cli/main.go`, `ex5-operational-knowledge-system/cmd/oks-cli/socket_transport.go`, `ex5-operational-knowledge-system/cmd/oks-cli/main_test.go`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/architecture.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`, `docs/thought-experiments/TE-nurek-ex5-cli-socket-fallback-honesty.md`

ID: DI-musav
Date: 2026-07-22 17:45:47 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track the remaining ex5 PromiseGrid alignment gap around terminal transport honesty and the last long-form doc drift.
Intent: Tighten the direct terminal embodiment contract so the shipped socket-first design is described honestly and does not silently hide transport demotion in normal operator use.
Constraints: Keep the browser on the current HTTP adapter; focus this pass on CLI socket-fallback semantics plus stale architecture/implementation/README wording.
Affects: `ex5-operational-knowledge-system/cmd/oks-cli/socket_transport.go`, `ex5-operational-knowledge-system/cmd/oks-cli/main.go`, `ex5-operational-knowledge-system/cmd/oks-cli/main_test.go`, `ex5-operational-knowledge-system/docs/architecture.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-musav-ex5-terminal-transport-alignment-pass.md`

## Goal

Align the remaining ex5 terminal-transport behavior and docs with the shipped
PromiseGrid embodiment split: browser on the HTTP adapter, CLI and Neovim on
the direct local Unix-socket contract with explicit compatibility fallback.

## Tasks

- [x] musav.1 Run a TE on CLI local-socket fallback semantics and lock the transport-honesty choice.
- [x] musav.2 Implement the chosen CLI behavior and add regression coverage.
- [x] musav.3 Update the remaining long-form docs so Neovim and terminal transport wording matches the shipped socket-first design.

## Status

- closed
- resolved by fail-closed CLI socket behavior with explicit `-socket=off` HTTP opt-in, plus terminal transport doc alignment
