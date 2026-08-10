# Ex1 durable observer evidence

TE ID: TE-golig

## Status

decided

## Decision status

Locked by DI-vihoz, DI-riguz, and DI-purum: retain append-only local
`ObservationRecord` entries in `observations.jsonl`, retain raw bytes at
ingress, give the kernel its own local data directory, and use
`Client.RecordObservation` for handler-local evidence.

## Decision under test

How should Ex1 durably retain local evidence of timeout, malformed or
unauthorized input, signed refusal, and unknown-pCID receipt without changing
what another agent promised or turning the kernel into a business authority?
This TE is the prerequisite for `lubav.4`.

## Assumptions and trust model

- Ex1 already retains valid sent and received envelopes as raw CBOR artifacts,
  with append-only message records. A malformed incoming byte currently fails
  before that retention path; an unknown pCID at the kernel can have no local
  subscriber and be dropped without evidence.
- A signed refusal is evidence of the signer’s own response. A timeout,
  malformed byte, invalid proof, invalid token, or unknown-pCID receipt is only
  the observer's local evidence; none proves another agent's intent.
- Alice, Bob, Carol, Dave, and Ellen are cooperative example roles. Mallory may
  send malformed, unauthorized, replayed, or unknown-pCID traffic.
- The kernel may retain local operational evidence needed to account for bytes
  it observed, but must still dispatch only by pCID and avoid business/promise
  assessment.
- The current local draft profile pCIDs and envelope shape remain unchanged.

## Alternatives

### A. Append-only local observation records plus raw-byte retention

Each observer writes one append-only local record for a material observation.
Records reference an exact raw-byte CID when bytes exist, expected/request CIDs
when a timeout is observed, and a narrow locally assessed reason. Agents retain
their observations; the kernel retains its own unknown/malformed ingress
observations. A signed refusal is retained as raw signed evidence and may have
an observer record that says only it was received and verified.

### B. Reuse `messages.jsonl` for all observations

Timeouts, malformed bytes, and signed refusals become new directions or status
values in the existing sent/received message log.

### C. Keep only normal logs and raw valid-envelope artifacts

The system reports errors or timeouts operationally but adds no durable local
observation evidence.

## Scenario analysis

### Normal signed refusal

Carol sends Bob a valid signed `pick_pack` result with `status = refused`.
Bob stores its raw bytes and can verify Carol's signature. Under A, Bob also
appends a local `refusal_observed` record referring to that CID; the record
does not claim that Carol broke a promise. B mixes a local assessment into a
message transport log. C preserves raw evidence but makes the local observed
outcome difficult to query and distinguish from other results.

### Timeout

Bob sends Ellen a shipment request, stores its exact CID, and receives nothing
before his deadline. Under A, Bob appends a `timeout_observed` record referring
to the request CID and the locally configured deadline. It does not say that
Ellen refused, failed, or broke a promise. B has no natural raw received
envelope to place in `messages.jsonl`. C loses the durable basis for later
assessment of Bob's signed top-level failure result.

### Malformed or unauthorized input

Mallory sends bytes that cannot parse, a valid envelope with a bad proof, or a
valid proof with an expired/mismatched token. Under A, the receiving agent (or
kernel at ingress) retains the exact bytes when present, computes a CID, and
appends a local observation such as `malformed_input`, `invalid_proof`, or
`invalid_capability`. No response is required. B overloads a log whose current
schema assumes a valid pCID and direction. C lets Mallory erase the forensic
trail merely by sending unparsable traffic.

### Unknown pCID

Mallory sends a syntactically valid envelope with an unknown pCID to the
kernel. Under A, the kernel stores the exact bytes and an `unknown_pcid`
observation before declining to dispatch them. An agent that receives an
unexpected pCID records its own local observation too. B gives no owner or
valid pCID field for an unrecognized profile. C silently drops the only
evidence that the kernel observed it.

### Concurrent actors, long horizon, and scale

Append-only records let Alice and Bob disagree locally without overwriting one
another. CID references avoid duplicating raw bytes. A profile revision yields
a new pCID, and earlier observations still name the old raw bytes and profile
text. At scale, A adds local storage proportional to observed exceptional
traffic and requires retention limits to remain host-local policy; it does not
create a central evidence service. B complicates every normal message-log
reader. C is cheap but undermines audit and migration.

## Conclusions

C is rejected because it cannot preserve the local evidence required to
distinguish timeout observation from a claim about another agent. B is rejected
because `messages.jsonl` describes valid sent/received envelopes and is not a
sound schema for malformed bytes or no-byte timeouts.

A survives and is recommended: use a distinct append-only local observation
record with raw-byte CID references where applicable. Retain the raw bytes
before parse/validation at every observing ingress. Have the kernel record only
its local ingress/dispatch observations, never a business outcome. Signed
refusal remains the signer’s evidence; the observer record states only receipt
and successful verification.

## Decisions still requiring DF

1. **Record location:** add a distinct append-only local observation log and
   raw-byte retention path (recommended), or extend `messages.jsonl`?
2. **Observer coverage:** have both agents and the kernel retain their own
   observations (recommended), or agents only?
3. **Refusal representation:** append `refusal_observed` pointing to the
   verified signed result (recommended), or retain only the signed result with
   no observation record?
4. **Naming and path:** use `observations.jsonl` beside each existing local
   role store and a shared `ObservationRecord` type (recommended), or choose a
   different record/file name and path before implementation?

## Implications for open work

- `lubav.4` can proceed after these four decisions are locked.
- The chosen record must not introduce a top-level PromiseGrid action kind; it
  is local durable evidence about observation.
- `lubav.5` can later document the resulting kernel unknown-pCID policy.
- `lubav.6` must add deterministic tests for each retained observation path.
