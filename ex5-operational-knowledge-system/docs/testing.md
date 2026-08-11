# Ex5 testing guide

Run the full Go verification suite from `ex5-operational-knowledge-system/`:

```bash
go test ./...
go vet ./...
errcheck ./...
```

The suite verifies frozen-family signed envelopes, CAS-backed replay/export,
peer and relay exchange, local socket embodiments, browser native-host bridge
contracts, and operational projections. Browser DOM smoke tests remain
stub-backed: they prove page rendering and mock-bridge behavior, not the
extension/native-host path.

## Real Chrome/native-host evidence

The external browser harness lives in
`../../grid-examples-browser-checks/ex5/`. It attaches to the prelaunched,
verified Google Chrome session at `http://127.0.0.1:9222`; it must not launch a
synthetic browser session. The demo launcher retains the Chrome DevTools-pipe
owner that loads the unpacked Ex5 extension for that live session, because
Chrome does not persist unpacked developer extensions across restarts. Its
assertions cover Draft queue, Problem hotspots,
Known-record search, and visible `Current Record` changes, followed by run,
evidence, and approval interactions. Chromium is not an acceptable substitute
until TODO 153 has its own passing verification.

Before running the harness, prepare and launch the checked-in demo path:

```bash
./scripts/setup-demo-browser.sh
./scripts/launch-demo-browser.sh
./scripts/verify-demo-browser.sh
cd ../../grid-examples-browser-checks/ex5
npm run test:demo
```

After a run, inspect `events.jsonl`, the eight `*-messages.jsonl` signed-family
logs, `cas/objects/`, `attachments/`, `drafts/`, and relay-feed records under
the demo runtime. These are durable/runtime evidence; the browser DOM itself
is presentation evidence only. The launcher replaces only the prior PID in
`/tmp/ex5-demo-browser/chrome.pid` before it prepares the next attachable demo
session. Source: `DI-sobek`; `DI-punek`; `DI-bahak`; `DI-danir`; `DI-bilim`;
`DI-sofol`.

For a live browser demonstration, follow the current browser setup and
walkthrough in [the user guide](user-guide.md), then inspect runtime evidence:
the signed-family logs and CAS objects, copied evidence blobs, and relay-feed
records. The exact frozen-family implementation claims and their scope are in
[PromiseGrid Implementation Claims](promisegrid-implementation-claims.md).
