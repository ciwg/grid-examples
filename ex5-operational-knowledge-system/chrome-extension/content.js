const pagePorts = new Map();

// Intent: Keep the browser UI in the page while bridging only carriage through
// the extension boundary, so `web/app.js` can move onto the direct embodiment
// contract without being rewritten as an extension-owned UI. Source: DI-punek
function postToPage(message) {
  window.postMessage({
    __oks_bridge: true,
    direction: "bridge->page",
    ...message,
  }, window.location.origin);
}

// Intent: Translate page-level bridge messages into extension runtime traffic
// without inventing new browser-specific semantics beyond the locked direct
// contract family. Source: DI-punek
window.addEventListener("message", (event) => {
  if (event.source !== window || !event.data || event.data.__oks_bridge !== true || event.data.direction !== "page->bridge") {
    return;
  }
  const message = event.data;
  if (message.kind === "handshake") {
    postToPage({
      kind: "handshake",
      request_id: message.request_id,
      ok: true,
    });
    return;
  }
  if (message.kind === "rpc") {
    chrome.runtime.sendMessage({
      kind: "rpc",
      request_id: message.request_id,
      socket_path: message.socket_path,
      request: message.request,
    }).then((response) => {
      if (response && response.error) {
        postToPage({
          kind: "error",
          request_id: message.request_id,
          error: response.error,
        });
        return;
      }
      postToPage({
        kind: "rpc-response",
        request_id: message.request_id,
        response: response ? response.response : null,
      });
    }).catch((error) => {
      postToPage({
        kind: "error",
        request_id: message.request_id,
        error: String(error),
      });
    });
    return;
  }
  if (message.kind === "live-open") {
    const port = chrome.runtime.connect({ name: `oks-live:${message.request_id}` });
    pagePorts.set(message.request_id, port);
    port.onMessage.addListener((response) => {
      if (response && response.error) {
        postToPage({
          kind: "error",
          request_id: message.request_id,
          error: response.error,
        });
        return;
      }
      postToPage({
        kind: "live-message",
        request_id: message.request_id,
        response: response ? response.response : null,
      });
    });
    port.onDisconnect.addListener(() => {
      pagePorts.delete(message.request_id);
    });
    port.postMessage({
      request_id: message.request_id,
      socket_path: message.socket_path,
      request: message.request,
    });
    return;
  }
  if (message.kind === "live-update" || message.kind === "live-close") {
    const port = pagePorts.get(message.request_id);
    if (!port) {
      postToPage({
        kind: "error",
        request_id: message.request_id,
        error: "browser live bridge is not open",
      });
      return;
    }
    port.postMessage({
      request_id: message.request_id,
      socket_path: message.socket_path,
      request: message.request,
    });
    if (message.kind === "live-close") {
      try {
        port.disconnect();
      } catch {
      }
      pagePorts.delete(message.request_id);
    }
  }
});
