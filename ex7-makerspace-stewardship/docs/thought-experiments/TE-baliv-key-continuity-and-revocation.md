# Participant Key Continuity and Revocation

TE ID: TE-baliv

## Status

needs DF

## Decision under test

How independently signing Ex7 agents publish key rotation and revocation evidence without a central identity authority.

## Alternatives

### A. Signed continuity and revocation records

Freeze separate pCID families for key-continuity and key-revocation promises. An active old key signs a rotation record naming the new public key; an active key signs a revocation record naming the revoked key and reason. Every agent retains the exact records and applies its own trust policy and conflict rules.

### B. Local configuration only

Each agent edits its trusted-key list with no portable evidence.

### C. Central account registry

One service declares current and revoked keys for all participants.

## Scenario analysis

A lets Alice rotate a key while peers retain a verifiable chain from the prior key. If Alice loses every active key, recovery is not silently invented; it needs a later recovery/witness protocol. A revocation is evidence that peers may assess under their own policy, not a universal command. B cannot carry continuity to another agent. C contradicts decentralization.

Under partitions, agents can retain and later exchange exact continuity/revocation bytes. Conflicting chains remain evidence rather than being overwritten by a registry. At scale, chains are compact and family-specific; local trust policies remain explicit.

## Conclusions

B and C are rejected. A is recommended, with one new frozen pCID per continuity and revocation family, local assessment rules, no automatic lost-key recovery, and no central registry.

## Decision status

needs DF
