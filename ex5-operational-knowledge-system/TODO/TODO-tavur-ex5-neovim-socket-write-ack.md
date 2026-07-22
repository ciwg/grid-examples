# TODO tavur - ex5 Neovim socket write acknowledgment

## Decision Intent Log

ID: DI-vatik
Date: 2026-07-22 16:43:04 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track Neovim live-socket write acknowledgment and fallback correctness as a dedicated follow-on after the `117A` rollout.
Intent: Prevent Neovim from treating a failed local-socket live update as success and silently skipping the compatibility fallback path.
Constraints: Keep the existing live-draft semantics intact; focus on write completion/error handling and fallback behavior only.
Affects: `ex5-operational-knowledge-system/nvim/lua/oks/init.lua`, `ex5-operational-knowledge-system/nvim/*_test.go`, `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-tavur-ex5-neovim-socket-write-ack.md`

ID: DI-sudik
Date: 2026-07-22 17:34:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Treat Neovim live-socket updates as successful only after the async write callback succeeds; otherwise reconnect and allow compatibility fallback immediately.
Intent: Make live-draft delivery truth honest so failed local-socket writes do not get reported as successful transport updates.
Constraints: Preserve the existing live-draft semantics and reconnect loop; do not add a new local retry queue.
Affects: `ex5-operational-knowledge-system/nvim/lua/oks/init.lua`, `ex5-operational-knowledge-system/nvim/snapshot_test.go`, `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-tavur-ex5-neovim-socket-write-ack.md`

## Goal

Make Neovim live-draft updates acknowledge Unix-socket write success honestly
before suppressing fallback behavior.

## Tasks

- [x] tavur.1 Define the required error/ack behavior for Neovim live-socket writes.
- [x] tavur.2 Implement write-success detection plus fallback/reconnect handling for failed writes.
- [x] tavur.3 Add regression coverage for dropped or broken live-socket writes.

## Status

- done
- fixed by the post-`117A` socket hardening pass
- Neovim now falls back honestly when a live-socket write does not complete
