#!/usr/bin/env bash
set -eu

repo_root="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
demo_root="/tmp/ex5-demo-browser"
runtime_root="${demo_root}/runtime"
profile_root="${demo_root}/chrome-profile"
native_host_root="${demo_root}/native-host"
host_binary="${native_host_root}/operational-browser-host"
host_manifest_template="${repo_root}/chrome-extension/native-host/operational_browser_host.json"
host_manifest_local="${native_host_root}/operational_browser_host.json"
host_manifest_chrome="/home/jj/.config/google-chrome/NativeMessagingHosts/operational_browser_host.json"

# Intent: Turn the browser demo from assumed ambient setup into one explicit,
# reproducible preflight that prepares sample data, browser-host artifacts, and
# Chrome registration before the one-sheet demo path is used. Source: DI-dabek

if ! command -v go >/dev/null 2>&1; then
  echo "go is required" >&2
  exit 1
fi

mkdir -p "${demo_root}" "${native_host_root}" "$(dirname -- "${host_manifest_chrome}")"

if [ ! -d "${runtime_root}" ]; then
  "${repo_root}/scripts/load-sample-data.sh" "${runtime_root}"
fi

if [ ! -d "${profile_root}" ]; then
  mkdir -p "${profile_root}"
fi

(
  cd "${repo_root}"
  go build -o "${host_binary}" ./cmd/operational-browser-host
)
chmod 755 "${host_binary}"

python3 - "${host_manifest_template}" "${host_binary}" "${host_manifest_local}" "${host_manifest_chrome}" <<'PY'
import pathlib
import sys

template_path = pathlib.Path(sys.argv[1])
host_binary = pathlib.Path(sys.argv[2])
local_manifest = pathlib.Path(sys.argv[3])
chrome_manifest = pathlib.Path(sys.argv[4])

text = template_path.read_text()
text = text.replace("__BROWSER_HOST_PATH__", str(host_binary))
local_manifest.write_text(text)
chrome_manifest.write_text(text)
PY

echo "ex5 demo browser setup is ready"
echo "runtime root: ${runtime_root}"
echo "chrome profile: ${profile_root}"
echo "native host binary: ${host_binary}"
echo "chrome native-host manifest: ${host_manifest_chrome}"
