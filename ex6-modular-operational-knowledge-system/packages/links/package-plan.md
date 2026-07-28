# Links Package Plan

## Intended Families

- `moks.links.typed.v1`

## Intended Claims

- `pcid:moks.links.typed.v1` / `family-validator`

## Intended Commands

- `links create`
- `links inspect`

## Notes

This package should keep cross-record relationships explicit and portable
across eggs.

## Current State

Phase 1 is implemented as a first-party built-in package:

- typed-link create
- typed-link inspect
- family validator for the typed-link family

This package is built-in for now because the installed-package mutation
contract is still too shallow for serious durable runtime writes. Source:
`DI-figar`.
