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

ID: DI-noruv
Date: 2026-07-22 09:58:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Implement TODO `005` as a shared websocket-preferred live-draft transport for both browser and Neovim, while keeping the existing `/api/items/{id}/live` HTTP routes as fallback and compatibility paths.
Intent: Remove HTTP polling as the primary live collaboration transport in `ex5` without reopening the durable PromiseGrid family design or forcing a flag-day client migration.
Constraints: Websocket remains carriage for shared draft state and presence only; durable revisions, runs, approvals, evidence, and peer exchange stay on their existing routes; browser and Neovim must both prefer websocket when available; HTTP live routes stay working for fallback, bootstrap fetch, and older clients.
Affects: `ex5-operational-knowledge-system/service/**`, `ex5-operational-knowledge-system/web/app.js`, `ex5-operational-knowledge-system/nvim/lua/oks/**`, `ex5-operational-knowledge-system/docs/**`, `ex5-operational-knowledge-system/TODO/TODO.md`, `docs/thought-experiments/TE-tivok-ex5-websocket-collaboration-transport-scope.md`

## Goal

Port `ex5` from its current local HTTP live-draft transport to a real
websocket/relay collaboration path if that direction is confirmed.

## Tasks

- [x] masad.1 Copy or adapt the standalone websocket collaboration machinery needed for `ex5` without creating runtime dependencies on `ex3`.
- [x] masad.2 Replace HTTP-polled live draft updates with websocket sync, presence, and peer-state transport for browser and Neovim collaboration.
- [x] masad.3 Add tests and docs for websocket collaboration behavior, failure modes, and operator setup.

## Result

`ex5` now ships websocket-preferred shared live-draft transport for both the
browser authoring surface and the first-phase Neovim embodiment, while keeping
`GET/POST /api/items/{id}/live` as fallback and compatibility routes under the
same local adapter. Source: `DI-noruv`.
