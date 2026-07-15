import * as Automerge from "@automerge/automerge";
import { canUseWebSocket, toWebSocketURL } from "./live-transport.js";

const AutomergeNext = Automerge.next;
const EMPTY_DOCUMENT_BYTES = base64ToBytes("hW9Kg8HDZmEAdQEQUDnUuZsuTLOKK6EtAqSUwAF91ThR16b5XY1P61eTHXkwnJNTicqZ35V+jMImBQWmigYBAgMCEwIjBkACVgIHFQkhAiMCNAFCAlYCgAECfwB/AX8Bf8660dIGfwB/B38HY29udGVudH8AfwEBfwR/AH8AAA==");

export class AutomergeRelayClient {
  constructor(options) {
    this.basePath = options.basePath;
    this.participantID = options.participantID;
    this.documentID = options.documentID;
    this.awareness = options.awareness;
    this.doc = ensureDocument(null);
    this.listeners = new Map();
    this.offset = 0;
    this.pollTimer = null;
    this.socket = null;
    this.seenEnvelopes = new Set();
    this.connected = false;
    this.pendingChanges = 0;
    this.transport = "polling";
  }

  on(event, callback) {
    if (!this.listeners.has(event)) {
      this.listeners.set(event, new Set());
    }
    this.listeners.get(event).add(callback);
  }

  emit(event, value) {
    this.listeners.get(event)?.forEach((callback) => callback(value));
  }

  getText() {
    return this.doc.content?.toString() || "";
  }

  getReplicaCID() {
    return bytesToBase64(Automerge.save(this.doc)).slice(0, 16);
  }

  getReplicaBytes() {
    return Automerge.save(this.doc);
  }

  transportMode() {
    return this.transport;
  }

  async connect() {
    this.emitStatus("connecting");
    if (canUseWebSocket()) {
      await this.connectWebSocket();
      return;
    }
    this.transport = "polling";
    await this.poll();
    this.emitStatus(this.pendingChanges > 0 ? "syncing" : "ready");
    this.pollTimer = window.setInterval(() => {
      this.poll().catch((error) => {
        this.connected = false;
        this.emitStatus("disconnected", error.message);
        this.emit("error", error);
      });
    }, 250);
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
    this.connected = false;
    this.emitStatus("disconnected");
  }

  async connectWebSocket() {
    this.transport = "websocket";
    const url = toWebSocketURL(this.basePath, "sync-socket", { since: this.offset });
    await new Promise((resolve, reject) => {
      let settled = false;
      const socket = new WebSocket(url);
      this.socket = socket;
      socket.addEventListener("open", () => {
        this.connected = true;
      });
      socket.addEventListener("message", (event) => {
        this.handleSocketMessage(event.data)
          .then(() => {
            if (!settled) {
              settled = true;
              resolve();
            }
          })
          .catch((error) => {
            if (!settled) {
              settled = true;
              reject(error);
              return;
            }
            this.emitStatus("error", error.message);
            this.emit("error", error);
          });
      });
      socket.addEventListener("error", () => {
        const error = new Error("sync websocket failed");
        if (!settled) {
          settled = true;
          reject(error);
          return;
        }
        this.connected = false;
        this.emitStatus("error", error.message);
        this.emit("error", error);
      });
      socket.addEventListener("close", () => {
        this.connected = false;
        this.socket = null;
        this.emitStatus("disconnected", "sync websocket closed");
      });
    });
    this.emitStatus(this.pendingChanges > 0 ? "syncing" : "ready");
  }

  async poll() {
    // Intent: Drain bounded relay pages until caught up so late joiners do not
    // silently skip older changes when the server applies feed-size limits.
    // Source: DI-rabod
    while (true) {
      const response = await fetch(`${this.basePath}/sync?since=${this.offset}&limit=256`);
      if (!response.ok) {
        throw new Error(`sync GET failed: ${response.status}`);
      }
      this.connected = true;
      const payload = await response.json();
      for (const record of payload.messages || []) {
        this.seenEnvelopes.add(record.envelope_cid);
        // Intent: Keep the relay feed peer-oriented so the browser replica does
        // not re-apply its own signed sync messages after they round-trip through
        // the relay. Source: DI-ramuv; DI-zegov
        if (record.participant_id === this.participantID) {
          continue;
        }
        if (record.recipient_id && record.recipient_id !== this.participantID) {
          continue;
        }
        await this.receive(record);
      }
      const nextOffset = payload.next_offset || this.offset;
      if (nextOffset <= this.offset) {
        break;
      }
      this.offset = nextOffset;
      if ((payload.messages || []).length < 256) {
        break;
      }
    }
    this.emitStatus(this.pendingChanges > 0 ? "syncing" : "ready");
  }

  observePeers(_peers) {}

