# TODO zibek - ex6 receiving package

## Decision Intent Log

ID: DI-zibek
Date: 2026-07-28 17:40:00
Status: active
Decision: Implement the first receiving package as a first-party built-in egg with explicit receiving item, receipt, and disposition families over shared context, knowledge, and runs records.
Intent: Add the next real domain egg without collapsing receiving behavior into the basket or hiding receipt/disposition truth inside generic run notes.
Constraints: Keep the package under `packages/receiving/`; preserve basket-mediated durable history; compose with `context`, `knowledge`, and `runs` instead of duplicating their cores.
Affects: `packages/receiving`, `cmd/moks`, runtime tests, and ex6 status docs.

## Goal

Build the first-party `receiving` package.

## Scope

- receiving item kind
- receipt recording
- disposition recording
- receiving inspection and listing
- durable family declarations

## Why

Receiving is one of the next real ex5-derived domain surfaces and should
exist as its own egg instead of waiting behind later embodiments.
