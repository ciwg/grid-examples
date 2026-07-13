import * as Automerge from "@automerge/automerge";

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
    this.offset = payload.next_offset || this.offset;
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
      this.postChange(change).catch((error) => this.emit("error", error));
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
    await fetch(`${this.basePath}/sync`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        participant_id: this.participantID,
        recipient_id: "",
        message_base64: bytesToBase64(change),
        embodiment: "browser",
      }),
    });
  }

  setDoc(nextDoc) {
    this.doc = ensureDocument(nextDoc);
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
