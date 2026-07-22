# TODO vupok - sync ex5 HTTP API guide with the shipped PromiseGrid adapter contract

## Decision Intent Log

ID: DI-vupok
Date: 2026-07-22 14:18:44 -0700
Status: active
Decision: Update the HTTP API guide to document the shipped `/api/meta`
capability fields and `/api/peer-exchange/*` routes while keeping the guide
explicit that HTTP is the local embodiment adapter, not the PromiseGrid peer
contract.
Intent: Remove adapter-contract drift without reopening transport design or
implying that route names are frozen wire semantics.
Constraints: Keep wording consistent with the shipped runtime and the current
PromiseGrid implementation claims.
Affects: docs/http-api-guide.md; TODO/TODO.md

## Goal

Bring `docs/http-api-guide.md` up to the current shipped adapter contract so it
describes the real `/api/meta` capability surface and the shipped
`/api/peer-exchange/*` routes.

## Why this exists

The runtime now exposes peer-exchange format/family metadata, CAS draft-body
support, and the `local_http` embodiment adapter signal through `/api/meta`,
and it ships peer exchange export/import routes. The HTTP API guide still
documents only the older metadata surface and does not document the peer
exchange endpoints.

## Tasks

- [x] vupok.1 Update the `/api/meta` section to include the shipped capability
  fields and their meaning.
- [x] vupok.2 Add the shipped `/api/peer-exchange/export` and
  `/api/peer-exchange/import` routes to the guide.
- [x] vupok.3 Align wording so the guide clearly treats HTTP as the local
  embodiment adapter, not the PromiseGrid peer contract.

## Status

- closed
- created from the post-112 PromiseGrid alignment pass
