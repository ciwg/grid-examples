# Ex4 testing guide

## Automated checks

From `ex4-bug-tracker/`, run:

```bash
bash scripts/verify.sh
```

For an individual layer, run `go test ./...`, `go vet ./...`, `errcheck ./...`,
or `node scripts/verify-browser-webcrypto.mjs` directly.

The service regression suite covers explicit enrollment, canonical prepare/finalize signing bytes, tag-42 pCID envelope parsing, accepted CAS retention, rejected ingress, unsigned-write rejection, browser signer-asset serving, and signed attachment-reference acceptance/download. The CLI integration suite verifies that its mutable workflow commands enroll an explicit local `--agent-key`, produce a signed lifecycle promise, and update the local projection only after server acceptance.

## Browser check

Start an isolated runtime root:

```bash
go run ./cmd/bug-tracker --data-root /tmp/ex4-check
```

Open the displayed URL, select `reporter`, and create an issue. The browser creates a non-extractable IndexedDB WebCrypto key for that role, completes explicit service-scoped enrollment, asks the service to prepare canonical signing bytes, signs them locally, finalizes, and posts the returned raw-CBOR envelope. Confirm the issue appears, then inspect `/tmp/ex4-check/cas/`, `accepted-promises.jsonl`, `observations.jsonl`, and `agent-bindings.jsonl`. Stop and restart with the same runtime root; the issue projection is replayed from `events.jsonl` and the acceptance/evidence records remain available.

For an attachment, upload the bytes first. They receive a local CAS CID but do not affect an issue. The browser then signs an `issue-attachment-reference` promise naming that CID; only acceptance of that promise adds the attachment to the timeline. The browser's identity picker is demonstration-local role policy, not a claim of global identity, delegation, or federation. Source: `DI-kolaf`; `DI-ninul`.

## Adapter and artifact boundary

`/api/promises/prepare` and `/api/promises/finalize` are bounded JSON adapter requests. They do not carry a pCID or create a PromiseGrid artifact. The service resolves the selected local-draft profile, constructs canonical bytes, and the browser or CLI signs those exact bytes with its own key. `/api/promises` and the CLI enrollment endpoint require `application/cbor` and carry the final `grid([42(pCID), payload, proof])` artifact. Source: `DI-kolaf`; `DI-rutul`.

The current local-draft profiles are [issue report](../protocols/issue-report.md), [lifecycle update](../protocols/issue-lifecycle-update.md), and [attachment reference](../protocols/issue-attachment-reference.md). Their pCIDs are exposed at runtime from `GET /api/meta`; derive them from those exact source files before treating them as stable references. They are not frozen upstream PromiseGrid specifications or claims of cross-tracker interoperability.
