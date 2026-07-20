import test from "node:test";
import assert from "node:assert/strict";
import { readFileSync } from "node:fs";

const indexHTML = readFileSync(new URL("../index.html", import.meta.url), "utf8");

test("documented sidebar sections exist in the page shell", () => {
  const expectedHeadings = [
    "<h2>PromiseGrid Flow</h2>",
    "<h2>Relay</h2>",
    "<h2>Document</h2>",
    "<h2>You</h2>",
    "<h2>Workspace</h2>",
    "<h2>Metadata</h2>",
    "<h2>Peers</h2>",
    "<h2>Review</h2>",
  ];
  for (const heading of expectedHeadings) {
    assert.match(indexHTML, new RegExp(escapeRegExp(heading)));
  }
});

test("top-left sidebar order puts PromiseGrid flow and relay ahead of document controls", () => {
  const flowIndex = indexHTML.indexOf('<h2>PromiseGrid Flow</h2>');
  const relayIndex = indexHTML.indexOf('<h2>Relay</h2>');
  const documentIndex = indexHTML.indexOf('<h2>Document</h2>');
  assert.ok(flowIndex >= 0, "expected PromiseGrid Flow heading");
  assert.ok(relayIndex >= 0, "expected Relay heading");
  assert.ok(documentIndex >= 0, "expected Document heading");
  assert.ok(flowIndex < relayIndex, "expected PromiseGrid Flow before Relay");
  assert.ok(relayIndex < documentIndex, "expected Relay before Document");
});

test("page includes the documented editor and PromiseGrid controls", () => {
  const expectedIDs = [
    "search-button",
    "preview-button",
    "split-button",
    "comment-button",
    "summary-button",
    "debug-button",
    "transport-mode",
    "trace-caption",
    "message-trace",
    "peer-badges",
    "metadata-results",
    "published-list",
  ];
  for (const id of expectedIDs) {
    assert.match(indexHTML, new RegExp(`id="${escapeRegExp(id)}"`));
  }
});

test("page includes the documented hidden overlays", () => {
  const overlays = [
    "settings-panel",
    "help-panel",
    "search-panel",
    "export-panel",
    "comment-panel",
    "summary-panel",
    "debug-panel",
  ];
  for (const id of overlays) {
    assert.match(indexHTML, new RegExp(`id="${escapeRegExp(id)}"`));
  }
});

function escapeRegExp(value) {
  return value.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
}
