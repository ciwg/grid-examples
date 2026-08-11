# Testing

Run `go test ./...` for codec, pCID, policy, framed replay, projection, and
HTTP ingress coverage. Run `scripts/run-record-ingress-proof.sh` for the real
command proof: recognized observation projection, unrecognized exact-byte
retention without projection, and rejection of a recognized non-steward
clearance. The fixture generator is test-only and never supplies browser or
account authoring. Source: DI-tohak; DI-piruf; DI-likoh.
