# Receiving Package Plan

## Intended Families

- `moks.receiving.item.v1`
- `moks.receiving.receipt.v1`
- `moks.receiving.disposition.v1`

## Intended Claims

- `pcid:moks.receiving.item.v1` / `domain-behavior`
- `pcid:moks.receiving.receipt.v1` / `domain-behavior`
- `pcid:moks.receiving.disposition.v1` / `domain-behavior`

## Intended Commands

- `receiving create`
- `receiving list`
- `receiving inspect`
- `receiving record-receipt`
- `receiving record-disposition`

## Notes

This package should express receiving work as domain behavior over shared
`knowledge`, `runs`, and `context` records instead of inventing runtime-private
receiving state.

## Current State

Phase 1 is implemented as a first-party built-in package:

- receiving create / list / inspect
- receipt recording
- disposition recording
- validators for receiving item, receipt, and disposition records

This package is built-in for now because the installed-package mutation
contract is still too shallow for serious domain packages with multiple durable
record families. Source: `DI-zibek`.
