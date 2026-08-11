# moks.inventory.count.v1

Status: frozen

## Record contract

This family uses the canonical Grid package-record carriage in `../package-record-v1.md`. Its payload is canonical JSON with required non-empty strings `inventory_id`, `run_id`, `place_id`, `counter`, and `quantity`. It records an inventory count; local policy decides its consequences.

## Evolution

This file is immutable. A changed payload meaning requires a new versioned family specification and pCID. Source: DI-jusij.
