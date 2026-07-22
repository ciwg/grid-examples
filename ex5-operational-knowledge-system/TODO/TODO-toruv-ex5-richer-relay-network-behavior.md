# TODO toruv - ex5 richer relay-network behavior

## Decision Intent Log

ID: DI-toruv
Date: 2026-07-22 15:11:29 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Open a new ex5 follow-on wave for richer relay-network behavior after the current PromiseGrid scope is complete.
Intent: Keep the next PromiseGrid expansion explicit and separate from the now-closed local-adapter, CAS, websocket, and Neovim embodiment waves.
Constraints: Treat the current shipped peer exchange as the starting point; do not reopen frozen-family, canonical-ID, or CAS-authority decisions unless the relay TE proves a new dependency; keep direct non-HTTP embodiment contracts out of this TODO unless a later decision explicitly pulls them in.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-toruv-ex5-richer-relay-network-behavior.md`, `docs/thought-experiments/TE-relav-ex5-richer-relay-network-scope.md`

ID: DI-pazek
Date: 2026-07-22 15:18:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Implement TODO `115` as an incremental relay feed over origin-aware signed records with separate blob transfer by CID. The first slice uses per-origin cursor requests, `POST` relay-feed export/import routes over the current local HTTP adapter, and `GET`/`PUT` relay blob routes for raw CAS object transfer.
Intent: Advance `ex5` beyond local-adapter bundle exchange toward a cleaner PromiseGrid relay shape without reintroducing whole-bundle blob carriage or conflating durable feed exchange with live collaboration transport.
Constraints: Keep the existing peer-exchange bundle routes working for compatibility; require evidence blobs to be staged into local CAS before relay-feed import of evidence records succeeds; do not reopen frozen-family, canonical-ID, websocket, or non-HTTP embodiment decisions in this slice.
Affects: `ex5-operational-knowledge-system/service/**`, `ex5-operational-knowledge-system/docs/**`, `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-toruv-ex5-richer-relay-network-behavior.md`, `docs/thought-experiments/TE-relav-ex5-richer-relay-network-scope.md`

## Goal

Extend `ex5` beyond local-adapter peer exchange into a richer relay-network
shape that remains honest to the current signed-family, origin-aware, and
CAS-backed runtime.

## Tasks

- [x] toruv.1 Open the relay-network expansion wave and capture the scope question as its own TODO.
- [x] toruv.2 Run the required TE for the next relay-network step beyond local-adapter export/import.
- [x] toruv.3 Lock the surviving relay-network design and implementation slice.
- [x] toruv.4 Implement the chosen relay-network slice with matching tests, claims, and docs.

## Status

- closed
- `TE-relav` completed
- `DI-pazek` locks incremental feed plus CID-separated blob transfer
- `toruv.4` now ships `/api/relay/feed/export`, `/api/relay/feed/import`,
  and `/api/relay/blobs/{cid}` over the current local adapter
- starts from the shipped eight-family origin-aware peer-exchange runtime
- implementation shipped
