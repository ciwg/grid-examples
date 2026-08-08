# TODO kuzan — Document root handle-minting workflow

## Decision Intent Log

### DI-hidiz

- ID: DI-hidiz
- Date: 2026-08-07 21:35:37 -0700
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Interpret a user request to "mint a handle" as an instruction to run `tools/mint-handle` from the repository root and use its exact result for the new coordination record.
- Intent: Agents should consistently use the shared allocator rather than manually choosing a proquint.
- Constraints: The ledger remains append-only and must be committed with the consuming artifact; do not reuse historical handles.
- Affects: `AGENTS.md`, `TODO/handle-namespace.tsv`, `TODO/TODO.md`, and this TODO record.

## Tasks

- [x] kuzan.1 Add the repo-wide handle-minting instruction to `AGENTS.md`.
- [x] kuzan.2 Commit and push the instruction.
