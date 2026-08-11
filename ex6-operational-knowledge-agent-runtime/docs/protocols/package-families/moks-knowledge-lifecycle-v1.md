# moks.knowledge.lifecycle.v1

Status: frozen

## Record contract

This family uses the canonical Grid package-record carriage in `../package-record-v1.md`. Its payload is canonical JSON with required non-empty strings `item_id` and `status`, plus optional string `notes`. It records an item lifecycle statement; local policy interprets status effects.

## Evolution

This file is immutable. A changed payload meaning requires a new versioned family specification and pCID. Source: DI-jusij.
