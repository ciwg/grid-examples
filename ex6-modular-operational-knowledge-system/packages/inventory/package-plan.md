# Inventory Package Plan

## Intended Families

- `moks.inventory.item.v1`
- `moks.inventory.count.v1`
- `moks.inventory.reconcile.v1`

## Intended Claims

- `pcid:moks.inventory.item.v1` / `domain-behavior`
- `pcid:moks.inventory.count.v1` / `domain-behavior`
- `pcid:moks.inventory.reconcile.v1` / `domain-behavior`

## Intended Commands

- `inventory create`
- `inventory list`
- `inventory inspect`
- `inventory record-count`
- `inventory record-reconcile`

## Notes

This package should express inventory work as domain behavior over shared
`knowledge`, `runs`, and `context` records instead of inventing basket-private
inventory state.

## Current State

Phase 1 is implemented as a first-party built-in package:

- inventory create / list / inspect
- count recording
- reconciliation recording
- validators for inventory item, count, and reconciliation records

This package is built-in for now because the installed-package mutation
contract is still too shallow for serious domain eggs with multiple durable
record families. Source: `DI-lavom`.
