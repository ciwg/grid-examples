# OKAR Package Author Guide

This guide explains how a new package fits into OKAR.

Current rule: package authors extend the runtime through explicit declarations
and runtime contracts. They do not bypass the runtime and make private package
systems the source of truth. Source: `DI-moksu`; `DI-lupok`.

## Where Packages Live

There are three important locations:

- source packages in the repo:
  `ex6-operational-knowledge-agent-runtime/packages/<package-id>/`
- package template:
  `ex6-operational-knowledge-agent-runtime/templates/package/`
- runtime-installed copies:
  `ex6-operational-knowledge-agent-runtime/.moks/packages/<package-id>/`

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
claim.

Current routing rule:

- those claims are also the runtime's intermediate route registrations
- if a package declares a family, it must declare a matching
  `family-validator` claim for that family's `protocol_pcid`
- relay export carries the derived route registrations so peers can inspect
  routing roles from batch metadata too
- parser or transform claims must declare `route_type` plus `emits_protocols`
- plain direct handlers can omit `route_type` and default to `direct`
- local consumers can query one input protocol with
  `moks route inspect <protocol-pcid>`
- local consumers can ask for the kernel's preferred executable route plan with
  `moks route plan <protocol-pcid>`
- operators can tune route-plan ordering with
  `moks route policy show` and `moks route policy set ...`

Source: `DI-lupok`; `DI-rutom`; `DI-ruvot`; `DI-lafek`; `DI-fotav`; `DI-pabut`; `DI-matek`.

## Workflow-Adapter Extension

An executable workflow adapter is declared in the active package's
`moks-package.json`, rather than in a separate adapter file. The declaration
names the adapter, Docker image/command, exact input/output pCIDs, and required
CPU, memory, PID, and timeout limits.

The first package using this contract is
[`examples/procedure-execution-adapter/`](../examples/procedure-execution-adapter/).
A workflow artifact may use an adapter only when its active manifest and the
active package declaration agree exactly on the adapter name and input/output
pCIDs.

An adapter is an executable agent, not an ordinary package command. The
runtime sends exact CBOR input to a Docker-confined worker through stdin.
The worker may return typed CBOR plus proposed CAS and record writes, but it
cannot write runtime state itself. The runtime validates the declared output
pCID before applying any proposed durable writes. The worker result uses the
frozen [`workflow-adapter-result-v1.json`](./protocols/workflow-adapter-result-v1.json)
protocol and contains the output CBOR plus explicit CAS/record proposals.
The runtime accepts at most 16 MiB on either worker stream.

During this first confined-worker slice, a proposal may write an unknown family
or a runtime/built-in validated family. A proposal for an externally validated
installed-package family is rejected: validating it would otherwise start that
package directly on the host. A future confined validator protocol can extend
this safely. Source: `DI-fofuh`; `TE-dovek`.

Workers will have no runtime-root, CAS, history, peer-key, Docker-socket, or
ambient-secret access; no network; and no direct host-process fallback. The
existing package `run` command remains its current separate contract and does
not gain workflow-adapter authority merely by being installed. Source:
`DI-fofuh`; `TE-dovek`.

The repository now ships the source for a locally built procedure-execution
image; its package declaration uses the immutable local Docker image ID, not
the convenient mutable build tag. It is not registry-published.
Existing built-in adapters remain available; an installed package becomes an
adapter supplier only after its manifest and `describe` self-check agree.

Build and install the first adapter from the ex6 root:

```bash
docker build -f examples/procedure-execution-adapter/Dockerfile -t moks/procedure-execution-adapter:dev .
go run ./cmd/moks package install ./examples/procedure-execution-adapter
```

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

Packages that declare `workflow_adapters` additionally supply the Docker image
and command named by that declaration. The worker itself receives CBOR on
standard input and writes one CBOR adapter-result envelope to standard output;
it does not implement a host-process `run` fallback. Source: `DI-fofuh`;
`TE-dovek`.

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
