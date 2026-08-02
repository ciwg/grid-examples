# Adapter image acquisition boundary

TE ID: TE-pohug

## Status

decided

## Decision under test

When EX6 acquires a missing, allow-listed, digest-pinned adapter image for a
workflow artifact.

This refines TE-sobir / DI-harib and TE-malob / DI-hapak.

## Assumptions

- `workflow verify` is an operator inspection command and currently does not
  alter lifecycle state.
- A workflow can execute only when its artifact, active package declaration,
  pCID contracts, registry policy, and exact image are locally ready.
- Pulling an image changes local Docker state and may contact an allowed
  network registry.

## Alternatives

1. Pull automatically during `workflow verify`.
2. Pull automatically during `workflow run start`.
3. Require `moks workflow image pull <alias>` as an explicit operation.

## Scenario analysis

In normal operation, Alice transfers an artifact and Bob runs verification.
Option 1 surprises Bob by creating Docker state and contacting a registry in a
readiness check. Option 2 hides a slow or failing network operation inside a
state-transition command. Option 3 lets Bob inspect the missing dependency,
approve acquisition explicitly, then verify or start separately.

If a registry is unreachable, option 1 makes inspection unreliable and option
2 leaves start behavior coupled to network timing. Option 3 reports an
acquisition failure without changing workflow lifecycle state. With concurrent
operators or mixed-version nodes, an explicit command creates a clear audit
point and avoids a newer node making `verify` unexpectedly mutating.

At scale, explicit acquisition supports prefetching and controlled maintenance
windows. It keeps the registry allow-list meaningful because each network
contact is deliberate. All alternatives still require post-pull digest
verification; only option 3 preserves the read-only meaning of verification.

## Conclusion

Choose option 3. `moks workflow image pull <alias>` resolves the image only
from the verified artifact and its matching installed package declaration,
checks the local registry allow-list, pulls the exact digest, and verifies the
local Docker result. It does not activate the workflow or start a run.

## Decision status

Locked by DI-zivut.

## Implications for open work

- Add image availability to workflow verification output.
- Add explicit acquisition and post-pull digest validation.
- Keep `workflow verify` read-only and `workflow run start` fail-closed when
  the required image is unavailable.
