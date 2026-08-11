# OKAR Architecture

`OKAR` is a runtime-first system.

That means the runtime is real:

- packages register with the runtime
- packages declare families and protocol claims
- packages route commands through the runtime
- relay export/import happens through the runtime
- unknown future families can still be stored and relayed as exact bytes

Packages are not supposed to become their own independent little systems with
private source-of-truth channels. Source: `DI-moksu`; `DI-lupok`.

## Topology

```text
       frozen built-in pCID registry / external package pCIDs
                               │
                      package manifests + self-check
                               │
                    OKAR runtime kernel
          /                 /          \                 \
         /                 /            \                 \
  command routing  canonical history    CAS       relay export/import
         \                 \            /                 /
          \                 \          /                 /
                  built-in and installed packages
```

Source: `DI-moksu`; `DI-lupok`.

## Core Responsibilities

The runtime currently owns:

- package activation
- manifest validation
- installed-package self-check
- command routing
- family registration
- protocol claim publication
- fixed lookup for the 27 built-in family pCIDs
- canonical Grid record encoding and append-only durable history
- CAS storage
- relay batch export/import

Source: `DI-moksu`; `DI-lupok`.

## Package Responsibilities

A package currently owns:

- its manifest
- its protocol implementation claims
- its durable families
- its command surface
- validation for the families it owns

An installed third-party package also owns the immutable specifications for its
own external families. It provides explicit pCIDs; the runtime never computes a
new pCID from a family label. Workflows compose package capabilities and family
pCIDs, so they are not package protocols by default. Source: `DI-jusij`;
`DI-solan`.

Built-in and installed packages follow the same conceptual contract. Source:
`DI-moksu`; `DI-lupok`.

## Protocol-First Rule

The current PromiseGrid alignment rule in ex6 is:

- each built-in family is registered with the fixed CIDv1 pCID of its immutable
  specification; the [registry](protocols/package-family-pcid-registry.md)
  publishes the mapping
- an external family supplies its own explicit frozen pCID in its package
  manifest and record bytes
- a package must publish an explicit implementation claim for that protocol
- a known-family record is rejected if its `protocol_pcid` does not match the
  registered family contract
- an unknown-family record is still retained as exact bytes

This keeps ex6 closer to protocol work than to app-local RPC work. Source:
`DI-lupok`; `DI-jusij`; `DI-solan`.

## Evidence and storage map

| Evidence or state | Representation | Boundary |
| --- | --- | --- |
| Durable package record | canonical Grid CBOR bytes; canonical JSON payload slot | semantic family pCID and local validator/policy |
| Semantic author evidence | author key ID, public key, and signature slots | required by Ex6 local append/import policy; evidence of authorship, not automatic authority |
| Append-only history | base64 line framing of exact canonical record bytes | local durable storage |
| CAS and blobs | content-addressed local objects | local runtime storage |
| Relay batch | exact record bytes, batch metadata, relay signatures and digest proofs | carriage evidence, separate from author evidence |
| Route promises and bindings | local durable promise/binding records | local availability and route-selection policy |
| Workflow artifacts and projections | immutable CAS artifacts plus local lifecycle/projection state | local import, activation, and execution decisions |

Unknown pCID records remain exact bytes in history, CAS, and relay carriage
without being treated as known semantics. Role effects, approvals, and route
selection are local-policy interpretations, not protocol-granted authority.
Ex6 also applies its [strict semantic author-signature admission policy](author-signature-admission-policy.md)
without changing the frozen family pCIDs. Source: `DI-sidoh`; `DI-jusij`;
`DI-solan`; `DI-damab`.

## Current Limits

This is still an early foundation.

It is not yet:

- a complete ex5 replacement
- a complete PromiseGrid implementation
- a full package ecosystem

The current code should be read as kernel groundwork, not as a finished modular
product. Source: `DI-moksu`; `DI-lupok`.
