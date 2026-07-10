
## Goal

Build a PromiseGrid example in this repo: a order-fulfillment demo that makes
multi-agent grid messaging visible.

Refer to: 
- PromiseGrid dev guide: ~/lab/cswg/promisegrid-dev-guide/README.md
- Wire-lab repo POCs 16, 18, 19, 20.  ~/lab/wire-lab 

This example should:

- demonstrate multiple business agents exchanging real grid messages
- run with one container per agent plus one kernel container
- keep the PromiseGrid contract in the message/spec layer, not in hidden
  in-process calls or a shared database
- stay aligned with the PromiseGrid dev guide, especially the current `App
  Devs` guidance around `pCID`-selected envelopes, local promise judgment,
  append-only evidence, and `poc12` hybrid fulfillment orientation

## Why This Example

The fulfillment demonstrates why the mechanics matter:

- several independent agents participate
- multiple `pCID`s are active in one workflow
- business refusal and protocol failure are different things
- timeouts, duplicate delivery, and local judgment all matter
- durable evidence and raw message artifacts are obviously useful
- per-agent containers feel like a real distributed deployment

The dev guide's `poc12` guidance is the closest upstream orientation point:
multi-`pCID` handler routing, hybrid fulfillment, postal-scale/label/accounting
style subflows, and app-local evidence without turning the kernel into the
business workflow owner.

## Non-Negotiable Alignment Points

The next Codex session should treat these as hard constraints:

- All inter-agent communication must be real grid traffic:
  `grid([42(pCID), payload, ...protocol-owned-slots])`.
- The kernel parses slot `0` only, routes by `pCID` only, and forwards exact
  bytes only.
- The kernel is not a router-with-business-logic, service registry, or RPC
  authority.
- Promise judgment stays in app agents, not in the kernel.
- Malformed bytes are quarantined or declined locally, never silently
  reinterpreted.
- Evidence, journals, and raw message artifacts are append-only.
- Process shape is host-local. The protocol contract is the stable thing.
- Do not copy `poc12` names, route strings, or toy payloads as if they were
  final API.
- No hidden direct calls between business agents. If one agent needs another,
  it sends a message through the kernel boundary.
- The MVP should be deterministic. Do not require an LLM for the first working
  slice.

## Repo Strategy

Recommended direction:

- prefer an isolated namespace such as:
  - `cmd/pg-order-agent`
  - `cmd/pg-order-submit`
  - `internal/orderdemo/...`
  - `docs/order-fulfillment/...`
  - `deploy/order-fulfillment/docker-compose.yml`

Recommended binary shape:
- a different binary for each agent.

## Runtime Topology

Rules:

- every business agent dials the kernel only
- no direct seller->warehouse TCP socket outside the kernel path
- the submitter/intake agent may run as a short-lived container invocation
- each agent has its own local CAS store 
- shared storage may be used for artifact persistence, but not as the live
  communication mechanism
- use the microkernel style supervisor, along with the stdout collector and analyzer from POC16

## Agent Roles And Promises

### Intake agent

Responsibility:

- accept one order fixture from the operator
- send the top-level order request
- wait for one final order result
- judge the seller's promise locally
- submit top-level evidence to the evidence agent

Promise:

> I promise to submit this order for fulfillment and record the resulting
> promise judgment as evidence.

### Seller agent

Responsibility:

- receive `order_request_v1`
- validate order shape and local policy
- orchestrate warehouse, carrier, and accounting substeps by message
- send one final `order_result_v1`

Promise:

> I promise to process this order request under the selected protocol and send
> one conforming final result describing fulfillment, refusal, or failure.

Important:

- `fulfilled`, `refused`, and `failed` are business outcomes
- a conforming `refused` result can still mean the seller kept its protocol
  promise

### Warehouse agent

Responsibility:

- receive `pick_pack_request_v1`
- deterministically decide whether items can be picked/packed
- emit one `pick_pack_result_v1`

Promise:

> I promise to return one conforming pick/pack result for this work request.

### Carrier agent

Responsibility:

- receive `shipment_request_v1`
- deterministically book or refuse shipment
- emit one `shipment_result_v1`

Promise:

> I promise to return one conforming shipment-booking result for this shipment
> request.

### Accounting agent

Responsibility:

- receive `accounting_request_v1`
- deterministically record or refuse accounting state
- emit one `accounting_result_v1`

Promise:

> I promise to return one conforming accounting result for this accounting
> request.

### Evidence agent

Responsibility:

- receive top-level `order_evidence_v1`
- append it durably and idempotently

Promise:

> I promise to record the order-result judgment as evidence.

### Kernel

Responsibility:

- receive registration messages
- parse slot `0` only
- route exact bytes
- emit kernel-local operational events only

The kernel never claims whether the seller, warehouse, carrier, or accounting
promises were kept.

## MVP Scenario

The first working slice should be intentionally narrow:

