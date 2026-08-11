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
- [ ] 1. Verify the exact immutable makerspace spec bytes, pCIDs, and cited
  record-profile document identity.
- [ ] 2. Replace legacy event persistence with exact-record validation,
  storage, replay, and projection under the four family specs.
- [ ] 3. Lock and implement the participant signing/ingress embodiment needed
  for semantic author evidence; keep browser/account behavior distinct.
- [ ] 4. Specify and implement byte carriage only after agents own verified
  durable records.
- [x] 5. Add the browser/account embodiment and an opt-in multi-agent proof.
  `scripts/run-two-agent-browser-proof.sh` owns a disposable Alice/Bob/Chrome
  session and emits an inspectable per-run evidence directory. The browser
  transports unsigned requests only; author evidence remains Alice's returned
  signed record. Source: DI-fuzar; DI-hibok; DI-kasaz.
- [ ] 6. Publish implementation claims, testing evidence, and all deferrals.

## Completion condition

Ex7 is not described as PromiseGrid-complete until every claimed family has
its named spec identity, live conformance evidence, scoped implementation
claim, and explicit embodiment/dependency limits. Unimplemented identity,
carriage, account, and recovery work remains explicit.
