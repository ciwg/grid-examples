# Local Recognition-Policy Configuration Embodiment

TE ID: TE-bilad

## Status

decided

## Decision under test

How an Ex7 operator supplies the local public-key recognition policy required
by DI-piruf without turning that configuration into a participant credential,
membership protocol, browser capability, or shared authority.

## Assumptions and scope

- An Ex7 record's signature is author evidence; the agent still needs local
  policy to decide whether a public key and label are recognized for its own
  projection.
- This embodiment configures public information only: label, Ed25519 public
  key, and optional local roles. It never accepts, creates, stores, exports, or
  uses a participant private key.
- Alice may have an independently held key. An operator may choose to
  recognize Alice's public key in one Ex7 agent view while Bob chooses not to
  recognize it in another. Neither configuration changes Alice's records.
- Mallory may alter a config file, submit malformed keys, claim a recognized
  label in a record, or submit unknown-family bytes. A malformed policy must
  fail startup rather than silently widening recognition.
- Key continuity, revocation, portable membership, accounts, and browser
  administration remain separate work.

## Alternatives

### A. Agent-local public recognition file

Use `<runtime-root>/recognition.json`, a mode-`0600` JSON configuration file
provided by the local operator. It contains an explicit version and an array
of `{label, ed25519_public_key_base64}` entries. Startup parses every key,
rejects duplicates, empty labels, malformed base64, non-Ed25519 key lengths,
and policy/record label mismatches. The file is read-only to the running agent;
there is no HTTP or browser mutation route.

### B. Repeated command-line recognition flags

Require one `-recognize label:base64-public-key` flag per recognized key on
each agent start. Nothing is persisted beside ordinary shell history or service
configuration.

### C. Browser/account policy editor

Let an account-authenticated browser page add, remove, or approve recognized
keys through the running HTTP service.

## Scenario analysis

### Normal operation

Under A, an operator places Alice's and Carol's public keys in the agent's
private local configuration before launch. The agent verifies records against
those exact keys; its existing local area-role data decides whether a
recognized signer may clear a hold. The browser can show the result but cannot
change recognition.

B can produce the same immediate behavior but makes a restart-sensitive
security boundary easy to omit or misquote. C turns a web session into a
recognition-management capability whose authentication and audit contract has
not been designed.

### Failure and corruption

With A, a missing file means an empty policy only when the operator explicitly
chooses that bootstrap mode; a present malformed file fails startup unchanged.
The agent never falls back to label-only recognition. A write interruption is
an operator-visible configuration failure rather than an automatic repair.

B fails through lost flags or shell history exposure. C adds request races,
session compromise, and an unplanned audit/membership surface.

### Different participants and long-horizon change

A allows Alice's and Bob's agents to recognize different keys locally. An
operator may replace the file during a controlled restart, while existing
records remain exact historical bytes. A future key-continuity protocol may
inform a later policy loader without changing this first embodiment.

B repeats configuration at every service manager, container, and restart. C
invites an unplanned shared roster semantics and leaves key changes tied to a
web account system.

### Trust and operational scale

A is deliberately modest: it is a local operator input, not a promise record
or claim about all agents. Mode `0600` reduces casual exposure even though the
contents are public keys. Its obligations are explicit parsing, fail-closed
startup, deterministic tests, and documentation of the file lifecycle.

B is smaller in code but weakens repeatability. C is convenient-looking but
would need its own browser administration, account, and durable audit
protocols before it could be trustworthy.

## Conclusions

B is rejected because a security-relevant recognition boundary must survive
restart as an inspectable agent-local artifact. C is rejected because it would
make an unplanned browser/account control surface responsible for a trust
decision.

A survives and is recommended: a mode-`0600`, agent-local, public-key-only
recognition file with no browser mutation route.

## Output to decision framing

**A — load `<runtime-root>/recognition.json` as a versioned, public-key-only
agent-local recognition configuration; fail closed on malformed content; keep
the browser and account surfaces read-only with respect to recognition.**

## Decision status

locked: Alternative A, agent-local public recognition file, by DI-likoh in
`../../TODO/TODO-bubuz-canonical-makerspace-records.md`.
