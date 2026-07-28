# Inventory Package

Planned ownership:

- inventory item kind
- inventory count session records
- inventory reconciliation records

This package should carry the first dedicated inventory domain behavior above
shared `knowledge`, `runs`, and `context` families.

Current state: first built-in implementation now exists in `package.go`, with
commands for inventory creation, listing, inspection, count recording, and
reconciliation recording. Source: `DI-lavom`.
