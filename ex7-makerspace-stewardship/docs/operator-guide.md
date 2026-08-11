# Operator guide

Run an isolated demo:

```bash
go run ./cmd/makerspace-stewardship -runtime-root /tmp/ex7-makerspace-demo
```

Then open `http://127.0.0.1:7037/`. The agent retains exact record frames in
`/tmp/ex7-makerspace-demo/records.frames`. Submit only externally signed,
base64-encoded Ex7 records through the page or `POST /api/records`; a browser
selection or account session cannot create an author signature. Stop and
restart the same command to confirm record replay. Do not edit a corrupt frame
file: the service fails closed and leaves it for investigation. For a fresh
demo, choose a new empty runtime root. Before normal startup, create
`<runtime-root>/recognition.json` with mode `0600`:

```json
{"version":1,"keys":[{"label":"alice","ed25519_public_key_base64":"<base64-32-byte-ed25519-public-key>"}]}
```

The file contains public keys only and is read at startup. It is not editable
through the browser or an account. Missing, malformed, duplicate, or insecure
files fail startup; use `-allow-empty-recognition` only for retention-only
bootstrap. Source: DI-tohak; DI-piruf; DI-likoh.

## Repeatable two-agent browser proof

Run the complete browser embodiment proof from this directory:

```bash
scripts/run-two-agent-browser-proof.sh
```

The harness owns one disposable process session: it bootstraps Alice and Bob,
starts Alice's participant agent on `127.0.0.1:7037` and Bob's terminal agent
on `127.0.0.1:7038`, launches Chrome with DevTools, submits Bob's unsigned
request to Alice's explicit loopback target, approves it locally at Alice, and
waits for Bob to ingest and project Alice's returned signed record. The
browser transports a request and displays evidence; it does not create author
evidence, and neither an account nor the request token is an author key.

On success the harness prints `browser proof evidence: <path>`. That path is a
new, disposable `/tmp/ex7-browser-proof.XXXXXX` directory for that run.
Inspect `bob-final.png`, `alice-approval-response.json`, `alice.log`,
`bob.log`, and `chrome.log` there before removing it. `bob-final.png` is the
browser-visible final projection; the approval response and framed agent state
remain the exact-record evidence. Source: DI-fuzar; DI-hibok; DI-kasaz.
