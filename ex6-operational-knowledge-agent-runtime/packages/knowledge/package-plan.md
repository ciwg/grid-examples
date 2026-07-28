# Knowledge Package Plan

## Intended Families

- `moks.knowledge.item.v1`
- `moks.knowledge.revision.v1`
- `moks.knowledge.lifecycle.v1`

## Intended Claims

- `pcid:moks.knowledge.item.v1` / `family-validator`
- `pcid:moks.knowledge.revision.v1` / `family-validator`
- `pcid:moks.knowledge.lifecycle.v1` / `family-validator`

## Intended Commands

- `knowledge item create`
- `knowledge item inspect`
- `knowledge revision snapshot`
- `knowledge item approve`
- `knowledge item supersede`

## Notes

This package should own the modular replacement for the ex5 knowledge-item and
revision surface.

## Current State

Phase 1 is implemented as a first-party built-in package:

- item create / list / inspect
- revision snapshot with CAS-backed body storage
- approve / supersede lifecycle events
- family validators for item, revision, and lifecycle records

This package is built-in for now because the installed-package mutation
contract is still too shallow for serious durable runtime writes. Source:
`DI-vakod`.
