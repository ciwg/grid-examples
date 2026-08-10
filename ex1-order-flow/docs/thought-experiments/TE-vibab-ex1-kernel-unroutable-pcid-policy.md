# Ex1 kernel unroutable-pCID policy

TE ID: TE-vibab

## Status

decided

## Decision status

Locked by DI-zosiz: record `no_registered_recipient` as a kernel-local dispatch
observation without maintaining an allow-list or claiming global pCID validity.

## Decision under test

What should Ex1's kernel claim and record when it receives a syntactically
valid pCID-selected envelope with no registered local recipient? This TE is the
prerequisite for `lubav.5`.

## Assumptions and trust model

- The kernel parses only enough envelope structure to select a pCID and route
  exact bytes. It does not interpret the pCID payload or assess business
  promises.
- Registration declares current local recipients, not a global protocol
  registry. A profile can be valid yet have no current local subscriber.
- The kernel already retains raw ingress bytes and a local observation record.
- Mallory may send an unknown pCID, a known local profile pCID while its role
  is offline, or a future profile pCID unsupported by this checkout.

## Alternatives

### A. Record `no_registered_recipient` and decline local dispatch

The kernel makes only the local statement that no current registration matched
the selected pCID. It retains exact bytes and makes no claim whether the pCID
is globally valid, unsupported, or meaningful elsewhere.

### B. Maintain a kernel allow-list of Ex1 profile pCIDs and record `unknown_pcid`

The kernel distinguishes known local profile identities from other pCIDs before
checking recipients.

### C. Drop unmatched traffic without a durable kernel record

The kernel avoids maintaining an observation for unroutable ingress.

## Scenario analysis

### Offline valid recipient

Alice sends a valid `shipment` envelope while Ellen is offline. A records only
that no current registration can receive it. B calls the pCID known but still
cannot deliver it, adding a second state with no routing benefit. C loses the
local delivery evidence.

### New or foreign pCID

Mallory sends bytes under a pCID this checkout does not recognize. A retains
the bytes without pretending the kernel can decide global validity. B turns the
kernel into a local profile registry and risks treating a future compatible
handler as semantically special. C loses evidence.

### Long-horizon evolution

When an Ex1 profile document changes, its pCID changes. A continues to state
only current recipient availability, so old and new pCIDs remain equally
carriage-safe. B requires every profile revision to update kernel knowledge and
invites a false distinction between local version inventory and protocol truth.

## Conclusion

C is rejected because it discards local evidence. B is rejected because the
kernel's job is local pCID dispatch, not profile validation or registry
authority. A survives and is recommended: rename the current observation to
`no_registered_recipient`, retain the raw envelope, make no signed reply, and
document that an agent receiving an unexpected pCID separately records its own
`unexpected_pcid` observation.

## Decision still requiring DF

Adopt A, `no_registered_recipient` as the kernel-local policy and observation
name (recommended), or B, an allow-list-backed `unknown_pcid` policy?

## Implications for open work

- `lubav.5` updates the kernel record name, design/spec documentation, and
  regression tests after the choice is locked.
- `lubav.6` can rely on the policy-specific test coverage.
