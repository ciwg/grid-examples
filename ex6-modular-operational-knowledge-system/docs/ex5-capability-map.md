# EX5 Capability Map

This is the current first cut for sorting ex5 into ex6.

The goal is not to preserve ex5 as one app. The goal is to identify what
belongs in the basket, what belongs in eggs, and what should be left behind.
Source: `DI-moksu`.

## Belongs In EX6 Core

- package install / list / inspect
- package manifest validation
- package self-check activation
- command routing
- family registration
- protocol implementation claims
- append-only history
- CAS
- relay import / export
- unknown-family exact-byte retention

Source: `DI-moksu`; `DI-lupok`.

## Concrete EX5 To EX6 Mapping

| EX5 capability | EX6 destination |
| --- | --- |
| package install / activation | core |
| append-only operational history | core |
| CAS blob storage | core |
| relay export / import | core |
| family registry | core |
| protocol implementation claims | core |
| responsibilities | `context` package |
| places | `context` package |
| resources | `context` package |
| knowledge items | `knowledge` package |
| revision snapshots | `knowledge` package |
| item lifecycle (`draft`, `approved`, `superseded`) | `knowledge` package |
| runs | `runs` package |
| evidence metadata | `runs` package |
| approvals | `runs` package first, maybe later split |
| typed links | `links` package |
| procedure-specific behavior | `procedures` package |
| training-specific behavior | `training` package |
| maintenance-specific behavior | later `maintenance` package |
| receiving-check behavior | later `receiving` package |
| inventory-audit behavior | later `inventory` package |
| browser embodiment | not core; later optional embodiment |
| Neovim embodiment | not core; later optional embodiment |

Source: `DI-moksu`; `DI-lupok`.

## Should Become EX6 Packages

- responsibilities
- places
- resources
- knowledge items
- revisions
- runs
- evidence
- approvals
- typed links
- procedures
- training
- maintenance
- receiving checks
- inventory audits

Source: `DI-moksu`.

## Should Be Left Behind As Core Assumptions

- browser as required embodiment
- Neovim as required embodiment
- ex5 command naming
- ex5’s single-monolith workflow structure
- ex5’s local embodiment assumptions as the definition of the product

Source: `DI-moksu`.

## Recommended First Package Set

Start with fewer larger eggs:

1. `context`
2. `knowledge`
3. `runs`
4. `links`
5. `procedures`

Later domain eggs can extend from there. Source: `DI-moksu`.

## Next Domain Eggs

With `training` now implemented, the next domain eggs are:

1. `maintenance`
2. `receiving`
3. `inventory`

Source: `DI-sivuk`.

## First Five Package Ownership

### `context`

- responsibilities
- places
- resources
- context drilldown identifiers and family ownership

### `knowledge`

- knowledge items
- revisions
- lifecycle changes
- knowledge-family validation

### `runs`

- run records
- evidence
- approvals
- actor/outcome context around performed work

### `links`

- typed links
- link validation
- cross-record relation surface

### `procedures`

- procedure item kind
- procedure-flavored commands
- first serious domain example built on top of `knowledge`, `runs`, and `context`

Source: `DI-moksu`.
