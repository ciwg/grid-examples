# Completeness verification

In the browser, verify that the page labels itself as exact signed-record
ingress, not a member-signing form. Submit a base64-encoded externally signed
record only when its signing key is linked through signed participant root and
device history and is locally recognized for the relevant makerspace role. The
stock command has no recognized keys and therefore demonstrates exact
retention only. The Go suite provides the configured-policy proof: a known
observation projects, an unknown pCID is retained without projection, a
revoked device is rejected, one recovery witness is insufficient, two matching
witnesses permit a replacement root, and a recognized non-steward key cannot
clear a safety hold. Restart with the same runtime root to verify replay.

Equivalent API checks, while the local server is running:

```bash
curl -sS -X POST http://127.0.0.1:7037/api/records -H 'Content-Type: application/json' -d '{"records":["<base64-canonical-signed-record>"]}'
```

For API checks, run `go test ./...`. A malformed `records.frames` file must
make startup fail without changing its bytes. This embodiment provides no
browser signing, account-based authoring, a running peer-discovery or relay
service, blob retrieval, or portable governance. Source: DI-tohak; DI-piruf;
DI-kasaz; DI-sisad.

For the complete two-agent browser evidence, run
`scripts/run-two-agent-browser-proof.sh`. It creates a disposable Alice/Bob
session, performs Bob's unsigned terminal request and Alice's local approval,
then requires Bob's final browser projection. The runner prints its temporary
evidence root (`/tmp/ex7-browser-proof.XXXXXX`); inspect `bob-final.png`,
`alice-approval-response.json`, `alice.log`, `bob.log`, and `chrome.log` in
that root. This check does not alter the non-claim above: browser controls and
accounts are not signing embodiments. Source: DI-fuzar; DI-hibok; DI-kasaz.
