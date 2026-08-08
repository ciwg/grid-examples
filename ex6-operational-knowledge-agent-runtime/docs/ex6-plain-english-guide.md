# EX6 in plain English

## The short version

EX6 is a small runtime for keeping operational knowledge. It is not yet a full
application with screens and buttons. Think of it as the reliable machinery
underneath an application: it records facts, checks known rules, keeps exact
files, and can exchange those files with another independent computer.

The important promise is simple:

> A workflow can be packaged exactly once, identified by its contents, sent to
> another node, verified there, and kept there without that receiving node
> silently deciding to run it.

Source: DI-lovek; DI-bavuk; DI-novuk.

## The three things EX6 keeps

### 1. Operational records

Packages such as `context`, `procedures`, and `runs` create durable records:
places, procedures, evidence, approvals, and so on. EX6 stores the exact bytes
and can relay them to a peer.

### 2. Workflow artifacts

A workflow is a directory of instructions, policies, examples, and a
`workflow.json` manifest. EX6 turns that directory into a deterministic tar
archive: the same valid directory always produces the same bytes.

Those bytes are stored in CAS (content-addressed storage). CAS gives the
archive a CID, like this:

```
bafkreiahdp34nto2rnnqde26jw3xnkd6xnlalnr72sug3w7tjb3bhhoj4q
```

That CID is the identity of the exact workflow version. If one character in
the workflow changes, the CID changes too.

### 3. Local lifecycle decisions

Each computer decides for itself whether it has merely imported a workflow,
made it active, deactivated it, or revoked it. Those decisions are append-only
lifecycle events in local CAS.

This distinction matters:

```
has the workflow bytes       !=       has decided to use the workflow
```

Receiving an artifact does not mean receiving authority to run it.

## Workflow loading

The workflow loader retains exact workflow artifacts in local CAS.

1. Pick up a workflow directory.
2. Put an exact archived copy into local CAS.
3. Give that copy a CID.
4. Optionally mark it active on this computer.
5. Take out the exact copy later for inspection.

Useful commands:

```bash
go run ./cmd/moks workflow demo procedure-execution
go run ./cmd/moks workflow list
go run ./cmd/moks workflow status procedure-execution
go run ./cmd/moks workflow extract procedure-execution ./inspect-here
```

The `demo` command is the easiest local walkthrough. Its four `[ok]` lines
show capture, verification, explicit activation, and exact extraction.

## What the relay is for

The relay lets two separate EX6 runtimes communicate. Each runtime has:

- a separate local state directory;
- separate CAS storage;
- its own peer identity and signing key;
- its own trust choices;
- its own workflow lifecycle decisions.

When Alice sends a workflow to Bob:

```
Alice workflow directory
        ↓ deterministic archive
Alice CAS artifact CID
        ↓ signed relay transfer
Bob verifies signature and CID
        ↓
Bob CAS stores exact artifact bytes
        ↓
Bob stores Alice's lifecycle statement as evidence
        ↓
Bob still has no locally active workflow
```

Bob must explicitly import and activate the artifact if Bob chooses to use it.
That is deliberate: Alice cannot remotely turn on software on Bob's node just
by sending it.

## The public peer-card page you saw

The JSON page is not a user interface. It is a machine-readable business card
for another EX6 node.

```json
{
  "peer_id": "...",
  "public_key": "...",
  "workflow_import_url": "..."
}
```

It tells another node:

- who this node claims to be;
- which public key to use to verify its signatures;
- where to send normal records;
- where to send a workflow artifact.

For a human demo, say:

> “This is the node's public identity and inbox address. Another independent
> node can discover it, verify who it is, and send it an exact workflow
> artifact. This node still decides locally whether to use that workflow.”

## What EX6 does not do yet

EX6 executes an active workflow only after an operator explicitly starts it;
it does not automatically start imported artifacts or provide a finished
human-facing workflow screen. The current trusted built-in adapters invoke
declared package commands and retain the lifecycle in CAS. It is still proving
the safer lower-level parts first:

- exact workflow packaging;
- reproducible identity with CIDs;
- durable history;
- validation;
- explicit local activation;
- independently verified node-to-node transfer.

Docker-confined workers are now an available execution backend alongside the
explicitly selected built-in adapter path. A package declares its adapter in
`moks-package.json`; the artifact and package must agree on its name and pCID
contracts before it can run. The locked-down worker cannot see the runtime's
files, network, secrets, or keys, and it can only propose work for the runtime
to validate and save. EX6 now ships source for one local Docker image, pinned
in its package declaration by immutable image ID, which performs procedure
execution.
Source: `DI-lumek`; `DI-fofuh`; `TE-dovek`.

That is valuable because a user interface can change, but those trust and
evidence boundaries must stay correct underneath it.

## A simple boss-meeting walkthrough

1. Run `workflow demo procedure-execution`.
2. Explain: “The workflow directory became an immutable artifact with a CID.”
3. Show the four `[ok]` lines.
4. Explain: “Activation is an explicit local decision; loading does not run
   anything.”
5. Show the peer card only if discussing decentralized exchange.
6. Explain: “A remote node can receive the exact artifact, but it cannot be
   remotely activated.”

## The one sentence to remember

EX6 makes operational workflows portable and verifiable without letting a
sender take control of the receiver.
