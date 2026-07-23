# Ex5 User Guide

This is the primary operator guide for `ex5-operational-knowledge-system`. It
is written for a newcomer who wants to understand what the system is for, load
real sample data, and move through the current browser, CLI, and Neovim
surfaces without guessing. It describes the shipped product surface, not future
work. Source: `DI-movar`; `DI-rubav`.

## What Ex5 Is

`ex5` is a local operational memory system. It keeps operational context,
versioned knowledge, runs, evidence, approvals, links, and live drafts in one
runtime so a team can answer practical questions later:

- what revision existed
- what run used it
- what evidence was captured
- who reviewed it
- what place, resource, or responsibility it belonged to

The same model covers procedures, receiving checks, inventory audits, training,
and maintenance because the core problem is durable operational memory, not one
workflow-specific form. Source: `DI-radok`; `DI-vemok`; `DI-fudok`.

## Start Here

The fastest honest newcomer path is:

1. Load the checked-in sample corpus into a fresh runtime root.
2. Start the ex5 server against that runtime root.
3. Use the CLI to inspect the sample world first.
4. Use the browser if you have the shipped Chrome/Chromium embodiment set up.
5. Use Neovim if you want live draft editing inside the editor.

The sample corpus is not mock data. It is a checked-in runtime with real
append-only history, a real attachment, and one persisted live draft. Source:
`DI-rubav`.

## Core Concepts

- Responsibility: who owns or reviews a slice of work.
- Place: where work happens.
- Resource: which machine, tool, bin, station, or document is involved.
- Knowledge item: a versioned operational record such as a receiving check,
  inventory audit, training item, or maintenance item.
- Run: a performed event against an exact item revision.
- Evidence: facts and an optional immutable attachment captured for a run.
- Approval: a named review record for an item or run.
- Link: a typed connection between records.
- Live draft: the current collaborative working body for one item, separate
  from durable revision history.

Source: `DI-radok`; `DI-zuvob`; `DI-luzaf`; `DI-zanub`; `DI-zoruk`.

## Load The Sample Corpus

Choose a fresh runtime root and copy the checked-in newcomer corpus into it:

```bash
./scripts/load-sample-data.sh /tmp/ex5-newcomer-runtime
```

The loader fails closed if the target already contains data. That is
intentional: the newcomer corpus should load into a fresh root instead of
silently mutating an existing runtime. Source: `DI-rubav`.

## Start The Runtime

Run the server against the loaded sample root:

```bash
go run ./cmd/operational-knowledge -data-root /tmp/ex5-newcomer-runtime
```

The server exposes:

- the browser shell on `http://127.0.0.1:7045/`
- the direct local Unix socket at the runtime root
- the same projected state to browser, CLI, and Neovim

Source: `DI-zorav`; `DI-favel`; `DI-fudok`.

## Know The Sample World

The sample corpus keeps all four storylines in one shared runtime:

- Receiving: inbound pallet inspection with a short-count and torn-wrap problem
- Inventory: RJ45 bin cycle count with a positive variance
- Training: new receiver onboarding and signoff
- Maintenance: daily heat sealer check with one active draft follow-up

The main sample records are:

- Responsibilities:
  - `RESP-0001` Receiving lead
  - `RESP-0002` Inventory steward
  - `RESP-0003` Training coordinator
  - `RESP-0004` Maintenance lead
- Places:
  - `PLACE-0001` North Warehouse
  - `PLACE-0002` Receiving Dock A
  - `PLACE-0003` Stock Cage 3
  - `PLACE-0004` Training Bay
  - `PLACE-0005` Line A
- Resources:
  - `RES-0001` Inbound Pallet Camera
  - `RES-0002` RJ45 Bin
  - `RES-0003` Training Binder
  - `RES-0004` Heat Sealer 7
- Knowledge items:
  - `RECV-0001` Inspect inbound connector pallet
  - `INV-0001` Count RJ45 bin
  - `TRAIN-0001` New receiver onboarding checkoff
  - `MAINT-0001` Daily heat sealer check
- Runs:
  - `RUN-0001` receiving problem thread
  - `RUN-0002` inventory discrepancy thread
  - `RUN-0003` training completion
  - `RUN-0004` maintenance startup check

The corpus also includes:

- one real evidence attachment on `RUN-0001`
- one persisted live draft on `MAINT-0001`
- one draft-status item in the pending-review queue
- two problem hotspots visible in problem review

Source: `DI-rubav`.

## First CLI Tour

Start with the high-level projections:

```bash
go run ./cmd/oks-cli -socket /tmp/ex5-newcomer-runtime/embodiment.sock dashboard
go run ./cmd/oks-cli -socket /tmp/ex5-newcomer-runtime/embodiment.sock pending-review
go run ./cmd/oks-cli -socket /tmp/ex5-newcomer-runtime/embodiment.sock problem-review
```

