## Goal

Build a PromiseGrid example in this repo: an order-fulfillment demo that makes
multi-agent grid messaging visible.

Refer to:
- PromiseGrid dev guide: `~/lab/cswg/promisegrid-dev-guide/README.md`
- Wire-lab repo POCs 16, 18, 19, 20: `~/lab/wire-lab`

This example should:

- demonstrate multiple business agents exchanging real grid messages
- run with one container per agent plus one kernel container
- keep the PromiseGrid contract in the message/spec layer, not in hidden
  in-process calls or a shared database
- stay aligned with the PromiseGrid dev guide, especially the current `App
  Devs` guidance around pCID-selected envelopes, local promise interpretation,
  append-only artifacts, and `poc12` hybrid fulfillment orientation

## Why This Example

The fulfillment demo demonstrates why the mechanics matter:

- several independent agents participate
- multiple pCIDs are active in one workflow
- business refusal and protocol failure are different things
- timeouts, duplicate delivery, signatures, and local promise interpretation
  all matter
- durable artifacts and raw message bytes are obviously useful
- per-agent containers feel like a real distributed deployment

The dev guide's `poc12` guidance is the closest upstream orientation point:
multi-pCID handler routing, hybrid fulfillment, postal-scale/label/accounting
style subflows, signed traffic, capability-token promises, and app-local
artifacts without turning the kernel into the business workflow owner.

## Non-Negotiable Alignment Points

The next implementation session should treat these as hard constraints:

- All inter-agent communication must be real grid traffic:
  `grid([42(pCID), payload, proof])`.
- The kernel parses slot `0` only, routes by pCID only, and forwards exact
  bytes only.
- The kernel is not a router-with-business-logic, service registry, or RPC
  authority.
- Promise interpretation stays in app agents, not in the kernel.
- Malformed bytes are quarantined or declined locally, never silently
  reinterpreted.
- Raw message artifacts, signed capability tokens, and local journals are
  append-only.
- Process shape is host-local. The protocol contract is the stable thing.
- Payloads must not contain protocol names.
- Do not copy `poc12` names, route strings, or toy payloads as if they were
  final API.
- No hidden direct calls between business agents. If one agent needs another,
  it sends a message through the kernel boundary.
- MVP must use actual signatures on every message path, and capability tokens
  wherever the selected protocol profile requires them.
- The MVP should be deterministic. Do not require an LLM for the first working
  slice.

## Repo Strategy

Recommended direction:

- prefer an isolated namespace such as:
  - `cmd/pg-order-agent`
  - `cmd/pg-order-submit`
  - `orderdemo/...`
  - `docs/order-fulfillment/...`
  - `deploy/order-fulfillment/docker-compose.yml`

Recommended binary shape:
- a different binary for each agent role.

## Runtime Topology

Rules:

- every business agent dials the kernel only
- no direct seller -> warehouse TCP socket outside the kernel path
- the submitter/intake agent may run as a short-lived container invocation
- each agent has its own local CAS store
- shared storage may be used for artifact persistence, but not as the live
  communication mechanism
- use the microkernel-style supervisor, along with the stdout collector and
  analyzer from POC16

## Agent Roles And Promises

### Intake agent

Responsibility:

- accept one order fixture from the operator
- send the top-level order submit message
- wait for one final order message
- interpret the seller's promise locally from the signed traffic it observed

Promise:

> I promise to submit this order for fulfillment and report what I observed
> about the resulting top-level order promise.

### Seller agent

Responsibility:

- receive `order` messages with `kind = "submit"`
- validate order shape, signature, capability token, and local policy
- orchestrate warehouse, accounting, and carrier substeps by message
- send one final `order` message with `kind = "final"`

Promise:

> I promise to process this order request under the selected protocol and send
> one conforming final message describing fulfillment, refusal, or failure.

Important:

- `fulfilled`, `refused`, and `failed` are business outcomes
- a conforming `refused` result can still mean the seller kept its protocol
  promise

### Warehouse agent

Responsibility:

- receive `pick_pack` messages with `kind = "request"`
- deterministically decide whether items can be picked and packed
- emit one `pick_pack` message with `kind = "result"`

Promise:

> I promise to return one conforming pick-and-pack result for this work
> request.

### Accounting agent

Responsibility:

- receive `accounting` messages with `kind = "request"`
- deterministically record or refuse accounting state
- emit one `accounting` message with `kind = "result"`

Promise:

> I promise to return one conforming accounting result for this accounting
> request.

### Carrier agent

Responsibility:

- receive `shipment` messages with `kind = "request"`
- deterministically book or refuse shipment
- emit one `shipment` message with `kind = "result"`

