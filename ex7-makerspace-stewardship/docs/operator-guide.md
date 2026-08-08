# Operator guide

Run `go run ./cmd/makerspace-stewardship`, then open `http://127.0.0.1:7037/`.
The local evidence log is `.makerspace-stewardship/events.jsonl`; use
`-runtime-root` for an isolated demo location. Do not edit a corrupt log: the
service fails closed and leaves it for investigation. For a fresh demo, choose
a new empty runtime root. Source: DI-dapod; DI-damod.
