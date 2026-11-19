export function canUseWebSocket() {
  return typeof globalThis.WebSocket === "function";
}

export function canUseWebRTC() {
  return typeof globalThis.RTCPeerConnection === "function";
}

export function preferredBrowserTransport() {
  const selected = globalThis.window?.__GRID_BROWSER_TRANSPORT__;
  return selected === "websocket" || selected === "polling" || selected === "webrtc" ? selected : "webrtc";
}

// Intent: Use the browser's WebRTC DataChannel API for live local-node
// carriage because browsers do not expose raw DTLS/UDP sockets. The existing
// signed application envelopes remain unchanged; WebRTC supplies the
// DTLS-backed connection beneath this embodiment-local channel. Source:
// DI-higom; DI-pirov.
export async function openWebRTCChannel(basePath, channelName, params = {}, capability = "") {
  const connection = new RTCPeerConnection({ iceServers: [] });
  const channel = connection.createDataChannel(channelName, { ordered: true });
  const ready = new Promise((resolve, reject) => {
    channel.addEventListener("open", () => resolve());
    channel.addEventListener("error", () => reject(new Error(`${channelName} webrtc channel failed`)));
  });
  const offer = await connection.createOffer();
  await connection.setLocalDescription(offer);
  await waitForICEGathering(connection);
  const origin = globalThis.window?.location?.origin || "http://127.0.0.1";
  const url = new URL(`${basePath.replace(/\/$/, "")}/webrtc-offer`, origin);
  for (const [key, value] of Object.entries(params)) {
    if (value !== undefined && value !== null && value !== "") {
      url.searchParams.set(key, String(value));
    }
  }
  const response = await fetch(url.toString(), {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      ...bearerHeaders(capability),
    },
    body: JSON.stringify({
      channel: channelName,
      offer: connection.localDescription,
    }),
  });
  if (!response.ok) {
    connection.close();
    throw new Error(`${channelName} webrtc offer failed: ${response.status}`);
  }
  const payload = await response.json();
  await connection.setRemoteDescription(payload.answer);
  return { channel, connection, ready };
}

export function isLiveConnectionOpen(connection) {
  return connection?.readyState === globalThis.WebSocket?.OPEN || connection?.readyState === "open";
}

function waitForICEGathering(connection) {
  if (connection.iceGatheringState === "complete") {
    return Promise.resolve();
  }
  return new Promise((resolve) => {
    connection.addEventListener("icegatheringstatechange", () => {
      if (connection.iceGatheringState === "complete") {
        resolve();
      }
    });
  });
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
  const url = new URL(`${basePath}/session`, globalThis.window?.location?.origin || "http://127.0.0.1");
  url.searchParams.set("access_token", accessToken);
  // Intent: Keep the long-lived share token out of ex3's steady-state live
  // mutation traffic by exchanging it once for short-lived document-scoped
  // capabilities before HTTP or websocket mutation begins. Source: DI-povip
  const response = await fetch(url.toString(), {
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
