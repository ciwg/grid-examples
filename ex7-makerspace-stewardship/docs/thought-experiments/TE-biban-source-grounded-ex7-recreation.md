# Source-Grounded Ex7 Recreation

TE ID: TE-biban

## Status

decided

## Decision under test

What Ex7 may truthfully claim as its first real PromiseGrid implementation,
and what must be rebuilt before it makes that claim.

This TE responds to the PromiseGrid development guide and to the observed
runtime, rather than treating earlier Ex7 design prose as implementation
evidence. It concerns TODO `bubuz` and TODO `giman`.

## Assumptions and scope

- The PromiseGrid development guide requires an explicit frozen protocol spec,
  its pCID, and a scoped implementation claim against the exact frozen
  document. It does not provide a universal app SDK, mandatory storage layout,
  identity product, envelope, browser model, or transport contract.
- Ex7 already has four frozen makerspace family documents: equipment
  observation, safety disposition, off-site loan, and off-site return. They
  are protocol-design artifacts until the running code demonstrably speaks
  them. Their registry values must be independently recalculated from the
  exact immutable files before any implementation claim cites them.
- The running Ex7 program instead accepts local JSON requests, writes
  `events.jsonl`, and replays `Event{Type: ...}` values. Its CBOR record helper
  and framed-record writer are not on that path. Therefore it currently proves
  a local demo, not implementation of the four family protocols.
- The current guide presents `grid([42(pCID), ...protocol-defined-slots])` and
  a signed-message profile as useful current direction, while universal
  envelope and signature conventions remain provisional. Any Ex7 record
  profile must be named as an Ex7 application contract, not called a final
  PromiseGrid-wide rule. The provisional tag value in the current helper must
  likewise not be described as an official universal allocation.
- A participant's signature is author evidence. A makerspace account, browser
  session, display name, relay, or local HTTP request can be an embodiment
  convenience or a local recognition input, but cannot be substituted for
  author evidence.
- A safety disposition, loan, or return is evidence of a promise or
  observation. Each agent's local policy may decide what it displays or acts
  on; no record is a global command, and no shared process decides everyone's
  trust.

## Alternatives

### A. Spec-and-claim-first recreation

Keep the four existing makerspace family specs only if exact-byte pCID checks
pass. Treat `makerspace-record-v1` as Ex7's explicitly scoped common record
profile, with its own frozen document identity in implementation claims.

Rebuild the live action path so a participant embodiment submits a proposed
action to that participant's Ex7 agent. The agent produces or verifies an
exact record selected by one family pCID, stores the exact accepted bytes,
and derives its own projection from those bytes and stated local policy. A
hold prompted by an observation becomes two linked records when both meanings
are asserted; it is no longer hidden inside one local event. A return links to
the exact loan record it addresses.

The first claim is deliberately narrow: a named Ex7 agent embodiment implements
the four family contracts and the scoped record profile, subject to listed
local bootstrap and trust assumptions. Browser/account access and byte carriage
are separately documented embodiments. Relay carriage is added only after its
own named contract and tests exist.

### B. Signed single-process retrofit

Keep one shared HTTP process and retrofit signatures and pCID fields into its
JSONL actions. The process chooses a signer from browser-supplied member IDs,
keeps the authoritative current projection, and later exposes the log to other
participants.

### C. Local-demo claim only

Keep the JSONL event model, document its limits, and defer all family-record
implementation indefinitely while continuing to call the application a
PromiseGrid example.

## Scenario analysis

### Normal operation

Alice observes that a table-saw guard is loose and asks her Ex7 agent to make
an observation promise. If she also makes a safety disposition, the agent
creates the distinct record specified by the safety family, with a link to the
observation where the spec permits it. Carol's agent later signs a clearance
disposition after Carol's inspection. Alice's and Carol's agents each retain
the exact records and each uses its own local recognition policy to decide the
displayed condition.

Under A, the durable meanings match their named family specs. The browser can
remain useful without being mistaken for the author. Under B, the shared
process remains the effective authoring and current-state decision point even
if it writes signatures. Under C, the demonstrated behavior remains useful but
does not implement the frozen protocol documents.

For an off-site loan, Alice's agent signs the loan commitment with the accepted
policy snapshot. A later return record carries the loan-record link required by
the return family. A current HTTP route cannot stand in for either durable
meaning; it is merely one way to request the agent action.

### Failure, corruption, and incomplete writes

