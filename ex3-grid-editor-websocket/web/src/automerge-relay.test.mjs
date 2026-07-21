import test from "node:test";
import assert from "node:assert/strict";
import * as Automerge from "@automerge/automerge";

if (!globalThis.window) {
  globalThis.window = globalThis;
}
if (!globalThis.window.btoa) {
  globalThis.window.btoa = (value) => Buffer.from(value, "binary").toString("base64");
}
if (!globalThis.window.atob) {
  globalThis.window.atob = (value) => Buffer.from(value, "base64").toString("binary");
}

test("primeFromRelayState hydrates browser replica from relay snapshot", async () => {
  const { AutomergeRelayClient } = await import("./automerge-relay.js");
  const source = new AutomergeRelayClient({
    basePath: "/api/local/documents/demo",
    participantID: "browser-source",
    documentID: "demo",
    awareness: { on() {} },
    capabilities: {},
  });
  source.postChange = async () => {};
  source.initialSyncReady = true;
  source.replaceText("# Live Demo Script\n\nShared manual");
  await Promise.resolve();

  const target = new AutomergeRelayClient({
    basePath: "/api/local/documents/demo",
    participantID: "browser-target",
    documentID: "demo",
    awareness: { on() {} },
    capabilities: {},
  });
  const replicaBase64 = Buffer.from(source.getReplicaBytes()).toString("base64");
  const didPrime = target.primeFromRelayState({
    snapshot_present: true,
    replica_base64: replicaBase64,
    snapshot_offset: 42,
  });

  assert.equal(didPrime, true);
  assert.equal(target.getText(), "# Live Demo Script\n\nShared manual");
  assert.equal(target.offset, 42);
});

test("recoverFromRelayHistory replays sync feed when startup stayed empty", async () => {
  const { AutomergeRelayClient } = await import("./automerge-relay.js");
  const source = new AutomergeRelayClient({
    basePath: "/api/local/documents/demo",
    participantID: "browser-source",
    documentID: "demo",
    awareness: { on() {} },
    capabilities: {},
  });
  source.postChange = async () => {};
  source.initialSyncReady = true;
  source.replaceText("# Live Demo Script\n\nRecovered over HTTP fallback");
  await Promise.resolve();

  const record = {
    participant_id: "browser-source",
    recipient_id: "",
    envelope_cid: "cid-1",
    message_base64: Buffer.from(Automerge.getLastLocalChange(source.doc)).toString("base64"),
  };

  const originalFetch = globalThis.fetch;
  globalThis.fetch = async () => ({
    ok: true,
    async json() {
      return {
        document_id: "demo",
        messages: [record],
        next_offset: 1,
      };
    },
  });

  try {
    const target = new AutomergeRelayClient({
      basePath: "/api/local/documents/demo",
      participantID: "browser-target",
      documentID: "demo",
      awareness: { on() {} },
      capabilities: {},
    });
    await target.recoverFromRelayHistory({ snapshot_present: false, message_count: 1 });
    assert.equal(target.getText(), "# Live Demo Script\n\nRecovered over HTTP fallback");
    assert.equal(target.offset, 1);
  } finally {
    globalThis.fetch = originalFetch;
  }
});

test("recoverFromRelayHistory replays from offset zero when a primed snapshot is blank", async () => {
  const { AutomergeRelayClient } = await import("./automerge-relay.js");
  const source = new AutomergeRelayClient({
    basePath: "/api/local/documents/demo",
    participantID: "browser-source",
    documentID: "demo",
    awareness: { on() {} },
    capabilities: {},
  });
  source.postChange = async () => {};
  source.initialSyncReady = true;
  source.replaceText("# Live Demo Script\n\nRecovered after stale snapshot");
  await Promise.resolve();

  const record = {
    participant_id: "browser-source",
    recipient_id: "",
    envelope_cid: "cid-stale",
    message_base64: Buffer.from(Automerge.getLastLocalChange(source.doc)).toString("base64"),
  };

  const originalFetch = globalThis.fetch;
  const seenQueries = [];
  globalThis.fetch = async (url) => {
    const parsed = new URL(url, "http://example.test");
    seenQueries.push(parsed.searchParams.get("since"));
    return {
      ok: true,
      async json() {
        return {
          document_id: "demo",
          messages: [record],
          next_offset: 1,
        };
      },
    };
  };

  try {
    const target = new AutomergeRelayClient({
      basePath: "/api/local/documents/demo",
      participantID: "browser-target",
      documentID: "demo",
      awareness: { on() {} },
      capabilities: {},
    });
    const blankReplicaBase64 = Buffer.from(target.getReplicaBytes()).toString("base64");
    target.primeFromRelayState({
      snapshot_present: true,
      replica_base64: blankReplicaBase64,
      snapshot_offset: 99,
    });
    await target.recoverFromRelayHistory({ snapshot_present: true, snapshot_offset: 99, message_count: 1 });
    assert.deepEqual(seenQueries, ["0"]);
    assert.equal(target.getText(), "# Live Demo Script\n\nRecovered after stale snapshot");
    assert.equal(target.offset, 1);
  } finally {
    globalThis.fetch = originalFetch;
  }
});
