# Finished Ex7 Participant-Agent Product

TE ID: TE-bijad

## Status

decided

## Decision under test

What complete, end-to-end PromiseGrid product Ex7 implements: durable
participant identity, authorized devices, recovery, shared-terminal use,
account/bootstrap behavior, discovery, byte carriage, local assessment, and
makerspace evidence.

## Assumptions and scope

- Ex7 is an application with its own pCID-selected contracts. The guide does
  not supply a universal identity, recovery, browser, relay, or runtime API.
- Alice must be able to act from a home device and a makerspace terminal
  without becoming different identities. Carol must make a separately authored
  steward assertion. Mallory can compromise a device, submit malformed or
  unknown bytes, impersonate a label, replay records, or operate an unreliable
  relay.
- Exact signed records, local projection, and `recognition.json` are current
  infrastructure; they are not yet the finished product claim.
- A makerspace account is local UI/bootstrap policy. It never creates author
  evidence. No universal identity registry exists.

## Alternatives

### A. Participant-root, authorized-device, threshold-recovery product

Each participant has an offline-capable root signing identity. Root/history
records authorize active device keys, replace keys, revoke keys, and define a
participant-selected threshold recovery set. Authorized devices sign ordinary
makerspace records. A terminal creates a request and obtains an exact signature
from an authorized device or leaves an unsigned draft. Signed peer cards enable
local discovery. A relay or direct session carries exact bytes under its own
pCID contract; every agent stores, verifies, retains unknown bytes, and makes
its own local projection/trust decision.

### B. Browser/account identity product

The account or browser profile owns the participant key and recovers it,
authorizes terminals, and supplies discovery.

### C. Local record product only

Keep the current externally signed ingress and local recognition file, with no
cross-device continuity, recovery, discovery, or carriage contract.

## Scenario analysis

### Normal cross-device operation

Under A, Alice's root history authorizes her phone and laptop. Either device
signs the same pCID-selected loan promise. At the makerspace terminal, Alice's
account opens a request UI; her phone reviews and signs the exact bytes. Carol
uses her own authorized device to sign a safety disposition. Agents retain and
assess those records independently.

B makes account/browser custody equal authorship. C cannot explain the
terminal or relationship between Alice's devices.

### Loss, compromise, and recovery

Under A, an active key revokes a lost device; a replacement is authorized by
the existing history. If every active key is unavailable, the participant's
predeclared threshold recovery keys issue the recovery evidence that activates
a new key. If threshold recovery is unavailable, the result is honestly a new
identity, not a silent reset. Each agent evaluates conflicting or late recovery
evidence under its own policy.

B makes account recovery the identity authority. C has no durable answer to
loss or replacement.

### Offline, discovery, relay outage, and duplicate traffic

Under A, Alice's agent retains her signed record while offline. Signed peer
cards provide locally chosen contact hints. Direct or relay carriage moves
opaque exact bytes, tolerates retries and duplicates, and cannot change their
meaning. A relay outage delays exchange; it does not change identity or local
history. Unknown pCIDs remain exact bytes without known semantics.

B makes the account/discovery service a required authority. C has no completed
multi-participant exchange story.

### Local makerspace assessment

Under A, signature validity, participant identity history, device status, and
local recognition are distinct. A valid Alice device key does not make Alice a
steward; Carol's recognized, active device evidence may clear a hold only under
that agent's local makerspace policy. No record creates a universal command or
reputation score.

### Long-horizon evolution and scale

A freezes narrow pCID families for root history, device authorization,
revocation, threshold recovery, peer cards, and carriage. New workflows
compose existing families; new durable meanings receive new immutable specs.
Agents can use different relays and policy sources without rewriting record
bytes. Storage remains exact records/blobs plus rebuildable local projections.

## Conclusions

B is rejected because account/browser custody would be confused with durable
participant authorship. C is rejected because it cannot deliver the required
finished multi-device, recoverable, multi-participant product.

A survives and is recommended as the finished Ex7 architecture.

## Output to decision framing

The complete DF lock must select A and additionally fix, as one coordinated
decision set: recovery threshold and key roles; device/terminal signing flow;
peer-card and carriage payloads; local policy/configuration lifecycle;
record/blob retention; runtime paths; exact pCID/doc-CID registry paths; UI
states; and finished E2E acceptance scenarios.

## Decision status

locked: Alternative A, participant-root authorized-device threshold-recovery
product, by DI-zodah in
`../../TODO/TODO-bubuz-canonical-makerspace-records.md`. Exact protocol and
embodiment details remain coordinated DF work before implementation.
