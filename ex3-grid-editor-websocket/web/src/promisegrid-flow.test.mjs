import test from "node:test";
import assert from "node:assert/strict";

import { formatTransportSummary, traceCaption, traceProtocolClass } from "./promisegrid-flow.js";

test("formatTransportSummary reports relay websocket modes explicitly", () => {
  assert.equal(
    formatTransportSummary("websocket", "websocket"),
    "browser sync: websocket · awareness: websocket · path: relay",
  );
});

test("formatTransportSummary falls back to dashes when modes are missing", () => {
  assert.equal(
    formatTransportSummary("", ""),
    "browser sync: - · awareness: - · path: relay",
  );
});

test("traceCaption explains live traffic when entries exist", () => {
  assert.equal(
    traceCaption("demo", 3),
    "Live relay-observed PromiseGrid traffic for demo. Click a message for decoded payload and raw CBOR base64.",
  );
});

test("traceCaption explains the empty state when no relay traffic exists", () => {
  assert.equal(
    traceCaption("demo", 0),
    "No relay traffic yet for demo. Start typing to watch signed messages flow.",
  );
});

test("traceProtocolClass maps known PromiseGrid protocols to stable UI classes", () => {
  assert.equal(traceProtocolClass("live-document"), "document");
  assert.equal(traceProtocolClass("live-awareness"), "awareness");
  assert.equal(traceProtocolClass("document-metadata"), "metadata");
  assert.equal(traceProtocolClass("publish-document"), "publish");
});
