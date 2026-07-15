export function canUseWebSocket() {
  return typeof globalThis.WebSocket === "function";
}

export function toWebSocketURL(basePath, action, params = {}) {
  const origin = globalThis.window?.location?.origin || "http://127.0.0.1";
  const url = new URL(`${basePath.replace(/\/$/, "")}/${action}`, origin);
  url.protocol = url.protocol === "https:" ? "wss:" : "ws:";
  for (const [key, value] of Object.entries(params)) {
    if (value === undefined || value === null || value === "") {
      continue;
    }
    url.searchParams.set(key, String(value));
  }
  return url.toString();
}
