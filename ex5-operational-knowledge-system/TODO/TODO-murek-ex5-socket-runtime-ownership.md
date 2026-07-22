# TODO murek - ex5 socket runtime ownership

## Decision Intent Log

ID: DI-lusek
Date: 2026-07-22 16:43:04 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track the local Unix-socket runtime-ownership collision as a dedicated follow-on after the `117A` rollout.
Intent: Prevent a second `operational-knowledge` process from silently stealing or unlinking the active embodiment socket for the first runtime.
Constraints: Treat this as a runtime-safety fix inside the shipped socket embodiment slice; do not reopen the broader `117A` transport boundary.
Affects: `ex5-operational-knowledge-system/service/local_socket.go`, `ex5-operational-knowledge-system/cmd/operational-knowledge/main.go`, `ex5-operational-knowledge-system/service/*_test.go`, `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-murek-ex5-socket-runtime-ownership.md`

ID: DI-fegom
Date: 2026-07-22 17:34:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Fail startup when another live runtime already owns the embodiment socket, and remove only stale socket files before rebinding.
Intent: Keep local embodiment ownership explicit so a second runtime cannot silently steal the first runtime's direct contract endpoint.
Constraints: Use a dedicated ownership guard named `ensureOwnedSocketListener`; preserve the existing single-runtime contract and only treat clearly inactive socket files as reclaimable.
Affects: `ex5-operational-knowledge-system/service/local_socket.go`, `ex5-operational-knowledge-system/service/local_socket_test.go`, `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-murek-ex5-socket-runtime-ownership.md`

## Goal

Make the direct local Unix-socket contract robust against multiple local
runtime processes pointing at the same `data-root`.

## Tasks

 - [x] murek.1 Define the expected behavior when a second runtime starts against an already-active embodiment socket.
 - [x] murek.2 Implement collision-safe socket startup and shutdown semantics.
 - [x] murek.3 Add regression coverage for same-root double-start behavior.

## Status

- done
- fixed by the post-`117A` socket hardening pass
- same-root double-start now fails without stealing the active socket
