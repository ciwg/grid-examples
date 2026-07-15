# grid-editor remote mutation bootstrap for multi-machine ex3

TE ID: TE-takav

## Status

decided

## Decision under test

How `ex3-grid-editor-websocket` should admit remote browser and Neovim clients
for live mutation so the example can run across multiple machines and normal
Docker bridge networking, without collapsing protocol meaning into transport or
pretending upstream PromiseGrid has already frozen one universal app-auth API.

## Assumptions

- `ex3` already keeps protocol meaning in the repo-local `live-document`,
  `live-awareness`, `document-metadata`, and `publish-document` specs.
- The current websocket work is carriage only; it does not redefine those
  protocol boundaries.
- The current blocker is admission: mutation is loopback-only today.
- `ex1-order-flow` provides the closest local model for signed short-lived
  capability tokens.
- The latest upstream dev-guide snapshot keeps app-facing auth/API guidance
  provisional under `DR-tuhaz` and separates semantic-model work from DevOps
  work.

## Alternatives

### Alt A: keep loopback-only mutation and host-network Docker

Preserve the current model and treat multi-machine behavior as out of scope.

### Alt B: accept one shared bearer secret directly on every mutation request

Add one operator-configured access token and require clients to send it on
every HTTP or websocket mutation path.

### Alt C: use a bootstrap secret only to mint short-lived per-protocol mutation capabilities

Introduce one operator-configured bootstrap secret for a repo-local session
bootstrap endpoint. A remote client presents that secret once, bound to a
document and participant identity, and the relay returns short-lived signed
capability tokens for the relevant document protocols. The client then uses the
capabilities for HTTP mutation and websocket live transport.

### Alt D: require every remote client to present an external signing identity and full proof chain

Make browser and sidecar clients bring their own durable keys, prove identity
up front, and participate in a more fully PromiseGrid-native auth model before
the relay will accept mutations.

## Scope and systems affected

- `ex3-grid-editor-websocket/service/**`
- `ex3-grid-editor-websocket/web/src/**`
- `ex3-grid-editor-websocket/cmd/grid-nvim-sidecar/**`
- `ex3-grid-editor-websocket/cmd/grid-relay/main.go`
- `ex3-grid-editor-websocket/cmd/grid-editor/main.go`
- `ex3-grid-editor-websocket/compose.yaml`
- `ex3-grid-editor-websocket/README.md`
- `ex3-grid-editor-websocket/docs/docker-simulation.md`

## Scenario analysis

### S1 — Alice opens the browser from another laptop against Bob's relay

With Alt A, Bob must keep host networking and loopback illusions in place, so
Alice cannot honestly mutate from another machine. Alt B works mechanically,
but the relay sees one undifferentiated long-lived bearer secret on every live
request and websocket, which makes scope, expiry, and replay reasoning weak.
Alt C lets Bob share one admission secret out of band, then gives Alice
short-lived document-scoped capabilities that match the relay's existing
protocol split. Alt D is more principled long-term, but it is a much larger
product and UX jump than `ex3` needs right now.

### S2 — Alice edits in the browser while Carol edits through Neovim

Alt A still fails for non-loopback. Alt B allows both embodiments, but every
path keeps carrying the same bootstrap bearer forever. Alt C works for both
embodiments while preserving the repo-local split between `live-document`,
`live-awareness`, `document-metadata`, and `publish-document`; websocket stays
carriage only. Alt D can also preserve the split, but forces both embodiments
to own a much heavier identity bootstrapping story immediately.

### S3 — Mallory captures one websocket frame or one HTTP request

Alt B is weak here because the same long-lived bearer secret remains useful for
new sessions and unrelated actions until the operator rotates it. Alt C limits
damage: the bootstrap secret is only for minting capabilities, while the live
mutation capability can be short-lived, document-scoped, protocol-scoped, and
replay-checked by token ID. Alt D can be even stronger, but only with much more
implementation and operator burden.

### S4 — normal Docker bridge networking and published ports

Alt A depends on `network_mode: host`, so it does not prove the thing we want.
Alt B and Alt C both work over published ports. Alt C is better aligned with
the PromiseGrid guidance that capability-like behavior should be explicit
promise-bearing artifacts rather than transport quirks or hidden route knobs.
Alt D also works over published ports, but again with much higher complexity.

### S5 — future upstream migration

Alt B leaves a sticky ad hoc auth surface that will be harder to supersede
cleanly later. Alt C keeps the long-lived secret clearly local and provisional,
while the actual remote mutation tokens already look like short-lived relay
capabilities. That gives this repo a cleaner migration seam if upstream later
locks a stronger app-facing auth profile. Alt D might end up closer to that
future, but it would overfit `ex3` to an upstream contract that does not yet
exist.

## Conclusions

Rejected:

- Alt A because it fails the stated multi-machine goal.
- Alt B because it solves admission too bluntly and leaves the bootstrap secret
  on every live mutation path.
- Alt D because it is too large and presumptive for the current provisional
  upstream state.

Surviving alternative:

- Alt C.

## Implications for open TODOs and pending DIs

- `TODO-buvir` should implement one repo-local bootstrap endpoint that mints
  short-lived per-protocol mutation capabilities from an operator-configured
  bootstrap secret.
- `ex3` should move its Docker demo from host networking to normal published
  ports.
- The README and Docker docs should explicitly say that this admission model is
  PromiseGrid-aligned in direction but still repo-local and provisional.

## Decision status

locked by `DI-povip`; `DI-talih`
