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

test("normalizeState filters dangling recent and tab references", () => {
  const state = normalizeState({
    documents: { demo: { title: "Demo" } },
    recent: ["demo", "missing"],
    openTabs: ["demo", "missing"],
  });
  assert.deepEqual(state.recent, ["demo"]);
  assert.deepEqual(state.openTabs, ["demo"]);
});
