# Procedures Package Plan

## Intended Families

- `moks.procedures.item.v1`
- `moks.procedures.use.v1`

## Intended Claims

- `pcid:moks.procedures.item.v1` / `domain-behavior`
- `pcid:moks.procedures.use.v1` / `domain-behavior`

## Intended Commands

- `procedures create`
- `procedures inspect`
- `procedures record-use`

## Notes

This package should be the first serious domain egg built on top of `context`,
`knowledge`, `runs`, and `links`.

## Current State

Phase 1 is implemented as a first-party built-in package:

- procedure create
- procedure inspect
- procedure record-use
- validators for procedure item and procedure use records

This package is built-in for now because the installed-package mutation
contract is still too shallow for serious durable runtime writes. Source:
`DI-tusav`.
