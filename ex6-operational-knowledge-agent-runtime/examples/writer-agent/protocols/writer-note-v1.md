# writer.note.v1

Status: frozen

## Record contract

This external package family uses Ex6's canonical Grid package-record carriage defined by `../../../docs/protocols/package-record-v1.md`. Its payload is canonical JSON with required non-empty strings `writer_id` and `note`. It records a writer note. The record does not grant authority; local policy decides whether to interpret or act on it.

## Evolution

This file is immutable. A changed payload meaning requires a new versioned family specification and pCID. Source: DI-jusij, DI-solan.
