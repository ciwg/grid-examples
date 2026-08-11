# Completeness verification

In the browser, verify that the page labels itself as exact signed-record
ingress, not a member-signing form. Submit a base64-encoded externally signed
record only when its public key is locally recognized. The stock command has
no recognized keys and therefore demonstrates exact retention only. The Go
suite provides the configured-policy proof: a known observation projects, an
unknown pCID is retained without projection, and a recognized non-steward key
cannot clear a safety hold. Restart with the same runtime root to verify
replay.

Equivalent API checks, while the local server is running:

```bash
curl -sS -X POST http://127.0.0.1:7037/api/records -H 'Content-Type: application/json' -d '{"records":["<base64-canonical-signed-record>"]}'
```

For API checks, run `go test ./...`. A malformed `records.frames` file must
make startup fail without changing its bytes. This embodiment provides no
browser signing, account-based authoring, key continuity/recovery, relay
carriage, or portable governance. Source: DI-tohak; DI-piruf.
