#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
port="${GRID_EDITOR_DEMO_PORT:-7038}"
relay_url="http://127.0.0.1:${port}"
demo_root="$(mktemp -d "${TMPDIR:-/tmp}/grid-editor-nvim-demo.XXXXXX")"
relay_log="${demo_root}/relay.log"
relay_pid=""

cleanup() {
  trap - EXIT INT TERM
  if [ -n "${relay_pid}" ] && kill -0 "${relay_pid}" 2>/dev/null; then
    kill "${relay_pid}"
  fi
  if [ -n "${relay_pid}" ]; then
    if wait "${relay_pid}"; then
      :
    else
      relay_status=$?
      printf 'Relay exited with status %s during cleanup.\n' "${relay_status}" >&2
    fi
  fi
  rm -rf -- "${demo_root}"
}
trap cleanup EXIT INT TERM

for required_command in curl go nvim xterm; do
  if ! command -v "${required_command}" >/dev/null 2>&1; then
    printf 'Missing required command: %s\n' "${required_command}" >&2
    exit 1
  fi
done

cd "${repo_root}"
go run ./cmd/grid-relay --listen "127.0.0.1:${port}" --data-root "${demo_root}/relay" >"${relay_log}" 2>&1 &
relay_pid=$!

relay_ready=0
for attempt in $(seq 1 40); do
  if curl --fail --silent "${relay_url}/api/meta" >/dev/null; then
    relay_ready=1
    break
  fi
  sleep 0.25
done
if [ "${relay_ready}" -ne 1 ]; then
  printf 'Relay did not become ready. Log follows:\n' >&2
  cat "${relay_log}" >&2
  exit 1
fi

# Intent: Give a presenter two immediately legible, independent Neovim
# embodiments without requiring any manual editor setup. Source: DI-loril
xterm -T 'Ex3 Neovim First' -fa Monospace -fs 18 -e env \
  GRID_EDITOR_RELAY_URL="${relay_url}" \
  GRID_EDITOR_DISPLAY_NAME='First Neovim' \
  GRID_EDITOR_COLOR='#d66f1d' \
  GRID_EDITOR_OPEN_DASHBOARD=1 \
  "${repo_root}/scripts/grid-editor-nvim" demo &

xterm -T 'Ex3 Neovim Second' -fa Monospace -fs 18 -e env \
  GRID_EDITOR_RELAY_URL="${relay_url}" \
  GRID_EDITOR_DISPLAY_NAME='Second Neovim' \
  GRID_EDITOR_COLOR='#1d6fd6' \
  GRID_EDITOR_OPEN_DASHBOARD=1 \
  "${repo_root}/scripts/grid-editor-nvim" demo &

printf 'Opened two Ex3 Neovim dashboard windows. Press Ctrl-C here when finished.\n'
wait "${relay_pid}"
