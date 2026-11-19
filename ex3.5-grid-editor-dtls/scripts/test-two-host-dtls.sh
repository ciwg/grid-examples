#!/usr/bin/env bash
# Proves the documented private-network deployment against a second host.
# Required environment: DTLS_REMOTE_SSH and DTLS_REMOTE_ROOT.
set -euo pipefail

if [ -z "${DTLS_REMOTE_SSH:-}" ] || [ -z "${DTLS_REMOTE_ROOT:-}" ]; then
  echo "DTLS_REMOTE_SSH and DTLS_REMOTE_ROOT are required" >&2
  exit 64
fi

local_root=$(mktemp -d /tmp/ex3.5-dtls-two-host.XXXXXX)
remote_log="$DTLS_REMOTE_ROOT/.dtls-two-host.log"
cleanup() {
  local status=$?
  if kill -0 "${local_pid:-0}" 2>/dev/null; then kill "$local_pid"; fi
  if ! ssh "$DTLS_REMOTE_SSH" "pkill -f 'grid-relay.*17526' >/dev/null 2>&1"; then
    echo "remote cleanup did not find a running node" >&2
  fi
  rm -rf "$local_root"
  exit "$status"
}
trap cleanup EXIT

remote_ip=$(ssh "$DTLS_REMOTE_SSH" "tailscale ip -4")
ssh "$DTLS_REMOTE_SSH" "cd '$DTLS_REMOTE_ROOT' && nohup go run ./cmd/grid-relay --listen 127.0.0.1:17525 --dtls-listen :17526 --dtls-advertise '$remote_ip:17526' --data-root /tmp/ex3.5-dtls-remote >'$remote_log' 2>&1 &"

for _ in $(seq 1 40); do
  invite=$(ssh "$DTLS_REMOTE_SSH" "sed -n 's/.*dtls invitation=//p' '$remote_log' | head -1")
  if [ -n "$invite" ]; then break; fi
  sleep 1
done
test -n "${invite:-}"

go run ./cmd/grid-relay --listen 127.0.0.1:17625 --dtls-listen :17626 --dtls-peer "$invite" --data-root "$local_root" >"$local_root/local.log" 2>&1 &
local_pid=$!
sleep 2
curl -fsS -X POST "http://127.0.0.1:17625/api/local/documents/demo/sync" -H 'Content-Type: application/json' --data '{"participant_id":"two-host","recipient_id":"remote","message_base64":"AQID","embodiment":"browser"}' >/dev/null
sleep 2
ssh "$DTLS_REMOTE_SSH" "curl -fsS 'http://127.0.0.1:17525/api/local/documents/demo/sync?since=0' | grep -q AQID"
echo "two-host direct DTLS proof passed"
