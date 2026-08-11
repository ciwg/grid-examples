# Training Package Plan

## Intended Families

- `moks.training.item.v1`
- `moks.training.session.v1`
- `moks.training.completion.v1`

## Intended Claims

- `moks.training.item.v1` / `domain-behavior` (fixed pCID in the [registry](../../docs/protocols/package-family-pcid-registry.md))
- `moks.training.session.v1` / `domain-behavior` (fixed pCID in the registry)
- `moks.training.completion.v1` / `domain-behavior` (fixed pCID in the registry)

## Intended Commands

- `training create`
- `training list`
- `training inspect`
- `training record-session`
- `training certify`

## Notes

This package should express training work as domain behavior over the shared
`knowledge` and `runs` families instead of inventing runtime-private training
state.

## Current State

Phase 1 is implemented as a first-party built-in package:

- training create / list / inspect
- training session recording
- training completion recording
- validators for training item, session, and completion records

This package is built-in for now because the installed-package mutation
contract is still too shallow for serious domain packages with multiple durable
record families. Source: `DI-sivuk`.
