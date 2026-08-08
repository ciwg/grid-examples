# Completeness verification

In the browser, select Alice and the table saw to record a safety-hold
observation; the tool becomes unavailable. Select Carol to clear that hold with
an inspection assessment; the tool becomes available again. Select Alice and
the cordless drill to create a loan, then record its return. Confirm the shared
evidence timeline changes after each action. Restart the service with the same
runtime root to verify replay.

Equivalent API checks, while the local server is running:

```bash
curl -sS -X POST http://127.0.0.1:7037/api/tools/table-saw/observations -H 'Content-Type: application/json' -d '{"reporterId":"alice","text":"Guard is loose","safetyHold":true}'
curl -sS -X POST http://127.0.0.1:7037/api/tools/table-saw/clear-safety-hold -H 'Content-Type: application/json' -d '{"stewardId":"carol","assessment":"Fastener tightened"}'
curl -sS -X POST http://127.0.0.1:7037/api/tools/cordless-drill/loans -H 'Content-Type: application/json' -d '{"memberId":"alice","dueAt":"2030-01-02T15:04:05Z"}'
curl -sS -X POST http://127.0.0.1:7037/api/tools/cordless-drill/returns -H 'Content-Type: application/json' -d '{"memberId":"alice","condition":"Returned with charger"}'
```

For API checks, run `go test ./...`. A malformed local `events.jsonl` must make
startup fail without changing its bytes. This is a single-process local demo:
it provides no authentication, replication, signatures, or multi-writer safety.
Source: DI-damod.
