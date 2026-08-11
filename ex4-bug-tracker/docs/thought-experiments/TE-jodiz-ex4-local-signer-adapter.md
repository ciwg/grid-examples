# Ex4 local signer adapter

TE ID: TE-jodiz

## Status

needs DF

## Decision under test

How can a browser-held private key produce a signer-owned Ex4 promise when the
only verified canonical-CBOR envelope implementation is Go?

## Alternatives

### A. Local signer bridge with browser key custody

The browser keeps its non-extractable key and sends a signing request to a
narrow local bridge. The Go bridge constructs canonical bytes and returns the
exact signed envelope; the browser submits it unchanged.

### B. Browser-local CBOR implementation

The browser independently builds and signs canonical CBOR envelopes.

### C. Server signs after browser HTTP input

The server constructs and signs the promise from an unsigned browser command.

## Scenario analysis

Alice runs the local service. Under A, Bob's browser key remains browser-held,
but the bridge must have a defined signing interface. Under B, cross-language
canonical-byte compatibility becomes a permanent proof obligation; the failed
`cbor-x` probe demonstrates that ordinary CBOR conformance is insufficient.
Under C, the service signs on Bob's behalf and violates DI-muzal.

For Mallory, A keeps the signer boundary explicit and lets the server verify
the returned exact bytes. B risks different but semantically equivalent bytes
that fail signature verification. C centralizes every private key and loses
embodiment provenance.

At restart and over long evolution, A keeps one Go canonical implementation
while preserving browser key custody. It creates a bridge trust boundary and
requires explicit local-only documentation. B creates duplicated canonical
encoding maintenance. C creates a false client-promise claim.

## Conclusion

C is rejected. B is rejected for this slice because its canonical-byte proof
failed. A survives, but it still needs DF for how the bridge asks the browser
to sign canonical Go-built bytes without exporting the private key.

## Decisions still requiring DF

1. Define whether the bridge returns Go-built signable bytes to the browser for
WebCrypto signing, or whether the browser uses a separate local signing agent.
2. Define the local-only bridge transport and browser enrollment binding.
3. Define whether the bridge is in scope for Ex4 now or the browser adapter is
deferred while CLI signing proceeds.
