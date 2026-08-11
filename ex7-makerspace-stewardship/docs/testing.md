# Testing

Run `go test ./...` for codec, pCID registries, participant root/device
history, revocation containment, 2-of-3 recovery, peer-card linkage,
exact-byte carriage, local role policy, framed replay, projection, and HTTP
ingress coverage. The participant tests prove that an ordinary record needs a
root-authorized device, a revoked device is rejected, Alice's root cannot
revoke Bob's device, one recovery witness cannot activate a replacement, two
matching witnesses can, and peer-card-linked carriage retains and projects its
exact enclosed record. Run `scripts/run-record-ingress-proof.sh` for the real
command proof: recognized observation projection, unrecognized exact-byte
retention without projection, and rejection of a recognized non-steward
clearance. The fixture generator is test-only and never supplies browser or
account authoring. Source: DI-tohak; DI-piruf; DI-likoh; DI-kasaz; DI-sisad.

Run `go vet ./...` for static Go correctness checks and `errcheck ./...` to
ensure the Go implementation does not discard errors. Run
`git diff --check` before committing documentation or code changes. The frozen
pCID tests hash every active specification and confirm that its checked-in
registry declares the same CIDv1 value. Source: DI-tohak; DI-sisad.

Run `scripts/run-two-agent-browser-proof.sh` for the browser-level proof. It
starts disposable Alice and Bob agents in one process session with Chrome
DevTools, has Bob submit an unsigned request to Alice's explicit loopback
target, approves from Alice's local page, and waits for Bob to independently
ingest and project the returned exact signed record. The runner requires both
agent listeners, DevTools, the local approval response, and nonempty
`bob-final.png`; it fails closed if any condition is absent. It prints the
per-run evidence root as `/tmp/ex7-browser-proof.XXXXXX`. Inspect
`bob-final.png`, `alice-approval-response.json`, and the Alice, Bob, and
Chrome logs in that directory. This proves request transport and record
ingress, not browser or account signing. Source: DI-fuzar; DI-hibok; DI-kasaz.
