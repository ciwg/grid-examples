# TODO zavok - reconcile ex5 peer-visible entity namespaces across independent peers

## Goal

Let `ex5` accept peer-visible histories from independent runtimes even when
those runtimes minted overlapping local-facing IDs such as `RECV-0001`,
`RUN-0001`, or `RESP-0001`.

## Why this exists

TODO `107` and TODO `103` introduced peer-stable origin identity and
non-bootstrap import, but the runtime still honestly rejects imported
create-event IDs that collide with an already-populated local namespace.

## Tasks

- [ ] zavok.1 Run the required TE for peer-stable entity namespace
  reconciliation.
- [ ] zavok.2 Lock whether the runtime adopts globally stable entity IDs,
  origin-qualified display IDs, or some other reconciliation model.
- [ ] zavok.3 Implement the chosen namespace model across the peer-visible
  families.
- [ ] zavok.4 Update the browser, CLI, docs, and peer-exchange claims for the
  new cross-peer entity identity model.

## Status

- open
- created because non-bootstrap import now works, but still rejects colliding
  local-facing entity IDs across independent peers
