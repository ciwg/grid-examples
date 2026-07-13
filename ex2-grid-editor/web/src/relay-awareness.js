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

  emit(event) {
    this.listeners.get(event)?.forEach((callback) => callback());
  }

  getStates() {
    const states = new Map(this.remoteStates);
    states.set(this.participantID, {
      user: { ...this.user },
      selection: { ...this.selection },
      typing: this.typing,
    });
    return states;
  }

  async connect() {
    await this.broadcast();
    await this.poll();
    this.pollTimer = window.setInterval(() => {
      this.poll().catch(() => {});
    }, 350);
  }

  disconnect() {
    if (this.pollTimer !== null) {
      window.clearInterval(this.pollTimer);
      this.pollTimer = null;
    }
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
      });
    }
    this.emit("change");
  }

  async broadcast() {
    // Intent: Keep awareness rooted in the old collab-awareness state shape
    // while swapping the transport to the new relay HTTP surface. Source:
    // DI-zegov
    await fetch(`${this.basePath}/awareness`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
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
  }

  updateSelection(anchor, head) {
    this.selection = { anchor, head };
    this.broadcast().catch(() => {});
    this.emit("change");
  }

  updateCursor(anchor) {
    this.updateSelection(anchor, anchor);
  }

  setTyping(typing) {
    this.typing = typing;
    this.broadcast().catch(() => {});
    this.emit("change");
  }

  setName(name) {
    this.user.name = name;
    this.broadcast().catch(() => {});
    this.emit("change");
  }

  setColor(color) {
    this.user.color = color;
    this.broadcast().catch(() => {});
    this.emit("change");
  }
}
