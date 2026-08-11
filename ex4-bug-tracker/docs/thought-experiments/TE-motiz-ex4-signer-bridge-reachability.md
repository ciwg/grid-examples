# Ex4 signer-bridge reachability

TE ID: TE-motiz

## Status

needs DF

## Decision under test

Should Ex4's signer bridge be same-machine only, or network-reachable for browser embodiments on other machines?

## Alternatives

### A. Network-reachable service bridge with browser-local keys

The browser calls prepare/finalize over the normal service transport from any reachable machine. Its non-extractable key remains in that browser profile. The service applies only its own enrollment and acceptance policy.

### B. Same-machine-only bridge

Restrict bridge requests to loopback/local execution.

## Scenario analysis

Alice operates the service and Bob opens the browser on another machine. A lets Bob obtain canonical signable bytes, sign locally, and return proof without exposing a private key. B prevents decentralized/multi-machine use for no protocol reason.

Mallory can reach the service under A, but still cannot create an accepted promise without an enrolled public key and valid proof; service acceptance is local assessment, not global authority. Under B, network reachability is blocked rather than correctly authenticated and verified.

Over time, A permits independent browser embodiments and later peer exchange of exact envelopes. B makes the browser feature a same-host convenience and conflicts with Ex4's selected multi-machine PromiseGrid direction.

## Conclusion

B is rejected. A is required: local key custody does not imply local-only transport. This supersedes DI-tigid's local-only bridge constraint.

## Decisions still requiring DF

1. Use A: network-reachable prepare/finalize service endpoints, browser-local keys, and service-scoped enrollment/acceptance policy.
