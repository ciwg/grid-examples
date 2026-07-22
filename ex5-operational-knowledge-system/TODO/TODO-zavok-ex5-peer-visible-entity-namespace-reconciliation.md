# TODO zavok - reconcile ex5 peer-visible entity namespaces across independent peers

## Decision Intent Log

ID: DI-loruk
Date: 2026-07-22 12:03:05 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Lock TODO `109` to Alternative A with canonical durable entity IDs derived from the create-event envelope CID. The current short IDs become presentation aliases only.
Intent: Make `ex5` fully PromiseGrid-native at the durable entity-identity layer instead of preserving local counter-minted IDs as the real keys.
Constraints: Keep enough alias data for legacy replay compatibility and embodiment transition. Do not keep import-time namespace translation tables as the long-term model. Cross-record references, stored entities, and peer exchange must converge on the canonical ID.
Affects: `docs/thought-experiments/TE-loruk-ex5-peer-visible-entity-namespace-reconciliation.md`, `ex5-operational-knowledge-system/service/types.go`, `ex5-operational-knowledge-system/service/app.go`, `ex5-operational-knowledge-system/service/peer_exchange.go`, frozen-family envelope helpers, browser/CLI/Neovim/API docs, and TODO `109`

## Goal

Let `ex5` accept peer-visible histories from independent runtimes even when
those runtimes minted overlapping local-facing IDs such as `RECV-0001`,
`RUN-0001`, or `RESP-0001`.

## Why this exists

TODO `107` and TODO `103` introduced peer-stable origin identity and
non-bootstrap import, but the runtime still honestly rejects imported
create-event IDs that collide with an already-populated local namespace.

## Tasks

- [x] zavok.1 Run the required TE for peer-stable entity namespace
  reconciliation.
- [x] zavok.2 Lock whether the runtime adopts globally stable entity IDs,
  origin-qualified display IDs, or some other reconciliation model.
- [x] zavok.3 Implement the chosen namespace model across the peer-visible
  families.
- [x] zavok.4 Update the browser, CLI, docs, and peer-exchange claims for the
  new cross-peer entity identity model.

## Status

- closed
- canonical peer-visible entity IDs now come from create-envelope CIDs, short
  IDs are preserved as aliases, and non-bootstrap peer exchange no longer
  rejects cross-peer alias reuse
