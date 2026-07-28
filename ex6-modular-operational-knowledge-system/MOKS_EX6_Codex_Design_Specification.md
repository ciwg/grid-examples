# MOKS (EX6) --- Codex Design Specification

> **Working name:** Modular Operational Knowledge System (MOKS)
>
> **Codename:** EX6

## Executive Summary

Design EX6 (MOKS) as a completely new, standalone, independently
shippable, CLI-first product built on PromiseGrid.

This is **not** EX5 version 2.

EX5 is available only as: - a capability inventory - a source of ideas -
a source of reusable implementation - a reference implementation

Do **not** preserve EX5's architecture. Copy only the pieces that are
appropriate for the new architecture.

------------------------------------------------------------------------

# Product Vision

Organizations should be able to install independent operational
capability packages into a common runtime.

The runtime provides shared infrastructure.

Packages provide business behavior.

The runtime must not require prior knowledge of every future workflow.

------------------------------------------------------------------------

# Settled Requirements

-   Standalone product.
-   Independently shippable.
-   CLI-first.
-   Browser not required for the first release.
-   PromiseGrid-native.
-   Installable packages.
-   Packages may include executable code.
-   Packages may introduce agents.
-   Packages may register CLI commands.
-   Packages may introduce new durable record families.
-   Packages explicitly register with the runtime.
-   Business logic belongs in packages.
-   The runtime provides shared services and coordination.

------------------------------------------------------------------------

# Runtime Responsibilities (Current Direction)

Evaluate which belong in the runtime:

-   PromiseGrid connectivity
-   package installation
-   package registration
-   package discovery
-   lifecycle management
-   CLI routing
-   append-only storage
-   CAS
-   identity
-   search
-   approvals
-   evidence
-   places
-   resources
-   responsibilities
-   typed links
-   capability routing
-   inter-package communication

------------------------------------------------------------------------

# Packages

Packages may contribute:

-   agents
-   executables
-   CLI namespaces
-   protocols
-   specifications
-   schemas
-   durable record families
-   projections
-   validators
-   documentation

Packages are **not** configuration files. They may contain arbitrary
executable behavior.

------------------------------------------------------------------------

# CLI

Example only:

``` text
moks package install ./maintenance
moks package list
moks maintenance record-work
moks training assign
moks search
```

Recommend a scalable command-routing model.

------------------------------------------------------------------------

# Record Families

Packages must be able to introduce durable record families unknown to
the runtime when it was compiled.

Design:

-   declaration
-   validation
-   storage
-   search
-   replication
-   signing
-   versioning

------------------------------------------------------------------------

# Communication

All communication should align with PromiseGrid concepts.

Recommend:

-   capability discovery
-   package communication
-   local embodiment strategy
-   remote behavior

Avoid unrelated communication systems.

------------------------------------------------------------------------

# Installation

Evaluate:

-   local directory
-   archive
-   Git
-   CAS
-   registry

Recommend an initial implementation and an evolution path.

------------------------------------------------------------------------

# Process Model

Evaluate:

1.  Separate process per package.
2.  In-process.
3.  Hybrid.

Recommend one.

------------------------------------------------------------------------

# Dependencies

Evaluate package and capability dependencies.

Support:

-   optional dependencies
-   version constraints
-   provider selection
-   missing dependency handling

------------------------------------------------------------------------

# Initial Packages

Recommend a minimal set proving:

-   package installation
-   independent behavior
-   shared runtime services
-   CLI registration
-   PromiseGrid integration
-   durable records

------------------------------------------------------------------------

# EX5 Guidance

Inspect EX5 and produce:

-   capability inventory
-   reusable implementation candidates
-   components to leave behind
-   architectural differences

Never preserve EX5 merely because it already exists.

------------------------------------------------------------------------

# Deliverables

Before major coding produce:

1.  Product definition
2.  Architecture
3.  Core vs package matrix
4.  Manifest proposal
5.  CLI model
6.  Communication model
7.  Capability model
8.  Record model
9.  Installation model
10. Repository layout
11. Security analysis
12. Phased implementation plan
13. Open design questions

Clearly distinguish: - settled requirements - preferences -
recommendations - unresolved questions

The objective is a clean architecture, not a migration.
