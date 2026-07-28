# EX6 Package Author Guide

This guide explains how a new package fits into ex6.

Current rule: package authors extend the runtime through explicit declarations
and runtime contracts. They do not bypass the runtime and make private package
systems the source of truth. Source: `DI-moksu`; `DI-lupok`.

## Where Packages Live

There are three important locations:

- source packages in the repo:
  `ex6-modular-operational-knowledge-system/packages/<package-id>/`
- package template:
  `ex6-modular-operational-knowledge-system/templates/package/`
- runtime-installed copies:
  `ex6-modular-operational-knowledge-system/.moks/packages/<package-id>/`

The starter template now exists under:

- `templates/package/moks-package.json`
- `templates/package/package.sh`
- `templates/package/README.md`
- installed-package example:
  `examples/writer-agent/`

Source: `DI-moksu`; `DI-nupad`.

## What A Package Must Declare

A package manifest currently declares:

- package `id`
- package `version`
- optional `description`
- installed package `executable`
- command registrations
- family registrations
- protocol implementation claims

Every declared family must be covered by an explicit protocol implementation
claim. Source: `DI-lupok`.

## Current Activation Model

Installed packages are executable packages.

Current activation flow:

1. the runtime reads `moks-package.json`
2. the runtime validates the manifest
3. the runtime runs the package `describe` self-check
4. the runtime compares the described package shape to the manifest
5. if they match, the package is activated

Source: `DI-moksu`; `DI-lupok`.

## Minimal Installed Package Shape

Current installed packages should support:

- `describe`
  print the package manifest JSON used for self-check
- `validate`
  read one record from stdin and validate it if the package owns that family
- `run <command-key> ...`
  execute a registered command and print either plain text output or a JSON
  result describing runtime-mediated CAS writes and record appends

Source: `DI-moksu`.

## Minimal Manifest Example

```json
{
  "id": "helper-agent",
  "version": "0.1.0",
  "description": "Example package",
  "executable": "./helper-agent.sh",
  "commands": [
    {
      "path": ["helper", "echo"],
      "summary": "Echo a string"
    }
  ],
  "families": [
    {
      "name": "helper.echo.v1",
      "protocol_pcid": "pcid:helper.echo.v1"
    }
  ],
  "claims": [
    {
      "protocol_pcid": "pcid:helper.echo.v1",
      "role": "family-validator",
      "summary": "Validates helper echo envelopes."
    }
  ]
}
```

## Installing A Package

Current command:

```bash
moks package install ./path-to-package
```

That copies the package into the runtime package root and activates it through
the manifest plus self-check path. Source: `DI-moksu`.

## Runtime-Mediated Writes

Installed packages do not write the runtime state directly.

Instead, `run` may print a JSON result with:

- `output`
- `cas`
- `records`

The runtime performs those CAS and append-only history writes itself. This
keeps durable mutation runtime-mediated even for external packages. Source:
`DI-rovum`.

## Runnable Installed-Package Example

The repo now includes one small installed-package example:

- `examples/writer-agent/moks-package.json`
- `examples/writer-agent/writer-agent.sh`

Try it from the ex6 repo root:

```bash
go run ./cmd/moks package install ./examples/writer-agent
go run ./cmd/moks writer create writer-1
go run ./cmd/moks relay export ./writer-batch.json
```

That flow proves an outside package can install, request runtime-mediated CAS
writes, append a durable record, and then have that record show up in relay
export. Source: `DI-rovum`; `DI-nupad`.

## Starter Flow

Current recommended author flow:

1. copy `templates/package/` to a new directory
2. rename package IDs, family names, and command paths
3. update claims so every family protocol is covered
4. implement `describe`, `validate`, and `run`
5. install with `moks package install ./your-package`

Source: `DI-moksu`; `DI-lupok`.

## Current Limit

This is still a foundation contract, not a mature package SDK.

There is not yet:

- a package scaffolding command
- a richer package API
- a full PromiseGrid proof/signature integration story

Source: `DI-moksu`; `DI-lupok`.
