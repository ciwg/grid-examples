#!/usr/bin/env bash
set -eu

repo_root="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
demo_root="/tmp/ex5-demo-browser"
runtime_root="${demo_root}/runtime"
socket_path="${runtime_root}/embodiment.sock"
host_binary="${demo_root}/native-host/operational-browser-host"
host_manifest_chrome="/home/jj/.config/google-chrome/NativeMessagingHosts/operational_browser_host.json"
expected_extension_origin="chrome-extension://miagfmaampfgjkojhccdilogehbjijpe/"
browser_base_url="http://127.0.0.1:7045"
chrome_debug_port="${EX5_CHROME_DEBUG_PORT:-9222}"

# Intent: Fail closed unless the browser demo's native-host and runtime bridge
# are actually ready, so the one-sheet path only claims readiness after the
# direct browser embodiment dependencies are proven. Source: DI-dabek

if [ ! -S "${socket_path}" ]; then
  echo "runtime socket missing: ${socket_path}" >&2
  exit 1
fi

if [ ! -x "${host_binary}" ]; then
  echo "native host binary missing: ${host_binary}" >&2
  exit 1
fi

if [ ! -f "${host_manifest_chrome}" ]; then
  echo "chrome native-host manifest missing: ${host_manifest_chrome}" >&2
  exit 1
fi

if ! grep -F "\"${expected_extension_origin}\"" "${host_manifest_chrome}" >/dev/null 2>&1; then
  echo "chrome native-host manifest does not allow the shipped extension origin" >&2
  exit 1
fi

if ! grep -F "\"path\": \"${host_binary}\"" "${host_manifest_chrome}" >/dev/null 2>&1; then
  echo "chrome native-host manifest does not point at the demo host binary" >&2
  exit 1
fi

meta_json=""
for _ in $(seq 1 40); do
  if meta_json="$(curl -s "${browser_base_url}/api/meta" 2>/dev/null)"; then
    if [ -n "${meta_json}" ]; then
      break
    fi
  fi
  sleep 0.25
done

if [ -z "${meta_json}" ]; then
  echo "could not fetch ${browser_base_url}/api/meta" >&2
  exit 1
fi

python3 - "${meta_json}" <<'PY'
import json
import sys

meta = json.loads(sys.argv[1])
browser = meta.get("embodiments", {}).get("browser", {})
if browser.get("primary_adapter") != "chrome_native_messaging":
    raise SystemExit("runtime does not advertise chrome_native_messaging")
if meta.get("local_unix_socket_path") != "/tmp/ex5-demo-browser/runtime/embodiment.sock":
    raise SystemExit("127.0.0.1:7045 is not serving the demo runtime socket")
PY

python3 - "${host_binary}" "${socket_path}" <<'PY'
import json
import pathlib
import struct
import subprocess
import sys

host_binary = pathlib.Path(sys.argv[1])
socket_path = sys.argv[2]
payload = json.dumps({
    "request_id": "demo-verify-runtime-ready",
    "socket_path": socket_path,
    "request": {
        "type": "operation",
        "operation": "runtime_ready",
    },
}).encode()
message = struct.pack("<I", len(payload)) + payload
completed = subprocess.run([str(host_binary)], input=message, capture_output=True, check=True)
if len(completed.stdout) < 4:
    raise SystemExit("native host returned no framed response")
size = struct.unpack("<I", completed.stdout[:4])[0]
response = json.loads(completed.stdout[4:4 + size])
if response.get("error"):
    raise SystemExit(response["error"])
body = response.get("response", {}).get("body")
if body != '{"ready":true}':
    raise SystemExit(f"unexpected runtime_ready body: {body!r}")
PY

# Intent: Verify the actual Chrome extension/page/native-host path, not only
# the host binary, before claiming the browser demo is ready. Source: DI-vuduz.
node - "${chrome_debug_port}" "${socket_path}" <<'JS'
const [debugPort, socketPath] = process.argv.slice(2);

async function main() {
  const targets = await (await fetch(`http://127.0.0.1:${debugPort}/json/list`)).json();
  const page = targets.find((target) => target.type === "page" && target.url === "http://127.0.0.1:7045/");
  if (!page) throw new Error("demo page is not a Chrome DevTools target");
  const socket = new WebSocket(page.webSocketDebuggerUrl);
  await new Promise((resolve, reject) => { socket.addEventListener("open", resolve, { once: true }); socket.addEventListener("error", reject, { once: true }); });
  const expression = `new Promise((resolve) => { const requestID = "verify-page-bridge-" + Date.now(); const timer = setTimeout(() => done({ ok: false, error: "page bridge timed out" }), 2000); const done = (result) => { clearTimeout(timer); window.removeEventListener("message", listener); resolve(result); }; const listener = (event) => { const message = event.data; if (event.source === window && message && message.__oks_bridge === true && message.direction === "bridge->page" && message.request_id === requestID && message.kind === "handshake") done({ ok: !!message.ok, error: message.error || "" }); }; window.addEventListener("message", listener); window.postMessage({ __oks_bridge: true, direction: "page->bridge", kind: "handshake", request_id: requestID, socket_path: ${JSON.stringify(socketPath)} }, window.location.origin); })`;
  const result = await new Promise((resolve, reject) => {
    socket.addEventListener("message", (event) => { const message = JSON.parse(event.data); if (message.id === 1) resolve(message); });
    socket.send(JSON.stringify({ id: 1, method: "Runtime.evaluate", params: { expression, awaitPromise: true, returnByValue: true } }));
    setTimeout(() => reject(new Error("Chrome DevTools page probe timed out")), 5000);
  });
  socket.close();
  const value = result.result?.result?.value;
  if (!value?.ok) throw new Error(value?.error || "page bridge failed");
}

main().catch((error) => { console.error(error.message); process.exit(1); });
JS

echo "ex5 browser demo verification passed"
echo "browser shell: ${browser_base_url}/"
echo "runtime socket: ${socket_path}"
