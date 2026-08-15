import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import test from "node:test";
import { fileURLToPath } from "node:url";

const sourceRoot = new URL("./", import.meta.url);

async function source(path) {
  return readFile(new URL(path, sourceRoot), "utf8");
}

test("Chrome color control has a visible native target and in-page fallback preview", async () => {
  const [html, css, main] = await Promise.all([
    source("../index.html"),
    source("../style.css"),
    source("main.js"),
  ]);

  assert.match(html, /<input id="color" type="color" value="#1d6fd6">/);
  assert.match(html, /id="color-preview-swatch"/);
  assert.match(html, /id="color-preview-value"/);
  assert.match(css, /input\[type="color"\]\s*\{\s*min-height: 48px;/);
  assert.match(css, /input\[type="color"\]::\-webkit\-color\-swatch-wrapper/);
  assert.match(css, /input\[type="color"\]::\-webkit\-color\-swatch/);
  assert.match(main, /function renderColorPreview\(name, value\)/);
  assert.match(main, /colorInput\.addEventListener\("change"/);
});

test("peer cursor rendering retains the received peer color hook", async () => {
  const editor = await source("editor.js");

  assert.match(editor, /--grid-peer-color/);
  assert.match(editor, /this\.color/);
});

test("Phase 2 controls expose the documented workflow labels and handlers", async () => {
  const [html, main] = await Promise.all([
    source("../index.html"),
    source("main.js"),
  ]);

  assert.match(html, /id="new-doc" type="button">New Shared Doc/);
  assert.match(html, /id="search-button" class="toolbar-button" type="button">Search/);
  assert.match(html, /id="preview-button" class="toolbar-button" type="button">Preview/);
  assert.match(html, /id="split-button" class="toolbar-button" type="button">Split View/);
  assert.match(html, /id="export-button" class="toolbar-button" type="button">Export \/ Exchange/);
  assert.match(main, /function createNewDocument\(\)/);
  assert.match(main, /function openSearch\(\)/);
  assert.match(main, /function togglePreview\(\)/);
  assert.match(main, /function toggleSplit\(\)/);
  assert.match(main, /function openExport\(\)/);
});
