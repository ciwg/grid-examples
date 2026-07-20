import test from "node:test";
import assert from "node:assert/strict";

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
