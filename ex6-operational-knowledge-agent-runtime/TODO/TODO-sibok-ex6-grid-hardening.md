# TODO sibok - ex6 grid hardening

## Goal

Move ex6 beyond the current relay batch shell.

## Scope

- stronger relay metadata
- better import/export validation
- proof/signature discipline
- mixed-peer handling
- migration and compatibility rules

## Why

The current runtime is architecture-correct but still shallow on actual grid
depth.

## Workflow Artifact Hardening

- [ ] Preserve imported workflow artifact CIDs and provenance through deactivation, revocation, replacement, and restart. (DI-lovek; TE-gavuk)
- [ ] Support local package-directory capture into CAS and direct CAS workflow-artifact import with validation. (DI-lovek; TE-gavuk)
