# Testing Ex3.5 decentralized DTLS Grid Editor

This guide explains how to verify Ex3.5 and what each test layer proves. A
passing local suite, a single node's evidence, or a shared-host simulation is
not global PromiseGrid trust proof. Source: `DI-hosit`; `DI-hadil`; `DI-pazis`.

## Run the verification suite

From `ex3.5-grid-editor-dtls/`, run the Go checks:

```bash
go vet ./...
go test ./...
errcheck ./...
```

Run the browser checks from `web/`:

```bash
npm test
npm run build
```

Run the Neovim sidecar build from `cmd/grid-nvim-sidecar/`:

```bash
npm run build
```

The Charm terminal inspector is covered by the Go suite. Its tests exercise
the Bubble Tea model's deterministic relay-entry filters; `go test ./cmd/grid-charm`
is the focused command. It does not require a terminal, live relay, or Gum.

`go test ./tui` covers the terminal workspace's UTF-16 cursor mapping, remote
cursor rendering, menu navigation, relay HTTP dispatch, export confirmation,
Glow dependency guidance, stable typing-presence rendering, and the absence of
an idle animation tick. Those assertions keep
terminal awareness offsets compatible with the Automerge sidecar's JavaScript
string offsets. It also covers Bubbles Help key bindings, List/Table relay
results, File Picker import entry, scrollable Activity Viewport, and Spinner
activation only for real command work. Source: `DI-mutoh`; `DI-tubol`;
`DI-vujub`.

The checked-in VHS tape can be validated from the exercise root with:

```bash
"$(go env GOPATH)/bin/vhs" validate demos/grid-tui-collaboration.tape
```

Run `scripts/record-grid-tui-vhs.sh <tape>` to start an isolated relay and
record a Grid TUI flow. The Bubbles tapes are `grid-tui-help.tape`,
`grid-tui-document-results.tape`, `grid-tui-file-picker.tape`, and
`grid-tui-activity.tape`; each produces matching GIF and MP4 evidence under
`demos/`. VHS needs
its `ttyd` runtime dependency; the launcher reports a missing recorder before
it starts the relay. On this host the recorder uses its project-local
`scripts/ffmpeg` Docker wrapper because the system Snap ffmpeg cannot encode
VHS frame sequences. The GIFs are visual demonstration evidence, not protocol
conformance proof. Source: `DI-mutoh`; `DI-vujub`.

The Go checks cover static diagnostics, deterministic Go tests, and handled Go
errors. Browser tests cover embodiment-local JavaScript behavior; the browser
build verifies the checked-in browser bundle can be rebuilt. The sidecar build
verifies the Neovim helper bundle can be rebuilt. Source: `DI-hosit`.

### Direct DTLS peer transport

`go test ./transport ./service` includes deterministic two-node DTLS checks.
They prove a node can retain a stable certificate fingerprint across restarts,
reject an invite with no fingerprint, establish a real UDP DTLS session with a
pinned peer certificate, exchange an acknowledged message, and ingest a
signed document envelope from one independent node data root into another.
They do not prove Internet reachability through arbitrary NATs; use a LAN,
port forwarding, or a private network such as Tailscale for multi-host demos.
Source: `DI-higom`; `DI-pirov`.

### Direct-node acceptance suite

`service.TestBrowserNodesConvergeThroughDirectDTLS` proves two independent
browser harnesses converge through separate nodes after the nodes exchange
only direct DTLS envelopes. `service.TestBrowserAndNeovimNodesConvergeThroughDirectDTLS`
proves the same path from a browser on one node to the Neovim sidecar on the
other. `service.TestDirectDTLSRetriesAfterLossAndReceiverRestart` starts the
receiver after a sender has already attempted delivery, restarts that receiver,
and verifies persisted envelopes converge without duplication.

`service.TestHeadlessBrowserEditorEditPersistsDocument` uses a real headless
Chrome CodeMirror surface, inserts text only after its document client is
ready, and verifies that the application's persisted document state contains
that edit. `service.TestHeadlessBrowserEditorsConvergeThroughDirectDTLS` uses
two independent app data roots and real Chrome editor surfaces in both
directions: browser A edits, node A sends the signed envelope directly over
DTLS to node B, browser B renders it; then browser B edits and browser A
renders the returned change. The browser-to-its-own-node hop remains the
current browser WebRTC DataChannel implementation. WebRTC is the
DTLS-backed browser transport; the browser does not expose raw-DTLS sockets.
The test proves the inter-node path is direct raw DTLS/UDP.

`service.TestHeadlessBrowserAndNeovimConvergeThroughDirectDTLS` uses Chrome's
real CodeMirror surface and the real Neovim sidecar in both directions across
two direct-DTLS nodes. `service.TestHeadlessBrowserDirectDTLSRetriesAfterReceiverRestart`
routes the direct path through a test-only UDP proxy that drops a datagram,
starts the receiver only after the first browser edit, restarts that receiver,
and verifies the browser's persisted envelopes converge without duplication.

The real-browser tests use Chrome DevTools automation rather than Chrome's
virtual-time DOM dump: WebRTC needs wall-clock ICE and DTLS timers before its
DataChannels are usable. The harness waits for both WebRTC channels and the
editable CodeMirror surface, then exits Chrome cleanly before test cleanup.
Run the exact browser/direct-node acceptance group with:

