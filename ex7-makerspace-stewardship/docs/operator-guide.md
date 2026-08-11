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
