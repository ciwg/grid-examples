# Dedicated workflow relay endpoint

TE ID: TE-novuk

## Status

decided

## Decision under test

How independently rooted PromiseGrid runtimes exchange an exact workflow
artifact together with the lifecycle envelope that describes the sender's
availability decision.

## Assumptions

- Alice and Bob have separate runtime roots and separate CAS directories.
- Workflow lifecycle envelopes are local decisions. A received envelope is
  evidence of Alice's decision, never Bob's decision.
- Ordinary relay batches remain JSON carriage for records. Large binary CAS
  objects must not be embedded in them.
- Existing peer registration and Ed25519 keys are the local trust root.

## Alternatives

1. Put artifact bytes and lifecycle bytes into a normal relay batch.
2. Require a shared CAS or shared filesystem.
3. Publish a dedicated, peer-authenticated workflow endpoint carrying one
   canonical transfer bundle.

## Scenario analysis

In normal operation, option 1 makes a batch grow with binary artifacts and
mixes record carriage with artifact retention. Option 2 avoids a transfer but
does not model independent nodes. Option 3 lets Alice send exact bytes to Bob,
who independently hashes, validates, and stores the artifact.

If a transfer is truncated or corrupt, options 1 and 3 can validate CIDs and
signatures before retaining bytes; a shared CAS can hide the failure behind the
shared store. With concurrent or mixed-version nodes, option 3 has an explicit
endpoint and bundle boundary, whereas a batch extension changes the routine
relay format for every node.

Over time, a sender's active/deactivated state is evidence, not authority for
the receiver. A lifecycle envelope placed in Bob's lifecycle CAS would be
replayed as Bob's local decision. Option 3 therefore retains received evidence
in a separate evidence CAS and requires Bob to make a new local import and
activation decision. This preserves the trust boundary at scale while adding a
small, bounded second CAS namespace.

## Conclusion

Choose option 3: a dedicated peer-authenticated workflow relay endpoint. The
transfer contains exact artifact bytes, exact lifecycle-envelope bytes, sender
identity, and a signature over the canonical transfer payload. Bob verifies the
sender against its allow-list, verifies the signature and both CIDs, stores the
artifact in its artifact CAS, and stores lifecycle evidence only in its
separate workflow-evidence CAS. Receipt never imports or activates a workflow.

## Decision status

Locked by DI-novuk.

## Implications for open work

The endpoint needs peer-card metadata, signed transfer codec, independent-node
tests, and operator documentation. It must not add a new top-level PromiseGrid
action: endpoint transport and pCID payload semantics remain implementation
mechanics under the existing `promise` vocabulary.
