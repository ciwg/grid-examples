# Maintenance Package Plan

## Intended Families

- `moks.maintenance.item.v1`
- `moks.maintenance.service.v1`
- `moks.maintenance.finding.v1`

## Intended Claims

- `pcid:moks.maintenance.item.v1` / `domain-behavior`
- `pcid:moks.maintenance.service.v1` / `domain-behavior`
- `pcid:moks.maintenance.finding.v1` / `domain-behavior`

## Intended Commands

- `maintenance create`
- `maintenance list`
- `maintenance inspect`
- `maintenance record-service`
- `maintenance record-finding`

## Notes

This package should express maintenance work as domain behavior over shared
`knowledge`, `runs`, and `context` records instead of inventing basket-private
maintenance state.

## Current State

Phase 1 is implemented as a first-party built-in package:

- maintenance create / list / inspect
- maintenance service recording
- maintenance finding recording
- validators for maintenance item, service, and finding records

This package is built-in for now because the installed-package mutation
contract is still too shallow for serious domain eggs with multiple durable
record families. Source: `DI-ramek`.