Promise:

> I promise to return one conforming shipment-booking result for this shipment
> request.

### Kernel

Responsibility:

- receive registration messages
- parse slot `0` only
- route exact bytes
- emit kernel-local operational events only

The kernel never claims whether the seller, warehouse, accounting, or carrier
promises were kept.

## MVP Scenario

The first working slice should be intentionally narrow:

- one order at a time
- one seller
- one warehouse
- one accounting agent
- one carrier
- one kernel
- one intake submission per run
- deterministic business logic only
- actual signatures on every message
- cryptographic capability tokens on every message path
- no real external APIs
- no LLMs in MVP

Use a tiny local catalog and scripted behavior:

- `widget-1` succeeds
- `widget-oos` causes warehouse refusal
- `widget-dup-pay` causes accounting refusal before shipment is attempted
- `widget-carrier-timeout` causes seller-visible shipment timeout/failure after
  accounting succeeded

## Business Flow

Happy path:

1. Intake sends an `order` submit message to seller.
2. Seller validates the message shape, signature, capability token, and local
   policy.
3. Seller sends a `pick_pack` request to warehouse.
4. Warehouse returns a `pick_pack` result.
5. Seller sends an `accounting` request to accounting.
6. Accounting returns an `accounting` result.
7. Seller sends a `shipment` request to carrier.
8. Carrier returns a `shipment` result.
9. Seller sends an `order` final message to intake.
10. Intake derives the top-level promise status locally from the signed traffic
    it observed.

Failure examples:

- warehouse refusal:
  - warehouse returns a conforming refusal result
  - seller returns a conforming top-level `order` final message with
    `order_status = "refused"`
  - intake interprets seller promise status as `kept` if the final message
    conforms
- accounting refusal:
  - accounting returns a conforming refusal result
  - seller does not attempt shipment booking
  - seller returns a conforming top-level `order` final message with
    `order_status = "refused"` and `failure_stage = "accounting"` when the
    accounting agent refused for business reasons
  - if the accounting exchange failed because of malformed bytes, bad
    signatures, or bad capability tokens, seller returns `order_status =
    "failed"` and `failure_stage = "accounting"`
- shipment timeout:
  - seller times out waiting for carrier
  - seller may emit a conforming `order` final message with
    `order_status = "failed"` and `failure_stage = "carrier"`
  - if seller itself stays silent past the intake timeout, intake reports
    `promise_status = "not_promised"`
- invalid signature or missing capability token:
  - the receiving agent refuses the message locally
  - the refusing agent may emit a conforming final failure or refusal message if
    its own promise requires one

## Business Semantics Vs Promise Semantics

Do not collapse these:

- **Promise status**:
  - `kept`
  - `broken`
  - `not_promised`

- **Business outcome**:
  - `fulfilled`
  - `refused`
  - `failed`

Examples:

- a conforming refusal from seller can be `promise_status = kept` and
  `order_status = refused`
- a malformed top-level final message is `promise_status = broken`
- seller silence is `promise_status = not_promised`

This distinction is one of the main reasons this demo is worth building.

## Order Input And Final Output

Recommended input fixture shape:

```json
{
  "customer_order_ref": "demo-001",
  "requested_by": "customer-demo",
  "ship_to": {
    "name": "Alice Example",
    "address1": "123 Example St",
    "city": "Portland",
    "region": "OR",
    "postal_code": "97201",
    "country": "US"
  },
  "items": [
    {
      "sku": "widget-1",
      "quantity": 1
    }
  ],
  "service_level": "ground",
  "payment_ref": "pay-demo-001",
  "notes": "Leave at front desk"
}
```

Recommended final intake output shape:

```json
{
  "customer_order_ref": "demo-001",
  "order_status": "fulfilled",
  "promise_status": "kept",
  "failure_stage": null,
  "package_id": "PKG-...",
  "tracking_number": "TRACK-...",
  "ledger_entry_id": "LEDGER-...",
  "final_order_cid": "b...",
  "notes": "conforming signed final order message received for this request"
}
```

## Shared Signed Envelope And Capability Token Requirements

Use real cryptography in MVP.

Shared envelope shape:

```text
grid([
  42(pCID),
  payload,
  proof
])
```

Rules:

- slot `0` is the routing pCID
- slot `1` is the protocol payload as a CBOR item directly, not a `bstr`
  wrapper around encoded payload bytes
- slot `2` carries the proof bytes for the pCID-defined signable view
- the kernel still parses slot `0` only
- each receiving app validates proof, payload shape, any required capability
  token fields, and local policy before acting

Capability-token requirements:

