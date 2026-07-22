# TODO pobek - ex5 socket path robustness

## Decision Intent Log

ID: DI-roben
Date: 2026-07-22 16:43:04 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track local-socket path robustness and default-path clarity as a dedicated follow-on after the `117A` rollout.
Intent: Reduce accidental fallback-to-HTTP behavior caused by relative socket defaults, changed working directories, or non-default runtime roots.
Constraints: Keep the current socket-first embodiment choice; focus on path resolution, defaults, and operator-facing reliability rather than reopening transport selection.
Affects: `ex5-operational-knowledge-system/cmd/oks-cli/main.go`, `ex5-operational-knowledge-system/nvim/lua/oks/init.lua`, `ex5-operational-knowledge-system/scripts/oks-nvim`, `ex5-operational-knowledge-system/docs/*`, `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-pobek-ex5-socket-path-robustness.md`

ID: DI-vorag
Date: 2026-07-22 17:34:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Resolve the default terminal embodiment socket path to an absolute path rooted from the repo/runtime launch root, with explicit env or flag overrides still taking precedence.
Intent: Remove working-directory accidents from the direct local contract so terminal embodiments reach the intended runtime more reliably.
Constraints: Use a helper named `defaultSocketPath` for the CLI, preserve explicit overrides, and keep browser HTTP compatibility untouched.
Affects: `ex5-operational-knowledge-system/cmd/oks-cli/main.go`, `ex5-operational-knowledge-system/cmd/oks-cli/main_test.go`, `ex5-operational-knowledge-system/nvim/lua/oks/init.lua`, `ex5-operational-knowledge-system/nvim/snapshot_test.go`, `ex5-operational-knowledge-system/scripts/oks-nvim`, `ex5-operational-knowledge-system/docs/architecture.md`, `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-pobek-ex5-socket-path-robustness.md`

## Goal

Make the direct local Unix-socket path predictable across working directories,
launcher paths, and non-default runtime roots.

## Tasks

- [x] pobek.1 Define the canonical default socket-path resolution rule for CLI and Neovim.
- [x] pobek.2 Implement path resolution or discovery so terminal embodiments reach the intended runtime more reliably.
- [x] pobek.3 Add tests and docs for non-default root and changed-working-directory behavior.

## Status

- done
- fixed by the post-`117A` socket hardening pass
- terminal embodiments now resolve the default socket path from the repo/runtime launch root instead of the ambient working directory
