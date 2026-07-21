import * as Automerge from "@automerge/automerge";
import { bearerHeaders, canUseWebSocket, toWebSocketURL } from "./live-transport.js";

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
    this.capabilities = options.capabilities || {};
    this.initialSyncReady = false;
    this.queuedUpdate = null;
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

  primeFromRelayState(state) {
    if (!state?.snapshot_present || !state.replica_base64) {
      return false;
    }
    // Intent: Bootstrap the browser replica from the relay's current snapshot
    // before live websocket replay starts so an old local seed or empty startup
    // buffer cannot masquerade as the shared document. Source: DI-ramuv;
    // DI-lumek; DI-gafit
    this.setDoc(Automerge.load(base64ToBytes(state.replica_base64)));
    if (state.snapshot_offset && state.snapshot_offset > this.offset) {
      this.offset = state.snapshot_offset;
    }
    return true;
  }

  async connect() {
    this.emitStatus("connecting");
    this.initialSyncReady = false;
    this.queuedUpdate = null;
    if (canUseWebSocket()) {
      await this.connectWebSocket();
      return;
    }
    this.transport = "polling";
    await this.poll();
    this.initialSyncReady = true;
    this.flushQueuedUpdate();
    this.emitStatus(this.pendingChanges > 0 ? "syncing" : "ready");
    this.pollTimer = window.setInterval(() => {
      this.poll().catch((error) => {
        this.connected = false;
        this.emitStatus("disconnected", error.message);
        this.emit("error", error);
      });
    }, 250);
  }

  // Intent: Give startup one explicit HTTP replay escape hatch when a browser
  // still opens with blank text even though the relay already reports shared
  // history, so private/incognito sessions are not forced to rely on perfect
  // websocket catch-up on first load. Source: DI-sulor
  async recoverFromRelayHistory(state) {
    const startOffset = state?.snapshot_present && state?.snapshot_offset ? state.snapshot_offset : 0;
    await this.fetchSyncFeed(startOffset);
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
    this.initialSyncReady = false;
    this.queuedUpdate = null;
    this.emitStatus("disconnected");
  }

  async connectWebSocket() {
    this.transport = "websocket";
    const url = toWebSocketURL(this.basePath, "sync-socket", { since: this.offset });
    await new Promise((resolve, reject) => {
      let settled = false;
      let socketWork = Promise.resolve();
      const socket = new WebSocket(url);
      this.socket = socket;
      socket.addEventListener("open", () => {
        this.connected = true;
        if (this.capabilities.sync) {
          socket.send(JSON.stringify({
            type: "auth",
            capability: this.capabilities.sync,
          }));
        }
      });
      socket.addEventListener("message", (event) => {
        // Intent: Process websocket sync frames in arrival order so browser
        // startup cannot mark the replica ready before earlier replay batches
        // are applied. Source: DI-gafit
        socketWork = socketWork.then(() => this.handleSocketMessage(event.data));
        socketWork
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
    await this.fetchSyncFeed(this.offset);
    this.emitStatus(this.pendingChanges > 0 ? "syncing" : "ready");
  }

  async fetchSyncFeed(startOffset) {
    // Intent: Drain bounded relay pages until caught up so late joiners do not
    // silently skip older changes when the server applies feed-size limits.
    // Source: DI-rabod
    let cursor = startOffset;
    while (true) {
      const response = await fetch(`${this.basePath}/sync?since=${cursor}&limit=256`);
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
      const nextOffset = payload.next_offset || cursor;
      if (nextOffset <= cursor) {
        break;
      }
      cursor = nextOffset;
      if ((payload.messages || []).length < 256) {
        break;
      }
    }
    if (cursor > this.offset) {
      this.offset = cursor;
    }
  }

  observePeers(_peers) {}

  applyLocalUpdate(update) {
    if (!this.initialSyncReady) {
      // Intent: Keep browser-side local edits from overwriting relay content
      // with an empty startup buffer before the initial sync catch-up has
      // completed. Source: DI-gafit
      this.queuedUpdate = update;
      return;
    }
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

  flushQueuedUpdate() {
    if (!this.initialSyncReady || this.queuedUpdate === null) {
      return;
    }
    const update = this.queuedUpdate;
    this.queuedUpdate = null;
    this.applyLocalUpdate(update);
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

  async publishSnapshot() {
    const response = await fetch(`${this.basePath}/snapshot`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        ...bearerHeaders(this.capabilities.sync),
      },
      body: JSON.stringify({
        participant_id: this.participantID,
        text_base64: bytesToBase64(new TextEncoder().encode(this.getText())),
        replica_base64: bytesToBase64(this.getReplicaBytes()),
      }),
    });
    if (!response.ok) {
      throw new Error(`snapshot POST failed: ${response.status}`);
    }
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
    const replicaBase64 = bytesToBase64(this.getReplicaBytes());
    const textBase64 = bytesToBase64(new TextEncoder().encode(this.getText()));
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify({
        type: "post-sync",
        participant_id: this.participantID,
        recipient_id: "",
        message_base64: bytesToBase64(change),
        text_base64: textBase64,
        replica_base64: replicaBase64,
        embodiment: "browser",
      }));
      return;
    }
    const response = await fetch(`${this.basePath}/sync`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        ...bearerHeaders(this.capabilities.sync),
      },
      body: JSON.stringify({
        participant_id: this.participantID,
        recipient_id: "",
        message_base64: bytesToBase64(change),
        text_base64: textBase64,
        replica_base64: replicaBase64,
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
      this.initialSyncReady = true;
      this.flushQueuedUpdate();
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