- a capability token is a signed promise artifact granting bounded authority to
  request a specific action
- capability tokens belong in pCID-owned payload fields rather than a universal
  envelope slot
- capability tokens must be cryptographically verifiable, not placeholder
  strings
- capability tokens may scope authority by sender, receiver, protocol family,
  message kind, and expiry
- refusal behavior for missing, expired, malformed, or unauthorized capability
  tokens must be explicit in the app-local logic

Signature requirements:

- signatures must be real asymmetric signatures, not test-only hashes or
  marker fields
- signature verification failure is a local protocol failure, not a business
  refusal
- registration traffic is signed too, and any required capability tokens are
  defined by the registration profile

## Draft Protocol Set

Use multiple pCIDs on purpose, but only where the payload family differs
materially. The MVP should use four business-domain pCIDs plus one reused
kernel-registration pCID.

### 1. `order`

Purpose:

- intake asks seller to process one order
- seller reports one final top-level outcome back to intake

Payload:

```text
{
  "kind": "submit" / "final",
  "customer_order_ref": tstr,
  ? "requested_by": tstr,
  ? "capability_token": bstr,
  ? "ship_to": {
    "name": tstr,
    "address1": tstr,
    ? "address2": tstr,
    "city": tstr,
    "region": tstr,
    "postal_code": tstr,
    "country": tstr
  },
  ? "items": [1*8 {
    "sku": tstr,
    "quantity": uint
  }],
  ? "service_level": "ground" / "priority",
  ? "payment_ref": tstr,
  ? "notes": tstr,
  ? "parent_order_cid": 42(bstr),
  ? "order_status": "fulfilled" / "refused" / "failed",
  ? "failure_stage": "seller_validation" / "warehouse" / "accounting" / "carrier",
  ? "warehouse_result_cid": 42(bstr),
  ? "accounting_result_cid": 42(bstr),
  ? "shipment_result_cid": 42(bstr),
  ? "package_id": tstr,
  ? "tracking_number": tstr,
  ? "ledger_entry_id": tstr,
  ? "summary": tstr
}
```

Notes:

- `kind = "submit"` carries the order request fields
- `kind = "final"` carries the final outcome fields
- payloads do not contain the protocol name

### 2. `pick_pack`

Purpose:

- seller asks warehouse to pick and pack the order
- warehouse reports packed or refused

Payload:

```text
{
  "kind": "request" / "result",
  "customer_order_ref": tstr,
  "parent_order_cid": 42(bstr),
  ? "items": [1*8 {
    "sku": tstr,
    "quantity": uint
  }],
  ? "capability_token": bstr,
  ? "service_level": "ground" / "priority",
  ? "status": "packed" / "refused",
  ? "package_id": tstr,
  ? "weight_grams": uint,
  ? "package_count": uint,
  ? "notes": tstr
}
```

### 3. `accounting`

Purpose:

- seller asks accounting to record the order financially
- accounting confirms or refuses the record before shipment booking

Payload:

```text
{
  "kind": "request" / "result",
  "customer_order_ref": tstr,
  "parent_order_cid": 42(bstr),
  ? "payment_ref": tstr,
  ? "capability_token": bstr,
  ? "currency": "USD",
  ? "amount_cents": uint,
  ? "status": "recorded" / "refused",
  ? "ledger_entry_id": tstr,
  ? "invoice_ref": tstr,
  ? "notes": tstr
}
```

### 4. `shipment`

Purpose:

- seller asks carrier to book shipment after accounting succeeded
- carrier reports shipment booked or refused

Payload:

```text
{
  "kind": "request" / "result",
  "customer_order_ref": tstr,
  "parent_order_cid": 42(bstr),
  ? "parent_pick_pack_cid": 42(bstr),
  ? "parent_accounting_cid": 42(bstr),
  ? "capability_token": bstr,
  ? "package_id": tstr,
  ? "weight_grams": uint,
  ? "ship_to": {
    "name": tstr,
    "address1": tstr,
    ? "address2": tstr,
    "city": tstr,
    "region": tstr,
    "postal_code": tstr,
    "country": tstr
  },
  ? "service_level": "ground" / "priority",
  ? "status": "booked" / "refused",
  ? "carrier_name": tstr,
  ? "tracking_number": tstr,
  ? "label_artifact_cid": 42(bstr),
  ? "notes": tstr
}
```

### 5. `kernel_register`

Reuse the existing local app/kernel registration protocol semantics rather than
inventing another one, but refer to it in this doc without a version suffix and
apply the same signature and capability-token requirements.

## Message DAG Expectations

The message DAG should be easy to inspect:

