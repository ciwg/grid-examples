# TODO rusav - ex5 absolute socket path advertisement

## Decision Intent Log

ID: DI-rakuv
Date: 2026-07-22 17:35:03 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Canonicalize the runtime `data-root` once at startup and derive the advertised local socket path from that canonical root.
Intent: Make `/api/meta` authoritative for terminal socket discovery even when the server and clients start from different working directories.
Constraints: Preserve runtime-first discovery from TODO `121`; fix the root cause at runtime startup instead of teaching clients to reinterpret relative paths.
Affects: `ex5-operational-knowledge-system/service/app.go`, `ex5-operational-knowledge-system/service/local_socket.go`, `ex5-operational-knowledge-system/service/server_test.go`, `ex5-operational-knowledge-system/TODO/TODO.md`, `docs/thought-experiments/TE-zavuk-ex5-absolute-socket-path-advertisement.md`
Supersedes: DI-rusav

ID: DI-rusav
Date: 2026-07-22 17:23:59 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track the remaining gap where `/api/meta` can still advertise a relative local socket path that is wrong outside the runtime process cwd.
Intent: Make the runtime authoritative for terminal socket discovery in a way that stays correct across custom `-data-root` values and different client working directories.
Constraints: Preserve the current runtime-first discovery model from TODO `121`; focus on path canonicalization and correctness, not on changing the transport choice again.
Affects: `ex5-operational-knowledge-system/cmd/operational-knowledge/main.go`, `ex5-operational-knowledge-system/service/local_socket.go`, `ex5-operational-knowledge-system/service/app.go`, `ex5-operational-knowledge-system/service/server_test.go`, `ex5-operational-knowledge-system/cmd/oks-cli/main.go`, `ex5-operational-knowledge-system/nvim/lua/oks/init.lua`, `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-rusav-ex5-absolute-socket-path-advertisement.md`

## Goal

Ensure the runtime publishes one canonical absolute local socket path so
runtime-first terminal discovery remains correct even when the server and the
client start from different working directories.

## Tasks

- [x] rusav.1 Define where the runtime canonicalizes its `data-root` and local socket path before advertising them.
- [x] rusav.2 Implement absolute socket-path publication through `/api/meta` and any other affected terminal discovery paths.
- [x] rusav.3 Add regression coverage for server/client cwd mismatch with a relative runtime root.

## Status

- closed
- resolved by canonicalizing the runtime root before socket-path advertisement
