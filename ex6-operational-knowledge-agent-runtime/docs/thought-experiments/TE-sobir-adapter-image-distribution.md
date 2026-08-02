# Adapter image distribution trust boundary

TE ID: TE-sobir

## Status

decided

## Decision under test

How a receiving EX6 node obtains the exact Docker worker image required by an
already transferred workflow artifact, while preserving independent local
activation and the worker confinement boundary.

This extends the workflow-transfer boundary in TODO `sibok` / DI-novuk and the
Docker-worker boundary in TODO `puvok` / DI-fofuh.

## Assumptions and trust model

- Alice and Bob have separate runtime roots, CAS directories, Docker daemons,
  and peer keys.
- Alice transfers an exact workflow artifact and signed lifecycle evidence to
  Bob through the existing dedicated workflow relay endpoint.
- A package manifest binds an adapter name and pCID pair to an immutable image
  reference. A mutable tag is not an executable dependency.
- A Docker worker remains untrusted code: no runtime-root, CAS, history,
  peer-key, secret, network, Docker-socket, or host-process access.
- Bob decides independently whether to fetch an image, install a package,
  activate an artifact, and execute it. Alice's lifecycle evidence is never
  authority for those decisions.
- The image bytes may be large; routine relay batches must remain record
  carriage and must not acquire an unbounded image payload.

## Alternatives

1. **Digest-pinned OCI registry.** The manifest uses a registry-qualified
   `name@sha256:<digest>` reference. Bob fetches it through the local Docker
   daemon, verifies Docker resolved the requested digest, then treats it as a
   locally available dependency.
2. **Workflow-relay image bundle.** The workflow-transfer endpoint carries an
   additional OCI image archive, addressed by its digest, beside the artifact
   and lifecycle evidence.
3. **Manual local preload.** Bob obtains the image out of band and manually
   updates a local package declaration to a Docker image ID before execution.

## Scenario analysis

### Normal operation

Alice captures and transfers a procedure workflow. Under option 1, Bob's
operator approves a package whose image field contains an OCI digest. Bob's
runtime asks Docker to pull exactly that digest and verifies the local result
before the existing readiness check permits execution. Artifact transfer and
image acquisition are separate: receipt stays inactive.

Option 2 lets Bob receive artifact and image in one relay request. It is
convenient for air-gapped peers, but turns a narrow artifact endpoint into a
large-content distributor and requires bounded streaming, archive validation,
image import semantics, and durable image-retention policy.

Option 3 is what the current local-image-ID example effectively requires. It
works for a single developer machine, but does not let a transferred artifact
become reproducibly executable on a second node.

### Failure, corruption, and incomplete acquisition

With option 1, a failed pull leaves Bob's workflow inactive or unavailable;
the requested digest stays visible as an unmet dependency. Docker's content
verification and the runtime's post-pull digest comparison prevent a tag from
silently satisfying the requirement.

Option 2 must independently validate every blob and manifest in the archive,
handle interrupted streams without retaining a partial executable image, and
prevent archive expansion from exhausting local storage. Its failure surface is
materially larger than the current workflow-transfer bundle.

Option 3 cannot distinguish an incorrect local rebuild from the intended
worker without an operator comparing image IDs manually.

### Concurrent actors and mixed versions

Alice may upgrade an adapter while Bob still has an older image. Option 1
allows both immutable digests to coexist; the artifact manifest decides which
one is required. Bob can retain the old image while retained artifacts still
need it, then garbage-collect by local policy.

Option 2 needs equivalent multi-image retention plus relay quota and resumable
transfer behavior. Option 3 encourages replacing a mutable local tag, which
can make a retained artifact unexpectedly run different code.

### Long-horizon evolution and scale

Registry distribution is an established content-addressed transport with cache
layers, resumable blobs, and independent availability. EX6 only needs to bind
the required digest to the active package and expose readiness; it does not
need to invent an image protocol.

Relay bundles improve self-contained transfer but force every PromiseGrid relay
operator to provision image-bandwidth, storage, and abuse controls. They are a
possible later air-gap extension after the registry path proves the contract.

Manual preload keeps EX6 small but leaves the key cross-node execution claim
outside the system and cannot provide an operator-readable readiness result.

## Conclusions

Reject option 3 for cross-node execution: it remains a development-only
workaround, not a portable execution dependency.

Defer option 2: it is appropriate only after a separate bounded OCI-archive
transport thought experiment and explicit storage/abuse policy.

Option 1 survives. It keeps workflow relay focused on exact workflow evidence,
uses an existing digest-verifying image transport, and leaves every fetch,
package installation, activation, and execution decision local to Bob.

## Decision status

Locked by DI-harib. The first EX6 slice uses a digest-pinned OCI registry
dependency and permits acquisition only from an operator-configured registry
allow-list. A package manifest cannot authorize an arbitrary registry hostname.

## Implications for open work

- TODO `puvok`: add a package adapter image-availability/readiness state and
  verify Docker resolves the exact registry digest before execution.
- TODO `sibok`: retain the existing workflow endpoint as artifact/evidence
  carriage only; do not add image bytes to it.
- Future work: define registry hostname policy, authentication/credential
  boundary, offline cache behavior, and image garbage collection in a separate
  decision before Docker pull behavior is implemented.
