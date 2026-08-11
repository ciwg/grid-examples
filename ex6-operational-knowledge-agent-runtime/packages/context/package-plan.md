# Context Package Plan

## Intended Families

- `moks.context.responsibility.v1`
- `moks.context.place.v1`
- `moks.context.resource.v1`

## Intended Claims

- `moks.context.responsibility.v1` / `family-validator` (fixed pCID in the [registry](../../docs/protocols/package-family-pcid-registry.md))
- `moks.context.place.v1` / `family-validator` (fixed pCID in the registry)
- `moks.context.resource.v1` / `family-validator` (fixed pCID in the registry)

## Intended Commands

- `context responsibility create`
- `context responsibility inspect`
- `context place create`
- `context place inspect`
- `context resource create`
- `context resource inspect`

## Notes

This package should become the shared operational context layer for the rest of
ex6.

## Current State

Phase 1 is implemented as a first-party built-in package:

- responsibility create / list / inspect
- place create / list / inspect
- resource create / list / inspect
- family validators for the three context families

This package is built-in for now because the installed-package mutation
contract is not yet deep enough to let external packages write durable runtime
state cleanly. Source: `DI-lorup`.
