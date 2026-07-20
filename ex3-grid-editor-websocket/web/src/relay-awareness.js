import { bearerHeaders, canUseWebSocket, toWebSocketURL } from "./live-transport.js";

export class RelayAwarenessClient {
  constructor(options) {
    this.basePath = options.basePath;
    this.participantID = options.participantID;
    this.documentID = options.documentID;
    this.user = {
      name: options.displayName,
      color: options.color,
    };
    this.selection = { anchor: 0, head: 0 };
    this.typing = false;
    this.remoteStates = new Map();
    this.listeners = new Map();
    this.pollTimer = null;
    this.heartbeatTimer = null;
    this.socket = null;
    this.transport = "polling";
    this.capabilities = options.capabilities || {};
  }

  on(event, callback) {
    if (!this.listeners.has(event)) {
      this.listeners.set(event, new Set());
    }
    this.listeners.get(event).add(callback);
  }

  off(event, callback) {
    this.listeners.get(event)?.delete(callback);
  }

  emit(event, value) {
    this.listeners.get(event)?.forEach((callback) => callback(value));
  }

  getStates() {
    const states = new Map(this.remoteStates);
    states.set(this.participantID, {
      user: { ...this.user },
      selection: { ...this.selection },
      typing: this.typing,
      lastSeenAt: new Date().toISOString(),
    });
    return states;
  }

  transportMode() {
    return this.transport;
  }

  async connect() {
    if (canUseWebSocket()) {
      await this.connectWebSocket();
      return;
    }
    this.transport = "polling";
    await this.broadcast();
    await this.poll();
    this.startHeartbeat();
    this.pollTimer = window.setInterval(() => {
      this.poll().catch((error) => this.emit("error", error));
    }, 350);
  }

  disconnect() {
    if (this.pollTimer !== null) {
      window.clearInterval(this.pollTimer);
      this.pollTimer = null;
    }
    if (this.socket) {
      this.socket.close();
      this.socket = null;
    }
    if (this.heartbeatTimer !== null) {
      window.clearInterval(this.heartbeatTimer);
      this.heartbeatTimer = null;
    }
  }

  async connectWebSocket() {
    this.transport = "websocket";
    const url = toWebSocketURL(this.basePath, "awareness-socket");
    await new Promise((resolve, reject) => {
      let settled = false;
      const socket = new WebSocket(url);
      this.socket = socket;
      socket.addEventListener("open", () => {
        if (this.capabilities.awareness) {
          socket.send(JSON.stringify({
            type: "auth",
            capability: this.capabilities.awareness,
          }));
        }
        Promise.resolve(this.broadcast()).catch((error) => {
          if (!settled) {
            settled = true;
            reject(error);
          }
        });
      });
      socket.addEventListener("message", (event) => {
        try {
          this.handleSocketMessage(event.data);
          if (!settled) {
            settled = true;
            resolve();
          }
        } catch (error) {
          if (!settled) {
            settled = true;
            reject(error);
            return;
          }
          this.emit("error", error);
        }
      });
      socket.addEventListener("error", () => {
        const error = new Error("awareness websocket failed");
        if (!settled) {
          settled = true;
          reject(error);
          return;
        }
        this.emit("error", error);
      });
      socket.addEventListener("close", () => {
        this.socket = null;
      });
    });
    this.startHeartbeat();
  }

  async poll() {
    const response = await fetch(`${this.basePath}/awareness`);
    if (!response.ok) {
      throw new Error(`awareness GET failed: ${response.status}`);
    }
    const payload = await response.json();
    this.remoteStates.clear();
    for (const peer of payload.awareness || []) {
      this.remoteStates.set(peer.participant_id, {
        user: {
          name: peer.display_name || peer.author,
          color: peer.color || "#999999",
        },
        selection: {
          anchor: peer.cursor || 0,
          head: peer.head || peer.cursor || 0,
        },
        typing: Boolean(peer.typing),
        lastSeenAt: peer.last_seen_at || null,
        embodiment: peer.embodiment || "",
      });
    }
    this.emit("change");
  }

  async broadcast() {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify({
        type: "post-awareness",
        participant_id: this.participantID,
        cursor: this.selection.anchor,
        head: this.selection.head,
        typing: this.typing,
        display_name: this.user.name,
        color: this.user.color,
        embodiment: "browser",
      }));
      return;
    }
    if (this.transport === "websocket") {
      return;
    }
    // Intent: Keep awareness rooted in the old collab-awareness state shape
    // while swapping the transport to the new relay HTTP surface. Source:
    // DI-zegov
    const response = await fetch(`${this.basePath}/awareness`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        ...bearerHeaders(this.capabilities.awareness),
      },
      body: JSON.stringify({
        participant_id: this.participantID,
        cursor: this.selection.anchor,
        head: this.selection.head,
        typing: this.typing,
        display_name: this.user.name,
        color: this.user.color,
        embodiment: "browser",
      }),
    });
    if (!response.ok) {
      // Intent: Make rejected awareness writes visible to the browser instead
      // of leaving users with stale peer labels or cursor state. Source:
      // DI-rabod
      throw new Error(`awareness POST failed: ${response.status}`);
    }
  }

  startHeartbeat() {
    if (this.heartbeatTimer !== null) {
      window.clearInterval(this.heartbeatTimer);
    }
    // Intent: Refresh active browser presence periodically so current
    // collaborators remain visible while abandoned sessions age out quickly.
    // Source: DI-gafit
    this.heartbeatTimer = window.setInterval(() => {
      this.broadcast().catch((error) => this.emit("error", error));
    }, 5000);
  }

  updateSelection(anchor, head) {
    this.selection = { anchor, head };
    this.broadcast().catch((error) => this.emit("error", error));
    this.emit("change");
  }

  updateCursor(anchor) {
    this.updateSelection(anchor, anchor);
  }

  setTyping(typing) {
    this.typing = typing;
    this.broadcast().catch((error) => this.emit("error", error));
    this.emit("change");
  }

  setName(name) {
    this.user.name = name;
    this.broadcast().catch((error) => this.emit("error", error));
    this.emit("change");
  }

  setColor(color) {
    this.user.color = color;
    this.broadcast().catch((error) => this.emit("error", error));
    this.emit("change");
  }

  handleSocketMessage(raw) {
    const payload = JSON.parse(raw);
    if (payload.type === "awareness-state") {
      this.remoteStates.clear();
      for (const peer of payload.awareness || []) {
        this.remoteStates.set(peer.participant_id, {
          user: {
            name: peer.display_name || peer.author,
            color: peer.color || "#999999",
          },
          selection: {
            anchor: peer.cursor || 0,
            head: peer.head || peer.cursor || 0,
          },
          typing: Boolean(peer.typing),
          lastSeenAt: peer.last_seen_at || null,
          embodiment: peer.embodiment || "",
        });
      }
      this.emit("change");
      return;
    }
    if (payload.type === "error") {
      throw new Error(payload.message || "awareness websocket error");
    }
  }
}
