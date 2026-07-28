# EX6 Runnable Examples

These examples are the shortest real paths through the current ex6 system.

They are meant to be copied and run, not just read. Source: `DI-nupad`.

## Built-In Package Walkthrough

From the ex6 repo root:

```bash
go run ./cmd/moks context place create place-1 Receiving Inbound-area
go run ./cmd/moks context resource create res-1 Scale Bench-scale place-1
go run ./cmd/moks procedures create proc-1 DockCheck Check-the-dock
go run ./cmd/moks procedures record-use proc-1 run-1 alice ok followed-v1
go run ./cmd/moks runs evidence add run-1 ev-1 photo kind=image,shift=night blob payload
go run ./cmd/moks runs approve run-1 ap-1 accepted looks-good
go run ./cmd/moks procedures inspect proc-1
go run ./cmd/moks runs inspect run-1
go run ./cmd/moks relay export ./relay.json
```

What this proves:

- the basket persists durable state across separate CLI invocations
- the `procedures` egg composes with `runs`
- relay export sees the same append-only history the eggs produced

Source: `DI-moksu`; `DI-pamuk`; `DI-tusav`; `DI-nupad`.

## Installed-Package Walkthrough

The repo includes a small outside egg example in `examples/writer-egg/`.

From the ex6 repo root:

```bash
go run ./cmd/moks package install ./examples/writer-egg
go run ./cmd/moks writer create writer-1
go run ./cmd/moks relay export ./writer-relay.json
```

What this proves:

- installed executable eggs activate through manifest plus self-check
- installed eggs can request basket-mediated CAS writes and record appends
- relay export includes records created by installed eggs too

Source: `DI-moksu`; `DI-rovum`; `DI-nupad`.

## Shell Scripts

If you want a throwaway workspace instead of writing `.moks/` under the repo
root, run:

```bash
./examples/builtin-quickstart.sh
./examples/install-writer-egg.sh
```

Each script creates its own temporary working directory and prints it before
running the demo. Source: `DI-nupad`.
