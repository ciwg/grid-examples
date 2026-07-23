const HOST_NAME = "operational_browser_host";

// Intent: Keep the extension worker as a thin carriage layer that forwards
// direct browser embodiment traffic into the native host instead of inventing a
// browser-only semantic runtime. Source: DI-punek
chrome.runtime.onMessage.addListener((message, _sender, sendResponse) => {
  if (!message || message.kind !== "rpc") {
    return false;
  }
  try {
    chrome.runtime.sendNativeMessage(HOST_NAME, {
      request_id: message.request_id,
      socket_path: message.socket_path,
      request: message.request,
    }, (response) => {
      if (chrome.runtime.lastError) {
        sendResponse({
          request_id: message.request_id,
          error: chrome.runtime.lastError.message,
        });
        return;
      }
      sendResponse(response || {
        request_id: message.request_id,
        error: "native host returned no response",
      });
    });
  } catch (error) {
    sendResponse({
      request_id: message.request_id,
      error: String(error),
    });
  }
  return true;
});

// Intent: Carry browser live-draft traffic over one long-lived native host
// connection so the browser stays on the same direct contract family as the
// terminal embodiments. Source: DI-punek
chrome.runtime.onConnect.addListener((port) => {
  if (!port.name.startsWith("oks-live:")) {
    return;
  }
  const nativePort = chrome.runtime.connectNative(HOST_NAME);
  const disconnect = () => {
    try {
      nativePort.disconnect();
    } catch {
    }
    try {
      port.disconnect();
    } catch {
    }
  };
  nativePort.onMessage.addListener((message) => {
    port.postMessage(message);
  });
  nativePort.onDisconnect.addListener(() => {
    port.postMessage({
      request_id: port.name.slice("oks-live:".length),
      error: chrome.runtime.lastError ? chrome.runtime.lastError.message : "native host disconnected",
    });
    disconnect();
  });
  port.onMessage.addListener((message) => {
    nativePort.postMessage({
      request_id: message.request_id,
      socket_path: message.socket_path,
      request: message.request,
    });
  });
  port.onDisconnect.addListener(() => {
    disconnect();
  });
});
