# Workflow receipt inbox

TE ID: TE-rasih

## Status

decided

## Decision under test

How a receiving EX6 node exposes workflow artifacts and sender lifecycle
evidence retained by the dedicated relay endpoint without converting receipt
into a local workflow lifecycle decision.

This extends TE-novuk / DI-novuk. It applies to the received artifact CAS and
the separate `workflow-evidence` CAS already populated by relay receipt.

## Assumptions and trust model

- Alice transfers an exact artifact plus signed lifecycle evidence to Bob.
- Bob independently verifies the peer, signature, artifact CID, and lifecycle
  envelope before retention.
- Receipt is evidence of Alice's decision, never Bob's import or activation.
- Bob needs an operator-readable way to discover retained artifacts before
  choosing an alias and importing one locally.
- Existing `workflow import <alias> <artifact-cid>` remains the only action
  that creates Bob's local imported lifecycle state.

## Alternatives

1. **Direct CAS scan.** Add `workflow inbox list` and `workflow inbox inspect`
   as read-only derived views over the received artifact CAS plus
   `workflow-evidence` CAS; retain no new durable inbox log.
2. **Durable inbox ledger.** Append a new local receipt event whenever the
   relay endpoint accepts a transfer, then project list/inspect from that log.
3. **Automatic import.** Make receipt create a normal imported workflow entry,
   distinguished only by a source flag.

## Scenario analysis

### Normal operation

Under option 1, Bob runs `moks workflow inbox list`, sees an artifact CID,
sender identity/evidence summary, and whether it is already locally imported.
Bob then uses `moks workflow inbox import <artifact-cid> <alias>`, which calls
the existing explicit import path. The inbox is a derived observation over
already durable bytes.

Option 2 offers explicit receipt chronology, but adds a second durable
lifecycle-like record even though the evidence CAS already proves what arrived.
Option 3 is convenient but makes Alice's transfer alter Bob's workflow catalog
and blurs the core independent-activation boundary.

### Failure, corruption, and incomplete writes

Option 1 lists only artifact/evidence pairs that passed the existing import
verification. If a crash leaves an artifact without evidence or vice versa,
the inbox reports it as incomplete and refuses the inbox-import convenience
action. It does not need a recovery log because full CAS scanning reconstructs
the view.

Option 2 must atomically coordinate artifact CAS, evidence CAS, and a ledger
or define reconciliation semantics. Option 3 can leave an apparent local
workflow imported when evidence retention failed.

### Concurrent actors and mixed versions

Several peers may send the same artifact with distinct evidence. Option 1
groups by artifact CID and reports all valid evidence CIDs/senders. Bob's local
import remains one alias-to-artifact decision. Older nodes simply retain the
same bytes without inbox commands.

Option 2 needs deduplication and ordering rules for duplicate receipts.
Option 3 needs source-flag merge rules and may change an alias unexpectedly.

### Long-horizon evolution and scale

Direct CAS scanning is proportional to locally retained receipt data and
requires no new protocol or event family. A disposable cache may later improve
large inboxes, as long as full CAS remains authoritative. The evidence CAS
already creates a natural retention/garbage-collection boundary.

An inbox ledger is only justified if EX6 later needs human receipt chronology,
acknowledgements, or policy-driven expiry. Automatic import remains unsafe at
any scale because transport receipt and local availability are different facts.

## Conclusion

Reject automatic import. Defer a durable inbox ledger.

Option 1 survives: derive a read-only inbox from artifact/evidence CAS scans,
group all valid evidence by artifact CID, and expose one explicit convenience
import command that delegates to the existing local import operation.

## Decision status

Locked by DI-jifuk. The command is `moks workflow inbox import <artifact-cid>
<alias>` so the retained remote identity is presented before Bob assigns a
local alias.

## Implications for open work

- Add deterministic inbox scan tests for multiple senders, corrupt evidence,
  restart, and already-imported artifacts.
- Keep evidence ownership/source visible but do not infer remote authority.
- Do not add a new top-level PromiseGrid action or automatic lifecycle event.