```text
order submit
  -> pick_pack request
    -> pick_pack result
      -> accounting request
        -> accounting result
          -> shipment request
            -> shipment result
              -> order final
```

Important:

- every saved raw message should be a `.cbor` artifact named by its CID
- parent links should reference exact envelope CIDs, not local sequence numbers
- signed capability-token artifacts should be inspectable alongside the message
  DAG
- `pg-inspect` should remain able to decode the resulting artifacts

## Determinism Requirements

The MVP should be deterministic end to end.
We will follow up with LLM-based agents in the next phase.

Recommended deterministic behavior:

- warehouse:
  - `widget-1` succeeds with fixed package count and weight
  - `widget-oos` refuses
- accounting:
  - success returns deterministic ledger IDs derived from request CID
  - duplicate fixture refuses based on `payment_ref`
- carrier:
  - success returns deterministic tracking numbers derived from request CID
  - timeout fixture sleeps past seller sub-timeout

No live APIs in MVP:

- no UPS API
- no payment gateway
- no LLM decision point

If an LLM is added later, it should be for a clearly host-local, bounded,
non-protocol-critical decision such as operator-facing exception summaries.

## Containerization Requirements

Add a container deployment because the demo value depends on visible process
separation.

Recommended shape:

- one shared image
- one compose file
- one service per long-running role
- one short-lived intake submission command

Suggested compose services:

- `kernel`
- `seller`
- `warehouse`
- `accounting`
- `carrier`
- `intake` used via `docker compose run --rm intake ...`

Suggested environment:

- `ROLE=kernel|seller|warehouse|accounting|carrier|intake`
- `PG_KERNEL_ADDR=kernel:7000`
- `PG_DATA_DIR=/data/<role>`
- `PG_TIMEOUT=...`

Suggested volumes:

- bind mount `/tmp/grid-examples-data:/data`
- each role writes only under its own directory or namespace

Do not use a shared mutable DB as the coordination mechanism.

## Artifacts And Observability

Use the collector and analyzer from POC16 to observe the message DAG, inspect
signed traffic, and inspect capability-token artifacts.

## Recommended Implementation Order

### Phase 1 - Protocol docs and payload types

- add the order-fulfillment draft spec docs
- derive content-addressed draft pCIDs from exact doc bytes
- add payload validation code
- add signature and capability-token validation code

### Phase 2 - Deterministic local agents

- implement seller, warehouse, accounting, carrier, and intake agents
- keep the message flow and protocol families fixed
- make refusal paths explicit for malformed payloads, bad signatures, and bad
  capability tokens

### Phase 3 - Container entrypoints

- add role-based binary entrypoints
- add Dockerfile and compose file
- keep the same message flow and protocols

### Phase 4 - Operator UX

- add fixture order files
- add README for the order demo
- document how to run happy path and failure cases

### Phase 5 - Stretch alignment features

- printer-port or postal-scale split
- local trust accounting per promise type
- richer signed-envelope metadata

## Acceptance Criteria

The handoff should be considered complete only when the next session can build
all of the following:

- a deterministic `go test` path for the new order demo
- a happy-path run that produces:
  - `order_status = fulfilled`
  - `promise_status = kept`
  - tracking number
  - accounting record ID
- at least two failure fixtures:
  - warehouse refusal
  - carrier timeout or accounting refusal
- one container per long-running agent plus kernel
- all inter-agent communication over the grid envelope through the kernel
- only four business-domain pCIDs:
  - `order`
  - `pick_pack`
  - `accounting`
  - `shipment`
- no payload field that repeats the protocol name
- actual signature creation and verification on every message path
- actual cryptographic capability-token issuance, presentation, and verification
  wherever the selected protocol profile requires a token
- explicit refusal handling for missing, malformed, expired, or unauthorized
  capability tokens
- explicit refusal handling for signature verification failure
- raw message artifacts and message DAG available for inspection

## Explicit Non-Goals For MVP

- real carrier APIs
- real payment processing
- browser UI
- global trust score
- generic workflow-engine abstraction
- central DB as the source of truth for inter-agent state

## Stretch After MVP

If the first deterministic containerized demo works, the best aligned next
steps are:

- split warehouse/device concerns into postal-scale and label-printer agents
- add richer issuer-local capability-token patterns for hardware access
- optionally add one bounded LLM-powered exception-summary path that does not
  own protocol validity

## Instructions To The Next Codex Session

- build the deterministic in-process version first
- keep business outcome and promise outcome separate
- do not sneak in direct RPC just because the agents live in one repo
- preserve raw-message inspectability throughout
- do not add extra pCIDs unless the payload family is materially different
