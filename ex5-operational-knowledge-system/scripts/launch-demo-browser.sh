#!/usr/bin/env bash
set -eu

repo_root="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
demo_root="/tmp/ex5-demo-browser"
runtime_root="${demo_root}/runtime"
socket_path="${runtime_root}/embodiment.sock"
profile_root="${demo_root}/chrome-profile"
runtime_log="${demo_root}/runtime.log"
runtime_pid_file="${demo_root}/runtime.pid"
extension_root="${repo_root}/chrome-extension"
browser_base_url="http://127.0.0.1:7045"
browser_url="${browser_base_url}/"

# Intent: Launch the browser demo through one dedicated Chrome profile and the
# shipped extension so the one-sheet path does not depend on ambient browser
# state or manual extension loading. Source: DI-dabek

if [ ! -d "${runtime_root}" ] || [ ! -d "${profile_root}" ]; then
  echo "demo browser setup missing; run scripts/setup-demo-browser.sh first" >&2
  exit 1
fi

browser_bin=""
for candidate in google-chrome chromium chromium-browser; do
  if command -v "${candidate}" >/dev/null 2>&1; then
    browser_bin="${candidate}"
    break
  fi
done

if [ -z "${browser_bin}" ]; then
  echo "could not find Chrome or Chromium" >&2
  exit 1
fi

if [ -f "${runtime_pid_file}" ]; then
  existing_pid="$(cat "${runtime_pid_file}")"
  if [ -n "${existing_pid}" ] && kill -0 "${existing_pid}" >/dev/null 2>&1; then
    :
  else
    rm -f "${runtime_pid_file}"
  fi
fi

if [ ! -f "${runtime_pid_file}" ]; then
  rm -f "${socket_path}"
  (
    cd "${repo_root}"
    nohup go run ./cmd/operational-knowledge -data-root "${runtime_root}" >"${runtime_log}" 2>&1 &
    echo "$!" > "${runtime_pid_file}"
  )
fi

meta_json=""
for _ in $(seq 1 40); do
  if ! kill -0 "$(cat "${runtime_pid_file}")" >/dev/null 2>&1; then
    break
  fi
  if meta_json="$(curl -s "${browser_base_url}/api/meta" 2>/dev/null)"; then
    if [ -n "${meta_json}" ]; then
      break
    fi
  fi
  sleep 0.25
done

if ! kill -0 "$(cat "${runtime_pid_file}")" >/dev/null 2>&1; then
  echo "demo runtime did not stay running; see ${runtime_log}" >&2
  exit 1
fi

if [ -z "${meta_json}" ]; then
  echo "ex5 server did not become reachable at ${browser_base_url}/api/meta" >&2
  exit 1
fi

if ! python3 - "${meta_json}" "${socket_path}" <<'PY'
import json
import sys

meta = json.loads(sys.argv[1])
expected_socket = sys.argv[2]
if meta.get("local_unix_socket_path") != expected_socket:
    raise SystemExit(1)
PY
then
  echo "127.0.0.1:7045 is not serving the demo runtime at ${socket_path}" >&2
  exit 1
fi

"${browser_bin}" \
  --user-data-dir="${profile_root}" \
  --load-extension="${extension_root}" \
  --no-first-run \
  --no-default-browser-check \
  "${browser_url}" >/dev/null 2>&1 &

echo "launched ${browser_bin} for the ex5 browser demo"
echo "browser url: ${browser_url}"
echo "runtime log: ${runtime_log}"
