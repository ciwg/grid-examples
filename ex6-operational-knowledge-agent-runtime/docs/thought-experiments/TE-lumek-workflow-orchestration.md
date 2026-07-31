# Workflow Orchestration

TE ID: TE-lumek

## Status

decided

## Decision under test

How EX6 should let independently loaded workflow artifacts exchange durable,
typed handoffs through the shared runtime.

## Assumptions

Alice operates one runtime with the seven shipped workflow artifacts active.
The first execution backend is trusted built-in package adapters; Docker workers
remain a later, separately contained backend under TE-dovek. Mallory can submit
malformed handoff bytes or stale identifiers, but cannot obtain runtime write
authority except through validated runtime APIs. Workflow artifacts remain
content-addressed CAS archives.

## Alternatives

1. Keep manual package commands and document coordination.
2. Parse package command text to infer workflow output.
3. Add pCID-selected CBOR handoffs, durable CAS run events, and typed built-in
   operation adapters.

## Scenario analysis

### Normal operation

Alice starts a procedure-execution run. Alternative 1 leaves the next operator
to remember which command and workflow follows. Alternative 2 makes output
formatting an accidental API. Alternative 3 records input, output, and every
state transition, then offers an explicit handoff or policy-selected target.

### Malformed or incomplete input

Carol supplies invalid CBOR or omits an adapter-required field. Alternatives 1
and 2 either fail late or have no durable explanation. Alternative 3 validates
the pCID envelope before execution and retains `waiting-for-input` or `failed`
state without altering the source workflow's history.

### Mixed workflows and concurrent operators

Alice and Bob may run distinct workflows from the same runtime. Alternative 1
has no run identity. Alternative 2 couples every workflow to presentation text.
Alternative 3 roots each run in a CAS event CID and lets a local policy match
`(source workflow, output pCID)` to a target workflow/input pCID.

### Restart, corruption, and migration

After a restart or deleted cache, Alternatives 1 and 2 cannot reconstruct
coordination. Alternative 3 scans CAS pCID-selected run events and rebuilds its
disposable cache. Invalid objects are ignored during scanning; an invalid new
request is rejected before persistence.

### Trust boundary and future workers

Built-ins have the runtime's trusted adapter boundary. Docker workers cannot be
silently granted that trust: they must later receive the same CBOR contract and
return runtime-validated output under TE-dovek. This keeps the run ledger and
handoff semantics stable while execution backends evolve.

### Scale

CAS replay costs a full scan at open, favoring correctness and auditability for
this first release. The cache is an acceleration/readability artifact only.
Handoff payloads are bounded and deterministic; later indexing can optimize
scan cost without changing durable event meaning.

## Conclusion

Alternative 3 is selected. EX6 adds pCID-selected CBOR input/output envelopes,
CAS-backed run lifecycle events, trusted built-in adapters for all seven shipped
workflows, explicit and policy-selected handoffs, and no automatic retry.
Missing compatible input is `waiting-for-input`; adapter failures are durable
`failed` states. All-to-all handoff requests are allowed, but target adapters
validate their own accepted input before executing.

## Implications

DI-lumek records the locked implementation, naming, and paths. DI-sovuk is
superseded only as to its no-engine constraint; its dispatcher scenario remains
useful coverage. Docker worker dispatch stays deferred to TE-dovek.

## Decision status

Locked by DI-lumek.
