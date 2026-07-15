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

export async function bootstrapRemoteSession(basePath, participantID, accessToken) {
  if (!accessToken) {
    return null;
  }
  // Intent: Keep the long-lived share token out of ex3's steady-state live
  // mutation traffic by exchanging it once for short-lived document-scoped
  // capabilities before HTTP or websocket mutation begins. Source: DI-povip
  const response = await fetch(`${basePath}/session`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "X-Grid-Access-Token": accessToken,
    },
    body: JSON.stringify({
      participant_id: participantID,
    }),
  });
  if (!response.ok) {
    const detail = await response.text();
    throw new Error(`session bootstrap failed: ${response.status}${detail ? ` ${detail}` : ""}`);
  }
  return await response.json();
}

export function bearerHeaders(capability) {
  if (!capability) {
    return {};
  }
  return {
    Authorization: `Bearer ${capability}`,
  };
}
