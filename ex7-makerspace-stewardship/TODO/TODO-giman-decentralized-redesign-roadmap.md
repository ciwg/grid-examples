# Ex7 decentralized redesign roadmap

## Decision Intent Log

This roadmap sequences DI-tohak. It does not create an implementation claim,
freeze an embodiment, or settle any PromiseGrid-wide runtime boundary.

## Purpose

Re-create Ex7 from its named makerspace contracts outward. The current local
demo remains an honest baseline until exact pCID-selected records are on the
live path and the implementation evidence supports a narrow claim. The full
roadmap is published at
[`docs/ex7-decentralized-redesign-roadmap.md`](../docs/ex7-decentralized-redesign-roadmap.md).

## Open slices

- [x] 0. Complete the source-grounded app-boundary audit and lock the
  recreation scope in TE-biban / DI-tohak. The audit is evidence, not an
  implementation claim.
- [x] 1. Verify the exact immutable makerspace spec bytes, pCIDs, and cited
  record-profile document identity. `service/protocol_registry_test.go` hashes
  the frozen makerspace, participant, peer, and carriage documents and checks
  their declared registries. Source: DI-tohak; DI-sisad.
- [x] 2. Replace legacy event persistence with exact-record validation,
  storage, replay, and projection under the four family specs. The live store
  uses `records.frames`; replay and unknown-pCID exact retention are covered
  by `service/store_test.go` and `service/app_test.go`. Source: DI-tohak;
  DI-piruf; DI-likoh.
- [x] 3. Lock and implement the participant signing/ingress embodiment needed
  for semantic author evidence; keep browser/account behavior distinct. Signed
  root/device history, revocation, 2-of-3 recovery, and terminal approval are
  covered by `service/participant_test.go` and `service/server_test.go`.
  Source: DI-kasaz; DI-sisad; DI-hibok; DI-fuzar.
- [x] 4. Specify and implement byte carriage only after agents own verified
  durable records. Peer-card-linked exact-byte carriage is implemented and
  tested without giving the carrier semantic-author authority. Source:
  DI-sisad; DI-kasaz.
- [x] 5. Add the browser/account embodiment and an opt-in multi-agent proof.
  `scripts/run-two-agent-browser-proof.sh` owns a disposable Alice/Bob/Chrome
  session and emits an inspectable per-run evidence directory. The browser
  transports unsigned requests only; author evidence remains Alice's returned
  signed record. Source: DI-fuzar; DI-hibok; DI-kasaz.
- [x] 6. Publish implementation claims, testing evidence, and all deferrals.
  The README, operator/testing/completeness guides, implementation claims,
  CHANGELOG, and this roadmap name the frozen contracts, evidence paths, and
  non-claims. Source: DI-tohak; DI-sisad; DI-fuzar.

## Completion condition

Ex7 is not described as PromiseGrid-complete until every claimed family has
its named spec identity, live conformance evidence, scoped implementation
claim, and explicit embodiment/dependency limits. Unimplemented identity,
carriage, account, and recovery work remains explicit.
