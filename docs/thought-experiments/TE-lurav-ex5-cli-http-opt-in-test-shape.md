# ex5 CLI HTTP opt-in test shape

TE ID: TE-lurav
## Status
decided

## Decision under test

How `ex5` should cover the user-facing CLI `-socket=off` behavior in tests now
that the transport rule is: fail closed on local socket loss unless the
operator explicitly opts into HTTP compatibility transport.

The behavior decision itself is already locked by TODO `127`. The remaining
question is what test shape best protects that user-facing contract without
muddying the terminal embodiment boundary.

## Assumptions

- The shipped behavior stays the same: CLI prefers the direct local Unix
  socket and fails closed by default.
- `-socket=off` remains the explicit operator escape hatch into HTTP
  compatibility transport.
- Existing transport-level tests already prove that an empty `SocketPath`
  yields HTTP behavior and that a missing socket fails closed.
- Alice is the operator invoking the CLI from a shell.
- Bob is the runtime owner exposing `/api/meta` and the local socket.

## Alternatives

### Alternative A: factor CLI startup resolution into a pure helper and unit test it

Move the `-socket=off` normalization and startup transport resolution into a
small helper that tests can call directly.

### Alternative B: add a subprocess-style integration test around actual flag parsing

Spawn the CLI entrypoint with argv containing `-socket=off` and assert the
resulting transport choice through a black-box process test.

### Alternative C: keep only the current transport-object tests

Rely on the existing tests that exercise `SocketPath=""` and missing socket
paths, without adding any direct coverage for the actual flag semantics.

## Scenario analysis

### Scenario 1: guarding the exact operator-facing CLI contract

Alice types `oks-cli -socket=off dashboard`.

Alternative A covers the exact normalization rule in-process. It proves the
flag value is mapped into explicit HTTP compatibility mode.

Alternative B also covers the exact operator-facing contract and is the most
black-box shape.

Alternative C does not actually protect the flag contract itself. The flag
could drift while transport-object tests still pass.

### Scenario 2: maintenance cost and determinism

Bob changes unrelated CLI startup logic later.

Alternative A is focused and deterministic. It isolates the startup decision
surface without requiring subprocess orchestration.

Alternative B is heavier and more brittle. Process tests are closer to the
shell contract, but they create more harness cost for a very small decision.

Alternative C is easiest short-term, but it leaves the specific user contract
under-protected.

### Scenario 3: PromiseGrid embodiment boundary clarity

Carol wants the tests to reinforce the embodiment contract rather than obscure
it.

Alternative A keeps the boundary explicit: startup resolves which embodiment
transport to use, then the existing transport tests cover the chosen path.

Alternative B also preserves the boundary, but it mixes CLI process mechanics
with transport semantics in one heavier test.

Alternative C leaves the operator-visible embodiment selection under-specified
in tests.

### Scenario 4: long-horizon evolution

Dave later adds another explicit transport-mode flag or startup rule.

Alternative A scales well because the startup resolution logic has one testable
home.

Alternative B scales more slowly because each behavior change pushes more work
into subprocess harnessing.

Alternative C keeps the code simple today but gives less confidence when the
startup contract grows.

## Conclusions

Rejected:

- Alternative C. It leaves the exact `-socket=off` contract uncovered.

Surviving:

- Alternative A: extract a pure startup-resolution helper and unit test it
- Alternative B: add a subprocess-style integration test around actual flag parsing

## Implications for TODOs and pending DIs

- TODO `128` (`rutav`) is locked to Alternative `A`.
- The implementation should keep startup transport resolution as a small pure
  helper, then test `-socket=off` directly there while leaving socket/HTTP
  transport-path tests in their existing narrower scope.
