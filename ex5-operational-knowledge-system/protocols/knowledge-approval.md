# knowledge-approval

`knowledge-approval` is the durable family for named-role review outcomes over
durable records. Source: `DI-vosul`.

Frozen v1 scope in ex5:

- approvals for knowledge items
- approvals for performed runs
- named review roles such as reviewer or approver
- decisions such as approved, rejected, or noted

Each durable knowledge-approval message in this second frozen family owns:

- stable approval identity
- target type
- target id
- actor
- sequence and durable timestamp
- target revision when the approval applies to a knowledge-item revision
- named review role
- review decision
- review notes

The second runtime slice signs and verifies this event kind under the family:

- `approval_recorded`

Knowledge-item lifecycle changes remain outside this family even when an
approval later causes a knowledge-item status transition. The current code
models the family through `Approval` and the signed knowledge-approval
envelope log. Source: `DI-vosul`.
