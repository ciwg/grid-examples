import * as Automerge from "@automerge/automerge";

const AutomergeNext = Automerge.next;

export class AutomergeRelayClient {
  constructor(options) {
    this.basePath = options.basePath;
    this.participantID = options.participantID;
    this.documentID = options.documentID;
    this.awareness = options.awareness;
    this.doc = loadDocument(documentKey(this.documentID));
    this.syncStates = new Map();
    this.listeners = new Map();
    this.offset = 0;
    this.pollTimer = null;
    this.knownPeers = new Set();
    this.seenEnvelopes = new Set();
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

  async connect() {
    await this.poll();
    await this.flushKnownPeers();
    this.pollTimer = window.setInterval(() => {
      this.poll().catch((error) => this.emit("error", error));
    }, 250);
  }

  disconnect() {
    if (this.pollTimer !== null) {
      window.clearInterval(this.pollTimer);
      this.pollTimer = null;
    }
  }

  async poll() {
    const response = await fetch(`${this.basePath}/sync?since=${this.offset}`);
    if (!response.ok) {
      throw new Error(`sync GET failed: ${response.status}`);
    }
    const payload = await response.json();
    for (const record of payload.messages || []) {
      this.seenEnvelopes.add(record.envelope_cid);
      if (record.participant_id === this.participantID) {
        continue;
      }
      if (record.recipient_id && record.recipient_id !== this.participantID) {
        continue;
      }
      await this.receive(record);
    }
    this.offset = payload.next_offset || this.offset;
  }

  observePeers(peers) {
    let changed = false;
    for (const peer of peers) {
      if (!peer.participant_id || peer.participant_id === this.participantID) {
        continue;
      }
      if (!this.knownPeers.has(peer.participant_id)) {
        this.knownPeers.add(peer.participant_id);
        changed = true;
      }
    }
    if (changed) {
      this.flushKnownPeers().catch((error) => this.emit("error", error));
    }
  }

  applyLocalUpdate(update) {
    let nextDoc = this.doc;
    update.changes.iterChanges((fromA, toA, _fromB, _toB, inserted) => {
      const deleteCount = toA - fromA;
      const insertText = inserted.toString();
      // Intent: Preserve the old working Automerge splice semantics so local
      // edits stay character-precise and concurrent merges happen in the CRDT,
      // not in ad hoc browser snapshot code. Source: DI-ramuv; DI-zegov
      nextDoc = Automerge.change(nextDoc, (draft) => {
        AutomergeNext.splice(draft, ["content"], fromA, deleteCount, insertText);
      });
    });
    this.setDoc(nextDoc);
    this.flushKnownPeers().catch((error) => this.emit("error", error));
  }

  async receive(record) {
    const syncState = this.syncStates.get(record.participant_id) || Automerge.initSyncState();
    const message = base64ToBytes(record.message_base64);
    const previous = this.doc;
    const [nextDoc, nextState] = Automerge.receiveSyncMessage(previous, syncState, message);
    this.syncStates.set(record.participant_id, nextState);
    if (!Automerge.equals(previous, nextDoc)) {
      this.setDoc(nextDoc);
      this.emit("document", this.getText());
    }
    await this.flushPeer(record.participant_id);
  }

  async flushKnownPeers() {
    for (const peerID of this.knownPeers) {
      await this.flushPeer(peerID);
    }
  }

  async flushPeer(peerID) {
    if (!peerID || peerID === this.participantID) {
      return;
    }
    let state = this.syncStates.get(peerID) || Automerge.initSyncState();
    for (let attempt = 0; attempt < 8; attempt += 1) {
      let message;
      [state, message] = Automerge.generateSyncMessage(this.doc, state);
      this.syncStates.set(peerID, state);
      if (!message) {
        break;
      }
      await fetch(`${this.basePath}/sync`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          participant_id: this.participantID,
          recipient_id: peerID,
          message_base64: bytesToBase64(message),
          embodiment: "browser",
        }),
      });
    }
  }

  setDoc(nextDoc) {
    this.doc = ensureDocument(nextDoc);
    localStorage.setItem(documentKey(this.documentID), bytesToBase64(Automerge.save(this.doc)));
  }
}

function ensureDocument(doc) {
  if (doc?.content !== undefined) {
    return doc;
  }
  return Automerge.from({ content: new Automerge.Text() });
}

function loadDocument(key) {
  const encoded = localStorage.getItem(key);
  if (!encoded) {
    return ensureDocument(null);
  }
  try {
    return ensureDocument(Automerge.load(base64ToBytes(encoded)));
  } catch (_error) {
    return ensureDocument(null);
  }
}

function documentKey(documentID) {
  return `grid-editor-automerge-${documentID}`;
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