- one order at a time
- one seller
- one warehouse
- one carrier
- one accounting agent
- one evidence agent
- one kernel
- one intake submission per run
- deterministic business logic only
- no real external APIs
- no signatures
- no capability tokens in MVP
- no LLMs in MVP

Use a tiny local catalog and scripted behavior:

- `widget-1` succeeds
- `widget-oos` causes warehouse refusal
- `widget-carrier-timeout` causes seller-visible shipment timeout/failure
- `widget-dup-pay` causes accounting refusal

## Business Flow

Happy path:

1. Intake sends `order_request_v1` to seller.
2. Seller validates the request.
3. Seller sends `pick_pack_request_v1` to warehouse.
4. Warehouse returns `pick_pack_result_v1`.
5. Seller sends `shipment_request_v1` to carrier.
6. Carrier returns `shipment_result_v1`.
7. Seller sends `accounting_request_v1` to accounting.
8. Accounting returns `accounting_result_v1`.
9. Seller sends `order_result_v1` to intake.
10. Intake judges the seller's promise and sends `order_evidence_v1`.
11. Evidence agent records the evidence append-only.

Failure examples:

- warehouse refusal:
  - warehouse returns a conforming refusal result
  - seller returns a conforming top-level `order_result_v1` with
    `order_status = "refused"`
  - intake judges seller promise `kept` if the result conforms
- shipment timeout:
  - seller times out waiting for carrier
  - seller may emit a conforming `order_result_v1` with
    `order_status = "failed"` and `failure_stage = "carrier"`
  - if seller itself stays silent past the intake timeout, intake records
    `not_promised`
- accounting refusal:
  - accounting returns a conforming refusal result
  - seller returns a conforming failure/refusal result

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
- a malformed top-level result is `promise_status = broken`
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
  "request_id": "b...",
  "order_result_id": "b...",
  "evidence_recorded": true,
  "notes": "conforming order result received for this request envelope"
}
```

## Draft Protocol Set

Use multiple `pCID`s on purpose. The point is to demonstrate that one runtime
can route several business protocols through one kernel boundary.

### 1. `order_request_v1`

Purpose:

- intake asks seller to process one order

Payload:

```text
{
  "protocol": "order_request_v1",
  "customer_order_ref": tstr,
  "requested_by": tstr,
  "ship_to": {
    "name": tstr,
    "address1": tstr,
    ? "address2": tstr,
    "city": tstr,
    "region": tstr,
    "postal_code": tstr,
    "country": tstr
  },
  "items": [1*8 {
    "sku": tstr,
    "quantity": uint
  }],
  "service_level": "ground" / "priority",
  "payment_ref": tstr,
  ? "notes": tstr
}
```

### 2. `pick_pack_request_v1`

Purpose:

- seller asks warehouse to pick and pack the order

Payload:

```text
{
  "protocol": "pick_pack_request_v1",
  "parent_order_request_cid": 42(bstr),
  "customer_order_ref": tstr,
  "items": [1*8 {
    "sku": tstr,
    "quantity": uint
  }],
  "service_level": "ground" / "priority"
}
```

### 3. `pick_pack_result_v1`

Purpose:

- warehouse reports picked/packed or refused

Payload:

```text
{
  "protocol": "pick_pack_result_v1",
  "parent_pick_pack_request_cid": 42(bstr),
  "customer_order_ref": tstr,
  "status": "packed" / "refused",
  ? "package_id": tstr,
  ? "weight_grams": uint,
  ? "package_count": uint,
  "notes": tstr
}
```

### 4. `shipment_request_v1`

Purpose:

- seller asks carrier to book shipment

Payload:

```text
{
  "protocol": "shipment_request_v1",
  "parent_order_request_cid": 42(bstr),
  "parent_pick_pack_result_cid": 42(bstr),
  "customer_order_ref": tstr,
  "package_id": tstr,
  "weight_grams": uint,
  "ship_to": { ...same fields as order... },
  "service_level": "ground" / "priority"
}
```

### 5. `shipment_result_v1`

Purpose:

- carrier reports shipment booked or refused

Payload:

```text
{
  "protocol": "shipment_result_v1",
  "parent_shipment_request_cid": 42(bstr),
  "customer_order_ref": tstr,
  "status": "booked" / "refused",
  ? "carrier_name": tstr,
  ? "tracking_number": tstr,
  ? "label_artifact_cid": 42(bstr),
  "notes": tstr
}
```

### 6. `accounting_request_v1`

Purpose:

- seller asks accounting to record the order financially

Payload:

```text
{
  "protocol": "accounting_request_v1",
  "parent_order_request_cid": 42(bstr),
  "customer_order_ref": tstr,
  "payment_ref": tstr,
  "currency": "USD",
  "amount_cents": uint
}
```

### 7. `accounting_result_v1`

Purpose:

- accounting confirms or refuses the record

Payload:

```text
{
  "protocol": "accounting_result_v1",
  "parent_accounting_request_cid": 42(bstr),
  "customer_order_ref": tstr,
  "status": "recorded" / "refused",
  ? "ledger_entry_id": tstr,
  ? "invoice_ref": tstr,
  "notes": tstr
}
```

### 8. `order_result_v1`

Purpose:

- seller reports the final business outcome back to intake

Payload:

```text
{
  "protocol": "order_result_v1",
  "parent_order_request_cid": 42(bstr),
  "customer_order_ref": tstr,
  "order_status": "fulfilled" / "refused" / "failed",
  ? "failure_stage": "seller_validation" / "warehouse" / "carrier" / "accounting",
  ? "warehouse_result_cid": 42(bstr),
  ? "shipment_result_cid": 42(bstr),
  ? "accounting_result_cid": 42(bstr),
  ? "package_id": tstr,
  ? "tracking_number": tstr,
  ? "ledger_entry_id": tstr,
  "summary": tstr
}
```

### 9. `order_evidence_v1`

Purpose:

- intake records the judgment of the seller's promise

Payload:

```text
{
  "protocol": "order_evidence_v1",
  "request_id": 42(bstr),
  "order_result_id": 42(bstr) / null,
  "promiser": tstr,
  "promisee": tstr,
  "promise_status": "kept" / "broken" / "not_promised",
  "notes": tstr
}
```

### 10. `kernel_register_v1`

Reuse the existing local app/kernel registration protocol rather than inventing
another one.

## Message DAG Expectations

The message DAG should be easy to inspect:

```text
order_request
  -> pick_pack_request
    -> pick_pack_result
      -> shipment_request
        -> shipment_result
          -> accounting_request
            -> accounting_result
              -> order_result
                -> order_evidence
