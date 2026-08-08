# Completeness verification

Use the browser to record an observation, clear a hold as a steward, create a
loan, and record its return. Confirm the shared evidence timeline changes after
each action. Restart the service with the same runtime root to verify replay.

For API checks, run `go test ./...`. A malformed local `events.jsonl` must make
startup fail without changing its bytes. This is a single-process local demo:
it provides no authentication, replication, signatures, or multi-writer safety.
Source: DI-damod.
