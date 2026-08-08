# Makerspace Stewardship

`ex7-makerspace-stewardship` is a browser-first PromiseGrid example application for a single volunteer makerspace. It models equipment condition, area qualifications, voluntary authority, safety holds, and off-site lending without a checkout master.

In-space use is not a checkout workflow. Any member can record an observation and place a potential safety hold; only a recognized area steward can clear it after recording an inspection. Loanable portable tools have a separate handoff and return record that preserves the terms accepted at checkout.

Run it with:

```bash
go run ./cmd/makerspace-stewardship
```

Then open `http://127.0.0.1:7037/`.

The server writes its append-only local evidence log to
`.makerspace-stewardship/events.jsonl` by default. Use `-runtime-root` to
choose another location.

Read the [operator guide](docs/operator-guide.md),
[workflow and evidence guide](docs/workflow-and-evidence-guide.md), and
[completeness verification guide](docs/completeness-verification.md) before
evaluating the example.
