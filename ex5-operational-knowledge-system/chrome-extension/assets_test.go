package chromeextension

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type manifest struct {
	ManifestVersion int      `json:"manifest_version"`
	Permissions     []string `json:"permissions"`
	Background      struct {
		ServiceWorker string `json:"service_worker"`
	} `json:"background"`
	ContentScripts []struct {
		Matches []string `json:"matches"`
		JS      []string `json:"js"`
	} `json:"content_scripts"`
	Key string `json:"key"`
}

func TestChromeExtensionManifestPinsTheDirectBrowserEmbodiment(t *testing.T) {
	root := repoChromeExtensionRoot(t)
	payload, err := os.ReadFile(filepath.Join(root, "manifest.json"))
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	var manifest manifest
	if err := json.Unmarshal(payload, &manifest); err != nil {
		t.Fatalf("decode manifest: %v", err)
	}
	if manifest.ManifestVersion != 3 {
		t.Fatalf("expected mv3 manifest: %+v", manifest)
	}
	if manifest.Background.ServiceWorker != "background.js" {
		t.Fatalf("unexpected background worker: %+v", manifest)
	}
	if len(manifest.ContentScripts) != 1 || len(manifest.ContentScripts[0].JS) != 1 || manifest.ContentScripts[0].JS[0] != "content.js" {
		t.Fatalf("unexpected content script config: %+v", manifest.ContentScripts)
	}
	if !containsString(manifest.Permissions, "nativeMessaging") {
		t.Fatalf("expected nativeMessaging permission: %+v", manifest.Permissions)
	}
	if manifest.Key == "" {
		t.Fatal("expected pinned extension key")
	}
}

func TestChromeExtensionNativeHostManifestPinsTheAllowedOrigin(t *testing.T) {
	root := repoChromeExtensionRoot(t)
	payload, err := os.ReadFile(filepath.Join(root, "native-host", "operational_browser_host.json"))
	if err != nil {
		t.Fatalf("read native host manifest: %v", err)
	}
	text := string(payload)
	if !strings.Contains(text, `"operational_browser_host"`) {
		t.Fatalf("missing native host name: %s", text)
	}
	if !strings.Contains(text, `"chrome-extension://miagfmaampfgjkojhccdilogehbjijpe/"`) {
		t.Fatalf("missing fixed allowed origin: %s", text)
	}
	if !strings.Contains(text, `"__BROWSER_HOST_PATH__"`) {
		t.Fatalf("missing host path placeholder: %s", text)
	}
}

func TestChromeExtensionBackgroundScriptForwardsRuntimeRPCAndErrors(t *testing.T) {
	runChromeExtensionScenario(t, "background-rpc")
}

func TestChromeExtensionBackgroundScriptBridgesLivePortsAndDisconnects(t *testing.T) {
	runChromeExtensionScenario(t, "background-live")
}

func TestChromeExtensionContentScriptProvesReadinessAndForwardsRPC(t *testing.T) {
	runChromeExtensionScenario(t, "content-rpc")
}

func TestChromeExtensionContentScriptBridgesLivePortsAndDisconnects(t *testing.T) {
	runChromeExtensionScenario(t, "content-live")
}

func repoChromeExtensionRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	return wd
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func runChromeExtensionScenario(t *testing.T, scenario string) {
	t.Helper()
	nodePath, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node not available")
	}
	root := repoChromeExtensionRoot(t)
	// Intent: Exercise the real shipped extension scripts under a deterministic
	// Chrome-API harness so browser contract coverage reaches the actual
	// extension/native-host boundary without depending on host registration or
	// external browser state. Source: DI-vasem
	command := exec.Command(nodePath, "-e", chromeExtensionScenarioHarness, scenario, root)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("run chrome extension scenario %s: %v\n%s", scenario, err, string(output))
	}
}

