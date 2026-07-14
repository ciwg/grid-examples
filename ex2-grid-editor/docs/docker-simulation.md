# Docker simulation

This is the quickest way to simulate two separate relay machines for
`grid-editor` without needing two laptops.

## What it does

- runs `relay-a` on `127.0.0.1:7015`
- runs `relay-b` on `127.0.0.1:7016`
- gives each relay its own Docker volume for data
- makes each relay poll the other as a peer

This uses `network_mode: host` on Linux so the browser still reaches the relay
as a loopback client. That matters because local mutation endpoints are still
loopback-only in the current security model.

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
http://127.0.0.1:7015/?doc=demo
```

Open a second browser window to:

```text
http://127.0.0.1:7016/?doc=demo
```

Use the same document ID in both windows.

## What this simulates

- separate relay processes
- separate relay identities
- separate relay data roots
- peer-to-peer relay polling

## What this does not simulate perfectly

- WAN latency
- firewall problems
- cross-host browser networking
- a future authenticated remote-browser mutation mode

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

## Why host networking is used

The current relay only accepts browser mutation requests from loopback clients.
With normal Docker bridge networking, the relay would see the host browser as a
non-loopback client and reject edits. Host networking keeps the simulation
simple for the current demo slice.
