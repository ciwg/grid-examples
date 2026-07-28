# EX6 Modular Operational Knowledge System

`ex6-modular-operational-knowledge-system` is the new basket-first runtime for
the modular operational knowledge system effort.

It is not `ex5` version 2.

It is the start of a new standalone product where:

- the runtime owns shared grid-facing infrastructure
- packages own business behavior
- built-in and installed packages use the same basket contract

Current status: this repo contains the runtime foundation plus one small
built-in example package. It does **not** yet contain the full ex5 capability
set as ex6 packages. Source: `DI-moksu`; `DI-lupok`.

## What Exists Today

- standalone Go module
- `moks` CLI
- package manifest loading
- manifest plus self-check activation for installed packages
- built-in package registration
- command routing
- family registration
- explicit per-package protocol implementation claims
- append-only durable history
- CAS storage
- relay batch export and import
- unknown-family exact-byte retention

Source: `DI-moksu`; `DI-lupok`.

## What Does Not Exist Yet

- the full ex5 package set
- browser embodiment
- Neovim embodiment
- ex5 review/search/workflow parity
- full PromiseGrid proof/signature discipline
- richer peer discovery and relay behavior

Source: `DI-moksu`; `DI-lupok`.

## Layout

- `cmd/moks/`
  runtime CLI entrypoint
- `kernel/`
  runtime package activation, routing, family registry, relay export/import
- `packages/`
  package manifest and installed-package runner support
- `records/`
  ex6 durable envelope shape
- `store/`
  append-only history and CAS
- `grid/`
  relay batch types
- `builtin/`
  built-in example egg(s)
- `packages/`
  first-party package source trees and package planning docs
- `templates/package/`
  starter template for outside package authors
- `docs/`
  ex6 reader and author documentation

Source: `DI-moksu`.

## Start Here

- [Architecture](./docs/architecture.md)
- [Current State](./docs/current-state.md)
- [Package Author Guide](./docs/package-author-guide.md)
- [EX5 Capability Map](./docs/ex5-capability-map.md)