You should see:

- one draft item: `MAINT-0001`
- zero unreviewed runs
- two problem runs: `RUN-0001` and `RUN-0002`

The CLI now treats the local Unix socket as the primary terminal contract. Use
`-socket=off` only when you intentionally want the HTTP compatibility path.
Source: `DI-zorav`; `DI-monuv`; `DI-rubav`.

## End-To-End Walkthrough

This walkthrough shows how the sample world fits together.

### 1. Inspect The Draft Item

```bash
go run ./cmd/oks-cli -socket /tmp/ex5-newcomer-runtime/embodiment.sock show-item MAINT-0001
```

What to notice:

- `MAINT-0001` is still `draft`
- it has two durable revisions
- it already has one completed maintenance run
- it also has a newer persisted live draft body that has not been snapped yet

This is the clearest example of how ex5 separates durable revision history from
current working draft state. Source: `DI-zoruk`; `DI-lusov`; `DI-rubav`.

### 2. Review The Receiving Problem

```bash
go run ./cmd/oks-cli -socket /tmp/ex5-newcomer-runtime/embodiment.sock show-run RUN-0001
```

What to notice:

- the run points to `RECV-0001` revision `1`
- the outcome is `accepted_with_notes`
- the evidence facts show `expected_count=240` and `actual_count=238`
- the attachment `inbound-wrap-photo.txt` is present
- the run already carries a review note

This is the shortest path to understanding why ex5 stores runs, evidence, and
approvals separately instead of flattening them into one document blob. Source:
`DI-zuvob`; `DI-zanub`; `DI-rubav`.

### 3. Review The Inventory Problem

```bash
go run ./cmd/oks-cli -socket /tmp/ex5-newcomer-runtime/embodiment.sock show-run RUN-0002
go run ./cmd/oks-cli -socket /tmp/ex5-newcomer-runtime/embodiment.sock search connector problem=true
```

What to notice:

- the inventory run is `completed`, but it is still problem-worthy because the
  evidence facts carry a non-zero variance and mismatched counts
- `problem=true` is not just about free-text notes; it follows the shared
  problem review logic

Source: `DI-pogul`; `DI-vafuk`; `DI-ralek`; `DI-rubav`.

### 4. Trace Ownership And Context

```bash
go run ./cmd/oks-cli -socket /tmp/ex5-newcomer-runtime/embodiment.sock show-responsibility RESP-0001
go run ./cmd/oks-cli -socket /tmp/ex5-newcomer-runtime/embodiment.sock show-place PLACE-0002
go run ./cmd/oks-cli -socket /tmp/ex5-newcomer-runtime/embodiment.sock show-resource RES-0002
```

What to notice:

- responsibilities own or review real items and runs
- places and resources are not decorative metadata; they are part of the
  searchable operational graph
- typed links connect the storylines instead of leaving context implicit

Source: `DI-kovup`; `DI-luzaf`; `DI-salup`; `DI-rubav`.

### 5. Compare Training And Maintenance

```bash
go run ./cmd/oks-cli -socket /tmp/ex5-newcomer-runtime/embodiment.sock show-item TRAIN-0001
go run ./cmd/oks-cli -socket /tmp/ex5-newcomer-runtime/embodiment.sock show-run RUN-0003
go run ./cmd/oks-cli -socket /tmp/ex5-newcomer-runtime/embodiment.sock show-run RUN-0004
```

What to notice:

- training and maintenance use the same durable model as receiving and
  inventory
- the difference is in the item kind, run kind, evidence, approvals, and links
  around them, not in a different storage system

Source: `DI-vemok`; `DI-kovup`; `DI-rubav`.

## Browser Path

The browser is the broadest embodiment. For the demo path below, do not
assemble it manually. Use the exact setup, launch, and verify steps here first,
then use the live-demo sheet. Unsupported browsers do not silently fall back to
the older HTTP browser path. Source: `DI-punek`; `DI-fonuv`; `DI-fovek`;
`DI-dabek`.

Run these exact steps:

```bash
./scripts/setup-demo-browser.sh
./scripts/launch-demo-browser.sh
./scripts/verify-demo-browser.sh
```

What those do:

- create a disposable demo runtime under `/tmp/ex5-demo-browser/runtime`
- build and register the native host for `operational_browser_host`
- launch Chrome with the shipped unpacked extension already loaded
- fail closed unless the runtime, host, and registration path are actually
  ready

Use the live-demo sheet only after `verify-demo-browser.sh` prints:

```text
ex5 browser demo verification passed
```

After verification passes, the newcomer sample is easiest to read in this
order:

1. Review workspace draft queue: find `MAINT-0001`
2. Review workspace hotspots: inspect the receiving and inventory problems
3. Known-record search: search for `connector`, `RJ45`, `training`, or
   `sealer`
