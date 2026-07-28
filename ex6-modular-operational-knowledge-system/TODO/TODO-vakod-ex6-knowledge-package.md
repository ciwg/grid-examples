# TODO vakod - ex6 knowledge package

## Decision Intent Log

ID: DI-vakod
Date: 2026-07-28 13:30:00
Status: active
Decision: Implement the first knowledge package as a first-party built-in egg with item, revision, and lifecycle families while the installed-package mutation contract remains shallow.
Intent: Get the modular replacement for the ex5 knowledge-item model into ex6 code now instead of waiting for a later external package ABI pass.
Constraints: Keep the package under `packages/knowledge/`; keep revision bodies CAS-backed; preserve basket-first durable history writes.
Affects: `packages/knowledge`, `cmd/moks`, runtime tests, and ex6 current-state docs.

## Goal

Build the first-party `knowledge` package.

## Scope

- knowledge item creation
- revisions
- status lifecycle
- approval hooks
- durable family declarations

## Why

This package replaces the central ex5 document model in modular form.