Under A, the agent validates its selected family payload and author evidence
before it makes a record projectable. It durably appends exact bytes before
changing the derived local view. It fails closed on malformed framing or a
truncated durable write, retaining intact prior bytes for inspection. A
well-framed unknown-pCID record can be retained unchanged without receiving
known makerspace semantics.

B can improve the old JSONL path but leaves an implementation-dependent event
format between the real action and the claimed family semantics. C retains the
existing JSON corruption behavior but has no exact-record conformance path.

### Concurrent participants, partitions, and mixed versions

With A, Alice can keep her own signed loan record while offline. Carol can keep
her own inspection record. They may exchange exact bytes directly or through
later byte-carriage infrastructure. A delayed, duplicated, unknown-family, or
locally unrecognized record does not empower the carrier to decide either
agent's projection. A newer agent may retain a record for a family it does not
implement.

B concentrates durable history and live projection in one process, making an
outage or policy choice there disproportionately decisive. C has no peer
exchange contract at all. Neither B nor C supplies the guide's required
spec-to-implementation claim relationship for the four families.

### Long-horizon evolution and user-created workflows

Under A, a new workflow composes the existing family pCIDs whenever its durable
meanings are observation, safety disposition, loan commitment, and return
claim. It does not need a new pCID merely because the UI flow is new. A truly
new interoperable durable meaning receives a new immutable specification and
pCID, with an explicit compatibility story. Earlier records remain exact,
auditable bytes.

B tends to extend an ever-growing event-type union controlled by the process.
C has no interoperable extension path. All alternatives must avoid silently
changing the meaning of historical bytes.

### Trust and evidence boundary

Under A, signature verification says which key made the record; it does not
establish membership, stewardship, correctness, or a universally binding
outcome. An Ex7 agent may locally recognize Carol's valid safety-disposition
record for a specified makerspace context, retain it without projection, or
reject it under explicit local policy. Account login may help the UI locate an
agent, but it cannot produce Alice's author evidence.

B risks presenting the shared process's account/session checks as semantic
authorship. C has no semantic author evidence at all. None of the alternatives
licenses a universal reputation score, a permanent gatekeeper, or an
unqualified claim about portable identity continuity.

### Scale and operational complexity

A requires real work: exact-byte codec/store/replay integration; per-family
validation and projections; explicit local trust configuration; hostile-input
tests; pCID/doc-CID checks; and narrow implementation claims. It buys a clear
interoperability target and prevents a polished demo from impersonating a
protocol implementation.

B is cheaper in the short term but leaves the crucial shared contract obscured
by local event mechanics. C is cheapest but cannot meet the requested product
goal. A does not require Ex7 to solve every future concern before its first
claim: relay, key rotation/recovery, account discovery, device delegation,
blob retrieval, and live-sync each remain separate protocol or embodiment work
until specified and tested.

## Conclusions

C is rejected for the product goal: it preserves a local demonstration but
does not implement the frozen makerspace protocols.

B is rejected: adding signatures to a shared event process does not make its
local event model the pCID-defined contract, and it blurs author evidence with
host behavior.

A survives. It is the smallest source-grounded recreation path because it
starts from named Ex7 contracts, makes their live bytes and implementation
claims testable, and leaves the guide's still-open universal layers explicitly
outside the claim.

## Output to decision framing

The required lock is:

**A — recreate Ex7 around the existing four frozen makerspace families and an
explicitly Ex7-scoped record profile; make the first implementation claim only
after exact pCID verification, live exact-byte storage/replay, per-family
validation/projection, and adversarial conformance tests.**

The following decisions remain deliberately separate and must not be smuggled
into this lock:

- the participant-agent and browser/account embodiment contract;
- key continuity, revocation, recovery, or delegated-device protocols;
- relay/feed carriage and its own evidence rules;
- blob availability/retrieval;
- portable governance, member, qualification, and tool-catalogue protocols.

## Decision status

locked: Alternative A, spec-and-claim-first recreation, by DI-tohak in
`../../TODO/TODO-bubuz-canonical-makerspace-records.md`.

## Implications for current Ex7 material

- The local JSONL app remains a baseline demonstration and must be described
  that way until replaced.
- The frozen makerspace specs remain preserved historical design artifacts,
  pending exact CID verification and live conformance.
- Earlier participant-agent, identity, and relay planning remains unimplemented
  and cannot be cited as evidence of the running product.
- TODO `bubuz` and TODO `giman` need a source-grounded rewrite after the DF
  lock, with explicit deferrals and a test/claim gate before each claim.
