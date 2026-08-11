# Ex4 signer-bridge adapter carriage

TE ID: TE-nujup

## Status

decided

## Decision under test

How browser and CLI embodiments request server-produced canonical signing bytes without making a browser CBOR implementation part of Ex4's signed wire protocol.

## Assumptions and trust model

Alice's browser and Bob's CLI retain their own private keys. The network-reachable Ex4 service may validate service-scoped enrollment, construct canonical bytes, and verify proofs; it must never receive a private key or sign for either embodiment. The PromiseGrid artifact remains `grid([42(pCID), payload, proof])`; pCID is binary tag 42 in that artifact.

## Alternatives

1. Raw-CBOR bridge requests containing pCIDs.
2. A browser CBOR encoder.
3. JSON adapter requests carrying local profile names and payload fields; the service resolves the embedded profile source to its pCID and produces the sole canonical CBOR artifact.

## Scenario analysis

In normal multi-machine use, alternative 1 requires Alice's browser to encode the raw-CBOR request before it can obtain signable bytes. That recreates the canonical-CBOR compatibility problem that the bridge was selected to avoid. Alternative 2 makes the browser a second canonical encoder, with ongoing compatibility and migration obligations. Alternative 3 lets Alice send ordinary HTTPS adapter data, receive exact signing bytes, sign them locally, and submit the unchanged finalized CBOR envelope.

If a malformed request arrives, alternatives 1 and 2 require a second CBOR decoder/encoder boundary in the browser. Alternative 3 rejects ordinary bounded JSON before any artifact exists; the service records a bounded observation. If a draft is lost or expires, no promise was created and Alice can prepare again. Concurrent and mixed-version embodiments remain interoperable at the final envelope boundary because the service derives the pCID from the embedded profile source rather than accepting text or untagged pCID selectors.

Over time, alternative 3 keeps HTTP adapter fields local and replaceable while preserving the stable signed artifact. It adds a server-side short-lived draft registry and explicit expiry, but avoids browser CBOR maintenance and false claims that JSON is a PromiseGrid artifact.

## Conclusion

Use bounded JSON only for prepare/finalize adapter carriage. The browser and CLI send a local profile name, never a pCID. The service selects the pCID from its embedded source, produces exact tag-42 signing bytes, verifies the returned proof, and returns unchanged finalized envelope bytes. This supersedes the raw-CBOR bridge-request portion of DI-tofuf and DI-gugah, not the raw-CBOR final `/api/promises` envelope boundary.

## Implications and future work

`kakon.3` implements the expiring in-memory draft registry, browser and CLI adapters, rejection observations, and regression tests. `kakon.4` documents the adapter and final artifact boundary.
