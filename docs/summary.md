# Grid Exercises Demo Guide

This is a practical, print-friendly guide for demonstrating each exercise.
Run every command from the exercise folder shown in its section. Each exercise
has its own `Makefile`; begin with `make help` if you want to see the available
paths again.

## Before you start

Most exercises require Go. The Docker demonstrations require a working Docker
daemon. The browser demonstrations need a modern browser; Ex7's two-agent proof
specifically uses Google Chrome. Start one exercise at a time: several use
nearby local ports.

## Ex1 — Order flow

**Folder:** `ex1-order-flow/`
**Needs:** Docker and Docker Compose.

This is the best container-based demo. It starts the collector, kernel, seller,
warehouse, accounting, carrier, and intake roles; then it runs an analyzer over
the resulting artifacts.

```bash
cd ex1-order-flow
make demo
```

Wait for the analyzer result and the final `preserved run data:` line. The demo
stops its containers itself and leaves its evidence under
`/tmp/grid-examples-ex1-data`. To show exception behavior, run either:

```bash
make demo-refusal
make demo-timeout
```

Use the preserved role directories to explain that each role keeps its own
messages and observations.

## Ex2 — Grid Editor

**Folder:** `ex2-grid-editor/`
**Needs:** Go and a browser. Docker Compose is optional.

For the fastest single-relay demo, start the relay in one terminal:

```bash
cd ex2-grid-editor
make relay
```

Open `http://127.0.0.1:7015/?doc=demo` in a browser. Edit the document and use
the collaboration UI. Stop the relay with `Ctrl-C` when finished.

To explain the two-relay topology, leave the first relay running and start this
in a second terminal:

```bash
cd ex2-grid-editor
make second-relay
```

Then open `http://127.0.0.1:7016/?doc=demo`. The two relays use separate local
data roots and exchange data as peers. If Docker is preferable, use
`make docker-demo` instead; it starts the two-relay simulation in the foreground.

## Ex3 — WebSocket Grid Editor

**Folder:** `ex3-grid-editor-websocket/`
**Needs:** Go and a browser. Docker Compose is optional.

For the normal local demonstration, run:

```bash
cd ex3-grid-editor-websocket
make relay
```

Open `http://127.0.0.1:7025/?doc=demo`. This is the WebSocket variant of the
Grid Editor; the browser gives the quickest way to see live document behavior.
Stop the relay with `Ctrl-C`.

For its ready-made two-relay container setup, run:

```bash
make docker-up
```

Open either `http://127.0.0.1:7025/?doc=demo&access_token=ex3-demo-access` or
`http://127.0.0.1:7026/?doc=demo&access_token=ex3-demo-access`. When finished,
run `make docker-down`. The `remote-relay` target is the documented remote-access
demo; it uses a checked-in demonstration-only access token and should not be
treated as deployment configuration.

## Ex4 — Bug Tracker

**Folder:** `ex4-bug-tracker/`
**Needs:** Go and a browser.

Use the seeded demo rather than starting empty:

```bash
cd ex4-bug-tracker
make demo
```

Open `http://127.0.0.1:7035/`. The sample data demonstrates a new issue, work
in progress, a resolved issue, a reopened issue, an attachment, and the three
built-in roles. Change the identity picker between `reporter`, `triage`, and
`engineer` to show different workflow actions. Stop the server with `Ctrl-C`.

The normal empty-server path is `make serve`; the full automated verification,
including browser-signing coverage, is `make verify`.

## Ex5 — Operational Knowledge System

**Folder:** `ex5-operational-knowledge-system/`
**Needs:** Go and a browser. Google Chrome is needed only for the extension demo.

Use the newcomer data set for the best first walkthrough:

```bash
cd ex5-operational-knowledge-system
make newcomer
```

Open `http://127.0.0.1:7045/`. The loaded runtime includes receiving and
inventory problems, training history, a maintenance draft, an attachment, and a
walkthrough. The server continues to run until `Ctrl-C`.

`make newcomer` uses `/tmp/ex5-newcomer-runtime` by default. To use another
empty location, run `make newcomer RUNTIME_ROOT=/path/to/runtime`. The full
verified Chrome extension demonstration is available as `make browser-demo`;
use it only on a machine prepared for that Chrome/native-host path.

## Ex6 — Operational Knowledge Agent Runtime (OKAR)

**Folder:** `ex6-operational-knowledge-agent-runtime/`
**Needs:** Go.

This is a command-line runtime, not a web server. Its shortest demonstration
uses a throwaway temporary workspace:

```bash
cd ex6-operational-knowledge-agent-runtime
make quickstart
```

The script creates a place and resource, creates and uses a procedure, adds
evidence and approval, inspects the results, and exports relay data. It prints
the temporary workspace and relay-export path for inspection.

To demonstrate installation of an outside package, run `make writer-example`.
The Docker and two-runtime integration targets are opt-in proofs and are not
needed for the basic walkthrough.

## Ex7 — Makerspace Stewardship

**Folder:** `ex7-makerspace-stewardship/`
**Needs:** Go and a browser. Google Chrome is needed for the two-agent proof.

Start the local record server:

```bash
cd ex7-makerspace-stewardship
make serve
```

Open `http://127.0.0.1:7037/`. This server displays local state and accepts
already signed records; it does not make the browser an author of records.
Stop it with `Ctrl-C`.

For a repeatable non-interactive record demonstration, run:

```bash
make ingress-proof
```

For the full two-agent Chrome proof, run `make browser-proof`. It starts its own
temporary Alice and Bob agents and prints the proof directory containing the
final screenshot, approval response, and logs. It requires local Chrome and
DevTools support.

## Quick chooser

- Want a complete container story: **Ex1**.
- Want a visible collaborative editor: **Ex2** or **Ex3**.
- Want a ready-to-show application: **Ex4** or **Ex5**.
- Want a CLI/runtime walkthrough: **Ex6**.
- Want signed record recognition and a two-agent proof: **Ex7**.
