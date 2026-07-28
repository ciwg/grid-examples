# TODO moksu - ex6 foundation

## Decision Intent Log

ID: DI-moksu
Date: 2026-07-28 00:00:00
Status: active
Decision: Build ex6 as a PromiseGrid-native runtime with manifest plus self-check packages, runtime-mediated package communication, and unknown-family exact-byte store-and-relay.
Intent: Keep ex6 as the durable runtime that owns grid-facing coordination while letting built-in and installed packages extend it without direct package-to-package coupling.
Constraints: All artifacts stay inside `ex6-operational-knowledge-agent-runtime/`; ex5 is reference material only; browser is out of scope; installed packages use executables rather than Go plugin ABI.
Affects: `cmd/moks`, `builtin`, `grid`, `kernel`, `packages`, `records`, `store`, tests, and ex6-local docs.

## Notes

- The repo-local `mint-handle` helper was not present in the checked workspace during implementation, so this ex6-local TODO uses a manually chosen proquint-style handle to preserve the required DI linkage inside the allowed ex6 boundary.

ID: DI-lupok
Date: 2026-07-28 00:30:00
Status: active
Decision: Tighten ex6 around explicit PromiseGrid implementation claims and runtime-enforced protocol identity matching for registered families.
Intent: Keep ex6 protocol-first so packages do not imply conformance through family names alone and relay exports can state what protocols the active packages claim to implement.
Constraints: Preserve exact-byte unknown-family retention; do not introduce a fake universal handler ABI; keep all changes inside the ex6 tree.
Affects: `grid`, `kernel`, `packages`, builtin package manifests, tests, and CLI-visible package metadata.
Supersedes: DI-moksu

ID: DI-okar
Date: 2026-07-28 12:10:00
Status: active
Decision: Rename the ex6 project in docs and filesystem presentation to Operational Knowledge Agent Runtime (`OKAR`), and rename the folder to `ex6-operational-knowledge-agent-runtime`.
Intent: Align the project name with the newer agent-runtime framing while keeping the existing CLI/runtime command surface stable for now.
Constraints: Keep the current `moks` CLI name and existing DI handles; preserve buildability by updating the Go module path and imports together with the folder rename.
Affects: folder path, `go.mod`, Go imports, top-level docs, and project-facing titles inside the renamed tree.
