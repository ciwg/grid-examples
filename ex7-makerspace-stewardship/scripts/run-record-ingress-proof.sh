#!/usr/bin/env bash
set -euo pipefail

# Intent: Prove configured public-key recognition and exact signed ingress
# through the real command without browser or account authoring. Source: DI-likoh.
proof_root=$(mktemp -d /tmp/ex7-record-proof.XXXXXX)
server_log="$proof_root/server.log"
cleanup() {
  if [ "$server_pid" -gt 0 ] && kill -0 "$server_pid" 2>/dev/null; then
    kill "$server_pid"
    wait "$server_pid" || status=$?
    if [ "${status:-0}" -ne 0 ] && [ "${status:-0}" -ne 143 ]; then
      exit "$status"
    fi
  fi
  rm -rf "$proof_root"
}
server_pid=0
trap cleanup EXIT

record=$(go run ./cmd/makerspace-record-fixture -runtime-root "$proof_root")
go run ./cmd/makerspace-stewardship -runtime-root "$proof_root" >"$server_log" 2>&1 &
server_pid=$!
for attempt in $(seq 1 30); do
  if curl -fsS http://127.0.0.1:7037/api/state >/dev/null 2>&1; then
    break
  fi
  sleep 0.1
done
response=$(curl -fsS -X POST http://127.0.0.1:7037/api/records -H 'Content-Type: application/json' -d "{\"records\":[\"$record\"]}")
if ! printf '%s' "$response" | rg -q 'Guard is loose'; then
  echo "signed record was not projected" >&2
  exit 1
fi
echo "record ingress proof passed"
unrecognized=$(go run ./cmd/makerspace-record-fixture -runtime-root "$proof_root" -label mallory -seed mallory-proof -observation 'Mallory claim' -write-recognition=false)
response=$(curl -fsS -X POST http://127.0.0.1:7037/api/records -H 'Content-Type: application/json' -d "{\"records\":[\"$unrecognized\"]}")
if printf '%s' "$response" | rg -q 'Mallory claim'; then
  echo "unrecognized record was projected" >&2
  exit 1
fi
hold=$(go run ./cmd/makerspace-record-fixture -runtime-root "$proof_root" -kind safety-hold -observation 'Guard is loose')
clear=$(go run ./cmd/makerspace-record-fixture -runtime-root "$proof_root" -kind safety-clear -observation 'Alice says clear')
curl -fsS -X POST http://127.0.0.1:7037/api/records -H 'Content-Type: application/json' -d "{\"records\":[\"$hold\",\"$clear\"]}" >/dev/null
state=$(curl -fsS http://127.0.0.1:7037/api/state)
if ! printf '%s' "$state" | rg -q 'Safety hold'; then
  echo "recognized non-steward cleared safety hold" >&2
  exit 1
fi
echo "negative record ingress proofs passed"