4. Current Record: open one item or run and follow the next-step actions
5. Author workspace: inspect the live draft for `MAINT-0001`
6. Operate workspace: review how run logging, evidence, and approvals attach
   to the current record

### Live Demo Sheet

Use this when you need a short live intro for people who have never seen `ex5`.
It stays on the same sample world and browser order described above. Source:
`DI-luren`; `DI-dabek`.

Opening:

- `I built this because our operational truth gets split across too many places.`
- `We end up with the procedure in one tool, the work record in another, evidence somewhere else, review in chat, and the real context in people's heads.`
- `This keeps the procedure, the run, the evidence, the review, and the context together.`

Browser flow:

- Start in Review and open the draft queue.
- Open `MAINT-0001` and call out that live draft state is separate from durable
  revision history.
- Switch to hotspots and open `RUN-0001` to show one receiving problem with its
  run, evidence, review, and context together.
- Open `RUN-0002` to show that problem review comes from real recorded work,
  not just free-text notes.
- Point at Current Record and the next-step actions to show that review and
  action stay connected instead of bouncing across tools.

CLI proof:

```bash
go run ./cmd/oks-cli -socket /tmp/ex5-newcomer-runtime/embodiment.sock pending-review
```

Close with:

- `Same runtime. Same operational truth. Different embodiment.`

If `verify-demo-browser.sh` does not pass, do not use this browser sheet. Fall
back to the CLI tour earlier in this same guide instead. Source: `DI-dabek`.

For field-by-field browser behavior, use the
[Browser UI Guide](./browser-ui-guide.md). Source: `DI-nalor`; `DI-rubav`.

## Neovim Path

Neovim is the best newcomer path when you want to inspect and edit the active
draft without setting up the browser embodiment first.

Open the maintenance draft:

```bash
OKS_BASE_URL=http://127.0.0.1:7045 ./scripts/oks-nvim MAINT-0001
```

Useful commands inside Neovim:

- `:OksInfo`
- `:OksInspect`
- `:OksPending`
- `:OksProblemReview`
- `:OksSnapshot`
- `:OksApproveItem`
- `:OksApproveRun`

Neovim now prefers the direct local Unix-socket contract and keeps websocket or
HTTP only as explicit compatibility mode. Source: `DI-fonuv`; `DI-monuv`;
`DI-rubav`.

## Common Real Tasks After The Tour

Once you understand the sample world, these are the next useful things to try
in your own copied runtime:

- approve or supersede `MAINT-0001`
- snapshot the active maintenance draft into a new durable revision
- record a second receiving or inventory run and compare the new history
- add another typed link from a responsibility to a run or item
- search by `problem=true`, by `place_id`, or by `responsibility_id`
- inspect how the draft file, attachment files, and CAS objects changed on disk

Because you loaded the sample into your own target root, you can mutate that
copy freely without changing the checked-in corpus. Source: `DI-rubav`.

## Troubleshooting

### The loader refuses my target root

The loader intentionally fails if the target already contains files. Pick a new
empty directory or remove the old scratch runtime first. Source: `DI-rubav`.

### The CLI says the local socket is unavailable

Make sure:

- the server is running against the same data root you loaded
- you pointed `-socket` at `<data-root>/embodiment.sock`
- you only use `-socket=off` when you explicitly want HTTP compatibility

Source: `DI-zorav`.

### The browser says Chrome/Chromium is required

Run the three browser-demo commands in this guide exactly:

```bash
./scripts/setup-demo-browser.sh
./scripts/launch-demo-browser.sh
./scripts/verify-demo-browser.sh
```

If verification still does not pass, use the CLI or Neovim path instead.
Source: `DI-punek`; `DI-fovek`; `DI-dabek`.

### The launch script says `127.0.0.1:7045` is not serving the demo runtime

That means some other local server is already bound to `127.0.0.1:7045`, and
it is not the disposable demo runtime under `/tmp/ex5-demo-browser/runtime`.
Stop that other local server first, then rerun:

```bash
./scripts/launch-demo-browser.sh
./scripts/verify-demo-browser.sh
```

Do not use the browser live-demo sheet until verification passes. Source:
`DI-dabek`.

### Pending review or problem review looks empty

Use the checked-in newcomer corpus first. It already contains:

- one draft item
- two problem runs
- four approved/noted run reviews

If you mutate your copy heavily, the projections will change with it. Source:
`DI-rubav`.

## Where To Read Next

- [Product Overview](./product-overview.md)
- [Browser UI Guide](./browser-ui-guide.md)
- [Terminal Capability Matrix](./terminal-capability-matrix.md)
- [Features Guide](./features-guide.md)
- [HTTP API Guide](./http-api-guide.md)
- [Architecture](./architecture.md)
