import test from "node:test";
import assert from "node:assert/strict";

import { extractMentions, inferVersionName, summarizeDocument } from "./review.js";

test("summarizeDocument prefers headings and early prose", () => {
  const summary = summarizeDocument("# Title\n\n## Deep\n\nFirst sentence.\nSecond sentence.");
  assert.match(summary, /Title/);
  assert.match(summary, /First sentence/);
});

test("extractMentions resolves visible names to stable ids when known", () => {
  const mentions = extractMentions("Hi @mallory and @bob", new Map([
    ["mallory", "browser-1"],
  ]));
  assert.deepEqual(mentions, [
    { label: "mallory", participantID: "browser-1" },
    { label: "bob", participantID: "" },
  ]);
});

test("inferVersionName uses heading before fallback title", () => {
  assert.equal(inferVersionName("Demo", "# Better name\n\ntext"), "Better name");
  assert.equal(inferVersionName("Demo", "plain"), "Demo version");
});
