# OKAR Architecture Cheat Sheet

## One sentence

**OKAR is a node-local runtime that turns package claims into inspectable
protocol routes, preserves durable evidence, and exchanges exact bytes under
local trust policy.**

## Layers

| Layer | Responsibility |
| --- | --- |
| Operator / CLI | Installs packages; inspects routes; sets local trust and route policy. |
| Runtime / kernel | Activates packages, derives routes, validates known records, manages relay and workflow lifecycle. Its 27 built-in pCIDs are fixed checked-in values. |
| Packages | Own domain commands, protocol families, validators, implementation claims, and—when external—their own immutable family specs. |
| Durable stores | Append-only exact canonical record bytes plus CAS for immutable payloads and lifecycle events. |
| Grid / relay | Moves exact records and relay evidence; discovery does not create trust. |

## Key flow

    built-in registry or external package spec + manifest + self-check
                ↓
    claims → runtime route table → inspect / plan / policy
                ↓
    known record → package validation → append-only history
    unknown record ─────────────────→ exact-byte relay retention
                ↓
    workflow artifact → CAS lifecycle event → disposable local projection

## Important boundaries

- **Protocol identity:** pCID identifies one immutable shared record contract;
  it is not a workflow, package, executable, or message hash.
- **Workflow extension:** workflows compose existing pCIDs. New shared durable
  semantics require a new immutable spec and pCID; ordinary workflow changes do
  not require recompilation.
- **Trust:** peer discovery is not permission to pull or push.
- **Evidence:** semantic author signatures and relay-carriage signatures are
  separate evidence layers and are interpreted by local policy.
- **Authority:** importing a workflow artifact does not by itself grant worker
  or route-execution authority.
- **Durability:** CAS events and artifacts are authoritative; caches are not.
- **Current limit:** routing is planned and explained, but not yet a complete
  promise-based multi-agent dispatch system.

Source: DI-moksu, DI-puvok, DI-bavuk, DI-jusij, DI-solan.