```

Important:

- every saved raw message should be a `.cbor` artifact named by its CID
- parent links should reference exact envelope CIDs, not local sequence numbers
- `pg-inspect` should remain able to decode the resulting artifacts

## Determinism Requirements

The MVP should be deterministic end to end.
We will follow up with LLM based agents in the next phase. 

Recommended deterministic behavior:

- warehouse:
  - `widget-1` succeeds with fixed package count/weight
  - `widget-oos` refuses
- carrier:
  - success returns deterministic tracking numbers derived from request CID
  - timeout fixture sleeps past seller sub-timeout
- accounting:
  - success returns deterministic ledger ids derived from request CID
  - duplicate fixture refuses based on `payment_ref`

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
- `carrier`
- `accounting`
- `evidence`
- `intake` (used via `docker compose run --rm intake ...`)

Suggested environment:

- `ROLE=kernel|seller|warehouse|carrier|accounting|evidence|intake`
- `PG_KERNEL_ADDR=kernel:7000`
- `PG_DATA_DIR=/data/<role>`
- `PG_TIMEOUT=...`

Suggested volumes:

- bind mount `/tmp/grid-examples-data:/data`
- each role writes only under its own directory or namespace

Do not use a shared mutable DB as the coordination mechanism.

## Artifacts And Observability

We will use the collector and analyzer from POC16 to observe the message DAG and inspect artifacts.

## Recommended Implementation Order

### Phase 1 - Protocol docs and payload types

- add the order-fulfillment draft spec docs
- derive content-addressed draft pCIDs from exact doc bytes
- add payload validation code

### Phase 3 - Container entrypoints

- implement seller, warehouse, carrier, accounting, intake, evidence agents
- add role-based binary entrypoints
- add Dockerfile and compose file
- keep the same message flow and protocols

### Phase 4 - Operator UX

- add fixture order files
- add README for the order demo
- document how to run happy path and failure cases

### Phase 5 - Stretch alignment features

- printer-port or postal-scale split
- capability token issue/redemption for scarce device access
- signed envelopes / proof slot
- local per-context trust ledger

## Acceptance Criteria

The handoff should be considered complete only when the next session can build
all of the following:

- a deterministic `go test` path for the new order demo
- a happy-path run that produces:
  - `order_status = fulfilled`
  - `promise_status = kept`
  - tracking number
  - accounting record id
  - `evidence_recorded = true`
- at least two failure fixtures:
  - warehouse refusal
  - carrier timeout or accounting refusal
- one container per long-running agent plus kernel
- all inter-agent communication over the grid envelope through the kernel
- raw message artifacts and message DAG available for inspection
- signatures / proof slot
- capability tokens in the first slice
- add local trust accounting per promise type

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
- add issuer-local capability-token issue/redemption for hardware access
- optionally add one bounded LLM-powered exception-summary path that does not
  own protocol validity

## Instructions To The Next Codex Session

- build the deterministic in-process version first
- keep business outcome and promise outcome separate
- do not sneak in direct RPC just because the agents live in one repo
- preserve raw-message inspectability throughout
