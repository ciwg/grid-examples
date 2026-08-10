# ex1-order-flow

`ex1-order-flow` is the order-fulfillment example in this repo.
It demonstrates several independent agents exchanging signed messages through a
kernel, with one container per role in the easiest demo path.

## What You Need To Run

For the standard demo:

- Docker
- either the `docker compose` plugin or `docker-compose`
- a local shell that can run `bash`

For local development and tests without Docker:

- Go

You do not need a browser, Node, npm, or Neovim for this example.

## Quick Demo

The easiest way to run `ex1` is the Docker demo harness:

```bash
cd docker
bash run-demo.sh
```

That script:

- builds the demo images
- starts the collector, kernel, seller, warehouse, accounting, and carrier
  services
- runs one short-lived intake container with a fixture
- runs the analyzer over the collected artifacts
- shuts the containers down after the run

## Alternate Fixtures

The default fixture is `happy-path.json`. You can run another fixture by
passing its filename:

```bash
cd docker
bash run-demo.sh warehouse-refusal.json
```

Fixtures live under [fixtures/](fixtures/).

## Runtime Data

The default demo runtime root is:

```text
/tmp/grid-examples-ex1-data
```

You can override it with `EX1_DATA_ROOT` before running the demo script.

The script starts from an empty runtime root on each run, then leaves the final
artifacts on disk for inspection.

## Scope, Guided Demo, and Evidence

Ex1 implements five named **local draft** profiles. It does not claim frozen
PromiseGrid-spec conformance or independent-peer interoperability; see the
[implementation scope declaration](CHANGELOG.md) and the
[protocol inventory](docs/design.md#protocol-inventory-current-local-draft-profiles).
Source: DI-josir, DI-motiv.

From this directory, run the normal local demo:

```bash
bash docker/run-demo.sh
```

The happy path prints the analyzer result and preserves its artifacts under
`/tmp/grid-examples-ex1-data` unless `EX1_DATA_ROOT` selects a different root.
To inspect the two exceptional outcomes covered by the guide:

```bash
bash docker/run-demo.sh warehouse-refusal.json
bash docker/run-demo.sh carrier-timeout.json
```

For each completed run, inspect these host-local artifacts:

- `collector/message-dag.jsonl` and `collector/message-cas/` show the
  collector's retained sent-envelope graph and raw bytes.
- `<role>/message-cas/` and `<role>/messages.jsonl` show raw envelopes and
  message records retained by that individual role.
- `<role>/observations.jsonl` records that role's exceptional local evidence;
  for example, the seller records `refusal_observed` for the signed warehouse
  refusal and `timeout_observed` when its carrier deadline expires.
- `kernel/observations.jsonl` records kernel-local ingress and dispatch facts,
  including `no_registered_recipient`; it does not decide whether a pCID is
  globally valid or a business promise was kept.

These records are local evidence, not shared proof. A signed refusal remains
the responder's artifact, while a timeout is only the observing role's record;
neither lets the reader infer another agent's intent. Source: DI-vihoz,
DI-riguz, DI-purum, DI-zosiz, DI-motiv.

If `run-demo.sh` exits nonzero, do not treat the resulting files as a completed
demo claim. Select or reuse an `EX1_DATA_ROOT`, rerun the script from its fresh
root, and inspect only artifacts from the successful completed run. Source:
DI-rokol, DI-motiv.

## Direct Local Development

If you want to work on the binaries directly instead of using Docker, the
entrypoints live under [cmd/](cmd/):

- `pg-order-collector`
- `pg-order-kernel`
- `pg-order-seller`
- `pg-order-warehouse`
- `pg-order-accounting`
- `pg-order-carrier`
- `pg-order-intake`
- `pg-order-analyze`

Go is enough for that workflow:

```bash
go test ./...
```

The Docker path remains the intended quick demo because it brings up the full
multi-agent topology with the least manual setup.

## Docs

- [Design notes](docs/design.md)
- [Current local protocol inventory](docs/design.md#protocol-inventory-current-local-draft-profiles)
- [Implementation scope declaration](CHANGELOG.md)
- [Testing guide](docs/testing.md)
