# Changelog

## Unreleased

### PromiseGrid alignment

- Re-created package-defined durable records as canonical Grid CBOR envelopes with canonical JSON payload bytes and optional semantic author evidence. The unreleased symbolic JSON package-record compatibility path is not retained. Source: DI-sidoh.
- Froze 27 built-in package-family specifications and their fixed CIDv1 pCIDs. The compiled registry is verified against the immutable specification bytes and the published human registry. Source: DI-jusij, DI-solan.
- Defined the extension boundary: workflows compose existing family pCIDs; new shared durable semantics require a new frozen specification and pCID; third-party packages provide their own explicit pCIDs without requiring an Ex6 rebuild. Source: DI-jusij, DI-solan.
- Changed relay batches and external worker/package results to carry exact canonical record bytes, with JSON base64 `[][]byte` only as an outer process-transport wrapper. Source: DI-sidoh.
- Preserved unknown-family records as exact bytes without inferring known semantics; kept relay-carriage signatures distinct from semantic author evidence. Source: DI-sidoh, DI-jusij.

### Documentation and verification

- Added the package-family pCID registry, frozen family specifications, implementation claims, and testing guide.
- Added registry-integrity coverage and verified `go test ./...`, `go vet ./...`, and `errcheck ./...`.