  applyLocalUpdate(update) {
    const nextDoc = Automerge.change(Automerge.clone(this.doc), (draft) => {
      update.changes.iterChanges((fromA, toA, _fromB, _toB, inserted) => {
        const deleteCount = toA - fromA;
        const insertText = inserted.toString();
        AutomergeNext.splice(draft, ["content"], fromA, deleteCount, insertText);
      });
    });
    // Intent: Publish durable Automerge change packets so relay history is
    // replayable for late joiners instead of depending on per-peer sync
    // sessions. Source: DI-zegov; DI-larok
    this.setDoc(nextDoc);
    const change = Automerge.getLastLocalChange(nextDoc);
    if (change) {
      this.pendingChanges += 1;
      this.emitStatus("unsynced");
      this.postChange(change).catch((error) => {
        this.pendingChanges = Math.max(0, this.pendingChanges - 1);
        this.emitStatus("error", error.message);
        this.emit("error", error);
      });
    }
  }

  replaceText(text) {
    const previous = this.getText();
    if (text === previous) {
      return;
    }
    const prefix = commonPrefix(previous, text);
    const suffix = commonSuffix(previous, text, prefix);
    const deleteCount = previous.length - prefix - suffix;
    const insertText = text.slice(prefix, text.length - suffix);
    this.applyLocalUpdate({
      changes: {
        iterChanges(callback) {
          callback(prefix, prefix + deleteCount, 0, 0, {
            toString() {
              return insertText;
            },
          });
        },
      },
    });
  }

  async receive(record) {
    const previous = Automerge.clone(this.doc);
    const [nextDoc] = Automerge.applyChanges(previous, [base64ToBytes(record.message_base64)]);
    if (!Automerge.equals(previous, nextDoc)) {
      this.setDoc(nextDoc);
      this.emit("document", this.getText());
    }
  }

  async postChange(change) {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify({
        type: "post-sync",
        participant_id: this.participantID,
        recipient_id: "",
        message_base64: bytesToBase64(change),
        embodiment: "browser",
      }));
      return;
    }
    const response = await fetch(`${this.basePath}/sync`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        participant_id: this.participantID,
        recipient_id: "",
        message_base64: bytesToBase64(change),
        embodiment: "browser",
      }),
    });
    if (!response.ok) {
      // Intent: Surface rejected relay writes to the browser UI instead of
      // failing silently after a local editor update. Source: DI-rabod
      throw new Error(`sync POST failed: ${response.status}`);
    }
    this.pendingChanges = Math.max(0, this.pendingChanges - 1);
    this.emitStatus("ready");
  }

  setDoc(nextDoc) {
    this.doc = ensureDocument(nextDoc);
  }

  async handleSocketMessage(raw) {
    const message = JSON.parse(raw);
    if (message.type === "sync-feed") {
      for (const record of message.messages || []) {
        this.seenEnvelopes.add(record.envelope_cid);
        if (record.participant_id === this.participantID) {
          continue;
        }
        if (record.recipient_id && record.recipient_id !== this.participantID) {
          continue;
        }
        await this.receive(record);
      }
      const nextOffset = message.next_offset || this.offset;
      if (nextOffset > this.offset) {
        this.offset = nextOffset;
      }
      this.emitStatus(this.pendingChanges > 0 ? "syncing" : "ready");
      return;
    }
    if (message.type === "sync-posted") {
      this.pendingChanges = Math.max(0, this.pendingChanges - 1);
      this.emitStatus(this.pendingChanges > 0 ? "syncing" : "ready");
      return;
    }
    if (message.type === "sync-ready") {
      const nextOffset = message.next_offset || this.offset;
      if (nextOffset > this.offset) {
        this.offset = nextOffset;
      }
      this.emitStatus(this.pendingChanges > 0 ? "syncing" : "ready");
      return;
    }
    if (message.type === "error") {
      throw new Error(message.message || "sync websocket error");
    }
  }

  emitStatus(phase, message = "") {
    this.emit("status", {
      connected: this.connected,
      pendingChanges: this.pendingChanges,
      phase,
      message,
    });
  }
}

function ensureDocument(doc) {
  if (doc?.content !== undefined) {
    return doc;
  }
  // Intent: Give every replica the same canonical serialized Automerge base,
  // including the same initial actor/object identities, so relayed change
  // packets replay identically across browser and Neovim embodiments.
  // Source: DI-zegov; DI-larok
  return Automerge.load(EMPTY_DOCUMENT_BYTES);
}

function bytesToBase64(bytes) {
  let text = "";
  bytes.forEach((value) => {
    text += String.fromCharCode(value);
  });
  return window.btoa(text);
}

function base64ToBytes(value) {
  return Uint8Array.from(window.atob(value), (char) => char.charCodeAt(0));
}

function commonPrefix(left, right) {
  const size = Math.min(left.length, right.length);
  let index = 0;
  while (index < size && left[index] === right[index]) {
    index += 1;
  }
  return index;
}

function commonSuffix(left, right, prefix) {
  const leftRemain = left.length - prefix;
  const rightRemain = right.length - prefix;
  const size = Math.min(leftRemain, rightRemain);
  let index = 0;
  while (index < size && left[left.length - 1 - index] === right[right.length - 1 - index]) {
    index += 1;
  }
  return index;
}
