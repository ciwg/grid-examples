#!/usr/bin/env bash
set -eu

repo_root="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
demo_root="/tmp/ex5-demo-browser"
runtime_root="${demo_root}/runtime"
socket_path="${runtime_root}/embodiment.sock"
profile_root="${demo_root}/chrome-profile"
runtime_log="${demo_root}/runtime.log"
runtime_pid_file="${demo_root}/runtime.pid"
browser_pid_file="${demo_root}/chrome.pid"
extension_root="${repo_root}/chrome-extension"
extension_launcher="${repo_root}/scripts/install-demo-chrome-extension.mjs"
browser_base_url="http://127.0.0.1:7045"
browser_url="${browser_base_url}/"
chrome_debug_port="${EX5_CHROME_DEBUG_PORT:-9222}"

# Intent: Launch the browser demo through one dedicated Chrome profile and the
# shipped extension so the one-sheet path does not depend on ambient browser
# state or manual extension loading. Source: DI-dabek

if [ ! -d "${runtime_root}" ] || [ ! -d "${profile_root}" ]; then
  echo "demo browser setup missing; run scripts/setup-demo-browser.sh first" >&2
  exit 1
fi

browser_bin="google-chrome"
if ! command -v "${browser_bin}" >/dev/null 2>&1; then
  echo "could not find verified browser: google-chrome is required; Chromium remains unverified" >&2
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

# Intent: Replace only the prior Chrome process recorded by this disposable
# demo launcher so attach-only checks cannot inherit a stale demo-profile
# session. Source: DI-danir.
if [ -f "${browser_pid_file}" ]; then
  existing_browser_pid="$(cat "${browser_pid_file}")"
  if [ -n "${existing_browser_pid}" ] && kill -0 "${existing_browser_pid}" >/dev/null 2>&1; then
    kill "${existing_browser_pid}"
    for _ in $(seq 1 40); do
      if ! kill -0 "${existing_browser_pid}" >/dev/null 2>&1; then
        break
      fi
      sleep 0.25
    done
    if kill -0 "${existing_browser_pid}" >/dev/null 2>&1; then
      echo "previous demo Chrome process did not exit: ${existing_browser_pid}" >&2
      exit 1
    fi
  fi
  rm -f "${browser_pid_file}"
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

"${extension_launcher}" "${profile_root}" "${extension_root}" "${chrome_debug_port}" "${browser_url}" >"${demo_root}/chrome-launcher.log" 2>&1 &
echo "$!" > "${browser_pid_file}"

debug_endpoint="http://127.0.0.1:${chrome_debug_port}/json/version"
debug_metadata=""
for _ in $(seq 1 40); do
  if debug_metadata="$(curl -s "${debug_endpoint}" 2>/dev/null)"; then
    if [ -n "${debug_metadata}" ]; then
      break
    fi
  fi
  sleep 0.25
done
if [ -z "${debug_metadata}" ]; then
  echo "Chrome DevTools did not become reachable at ${debug_endpoint}" >&2
  exit 1
fi

echo "launched ${browser_bin} for the ex5 browser demo"
echo "browser url: ${browser_url}"
echo "runtime log: ${runtime_log}"
echo "Chrome DevTools: http://127.0.0.1:${chrome_debug_port}/json/version"
