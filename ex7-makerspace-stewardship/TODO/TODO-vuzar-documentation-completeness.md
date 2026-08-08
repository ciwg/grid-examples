# Ex7 documentation completeness

## Decision Intent Log

ID: DI-damod
Date: 2026-08-08 12:27:09
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Complete Ex7's local-demo documentation with separate operator, workflow/evidence, and verification guides; use documented browser/API steps and deterministic HTTP tests rather than browser automation.
Intent: Make Ex7 understandable and safely evaluable without overstating its local-only trust model.
Constraints: Ex7-only; cover normal observation/hold-clear/loan/return and corrupt-log failure; no authentication, replication, signatures, or new domain behavior.
Affects: README, docs guides, server tests.

## Status

- [x] Publish the locked documentation-first completeness package.
