# Ex4 Bug Tracker

`ex4-bug-tracker` is a browser-first bug tracker with a small CLI for engineers.
It is a real example application that stands on its own and also shows a
durable, append-only workflow example in this repo.

The app keeps one issue timeline per bug. Creation, comments, assignments,
status changes, and attachments are stored as append-only events, and the
current issue view is projected from that history. Source: `DI-nunit`;
`DI-ninuf`.

## Features

- queue view with status and assignee filters
- issue detail view with a merged timeline
- new issue form
- fixed workflow: `New`, `Triaged`, `In Progress`, `Resolved`
- single active assignee
- reopen from `Resolved` back to `Triaged`
- built-in identities for `reporter`, `triage`, and `engineer`
- real uploaded attachment storage under the local runtime root
- simple engineer CLI for assigned work

## Runtime

By default the server stores local runtime data under `.bug-tracker/`:

- `events.jsonl` for append-only issue history
- `attachments/` for uploaded file copies

Every issue also carries a built-in `team` field. V1 behavior is single-team,
and the current default team is `CORE`, but the storage model already keeps the
field so later multi-team work does not require a data rewrite. Source:
DI-gofub

## Run

Start the server:

```bash
go run ./cmd/bug-tracker
```

Or start the seeded demo:

```bash
bash scripts/run-demo.sh
```

Then open:

```text
http://127.0.0.1:7035/
```

Use the identity picker to switch between the built-in users:

- `reporter`
- `triage`
- `engineer`

## Demo

The seeded demo starts the normal server against a fresh temp runtime root and
preloads four issues:

- one new issue with an attachment
- one issue in progress
- one resolved issue
- one reopened issue back in triage

That gives you a fast way to show:

- queue filtering
- issue detail and timeline
- attachments
- assignment
- status changes
- reopen behavior
- the engineer CLI

The default demo runtime root is:

```text
/tmp/grid-examples-ex4-demo
```

You can override it with `EX4_DATA_ROOT`, and you can override the listen
address with `EX4_LISTEN`.

## CLI

The CLI targets the same server:

```bash
go run ./cmd/bug-tracker-cli --user engineer assigned
go run ./cmd/bug-tracker-cli --user engineer show BUG-0001
go run ./cmd/bug-tracker-cli --user engineer start BUG-0001
go run ./cmd/bug-tracker-cli --user engineer resolve BUG-0001
go run ./cmd/bug-tracker-cli --user engineer comment BUG-0001 "working on a fix"
```

## Notes

- The server uses simple built-in identities and role checks for this first
  slice.
- The app is durable-first and does not depend on websocket sessions.
- Attachments are copied into the runtime root instead of referencing original
  host paths.
- `--seed-demo` only seeds when the selected runtime root is empty. Source:
  `DI-zogof`

## Docs

- [Architecture notes](docs/architecture.md)
- [Practical implementation notes](docs/practical-implementation.md)
- [Browser UI example](docs/bug-tracker-ui-example.md)
- [In-progress features](INPROGRESS-FEATURES.md)
