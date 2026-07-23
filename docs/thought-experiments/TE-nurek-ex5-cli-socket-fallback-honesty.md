# ex5 CLI socket fallback honesty

TE ID: TE-nurek
## Status
decided

## Decision under test

How the `ex5` CLI should behave when the runtime advertises a direct local
Unix-socket embodiment contract but the local socket path is unavailable or
unreachable.

The current implementation silently demotes to HTTP compatibility transport on
socket miss or dial failure. The review question is whether that silent
demotion is still the right PromiseGrid-aligned behavior now that the direct
local socket contract is part of the shipped embodiment model.

## Assumptions

- Browser remains on the local HTTP adapter.
- CLI remains a terminal embodiment that should prefer the direct local socket
  contract when available.
- HTTP compatibility transport still exists and remains useful for mixed local
  deployments and recovery.
- Alice is a normal operator using the CLI intentionally.
- Bob is another operator starting the runtime with the expected local socket.
- Mallory is not the central actor here; the main risk is transport honesty
  and operator confusion, not signature forgery.

## Alternatives

### Alternative A: silent automatic HTTP fallback

Keep the current behavior. If the socket path is missing or dialing fails, the
CLI falls straight through to HTTP without surfacing a warning.

### Alternative B: explicit warning, then HTTP fallback

Keep HTTP compatibility fallback, but make the CLI tell the operator when the
preferred direct socket path was unavailable and the command is running over
HTTP instead.

### Alternative C: fail closed unless the user opts into HTTP explicitly

Treat local socket unavailability as an error for normal CLI calls. HTTP would
remain available only through an explicit flag or mode switch.

## Scenario analysis

### Scenario 1: normal local runtime with healthy socket

Alice runs the CLI on the same machine as the runtime and the socket exists.

Alternative A works.
Alternative B also works; no warning is emitted because the socket succeeds.
Alternative C also works; the socket succeeds and HTTP is not involved.

No alternative is meaningfully worse in the steady state.

### Scenario 2: runtime not started yet, but Alice still invokes the CLI

Alice runs a CLI command before starting `operational-knowledge`.

Alternative A falls through to HTTP. If the HTTP server is also absent, the
command still fails, but only after hiding the fact that the preferred socket
lane was missing.

Alternative B still gives Alice one truthful clue: the direct embodiment
contract was unavailable and the CLI is trying compatibility transport.

Alternative C is strictest. It fails immediately on the socket problem unless
Alice asked for HTTP explicitly. That is honest, but it is also harsher for
operators who just want the command to keep working against a reachable local
HTTP runtime.

### Scenario 3: custom or mixed local deployment

Bob runs a local ex5 runtime where the direct socket is unavailable for some
environmental reason, but the HTTP adapter is still intentionally reachable.

Alternative A preserves usability but hides the degradation.

Alternative B preserves usability and exposes the actual operating mode.

Alternative C treats the deployment as invalid for normal CLI use, even though
the compatibility lane still exists and may be a deliberate short-term choice.

### Scenario 4: PromiseGrid embodiment honesty

Carol reads the docs saying CLI prefers the direct local Unix-socket contract.

Alternative A is weakest. A command may actually run over HTTP with no visible
signal, so the embodiment contract is not operationally honest.

Alternative B is stronger. The CLI still prefers the socket contract, but
operators are told when compatibility transport is carrying the command.

Alternative C is strongest in purity terms because the CLI never pretends to
be on the direct contract when it is not.

### Scenario 5: migration and mixed-version local nodes

Dave upgrades one machine where some older wrappers or launchers still assume
HTTP-only operation.

Alternative A is easiest to roll through, but it keeps the ambiguity forever.

Alternative B is a staged tightening: existing flows keep working, but the
warning highlights where the environment is not on the preferred contract.

Alternative C is the sharpest cutover and creates the biggest migration
obligation because every mixed setup must be made explicit immediately.

### Scenario 6: long-horizon PromiseGrid direction

Ellen wants `ex5` to keep moving toward clearer embodiment contracts without
breaking the current shipped compatibility surface all at once.

Alternative A leaves the contract blurry.

Alternative B keeps the compatibility lane but makes the degradation explicit.
That supports future hardening later if the repo decides to make HTTP opt-in.

Alternative C is the most pure end state for terminal contract honesty, but it
conflates this alignment pass with a stronger product/runtime policy shift.

## Conclusions

Rejected:

- Alternative A. It keeps the CLI operationally ambiguous.

Surviving:

- Alternative B: explicit warning, then HTTP fallback
- Alternative C: fail closed unless the user opts into HTTP explicitly

## Implications for TODOs and pending DIs

- TODO `127` (`musav`) is locked to Alternative `C`.
- HTTP compatibility remains shipped, but the CLI now requires an explicit
  `-socket=off` opt-in before using it.
- The doc sweep should update long-form transport wording so Neovim reads as
  socket-first with websocket and HTTP fallback, and the CLI reads as
  fail-closed by default on local socket loss.
