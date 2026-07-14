import test from "node:test";
import assert from "node:assert/strict";

import { DocumentRegistry, normalizeState, templateCatalog } from "./documents.js";

class MemoryStorage {
  constructor() {
    this.data = new Map();
  }
  getItem(key) {
    return this.data.has(key) ? this.data.get(key) : null;
  }
  setItem(key, value) {
    this.data.set(key, value);
  }
}

test("document registry tracks recent docs and timestamps", () => {
  const registry = new DocumentRegistry(new MemoryStorage());
  const first = registry.touchViewed("demo");
  assert.equal(first.documentID, "demo");
  assert.ok(first.createdAt);

  registry.touchEdited("demo");
  registry.openTab("demo");
  const recent = registry.listRecent();
  assert.equal(recent.length, 1);
  assert.equal(recent[0].documentID, "demo");
  assert.ok(recent[0].lastEditedAt);
});

test("duplicate document seeds a new local copy", () => {
  const registry = new DocumentRegistry(new MemoryStorage());
  registry.updateTitle("demo", "Original");
  const duplicate = registry.duplicateDocument("demo", "demo-copy", "# copied");
  assert.equal(duplicate.documentID, "demo-copy");
  assert.equal(duplicate.title, "Original copy");
  assert.equal(registry.seedContent("demo-copy"), "# copied");
});

test("template catalog exposes phase 2 starter docs", () => {
  const labels = templateCatalog().map((entry) => entry.label);
  assert.ok(labels.includes("Meeting Notes"));
  assert.ok(labels.includes("Demo Sample"));
});

test("review metadata tracks comments, versions, and participants", () => {
  const registry = new DocumentRegistry(new MemoryStorage());
  registry.addComment("demo", {
    id: "comment-1",
    authorName: "Mallory",
    createdAt: "2026-07-13T10:00:00Z",
    body: "Looks good",
    reactions: [],
  });
  registry.addReaction("demo", "comment-1", {
    emoji: "👍",
    createdAt: "2026-07-13T10:01:00Z",
  });
  registry.toggleCommentResolved("demo", "comment-1", "Mallory");
  registry.addSavedVersion("demo", {
    id: "version-1",
    name: "Draft 1",
    createdAt: "2026-07-13T10:02:00Z",
    content: "# demo",
    replicaBase64: "AQID",
  });
  registry.noteParticipant("demo", {
    participantID: "browser-1",
    name: "Mallory",
    color: "#123456",
  });

  assert.equal(registry.listComments("demo").length, 1);
  assert.equal(registry.listComments("demo")[0].reactions[0].emoji, "👍");
  assert.equal(registry.listComments("demo")[0].resolved, true);
  assert.equal(registry.listSavedVersions("demo")[0].name, "Draft 1");
  assert.equal(registry.listSavedVersions("demo")[0].content, "# demo");
  assert.equal(registry.listRecentParticipants("demo")[0].name, "Mallory");
  assert.ok(registry.listActivity("demo").length >= 4);
});

test("normalizeState filters dangling recent and tab references", () => {
  const state = normalizeState({
    documents: { demo: { title: "Demo", comments: [{}], activity: [{}], recentParticipants: [{}], savedVersions: [{ content: "# doc", replicaBase64: "AQID" }] } },
    recent: ["demo", "missing"],
    openTabs: ["demo", "missing"],
  });
  assert.deepEqual(state.recent, ["demo"]);
  assert.deepEqual(state.openTabs, ["demo"]);
  assert.equal(state.documents.demo.comments.length, 1);
  assert.equal(state.documents.demo.savedVersions.length, 1);
  assert.equal(state.documents.demo.savedVersions[0].content, "# doc");
});