const chromeExtensionScenarioHarness = `
const fs = require("fs");
const path = require("path");
const vm = require("vm");

function assert(condition, message) {
  if (!condition) {
    throw new Error(message);
  }
}

function flush() {
  return new Promise((resolve) => setTimeout(resolve, 0));
}

function createEvent() {
  const listeners = [];
  return {
    api: {
      addListener(listener) {
        listeners.push(listener);
      },
    },
    listeners,
    emit(...args) {
      for (const listener of listeners) {
        listener(...args);
      }
    },
  };
}

function createFakePort(name) {
  const onMessage = createEvent();
  const onDisconnect = createEvent();
  return {
    name,
    posted: [],
    disconnected: false,
    onMessage: onMessage.api,
    onDisconnect: onDisconnect.api,
    postMessage(message) {
      this.posted.push(message);
    },
    disconnect() {
      this.disconnected = true;
    },
    emitMessage(message) {
      onMessage.emit(message);
    },
    emitDisconnect() {
      onDisconnect.emit();
    },
  };
}

function loadScript(scriptPath, contextValues) {
  const script = fs.readFileSync(scriptPath, "utf8");
  const context = {
    console,
    setTimeout,
    clearTimeout,
    Map,
    ...contextValues,
  };
  context.globalThis = context;
  vm.createContext(context);
  vm.runInContext(script, context, { filename: scriptPath });
  return context;
}

async function runBackgroundRPC(root) {
  const onMessage = createEvent();
  const onConnect = createEvent();
  let callCount = 0;
  let lastPayload = null;
  const chrome = {
    runtime: {
      lastError: null,
      onMessage: onMessage.api,
      onConnect: onConnect.api,
      sendNativeMessage(host, payload, callback) {
        assert(host === "operational_browser_host", "background should use the fixed host name");
        callCount += 1;
        lastPayload = payload;
        if (callCount === 1) {
          callback({
            request_id: payload.request_id,
            response: { status: 200, body: "{\"ready\":true}" },
          });
          return;
        }
        chrome.runtime.lastError = { message: "native host unavailable" };
        callback(undefined);
        chrome.runtime.lastError = null;
      },
      connectNative() {
        throw new Error("connectNative should not be used in the rpc scenario");
      },
    },
  };

  loadScript(path.join(root, "background.js"), { chrome });
  assert(onMessage.listeners.length === 1, "background should register one runtime message listener");

  let response = null;
  const asyncResult = onMessage.listeners[0]({
    kind: "rpc",
    request_id: "req-1",
    socket_path: "/tmp/oks.sock",
    request: { type: "operation", operation: "runtime_ready" },
  }, null, (payload) => {
    response = payload;
  });
  assert(asyncResult === true, "background rpc handler should stay async");
  assert(lastPayload.request.operation === "runtime_ready", "background should forward the typed runtime_ready operation");
  assert(response && response.response && response.response.status === 200, "background should return the native response");

  response = null;
  onMessage.listeners[0]({
    kind: "rpc",
    request_id: "req-2",
    socket_path: "/tmp/oks.sock",
    request: { type: "operation", operation: "inspect_item", item_id: "ITEM-1" },
  }, null, (payload) => {
    response = payload;
  });
  assert(response && response.error === "native host unavailable", "background should surface chrome.runtime.lastError");
}

async function runBackgroundLive(root) {
  const onMessage = createEvent();
  const onConnect = createEvent();
  const nativePort = createFakePort("native");
  const chrome = {
    runtime: {
      lastError: null,
      onMessage: onMessage.api,
      onConnect: onConnect.api,
      sendNativeMessage() {
        throw new Error("sendNativeMessage should not be used in the live scenario");
      },
      connectNative(host) {
        assert(host === "operational_browser_host", "background live path should use the fixed host name");
        return nativePort;
      },
    },
  };

  loadScript(path.join(root, "background.js"), { chrome });
  assert(onConnect.listeners.length === 1, "background should register one connect listener");

  const pagePort = createFakePort("oks-live:live-1");
  onConnect.listeners[0](pagePort);
  pagePort.emitMessage({
    request_id: "live-1",
    socket_path: "/tmp/oks.sock",
    request: { type: "live-open", item_id: "ITEM-1" },
  });
  assert(nativePort.posted.length === 1, "background should forward live-open to the native port");
  assert(nativePort.posted[0].request.type === "live-open", "background should preserve the live-open request");

  nativePort.emitMessage({
    request_id: "live-1",
    response: { type: "live-state", state: { version: 1 } },
  });
  assert(pagePort.posted.length === 1, "background should return native live messages to the page port");
  assert(pagePort.posted[0].response.type === "live-state", "background should preserve native live responses");

  chrome.runtime.lastError = { message: "native host disconnected" };
  nativePort.emitDisconnect();
  chrome.runtime.lastError = null;
  assert(pagePort.posted.some((message) => message.error === "native host disconnected"), "background should surface native disconnect errors");
  assert(pagePort.disconnected === true, "background should disconnect the page port after native disconnect");
}

async function runContentRPC(root) {
  let messageHandler = null;
  const posted = [];
  const sendMessageCalls = [];
  const windowObject = {
    location: { origin: "http://example.test" },
    postMessage(message) {
      posted.push(message);
    },
    addEventListener(type, listener) {
      if (type === "message") {
        messageHandler = listener;
      }
    },
  };
  let failHandshake = false;
  const chrome = {
    runtime: {
      sendMessage(payload) {
        sendMessageCalls.push(payload);
        if (failHandshake) {
          return Promise.reject(new Error("missing native host"));
        }
        if (payload.request.operation === "inspect_item") {
          return Promise.resolve({
            response: { status: 200, body: "{\"id\":\"ITEM-1\"}" },
          });
        }
        return Promise.resolve({
          response: { status: 200, body: "{\"ready\":true}" },
        });
      },
      connect() {
        throw new Error("connect should not be used in the rpc scenario");
      },
    },
  };

  loadScript(path.join(root, "content.js"), { window: windowObject, chrome });
  assert(typeof messageHandler === "function", "content script should register one message listener");

  messageHandler({
    source: windowObject,
    data: {
      __oks_bridge: true,
      direction: "page->bridge",
      kind: "handshake",
      request_id: "handshake-1",
      socket_path: "/tmp/oks.sock",
    },
  });
  await flush();
  assert(sendMessageCalls[0].request.operation === "runtime_ready", "content handshake should probe runtime_ready");
  assert(posted.some((message) => message.kind === "handshake" && message.ok === true), "content script should report successful readiness");

  failHandshake = true;
  messageHandler({
    source: windowObject,
    data: {
      __oks_bridge: true,
      direction: "page->bridge",
      kind: "handshake",
      request_id: "handshake-2",
      socket_path: "/tmp/oks.sock",
    },
  });
  await flush();
  assert(posted.some((message) => message.request_id === "handshake-2" && message.kind === "handshake" && message.ok === false), "content script should fail closed when readiness probing fails");

  failHandshake = false;
  messageHandler({
    source: windowObject,
    data: {
      __oks_bridge: true,
      direction: "page->bridge",
      kind: "rpc",
      request_id: "rpc-1",
      socket_path: "/tmp/oks.sock",
      request: {
        type: "operation",
        operation: "inspect_item",
        item_id: "ITEM-1",
      },
    },
  });
  await flush();
  assert(posted.some((message) => message.kind === "rpc-response" && message.request_id === "rpc-1" && message.response && message.response.status === 200), "content script should return one-shot rpc responses");
}

async function runContentLive(root) {
  let messageHandler = null;
  const posted = [];
  const port = createFakePort("oks-live:live-1");
  const windowObject = {
    location: { origin: "http://example.test" },
    postMessage(message) {
      posted.push(message);
    },
    addEventListener(type, listener) {
      if (type === "message") {
        messageHandler = listener;
      }
    },
  };
  const chrome = {
    runtime: {
      sendMessage() {
        throw new Error("sendMessage should not be used in the live scenario");
      },
      connect(options) {
        assert(options && options.name === "oks-live:live-1", "content script should open a named live port");
        return port;
      },
    },
  };

  loadScript(path.join(root, "content.js"), { window: windowObject, chrome });
  assert(typeof messageHandler === "function", "content script should register one message listener");

  messageHandler({
    source: windowObject,
    data: {
      __oks_bridge: true,
      direction: "page->bridge",
      kind: "live-open",
      request_id: "live-1",
      socket_path: "/tmp/oks.sock",
      request: { type: "live-open", item_id: "ITEM-1" },
    },
  });
  assert(port.posted.length === 1, "content script should forward live-open through the extension port");
  assert(port.posted[0].request.type === "live-open", "content script should preserve the live-open request");

  port.emitMessage({
    request_id: "live-1",
    response: { type: "live-state", state: { version: 1 } },
  });
  assert(posted.some((message) => message.kind === "live-message" && message.response && message.response.type === "live-state"), "content script should surface live messages back to the page");

  messageHandler({
    source: windowObject,
    data: {
      __oks_bridge: true,
      direction: "page->bridge",
      kind: "live-update",
      request_id: "live-1",
      socket_path: "/tmp/oks.sock",
      request: { type: "live-update", body: "# updated" },
    },
  });
  assert(port.posted.length === 2, "content script should forward live-update through the existing port");
  assert(port.posted[1].request.type === "live-update", "content script should preserve live updates");

  port.emitMessage({ error: "native host disconnected" });
  assert(posted.some((message) => message.kind === "error" && message.error === "native host disconnected"), "content script should surface live-port errors");
  port.emitDisconnect();

  messageHandler({
    source: windowObject,
    data: {
      __oks_bridge: true,
      direction: "page->bridge",
      kind: "live-update",
      request_id: "live-1",
      socket_path: "/tmp/oks.sock",
      request: { type: "live-update", body: "# after disconnect" },
    },
  });
  assert(posted.some((message) => message.kind === "error" && message.error === "browser live bridge is not open"), "content script should fail closed after the live port disconnects");
}

async function main() {
  const scenario = process.argv[1];
  const root = process.argv[2];
  if (!scenario || !root) {
    throw new Error("missing scenario args");
  }
  switch (scenario) {
    case "background-rpc":
      await runBackgroundRPC(root);
      return;
    case "background-live":
      await runBackgroundLive(root);
      return;
    case "content-rpc":
      await runContentRPC(root);
      return;
    case "content-live":
      await runContentLive(root);
      return;
    default:
      throw new Error("unknown scenario: " + scenario);
  }
}

main().catch((error) => {
  console.error(error && error.stack ? error.stack : String(error));
  process.exit(1);
});
`
