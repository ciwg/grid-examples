#!/usr/bin/env bash
set -euo pipefail

# Intent: Record the real Grid TUI against an isolated live relay so the demo
# shows the same Ex3 collaboration path used in development. Source: DI-mutoh.
repo_dir=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
vhs_bin=${VHS_BIN:-}
if [ -z "$vhs_bin" ]; then
  if command -v vhs >/dev/null 2>&1; then
    vhs_bin=$(command -v vhs)
  else
    go_bin=$(go env GOPATH)/bin/vhs
    if [ -x "$go_bin" ]; then
      vhs_bin=$go_bin
    else
      printf '%s\n' 'VHS is required. Install it from https://charm.land/vhs or set VHS_BIN.' >&2
      exit 1
    fi
  fi
fi

runtime_dir=$(mktemp -d "${TMPDIR:-/tmp}/grid-tui-vhs.XXXXXX")
relay_log=$runtime_dir/relay.log
relay_pid=''
cleanup() {
  status=$?
  if [ -n "$relay_pid" ]; then
    if kill -0 "$relay_pid" 2>/dev/null; then
      kill "$relay_pid"
      wait "$relay_pid" || status=$?
    fi
  fi
  rm -rf "$runtime_dir"
  exit "$status"
}
trap cleanup EXIT INT TERM

cd "$repo_dir"
go run ./cmd/grid-relay --listen 127.0.0.1:7049 --data-root "$runtime_dir/relay" >"$relay_log" 2>&1 &
relay_pid=$!
for _ in $(seq 1 40); do
  if grep -q 'grid-relay listening' "$relay_log"; then
    break
  fi
  if ! kill -0 "$relay_pid" 2>/dev/null; then
    cat "$relay_log" >&2
    exit 1
  fi
  sleep 0.25
done
if ! grep -q 'grid-relay listening' "$relay_log"; then
  printf '%s\n' 'Timed out waiting for isolated Grid relay.' >&2
  exit 1
fi
"$vhs_bin" demos/grid-tui-collaboration.tape
