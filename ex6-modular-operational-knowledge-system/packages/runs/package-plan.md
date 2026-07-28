# Runs Package Plan

## Intended Families

- `moks.runs.run.v1`
- `moks.runs.evidence.v1`
- `moks.runs.approval.v1`

## Intended Claims

- `pcid:moks.runs.run.v1` / `family-validator`
- `pcid:moks.runs.evidence.v1` / `family-validator`
- `pcid:moks.runs.approval.v1` / `family-validator`

## Intended Commands

- `runs record`
- `runs inspect`
- `runs evidence add`
- `runs approve`

## Notes

This package should carry the durable record of work, attached evidence, and
approval outcomes.

## Current State

Phase 1 is implemented as a first-party built-in package:

- run record
- evidence add
- run approve
- run inspect
- family validators for run, evidence, and approval records

This package is built-in for now because the installed-package mutation
contract is still too shallow for serious durable runtime writes. Source:
`DI-pamuk`.
