# Docker simulation

This is the quickest way to simulate two separate relay machines for
`grid-editor` without needing two laptops.

This copied example uses `127.0.0.1:7025` and `127.0.0.1:7026` so its Docker
simulation can run beside `ex2-grid-editor`. Source: `DI-vatub`.

## What it does

- runs `relay-a` on `127.0.0.1:7025`
- runs `relay-b` on `127.0.0.1:7026`
- gives each relay its own Docker volume for data
- makes each relay poll the other as a peer
- enables one shared demo bootstrap token: `ex3-demo-access`

This now uses normal published Docker ports instead of host networking. Remote
mutation is admitted through a repo-local bootstrap token that mints
short-lived relay-signed document capabilities for browser or Neovim clients.
Source: `DI-povip`.

The compose file also runs the demo relays as root inside the container so the
named Docker volumes are writable without extra permission setup. This keeps
the simulation quick and predictable for local testing.

## Start it

From the repo root:

```bash
docker-compose up --build
```

If your Docker install has the newer plugin form, this also works:

```bash
docker compose up --build
```

## Open it

Open one browser window to:

```text
http://127.0.0.1:7025/?doc=demo&access_token=ex3-demo-access
```

Open a second browser window to:

```text
http://127.0.0.1:7026/?doc=demo&access_token=ex3-demo-access
```

Use the same document ID in both windows.

## What this simulates

- separate relay processes
- separate relay identities
- separate relay data roots
- peer-to-peer relay polling
- published-port remote browser access
- repo-local remote session bootstrap and mutation capabilities

## What this does not simulate perfectly

- WAN latency
- firewall problems
- long-lived operator auth policy
- a frozen upstream PromiseGrid app-auth API

## Useful commands

Stop:

```bash
docker-compose down
```

Stop and remove relay data:

```bash
docker-compose down -v
```

Show logs:

```bash
docker-compose logs -f
```

## Why the token is still provisional

The newer upstream PromiseGrid snapshot still keeps app-facing auth/API
guidance provisional, and now separates `POC20` semantic-model work from the
`POC21` DevOps/bootstrap track. `ex3` therefore uses a repo-local bootstrap
token and short-lived relay capabilities as a practical demo slice, not as a
claim that upstream has already frozen one universal remote-editor auth shape.
Source: `DI-talih`; `DI-povip`.
