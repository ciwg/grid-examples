# Makerspace Stewardship

`ex7-makerspace-stewardship` is an Ex7 makerspace record-agent example. It
retains externally signed, pCID-selected makerspace records and derives a
local view only from known-family records signed by a root-authorized
participant device (or active root) and recognized by its local role policy.
It validates participant recovery, peer-card, and exact-byte-carriage records,
but does not yet provide browser signing, account identity, a running relay,
or portable governance. Source: DI-tohak; DI-piruf; DI-kasaz; DI-sisad.

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
participant-history and makerspace records without projecting makerspace
effects. A makerspace record must first be linked to signed root/device history;
`recognition.json` then controls only local role assessment. Create
`.makerspace-stewardship/recognition.json` with mode `0600`, or use
`-allow-empty-recognition` only for retention-only bootstrap. Source: DI-piruf;
DI-likoh; DI-kasaz.

Read the [operator guide](docs/operator-guide.md),
[workflow and evidence guide](docs/workflow-and-evidence-guide.md), and
[testing guide](docs/testing.md),
[implementation claims](docs/promisegrid-implementation-claims.md),
[completeness verification guide](docs/completeness-verification.md), and
[CHANGELOG](CHANGELOG.md) before evaluating the example.

For the repeatable Chrome embodiment proof, run
`scripts/run-two-agent-browser-proof.sh`. It prints the disposable evidence
directory containing `bob-final.png`, Alice's approval response, and process
logs. The proof transports an unsigned terminal request to a participant
agent; the browser and any account context never author the returned record.
Source: DI-fuzar; DI-hibok.
