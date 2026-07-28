# EX6 Agent/Kernel Alignment Note

This note records how the recent architecture notes relate to current `ex6`.

## Working interpretation

- the **runtime** is still the right top-level idea
- the **kernel** should be understood as the runtime expressed as **node services / agent roles**
- the **apps / packages** should be understood as **agents**

So:

- runtime = node-local coordinating system
- kernel = roles inside that runtime
- packages/apps = agents that do useful work through those roles

Source: `DI-moksu`; `DI-puvok`.

## What already fits current ex6

- `ex6` is already protocol-first and package-first rather than monolith-first
- `ex6` already treats protocol identity as the stable boundary
- `ex6` already assumes modular growth and later-added packages
- `ex6` already allows trust, routing, and relay behavior to stay runtime-owned instead of becoming package-private side channels

Source: `DI-moksu`; `DI-lupok`; `DI-rovum`; `DI-rumek`; `DI-puvok`.

## What current ex6 does not model yet

- the current `kernel` is mostly one runtime process, not a visible set of cooperating agents
- package activation is still manifest/self-check activation, not startup promises to a routing agent
- routing is not yet modeled as a dedicated agent promising delivery by `pCID`
- parser-agent variations are not first-class yet
- multiple routing agents on one node are not modeled yet

Source: `DI-moksu`; `DI-rovum`; `DI-puvok`.

## Adjustment pressure on ex6

These notes mean current ex6 should be read as a **useful intermediate embodiment**, not the final architectural shape.

The likely long-term shift is:

- from `runtime loads packages`
- toward `node roles and agents promise capabilities to each other`

That does **not** invalidate the runtime-centered idea.

It refines it:

- the runtime remains the center
- the center becomes more explicitly agent-shaped
- routing becomes one kernel role
- parsing can become another kernel role
- app packages become agents that register promise-based protocol handling

Current concrete step:

- package claims now register explicit protocol routes inside the kernel
- registered families must declare a matching `family-validator` route claim
- `moks route list` exposes the current route table so routing is visible as a kernel service
- relay export now carries those route registrations so routing roles are visible across runtimes too

Source: `DI-puvok`; `DI-rutom`; `DI-ruvot`.

## Practical effect on next design work

When we touch kernel/app boundaries next, prefer changes that move ex6 toward:

- role-based kernel services
- agent-shaped packages
- `pCID`-declared routing promises
- optional parser-agent hops
- organic protocol growth without forcing one giant protocol document

Source: `DI-puvok`; `DI-rutom`.
