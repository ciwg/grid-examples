# TODO masad - ex5 websocket collaboration transport

## Decision Intent Log

ID: DI-masad
Date: 2026-07-20 11:56:32 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Capture the missing websocket collaboration work as its own ex5 TODO instead of implying that the current HTTP live-draft layer is the end state.
Intent: Keep the gap between the current runnable ex5 collaboration model and the intended PromiseGrid-style live transport visible and actionable.
Constraints: This TODO stays pending until the open collaboration decision is resolved; it should not force implementation ahead of the product choice.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-masad-ex5-websocket-collaboration-transport.md`

## Goal

Port `ex5` from its current local HTTP live-draft transport to a real
websocket/relay collaboration path if that direction is confirmed.

## Tasks

- [ ] masad.1 Copy or adapt the standalone websocket collaboration machinery needed for `ex5` without creating runtime dependencies on `ex3`.
- [ ] masad.2 Replace HTTP-polled live draft updates with websocket sync, presence, and peer-state transport for browser collaboration.
- [ ] masad.3 Add tests and docs for websocket collaboration behavior, failure modes, and operator setup.
