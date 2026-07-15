import test from "node:test";
import assert from "node:assert/strict";

import { DocumentRegistry, generateDemoText, normalizeState, templateCatalog } from "./documents.js";

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

test("setTitle updates stored title without needing review activity", () => {
  const registry = new DocumentRegistry(new MemoryStorage());
  const document = registry.setTitle("demo", "Relay Title");
  assert.equal(document.title, "Relay Title");
  assert.equal(registry.listActivity("demo").length, 0);
});

test("document registry keeps recent docs in recency order", () => {
  const registry = new DocumentRegistry(new MemoryStorage());
  registry.touchViewed("first");
  registry.touchViewed("second");
  registry.touchViewed("first");
  assert.deepEqual(registry.listRecent().map((value) => value.documentID), ["first", "second"]);
});

test("document registry manages open tabs in most-recent order", () => {
  const registry = new DocumentRegistry(new MemoryStorage());
  registry.openTab("first");
  registry.openTab("second");
  registry.openTab("first");
  assert.deepEqual(registry.listOpenTabs().map((value) => value.documentID), ["first", "second"]);
  registry.closeTab("first");
  assert.deepEqual(registry.listOpenTabs().map((value) => value.documentID), ["second"]);
});

test("duplicate document seeds a new local copy", () => {
  const registry = new DocumentRegistry(new MemoryStorage());
  registry.updateTitle("demo", "Original");
  const duplicate = registry.duplicateDocument("demo", "demo-copy", "# copied");
  assert.equal(duplicate.documentID, "demo-copy");
  assert.equal(duplicate.title, "Original copy");
  assert.equal(registry.seedContent("demo-copy"), "# copied");
});

test("duplicate document resets review metadata and saved versions", () => {
  const registry = new DocumentRegistry(new MemoryStorage());
  registry.addComment("demo", {
    id: "comment-1",
    authorName: "Mallory",
    createdAt: "2026-07-13T10:00:00Z",
    body: "Looks good",
    reactions: [],
  });
  registry.addSavedVersion("demo", {
    id: "version-1",
    name: "Draft 1",
    createdAt: "2026-07-13T10:02:00Z",
    content: "# demo",
    replicaBase64: "AQID",
  });
  const duplicate = registry.duplicateDocument("demo", "demo-copy", "# copied");
  assert.equal(duplicate.comments.length, 0);
  assert.equal(duplicate.savedVersions.length, 0);
  assert.equal(duplicate.activity.length, 0);
});

test("template catalog exposes phase 2 starter docs", () => {
  const labels = templateCatalog().map((entry) => entry.label);
  assert.ok(labels.includes("Meeting Notes"));
  assert.ok(labels.includes("Demo Sample"));
});

test("generateDemoText includes checklist workflow hints", () => {
  assert.match(generateDemoText(), /export markdown/);
  assert.match(generateDemoText(), /Generated demo document/);
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

test("bookmarks and snapshots are tracked in activity order", () => {
  const registry = new DocumentRegistry(new MemoryStorage());
  registry.addBookmark("demo", {
    id: "bookmark-1",
    label: "Intro",
    createdAt: "2026-07-13T10:00:00Z",
  });
  registry.addSnapshot("demo", {
    id: "snapshot-1",
    title: "Snapshot 1",
    createdAt: "2026-07-13T10:01:00Z",
  });
  const document = registry.get("demo");
  assert.equal(document.bookmarks.length, 1);
  assert.equal(document.snapshots.length, 1);
  assert.equal(registry.listActivity("demo")[0].type, "snapshot");
});

test("touchExported records export timestamps", () => {
  const registry = new DocumentRegistry(new MemoryStorage());
  const document = registry.touchExported("demo");
  assert.ok(document.lastExportedAt);
  assert.equal(registry.listActivity("demo")[0].type, "exported");
});

test("noteParticipant replaces existing participant entries", () => {
  const registry = new DocumentRegistry(new MemoryStorage());
  registry.noteParticipant("demo", {
    participantID: "browser-1",
    name: "Mallory",
    color: "#123456",
    lastSeenAt: "2026-07-13T10:00:00Z",
  });
  registry.noteParticipant("demo", {
    participantID: "browser-1",
    name: "Mallory 2",
    color: "#654321",
    lastSeenAt: "2026-07-13T10:05:00Z",
  });
  const participants = registry.listRecentParticipants("demo");
  assert.equal(participants.length, 1);
  assert.equal(participants[0].name, "Mallory 2");
  assert.equal(participants[0].color, "#654321");
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

test("normalizeState restores default arrays and titles", () => {
  const state = normalizeState({
    documents: {
      demo: {},
    },
  });
  assert.equal(state.documents.demo.title, "Document demo");
  assert.deepEqual(state.documents.demo.comments, []);
  assert.deepEqual(state.documents.demo.bookmarks, []);
});

test("normalizeState keeps seed content when present", () => {
  const state = normalizeState({
    documents: {
      demo: { seedContent: "# seeded" },
    },
  });
  assert.equal(state.documents.demo.seedContent, "# seeded");
});
