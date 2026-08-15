# Chromium Snap support boundary

TE ID: TE-tilun

## Status

decided

## Decision under test

Whether `ex5-operational-knowledge-system` should add supported Chromium Snap
browser embodiment setup, support a non-Snap Chromium distribution instead, or
remain Chrome-only while retaining its Chromium-capable extension code.

This TE completes TODO `lavem.3` analysis after the recorded Chromium Snap
reproduction in TODO `lavem`.

## Assumptions and trust boundary

- Alice uses the verified browser embodiment against one local Ex5 runtime.
- Bob can use the working Google Chrome setup today.
- Chromium `151.0.7922.108 snap` loads the shipped unpacked Ex5 extension but
  the page handshake reports `Specified native messaging host not found`.
- The extension, page bridge, native host binary, and local Unix socket are
  separate hops. A browser loading the extension is not evidence that native
  messaging is available.
- A native host is local executable code. Its installation and discovery path
  are part of the browser embodiment's trust boundary and must be explicit.
- Ex5 must remain fail-closed: it may not fall back silently to the old HTTP
  browser path when native messaging is unavailable. Source: DI-bahak.

## Alternatives

### A. Retain verified Google Chrome support; defer Chromium Snap

Keep the shipped Chrome setup, the local regression harness, and the
attach-only external evidence as the only supported browser path. Retain the
Chromium-capable extension code but do not advertise Chromium support.

### B. Add a separate Chromium Snap installation and verification path

Create a Snap-specific installer that places the native-host manifest and host
binary where the confined Chromium build can discover and execute them, then
add a separate passing Chromium verification harness and documentation.

### C. Support a non-Snap Chromium distribution

Support Chromium only when installed through a distribution that follows the
ordinary Chromium native-messaging registry contract. The Snap package remains
unsupported.

## Scenario analysis

### Normal operation

Alice starts Ex5 in Chrome today. Alternative A is already verified end to end:
the extension is loaded by Chrome's supported DevTools API, the native host is
registered, and visible review interactions traverse the native-host route.

Alternative B would need a new exact Snap installation contract before it
could make the same claim. The completed reproduction shows that merely adding
the manifest to disposable profile and apparent Snap directories is not that
contract. A host path that cannot be discovered or executed is not a supported
embodiment.

Alternative C can use the ordinary Chromium registry model, but adds a second
browser distribution requirement to the example and does not help Alice using
the installed Snap package.

### Failure, corruption, and incomplete setup

Under A, a missing Chrome native host produces the existing explicit failure;
the supported setup and harness make recovery repeatable.

Under B, every registration, packaging, permission, and host-path failure must
be separately diagnosed. A partial installer must leave no stale broad host
authorization and must fail before the browser claims readiness.

Under C, package differences move the operational burden to installation
documentation and risk users assuming their Snap Chromium is covered when it
is not.

### Concurrent actors and mixed environments

Bob's verified Chrome profile and Alice's unsupported Chromium Snap profile
can share Ex5 data only if each browser's embodiment claim is explicit. A does
that cleanly. B introduces a second local installer and evidence set that must
not alter or terminate the Chrome demo. C similarly requires distribution
detection and clear support reporting.

### Long-horizon evolution and migration

A is stable and low-maintenance, but defers a user choice. B creates a
Snap-specific compatibility commitment tied to Snap interfaces, host paths,
and future Chromium updates. C creates a distribution-specific support matrix
instead of one browser contract.

### Trust and scale

All options keep the promise semantics unchanged: browser carriage remains a
local embodiment adapter and no new PromiseGrid action is introduced. B has
the broadest local trust surface because it requires a confined browser,
installer, host executable, and native-messaging registration to cooperate.
At scale, A has one documented setup and one regression surface; B and C each
need their own setup, diagnostics, test harness, and release support.

## Conclusions

The observed fault is neither Ex5 runtime behavior nor unpacked-extension
loading. It is Chromium Snap native-host discovery/confinement. No current
evidence establishes a safe Snap-specific registration contract, so B is not
ready to implement. C is plausible but only if a concrete non-Snap Chromium
distribution is placed in scope.

The surviving alternatives are A and C. A is recommended for the present
example because it preserves a verified, evidence-backed browser embodiment
without inventing a fragile Snap installer.

## Decision status

Locked: Alternative A. Ex5 supports the verified Google Chrome embodiment and
retains Chromium Snap as explicitly deferred. Source: DI-vamid.

## Implications for TODOs and pending DIs

- TODO `lavem.3` needs DF before any Chromium support code, installer, or
  support claim changes.
- Under A, TODO `lavem` remains deferred and Ex5 continues to document only
  Google Chrome as supported.
- Under C, a follow-up TODO must name the exact supported Chromium
  distribution, installer paths, verification harness, and support boundary.
