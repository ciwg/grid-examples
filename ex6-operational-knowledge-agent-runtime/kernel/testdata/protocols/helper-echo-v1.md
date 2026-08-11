# helper.echo.v1

Status: frozen test fixture

## Record contract

This test-only external family uses the canonical Grid package-record carriage. Its payload is canonical JSON with the fixture-defined `message` field. It exists to prove installed-package validation, routing, and unknown-family behavior use an explicit external pCID rather than a synthesized identifier.

## Evolution

This fixture is immutable. A changed test contract requires a new fixture specification and CID. Source: DI-solan.
