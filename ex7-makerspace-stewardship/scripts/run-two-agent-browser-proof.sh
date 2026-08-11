#!/usr/bin/env bash
# Intent: Keep Alice, Bob, Chrome, and DevTools in one process namespace so
# loopback evidence is reproducible. Source: DI-fuzar.
set -euo pipefail

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
proof_root=$(mktemp -d /tmp/ex7-browser-proof.XXXXXX)
passphrase='browser-proof-passphrase'
alice_pid=0
bob_pid=0
chrome_pid=0

cleanup() {
  status=$?
  if [ "$alice_pid" -gt 0 ]; then kill "$alice_pid" 2>/dev/null || status=$?; fi
  if [ "$bob_pid" -gt 0 ]; then kill "$bob_pid" 2>/dev/null || status=$?; fi
  if [ "$chrome_pid" -gt 0 ]; then kill "$chrome_pid" 2>/dev/null || status=$?; fi
  exit "$status"
}
trap cleanup EXIT INT TERM

printf '%s' "$passphrase" | go run "$repo_root/cmd/makerspace-two-agent-bootstrap" -root "$proof_root" -passphrase-stdin
printf '%s' "$passphrase" | go run "$repo_root/cmd/makerspace-stewardship" -runtime-root "$proof_root/alice" -identity-passphrase-stdin -listen 127.0.0.1:7037 >"$proof_root/alice.log" 2>&1 &
alice_pid=$!
printf '%s' "$passphrase" | go run "$repo_root/cmd/makerspace-stewardship" -runtime-root "$proof_root/bob" -identity-passphrase-stdin -listen 127.0.0.1:7038 >"$proof_root/bob.log" 2>&1 &
bob_pid=$!
agents_ready=false
for attempt in 1 2 3 4 5; do
  if curl -fsS http://127.0.0.1:7037/api/state >/dev/null && curl -fsS http://127.0.0.1:7038/api/state >/dev/null; then agents_ready=true; break; fi
  sleep 1
done
if [ "$agents_ready" = false ]; then
  printf 'Alice or Bob did not become ready; inspect %s/alice.log and %s/bob.log\n' "$proof_root" "$proof_root" >&2
  exit 1
fi
google-chrome --headless=new --no-sandbox --remote-debugging-port=9229 --user-data-dir="$proof_root/chrome-profile" about:blank >"$proof_root/chrome.log" 2>&1 &
chrome_pid=$!
chrome_ready=false
for attempt in 1 2 3 4 5; do
  if curl -fsS http://127.0.0.1:9229/json/version >/dev/null; then chrome_ready=true; break; fi
  sleep 1
done
if [ "$chrome_ready" = false ]; then
  printf 'Chrome DevTools did not become ready; inspect %s/chrome.log\n' "$proof_root" >&2
  exit 1
fi
EX7_PROOF_ROOT="$proof_root" node "$repo_root/scripts/two-agent-browser-proof.mjs"
test -s "$proof_root/bob-final.png"
printf 'browser proof evidence: %s\n' "$proof_root"
