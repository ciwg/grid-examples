# Receiving-exception workflow

This workflow opens a quarantine case after a failed inbound receiving
inspection. It composes `context`, `receiving`, `quarantine`, and `runs`.

It records the receiving item, failed receipt run, actor, evidence reference,
and exception in one durable opening event. Release and rejection are separate
quarantine transitions and are deliberately outside this workflow's scope.
Source: DI-hogid.
