import test from "node:test";
import assert from "node:assert/strict";

import { extractMentions, inferVersionName, summarizeDocument } from "./review.js";

test("summarizeDocument prefers headings and early prose", () => {
  const summary = summarizeDocument("# Title\n\n## Deep\n\nFirst sentence.\nSecond sentence.");
  assert.match(summary, /Title/);
  assert.match(summary, /First sentence/);
});

test("summarizeDocument falls back cleanly for empty input", () => {
  assert.equal(summarizeDocument(""), "No summary yet.");
});

test("summarizeDocument truncates long output", () => {
  const text = "# Title\n\n" + "word ".repeat(100);
  assert.ok(summarizeDocument(text).length <= 240);
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

test("extractMentions is case-insensitive for participant lookups", () => {
  const mentions = extractMentions("Hi @Mallory", new Map([
    ["mallory", "browser-1"],
  ]));
  assert.deepEqual(mentions, [
    { label: "Mallory", participantID: "browser-1" },
  ]);
});

test("extractMentions ignores text without mentions", () => {
  assert.deepEqual(extractMentions("plain text"), []);
});

test("inferVersionName uses heading before fallback title", () => {
  assert.equal(inferVersionName("Demo", "# Better name\n\ntext"), "Better name");
  assert.equal(inferVersionName("Demo", "plain"), "Demo version");
});

test("inferVersionName trims long headings", () => {
  const name = inferVersionName("Demo", `# ${"a".repeat(80)}`);
  assert.equal(name.length, 60);
});

test("inferVersionName prefers heading over title fallback", () => {
  assert.equal(inferVersionName("Ignored", "# Kept\n\nbody"), "Kept");
});
