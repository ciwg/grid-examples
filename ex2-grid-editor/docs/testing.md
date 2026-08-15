# Testing Grid Editor

This guide explains how to verify Grid Editor and what each test layer proves.
It does not turn a passing local run, one relay's evidence, or a shared host
into a claim of centralized or globally authoritative PromiseGrid behavior.
Source: `DI-bubab`; `DI-nilas`; `DI-todav`.

## Run the verification suite

From `ex2-grid-editor/`, run:

```bash
go vet ./...
go test ./...
errcheck ./...
cd web && npm test && npm run build
```

`go vet` checks likely Go-code mistakes, `go test` runs the deterministic unit
and integration coverage, and `errcheck` verifies that Go errors are handled.
Run all three after a behavior or documentation-supported contract change.
Source: `DI-bubab`.

The `web` command runs deterministic browser-source tests, then rebuilds the
checked-in `web/app.js` bundle from those sources. It does not contact a relay
or create presence evidence.

## Test layers and their claims

### Source-derived protocol inventory

`protocol.TestProfileInventoryMatchesPublishedDocuments` derives each pCID from
the exact bytes of the four files in `protocols/`, then verifies that both
`docs/architecture.md` and `CHANGELOG.md` publish the same values. It proves
that those two documents match the current repo-local draft sources; it does
not claim that the drafts are frozen upstream specifications or universally
interoperable profiles. Source: `DI-guros`; `DI-ralit`.

### Relay-local exception evidence

The malformed-input, invalid-proof, and no-supported-handler tests in
`service/app_test.go` verify the relay's bounded local evidence policy:

- the received bytes are retained in that relay's CAS;
- one `observations.jsonl` record is appended for each receipt; and
- rejected input stays out of accepted peer feed and replay paths.

These records show what this relay observed. They do not prove another
participant's intent or globally invalidate a pCID. Source: `DI-todav`;
`DI-guros`.

### Collaboration and interoperability

`service/app_test.go` exercises accepted envelopes, peer relay exchange, and
replay using independently rooted relay instances. `service/interoperability_test.go`
exercises browser and Neovim embodiments through the relay-facing contract.
Together, these tests cover the current decentralized collaboration paths; they
do not define key rotation, delegation, or cross-relay trust policy. Source:
`DI-guros`; `DI-nilas`.

### Local presence lifecycle

`web/src/presence.test.mjs` injects a clock and scheduler to prove the normal
profile's `live` → `stale` → `offline` → removed boundaries, nearest-boundary
timer selection, and timer cancellation. The browser and Neovim use that same
observer-local policy to refresh presentation without emitting expiry traffic
or asking a relay to declare membership. This proves only local rendering of
the latest accepted awareness observation; it does not prove that a peer left,
revoked a promise, or is unreachable. Source: `DI-dizut`; `DI-dazin`.

## Test data and topology

Go tests use `t.TempDir()` for relay roots, identities, CAS directories, and
observation logs. Each temporary root is isolated and automatically removed by
the Go test framework. It is test data, not shared production evidence.

One host can run a genuine one-node session with one relay, key, and data root,
or simulate multiple decentralized relay nodes with separate processes, keys,
and roots. The latter models logical nodes; a shared host does not merge their
local evidence or identities. Source: `DI-bubab`; `DI-nilas`.

## Related documents

- [Architecture overview](architecture.md)
- [Implementation scope declaration](../CHANGELOG.md)
- [README](../README.md)
