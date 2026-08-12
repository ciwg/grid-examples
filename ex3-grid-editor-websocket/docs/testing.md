# Testing Ex3 Grid Editor

This guide explains how to verify Ex3 and what each test layer proves. A
passing local suite, a single relay's evidence, or a shared-host simulation is
not global PromiseGrid trust proof. Source: `DI-hosit`; `DI-hadil`; `DI-pazis`.

## Run the verification suite

From `ex3-grid-editor-websocket/`, run the Go checks:

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

The Go checks cover static diagnostics, deterministic Go tests, and handled Go
errors. Browser tests cover embodiment-local JavaScript behavior; the browser
build verifies the checked-in browser bundle can be rebuilt. The sidecar build
verifies the Neovim helper bundle can be rebuilt. Source: `DI-hosit`.

## Test layers and their claims

### Source-derived protocol publication

`protocol.TestProfileInventoryMatchesPublishedDocuments` derives every pCID
from the exact four `protocols/*.md` source files and checks both
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

### Decentralized collaboration and embodiments

Existing `service/interoperability_test.go`, WebSocket/service tests, and
headless browser startup tests cover browser/Neovim behavior, relay exchange,
late join, live carriage, and private-session hardening. Browser JavaScript
tests cover storage fallback and startup recovery; the sidecar build covers its
embodiment-local helper.

The headless browser startup proof records a sync WebSocket `sync-ready` event
separately from the rendered transport label. In its normal late-join case it
requires zero HTTP sync recovery reads. In its injected stale-blank-snapshot
case it requires the ready event plus one bounded HTTP sync recovery read from
the full relay history. These are observations from an isolated test browser
and relay; they neither change pCID-selected payload meaning nor establish a
network-wide transport guarantee. Source: `DI-gofut`; `DI-raron`.

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