```bash
go test ./service -run 'TestHeadlessBrowser(EditorsConvergeThroughDirectDTLS|AndNeovimConvergeThroughDirectDTLS|DirectDTLSRetriesAfterReceiverRestart)' -count=1
```

Run `scripts/test-two-host-dtls.sh` with `DTLS_REMOTE_SSH` and
`DTLS_REMOTE_ROOT` to automate the Tailscale/private-network proof against a
second host. The script is intentionally not silently skipped: the acceptance
suite is incomplete until that command passes against a reachable second host.

## Test layers and their claims

### Source-derived protocol publication

`protocol.TestProfileInventoryMatchesPublishedDocuments` derives every pCID
from the exact five `protocols/*.md` source files and checks both
`docs/architecture.md` and `CHANGELOG.md`. It proves those documents match the
current local drafts; it does not claim frozen upstream specifications or
automatic independent-peer interoperability. Source: `DI-dilav`; `DI-hadil`.

### Relay-local evidence and admission boundaries

Focused service/store tests verify that a bounded malformed or otherwise
rejected peer envelope is retained in the observing relay's CAS, produces one
`observations.jsonl` record per receipt, and stays out of accepted message
replay and peer feed. They also verify that remote-admission denials produce
non-secret `admission-diagnostics.jsonl` records without accepted messages or
bearer/bootstrap retention.

These are relay-local facts. They do not prove another participant's intent or
globally invalidate a pCID. Source: `DI-pazis`; `DI-darif`; `DI-lozut`.

### Published-version restore

Restore tests require an exact published manifest for the same document and
one accepted `restore-published-version` envelope carrying the exact Automerge
change projected into the sync feed. They cover malformed input, capability
denial, restart replay, and the fact that concurrent changes can merge with a
restore-derived state. They do not claim owner/admin authority, destructive
rollback, byte-identical output under concurrency, or frozen-spec
interoperability. Source: `DI-hihok`; `DI-tibum`; `DI-nihiz`.

### Decentralized collaboration and embodiments

Existing `service/interoperability_test.go`, WebRTC/browser tests, WebSocket
sidecar tests, and
headless browser startup tests cover browser/Neovim behavior, relay exchange,
late join, live carriage, and private-session hardening. Browser JavaScript
tests cover storage fallback and startup recovery; the sidecar build covers its
embodiment-local helper.

The browser startup proof records a WebRTC offer plus DataChannel readiness
separately from the rendered transport label. Its normal late-join and
injected stale-blank-snapshot cases both require the document to render from
the live DataChannel after the initial relay state is prepared. These are
observations from an isolated test browser and relay; they neither change
pCID-selected payload meaning nor establish a network-wide transport
guarantee. Source: `DI-gofut`; `DI-raron`.

The fresh-document regression opens a previously unknown, empty document ID in
an isolated real browser and requires the document-sync WebRTC DataChannel plus a ready,
editable CodeMirror surface. A separate unit test proves the awareness client
does not wait for a first remote-awareness snapshot before it reports its
socket connection established.

`tui/integration_test.go` extends that proof through the actual Grid TUI
sidecar process. `TestBrowserAndGridTUIConvergeOnFreshDocument` requires a
browser-written document update plus browser name, color, typing state, and
cursor to reach the Grid TUI model and rendered terminal surface.
`TestGridTUIClientsConvergeOnFreshDocument` requires the same document and
awareness convergence between two independent Grid TUI sidecars. Both tests
use an isolated relay root and only loopback test processes. Run the focused
proof with:

```bash
go test ./tui -run 'Test(BrowserAndGridTUI|GridTUIClients)ConvergeOnFreshDocument' -count=1
```

These tests prove Ex3's current browser and terminal embodiments honor the
same local-draft relay path in the isolated test topology. They do not claim
frozen-spec interoperability, a network-wide transport guarantee, or identity
authority beyond the tested presentation fields. Source: `DI-gufuj`.

Together these tests cover Ex3's current decentralized collaboration paths.
They do not define general key rotation, delegation, cross-relay role
recognition, or person identity. Source: `DI-dilav`; `DI-hadil`.

## Test data, topology, and private-browser boundary

Go tests use `t.TempDir()` for relay roots, identities, CAS directories,
observations, and admission diagnostics. Each root is isolated and
automatically removed by the Go test framework. It is test data, not shared
production evidence.

One host can run one logical relay node or simulate several decentralized nodes
with separate relay processes, keys, and data roots. Separate roots preserve
each relay's local evidence and identity even when a test host is shared.

Automated browser storage fallback, late-join, and blank-snapshot recovery are
covered. Browser-level verification has also passed with isolated normal and
incognito Chrome sessions using one isolated relay, including convergence in
both directions. The check used local DevTools and native browser input rather
than a human-driven usability review; its exact evidence is in the completed
[TODO tamuk](../TODO/TODO-tamuk-grid-editor-private-browser-document-sync.md).
Source: `DI-hosit`; `DI-figak`; `DI-hadil`; `DI-sodoj`.

## Related documents

- [Implementation scope declaration](../CHANGELOG.md)
- [Architecture overview](architecture.md)
- [README](../README.md)
- [Private-browser TODO](../TODO/TODO-tamuk-grid-editor-private-browser-document-sync.md)
