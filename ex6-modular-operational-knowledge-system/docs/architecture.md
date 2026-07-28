# EX6 Architecture

`ex6` is a basket-first runtime.

That means the basket is real:

- packages register with the runtime
- packages declare families and protocol claims
- packages route commands through the runtime
- relay export/import happens through the runtime
- unknown future families can still be stored and relayed as exact bytes

Packages are not supposed to become their own independent little systems with
private source-of-truth channels. Source: `DI-moksu`; `DI-lupok`.

## Topology

```text
             package protocol claims / family declarations
                               │
                      package manifests + self-check
                               │
                    ex6 runtime kernel ("the basket")
          /                 /          \                 \
         /                 /            \                 \
  command routing   append-only log      CAS       relay export/import
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
- append-only durable history
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

Built-in and installed packages follow the same conceptual contract. Source:
`DI-moksu`; `DI-lupok`.

## Protocol-First Rule

The current PromiseGrid alignment rule in ex6 is:

- a family is registered with a declared `protocol_pcid`
- a package must publish an explicit implementation claim for that protocol
- a known-family record is rejected if its `protocol_pcid` does not match the
  registered family contract
- an unknown-family record is still retained as exact bytes

This keeps ex6 closer to protocol work than to app-local RPC work. Source:
`DI-lupok`.

## Current Limits

This is still an early foundation.

It is not yet:

- a complete ex5 replacement
- a complete PromiseGrid implementation
- a full package ecosystem

The current code should be read as kernel groundwork, not as a finished modular
product. Source: `DI-moksu`; `DI-lupok`.
