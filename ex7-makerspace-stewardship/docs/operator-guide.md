# Operator guide

Run an isolated demo:

```bash
go run ./cmd/makerspace-stewardship -runtime-root /tmp/ex7-makerspace-demo
```

Then open `http://127.0.0.1:7037/`. The evidence log is
`/tmp/ex7-makerspace-demo/events.jsonl`. Stop and restart the same command to
confirm replay; inspect the log with `sed -n '1,20p' /tmp/ex7-makerspace-demo/events.jsonl`. Do not edit a corrupt log: the service
fails closed and leaves it for investigation. For a fresh demo, choose a new
empty runtime root. Source: DI-dapod; DI-damod.
