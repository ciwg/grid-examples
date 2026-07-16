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

