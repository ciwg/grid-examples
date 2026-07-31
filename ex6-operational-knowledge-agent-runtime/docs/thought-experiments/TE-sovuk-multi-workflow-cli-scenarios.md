# Multi-workflow CLI scenario tests

TE ID: TE-sovuk

## Status

decided

## Decision under test

How to prove multiple workflow artifacts interact through the main `moks`
program without falsely claiming that artifacts are an execution engine.

## Assumptions

- Current workflow artifacts are loadable instructions, not executable workers.
- Package commands are the existing operational behavior.
- The test must use the command dispatcher used by the `moks` CLI.

## Alternatives

1. Add a workflow execution engine.
2. Unit-test each workflow directory in isolation.
3. Run a shared-state CLI scenario that loads artifacts and invokes real package
   commands in one runtime.

## Scenario analysis

Option 1 changes product behavior and needs a separate execution design. Option
2 proves manifest validity but not shared operational state. Option 3 loads all
workflow artifacts, activates them locally, then creates shared context and
records receiving, inventory, maintenance, training, and knowledge activity
through the main command dispatcher. It proves composition while preserving the
truth that the artifacts do not autonomously execute.

## Conclusion

Choose option 3: add a deterministic CLI scenario test. It verifies shared
runtime interaction, active workflow status, and cross-package records without
adding a workflow engine.

## Decision status

Locked by DI-sovuk.
