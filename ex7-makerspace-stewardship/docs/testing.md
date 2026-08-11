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
