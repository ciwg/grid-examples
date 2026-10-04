import test from "node:test";
import assert from "node:assert/strict";

import { RelayAwarenessClient } from "./relay-awareness.js";

test("awareness websocket startup does not wait for an empty document snapshot", async () => {
  const originalWebSocket = globalThis.WebSocket;
  const originalWindow = globalThis.window;
  const sockets = [];

  class MockWebSocket {
    static OPEN = 1;

    constructor() {
      this.readyState = 0;
      this.listeners = new Map();
      this.sent = [];
      sockets.push(this);
      queueMicrotask(() => {
        this.readyState = MockWebSocket.OPEN;
        this.emit("open");
      });
    }

    addEventListener(type, listener) {
      const listeners = this.listeners.get(type) || [];
      listeners.push(listener);
      this.listeners.set(type, listeners);
    }

    emit(type) {
      for (const listener of this.listeners.get(type) || []) {
        listener({ type });
      }
    }

    send(message) {
      this.sent.push(JSON.parse(message));
    }

    close() {}
  }

  globalThis.WebSocket = MockWebSocket;
  globalThis.window = {
    location: { origin: "http://relay.test" },
    setInterval,
    clearInterval,
  };
  try {
    const client = new RelayAwarenessClient({
      basePath: "/api/local/documents/fresh",
      participantID: "browser-a",
      documentID: "fresh",
      displayName: "Browser A",
      color: "#1d6fd6",
    });

    await client.connect();

    assert.equal(sockets.length, 1);
    assert.deepEqual(sockets[0].sent[0], {
      type: "post-awareness",
      participant_id: "browser-a",
      cursor: 0,
      head: 0,
      typing: false,
      display_name: "Browser A",
      color: "#1d6fd6",
      embodiment: "browser",
    });
    client.disconnect();
  } finally {
    globalThis.WebSocket = originalWebSocket;
    globalThis.window = originalWindow;
  }
});
