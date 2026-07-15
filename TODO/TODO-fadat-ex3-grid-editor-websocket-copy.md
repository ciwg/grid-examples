# TODO fadat - ex3 grid editor websocket copy

## Decision Intent Log

ID: DI-norov
Date: 2026-07-14 16:42:54 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Create `ex3-grid-editor-websocket/` as a copy of `ex2-grid-editor/` while leaving every file under `ex2-grid-editor/` unchanged.
Intent: Start the websocket-oriented follow-on example from the current grid editor codebase without disturbing the existing example.
Constraints: The work must not edit files under `ex2-grid-editor/`; the new example lives as a sibling directory under the repo root; the first step is a copy, not a websocket redesign.
Affects: `ex3-grid-editor-websocket/**`, `TODO/TODO.md`, `TODO/TODO-fadat-ex3-grid-editor-websocket-copy.md`

ID: DI-talat
Date: 2026-07-14 16:42:54 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Keep the copy minimal by changing only what is required for the new example to coexist cleanly with `ex2-grid-editor`, including the nested Go module path and matching Go imports.
Intent: Preserve the copied behavior and documentation shape while avoiding immediate collisions in local development and Go package resolution.
Constraints: Broader product wording, copied TODO history, and most human-facing `grid-editor` strings stay as inherited copy content for now; rename the nested module path to `github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket`.
Affects: `ex3-grid-editor-websocket/go.mod`, `ex3-grid-editor-websocket/**/*.go`, `ex3-grid-editor-websocket/README.md`, `ex3-grid-editor-websocket/docs/**`

ID: DI-vatub
Date: 2026-07-14 16:42:54 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Make the new copy runnable beside `ex2-grid-editor` by moving its default local and Docker relay ports to `127.0.0.1:7025` and `127.0.0.1:7026`.
Intent: Let both examples run on the same machine without port conflicts while preserving the two-relay local simulation model.
Constraints: Update default listen URLs, peer URLs, launch scripts, and docs that would otherwise point operators to `7015` or `7016`.
Affects: `ex3-grid-editor-websocket/cmd/**`, `ex3-grid-editor-websocket/compose.yaml`, `ex3-grid-editor-websocket/scripts/**`, `ex3-grid-editor-websocket/nvim/**`, `ex3-grid-editor-websocket/README.md`, `ex3-grid-editor-websocket/docs/**`

ID: DI-fohif
Date: 2026-07-14 16:42:54 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Copy the repo-local runtime directory into `ex3-grid-editor-websocket/.grid-editor/` and reseed only `identity_ed25519_seed` so the new example starts as a distinct local identity.
Intent: Preserve the current runtime and data layout as part of the copy without forcing `ex2` and `ex3` to share the same local signing identity.
Constraints: Keep the copied `.grid-editor/` tree; do not edit `ex2-grid-editor/.grid-editor/`; only reset the copied identity seed, letting the app recreate it on first start if needed.
Affects: `ex3-grid-editor-websocket/.grid-editor/**`, `ex3-grid-editor-websocket/service/app.go`

ID: DI-rokod
Date: 2026-07-14 16:42:54 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Keep the copied Go runtime verification-clean by fixing inherited unchecked close paths only inside `ex3-grid-editor-websocket/`.
Intent: Let the copied example pass the repo's required `errcheck` verification without changing `ex2-grid-editor/` or widening the copy into a refactor.
Constraints: Limit the cleanup to explicit close-result handling in the copied files that `errcheck` reports; do not use the cleanup as a pretext for broader behavioral redesign.
Affects: `ex3-grid-editor-websocket/service/app.go`, `ex3-grid-editor-websocket/store/log.go`

## Goal

Create a new sibling example named `ex3-grid-editor-websocket` by copying the
current `ex2-grid-editor` tree and then applying the minimum separation changes
required for safe parallel development and execution.

## Tasks

- [x] fadat.1 Copy `ex2-grid-editor/` to `ex3-grid-editor-websocket/`, including `.grid-editor/`.
- [x] fadat.2 Rename the nested Go module path and matching Go imports in the copied tree.
- [x] fadat.3 Move the copied example's default relay ports from `7015/7016` to `7025/7026` anywhere the defaults would collide with `ex2`.
- [x] fadat.4 Reset the copied runtime identity seed while preserving the rest of the copied `.grid-editor/` tree.
- [x] fadat.5 Verify the new copied example builds/tests without modifying `ex2-grid-editor/`.

## Evidence

- `ex3-grid-editor-websocket/` exists as a sibling copy of `ex2-grid-editor/`.
- The copied nested Go module path is `github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket`.
- The copied default relay ports now use `127.0.0.1:7025` and `127.0.0.1:7026` in the copied command defaults, launcher, compose file, and operator docs.
- At implementation time, `ex2-grid-editor/.grid-editor/` was not present in the working tree, so there was no runtime identity seed to copy or reseed in `ex3-grid-editor-websocket/`.
- Verification passed with `go test ./...` and `errcheck ./...` from `ex3-grid-editor-websocket/`.
