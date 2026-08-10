# Ex1 published-claim regression coverage

TE ID: TE-pokis

## Status

decided

## Decision status

Locked by DI-potoj: use layered hybrid coverage, derive profile pCIDs from the
source specs and compare them with the published inventory and scope record,
keep one E2E refusal path and one E2E timeout path, and place focused tests
beside their respective packages.

## Decision under test

How should Ex1 regression coverage prove the claims it publishes: that the
five local profile documents match their derived pCIDs, invalid proof or
capability input is locally rejected, a signed refusal is distinguished from a
promise assessment, a timeout is retained as observer-local evidence, and raw
bytes remain available for audit and replay?

This TE is the prerequisite for `lubav.6`.

## Assumptions and trust model

- Ex1 is a local-draft example, not a claim of frozen upstream-spec
  conformance. Its five profile documents and their pCIDs are published in the
  design inventory and implementation-scope declaration.
- Every valid application envelope is signed. Capability-token meaning remains
  pCID-owned payload semantics, and malformed or unauthorized input does not
  require a signed response. Source: DI-garis.
- `observations.jsonl` is append-only local evidence. A timeout, invalid
  input, or a kernel dispatch outcome records what its observer saw; it does
  not establish another agent's intent. Source: DI-vihoz, DI-riguz, and
  DI-purum.
- The kernel retains and routes exact bytes by pCID, records
  `no_registered_recipient` when no current local registration matches, and
  does not decide global pCID validity. Source: DI-zosiz.
- Alice, Bob, Carol, Dave, and Ellen are cooperative local roles. Mallory may
  send malformed bytes, altered signatures, expired or mismatched capability
  tokens, and envelopes for absent recipients.
- Tests must be deterministic, use temporary test roots, avoid network calls
  outside the loopback harness, and preserve the existing role-level and E2E
  suite boundaries.

## Alternatives

### A. Focused package tests only

Add deterministic tests next to the artifact, agent, protocol, token, and
kernel packages. Each test builds the smallest input that exercises one local
decision, then inspects raw artifacts or one observation log.

### B. End-to-end scenario tests only

Extend the E2E harness with fixtures for altered proof, invalid capability,
refusal, timeout, and retained-artifact inspection. Each test starts the full
local topology and assesses the resulting role directories.

### C. Layered hybrid coverage

Use focused package tests for precise rejection, pCID-inventory, raw-retention,
and observation-record assertions. Retain a small number of E2E scenarios for
the refusal and timeout paths, where the distinction depends on several roles
and persisted artifacts together. Do not duplicate every rejection case at
both layers.

## Scenario analysis

### Normal operation and profile evolution

Alice updates a local profile document but forgets to update the published
inventory. A focused profile-inventory test under A can derive the document
pCID and point directly to the stale entry. B can discover that normal routing
still works, but it may never establish that the human-facing inventory is
wrong. C detects the mismatch precisely while keeping a normal E2E path as a
separate integration check.

When a later profile revision deliberately changes a pCID, C makes the
maintenance obligation explicit: update the profile, inventory, scope
publication, and one focused assertion together. It does not make the kernel
an allow-list or turn a test fixture into a global protocol registry.

### Malformed input, invalid proof, and invalid capability

Mallory sends unparsable bytes, then a structurally valid envelope whose proof
does not verify, then a message bearing a capability that is expired or has the
wrong issuer or audience. Under A, tests can assert the exact receiving
boundary, retained raw CID, local observation kind, and absence of a signed
response without timing a full topology. B demonstrates the same behavior only
through more setup and can make a failure look like a routing race. C uses
focused tests for those distinct validation conditions and reserves E2E for
claims that need role interaction.

### Signed refusal versus local promise assessment

Carol sends Bob a valid, signed `pick_pack` or `accounting` result with
`status = refused`. The test must prove that Bob retains the signed result and
records `refusal_observed`; it must not call the refusal proof that Carol broke
or withdrew any promise. A can inspect the local handler result but can miss a
broken caller-to-handler artifact chain. B exercises that chain but is slower
and less exact about which layer made the assertion. C keeps an E2E refusal
scenario that checks the final business result plus Bob's local evidence, while
focused tests protect the record semantics.

### Timeout, incomplete writes, and restart audit

Ellen does not return a shipment result before Bob's deadline. A can reliably
construct or invoke the deadline path and inspect a `timeout_observed` record
whose `ExpectedCID` names the request; it must not assert that Ellen refused.
B tests the full timeout topology, but timing introduces flake risk and it does
not naturally test incomplete observation writes. C uses a bounded E2E timeout
scenario for the cross-role path, plus focused artifact tests that validate
append-only JSONL and retained raw-byte lookup after reopening a store.

### Concurrent actors and mixed-version nodes

Alice and Bob may independently retain different observations of the same
traffic. A verifies individual store behavior but cannot prove that role
directories remain distinct through the application topology. B can cover that
separation but should not multiply every malformed-input test across every
role. C uses one E2E assertion that role-local evidence is not centralized and
focused tests for per-store append behavior.

A mixed-version node may register a future pCID or be absent altogether. The
existing kernel regression establishes the local `no_registered_recipient`
policy. C leaves that test focused at the kernel boundary, avoiding a brittle
claim that current Ex1 knows every valid future pCID.

### Scale and maintenance cost

A is fast and localizes failures, but it can leave an incorrect wiring path
undetected. B has broad confidence in each scenario but increases suite time,
network timing sensitivity, and fixture maintenance. C gives the highest
signal per test: low-level claims fail close to their source, while only
cross-role promises and evidence chains pay the E2E cost. Its obligation is to
document which claim belongs to which layer so that redundant tests do not
silently drift.

## Conclusions

A is rejected as the sole strategy because the published refusal and timeout
claims depend on a real multi-role chain and persisted evidence at more than
one boundary. B is rejected as the sole strategy because it obscures precise
profile, proof, capability, raw-retention, and append-only-record failures
behind an expensive topology.

C survives and is recommended: add focused deterministic tests for profile
inventory, invalid proof/capability handling, raw-byte retention, and local
observation records; add or extend only the E2E refusal and timeout scenarios
needed to prove cross-role behavior. Assertions must say what the local
observer retained or assessed, not infer another agent's intent.

## Decisions still requiring DF

1. **Test strategy:** use layered hybrid coverage (recommended), focused
   package tests only, or E2E scenarios only?
2. **Inventory authority:** derive each pCID from its `specdocs/*.md` source
   and compare it with the design inventory and implementation-scope
   declaration (recommended), or test only the runtime `protocol` constants?
3. **E2E scope:** retain one refusal scenario and one timeout scenario with
   evidence assertions (recommended), or add E2E coverage for every rejection
   class too?
4. **Test placement:** add focused tests beside their respective packages and
   extend `ex1-order-flow/e2e/e2e_test.go` only for cross-role scenarios
   (recommended), or place all new coverage in a separate integration package?

## Implications for open work

- `lubav.6` can proceed after the four DFs are locked and the corresponding DI
  is recorded in `TODO/TODO-lubav-order-flow-protocol-consolidation.md`.
- The chosen tests must use temporary test roots and clean them through the Go
  test lifecycle; no production runtime path is introduced.
- `lubav.7` should link readers to the resulting reproducible verification
  commands and evidence-inspection expectations without representing the local
  draft as a frozen upstream specification.
