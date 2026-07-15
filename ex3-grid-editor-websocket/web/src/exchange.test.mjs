import test from "node:test";
import assert from "node:assert/strict";

import { buildExportArtifact, buildPublishSource, parsePublishedURL } from "./exchange.js";

test("buildExportArtifact keeps raw markdown text unchanged", () => {
  const markdown = "# Title\n\n**bold** and <u>underlined</u>";
  const artifact = buildExportArtifact("markdown", "Demo", markdown, new Uint8Array([1, 2, 3]), () => "<p>ignored</p>");
  assert.equal(artifact.extension, "md");
  assert.equal(artifact.mime, "text/markdown;charset=utf-8");
  assert.equal(artifact.body, markdown);
});

test("buildExportArtifact renders html separately from markdown text", () => {
  const artifact = buildExportArtifact("html", "Demo", "# Title", new Uint8Array(), () => "<h1>Title</h1>");
  assert.equal(artifact.extension, "html");
  assert.match(artifact.body, /<h1>Title<\/h1>/);
});

test("buildPublishSource returns current state by default", () => {
  const source = buildPublishSource("# demo", "AQID", "Demo", [], "");
  assert.deepEqual(source, {
    sourceKind: "current",
    sourceVersionID: "",
    sourceVersionName: "",
    title: "Demo",
    summary: "demo",
    text: "# demo",
    replicaBase64: "AQID",
  });
});

test("buildPublishSource resolves a named saved version", () => {
  const source = buildPublishSource("ignored", "AQID", "Demo", [{
    id: "version-1",
    name: "Draft 1",
    content: "# Saved\n\ntext",
    replicaBase64: "CQgH",
    summary: "Saved summary",
  }], "Draft 1");
  assert.deepEqual(source, {
    sourceKind: "saved_version",
    sourceVersionID: "version-1",
    sourceVersionName: "Draft 1",
    title: "Draft 1",
    summary: "Saved summary",
    text: "# Saved\n\ntext",
    replicaBase64: "CQgH",
  });
});

test("buildPublishSource rejects incomplete saved versions", () => {
  const source = buildPublishSource("ignored", "AQID", "Demo", [{
    id: "version-1",
    name: "Draft 1",
    content: "",
    replicaBase64: "",
  }], "Draft 1");
  assert.equal(source, null);
});

test("parsePublishedURL resolves published exchange URLs", () => {
  assert.deepEqual(
    parsePublishedURL("http://relay.example/api/published/bafy123", "http://local.test"),
    { origin: "http://relay.example", envelopeCID: "bafy123" },
  );
  assert.deepEqual(
    parsePublishedURL("/api/published/bafy456", "http://local.test"),
    { origin: "http://local.test", envelopeCID: "bafy456" },
  );
});

test("parsePublishedURL rejects non-published URLs", () => {
  assert.equal(parsePublishedURL("http://relay.example/api/local/documents/demo", "http://local.test"), null);
});

test("parsePublishedURL rejects empty strings", () => {
  assert.equal(parsePublishedURL("", "http://local.test"), null);
});

test("buildExportArtifact keeps raw text for plain text and Automerge export", () => {
  const textArtifact = buildExportArtifact("text", "Demo", "**bold**", new Uint8Array([1, 2]), () => "<p>ignored</p>");
  assert.equal(textArtifact.body, "**bold**");
  assert.equal(textArtifact.extension, "txt");

  const replicaBytes = new Uint8Array([1, 2, 3]);
  const automergeArtifact = buildExportArtifact("automerge", "Demo", "ignored", replicaBytes, () => "<p>ignored</p>");
  assert.equal(automergeArtifact.body, replicaBytes);
  assert.equal(automergeArtifact.extension, "automerge");
});

test("buildExportArtifact keeps markdown markers intact", () => {
  const artifact = buildExportArtifact("markdown", "Demo", "## Heading\n\n* item", new Uint8Array(), () => "<p>ignored</p>");
  assert.equal(artifact.body, "## Heading\n\n* item");
});

test("buildExportArtifact preserves combined bold and underline markdown", () => {
  const boldOuter = buildExportArtifact("markdown", "Demo", "**<u>word</u>**", new Uint8Array(), () => "<p>ignored</p>");
  assert.equal(boldOuter.body, "**<u>word</u>**");

  const underlineOuter = buildExportArtifact("markdown", "Demo", "<u>**word**</u>", new Uint8Array(), () => "<p>ignored</p>");
  assert.equal(underlineOuter.body, "<u>**word**</u>");
});
