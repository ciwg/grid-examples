# TODO rovum - ex6 package template and scaffold

## Decision Intent Log

ID: DI-rovum
Date: 2026-07-28 15:10:00
Status: active
Decision: Extend the installed-package contract so `run` may return basket-mediated CAS write requests and record append requests instead of limiting installed eggs to plain text output.
Intent: Let outside package authors create durable records through the runtime without bypassing the basket or waiting for a later full plugin ABI.
Constraints: The runtime remains the only writer to its CAS and append-only history; the contract stays narrow and explicit.
Affects: `packages/external`, `kernel`, installed-package tests, template package, and package author docs.

ID: DI-lised
Date: 2026-07-28 16:25:00
Status: active
Decision: Re-activate installed packages from `.moks/packages/` when the runtime opens so package installation survives later CLI invocations and restarts.
Intent: Make installed executable eggs behave like real installed packages instead of process-local registrations that disappear on the next command.
Constraints: The runtime package root remains the source of installed-package activation; startup should ignore non-package files and fail fast on broken installed manifests.
Affects: `kernel/runtime.go`, installed-package tests, runnable examples, and CLI package flows.

## Goal

Make it practical for other people to create new ex6 packages.

## Scope

- `templates/package/`
- starter manifest
- starter executable contract
- author guide examples
- optional scaffold command later

## Why

The product is not genuinely modular if only the repo author can add eggs.
