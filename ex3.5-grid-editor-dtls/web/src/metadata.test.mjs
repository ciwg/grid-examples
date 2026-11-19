import test from "node:test";
import assert from "node:assert/strict";

import { formatMetadataList, metadataDisplayTitle, normalizeMetadataRecord, parseMetadataList } from "./metadata.js";

test("parseMetadataList trims empties and deduplicates case-insensitively", () => {
  assert.deepEqual(parseMetadataList("grid, docs, Grid, , team "), ["grid", "docs", "team"]);
});

test("formatMetadataList joins labels for form display", () => {
  assert.equal(formatMetadataList(["grid", "docs"]), "grid, docs");
});

test("normalizeMetadataRecord restores arrays and ids", () => {
  const record = normalizeMetadataRecord("demo", { title: "Demo", tags: ["grid"], favorite: true });
  assert.equal(record.document_id, "demo");
  assert.deepEqual(record.tags, ["grid"]);
  assert.deepEqual(record.collections, []);
  assert.equal(record.favorite, true);
  assert.equal(record.archived, false);
});

test("metadataDisplayTitle prefers relay title over local fallback", () => {
  assert.equal(metadataDisplayTitle("demo", "Local", { title: "Relay" }), "Relay");
  assert.equal(metadataDisplayTitle("demo", "Local", { title: "" }), "Local");
  assert.equal(metadataDisplayTitle("demo", "", { title: "" }), "Document demo");
});
