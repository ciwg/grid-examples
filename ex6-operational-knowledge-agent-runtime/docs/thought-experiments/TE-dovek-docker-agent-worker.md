# Docker-Confined Agent Worker

TE ID: TE-dovek

## Status

decided

## Decision

Docker is OKAR's first executable-agent backend. The runtime invokes a declared
image as a one-shot JSON stdin/stdout worker with a read-only root filesystem,
`--network none`, no implicit host mounts, explicit read-only mounts only, dropped
Linux capabilities, no-new-privileges, CPU/memory/PID limits, and a timeout.

No direct host-process fallback is permitted. The worker cannot receive the runtime
root, CAS, history, peer keys, Docker socket, or ambient environment. Durable
writes remain runtime-mediated through validated worker output.

## Rationale

Docker is already the deployed, reachable containment substrate for this project.
This supersedes TE-safuk's generic OS-process-worker recommendation.

## Decision status
