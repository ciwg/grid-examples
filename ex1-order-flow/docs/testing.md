# Ex1 testing guide

Ex1's tests verify the published behavior of its local-draft profiles. They do
not claim conformance to a frozen upstream PromiseGrid specification. Source:
DI-garis, DI-potoj.

## Run the suite

From `ex1-order-flow/`:

```bash
go test ./...
errcheck ./...
```

## What the tests prove

- `protocol/profile_inventory_test.go` derives each pCID from the exact local
  profile source bytes and checks that both the design inventory and local
  implementation-scope declaration publish that identity.
- `artifact/store_test.go` proves raw bytes can be read back by their retained
  CID after reopening a local store, and that local observations append as
  separate JSONL records.
- `agent/client_test.go` proves malformed input, invalid proof, and unexpected
  pCID receipt are locally recorded without treating them as another agent's
  promise or refusal.
- `agent/message_test.go` proves capability validation rejects a token whose
  authority does not match the receiving context.
- `e2e/e2e_test.go` exercises the full local topology. The warehouse-refusal
  scenario checks the seller's `refusal_observed` evidence; the carrier-timeout
  scenario checks its `timeout_observed` evidence names an expected request.

## Evidence boundaries

Each test uses a Go-managed temporary root. The tests inspect only the local
role's retained raw bytes and `observations.jsonl`; a timeout or refusal
observation is not proof of another agent's intent. The kernel's separate
`no_registered_recipient` behavior is covered in `kernel/server_test.go`.
Source: DI-vihoz, DI-riguz, DI-purum, DI-zosiz, DI-potoj.
