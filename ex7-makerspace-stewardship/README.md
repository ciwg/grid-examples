# Makerspace Stewardship

`ex7-makerspace-stewardship` is an Ex7 makerspace record-agent example. It
retains externally signed, pCID-selected makerspace records and derives a
local view only from known-family records whose public keys its local policy
recognizes. It is not a browser-signing, account-identity, relay, or portable
governance implementation. Source: DI-tohak; DI-piruf.

The browser displays local state and submits already signed record bytes. A
selected name, browser session, or account does not author a promise. Unknown
or unrecognized records remain retained evidence without changing the local
projection.

Run it with:

```bash
go run ./cmd/makerspace-stewardship
```

Then open `http://127.0.0.1:7037/`.

The server writes append-only exact record frames to
`.makerspace-stewardship/records.frames` by default. Use `-runtime-root` to
choose another location. Existing browser workflow controls are intentionally
replaced by signed-record ingress until a separate signing embodiment exists.
The stock command starts with no recognized public keys, so it retains valid
records without projecting them; configuring a local recognition policy is the
next embodiment slice. Create `.makerspace-stewardship/recognition.json` with
mode `0600`, or use `-allow-empty-recognition` only for retention-only
bootstrap. Source: DI-piruf; DI-likoh.

Read the [operator guide](docs/operator-guide.md),
[workflow and evidence guide](docs/workflow-and-evidence-guide.md), and
[testing guide](docs/testing.md),
[implementation claims](docs/promisegrid-implementation-claims.md),
[completeness verification guide](docs/completeness-verification.md), and
[CHANGELOG](CHANGELOG.md) before evaluating the example.
