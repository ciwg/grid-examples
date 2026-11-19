import test from "node:test";
import assert from "node:assert/strict";
import { describePaneMode, nextPaneState } from "./panes.js";

test("preview turns on preview-only from editor-only mode", () => {
  assert.deepEqual(
    nextPaneState({ previewEnabled: false, splitEnabled: false }, "preview"),
    { previewEnabled: true, splitEnabled: false },
  );
});

test("preview turns preview off when already in preview-only mode", () => {
  assert.deepEqual(
    nextPaneState({ previewEnabled: true, splitEnabled: false }, "preview"),
    { previewEnabled: false, splitEnabled: false },
  );
});

test("preview exits split mode into preview-only mode", () => {
  assert.deepEqual(
    nextPaneState({ previewEnabled: true, splitEnabled: true }, "preview"),
    { previewEnabled: true, splitEnabled: false },
  );
});

test("split turns on split mode from editor-only mode", () => {
  assert.deepEqual(
    nextPaneState({ previewEnabled: false, splitEnabled: false }, "split"),
    { previewEnabled: true, splitEnabled: true },
  );
});

test("split turns split mode off into preview-only mode", () => {
  assert.deepEqual(
    nextPaneState({ previewEnabled: true, splitEnabled: true }, "split"),
    { previewEnabled: true, splitEnabled: false },
  );
});

test("preview-only mode hides the editor pane and shows preview", () => {
  assert.deepEqual(
    describePaneMode({ previewEnabled: true, splitEnabled: false }),
    { showEditor: false, showPreview: true, splitEnabled: false },
  );
});

test("split mode shows both panes", () => {
  assert.deepEqual(
    describePaneMode({ previewEnabled: true, splitEnabled: true }),
    { showEditor: true, showPreview: true, splitEnabled: true },
  );
});

test("editor-only mode hides preview", () => {
  assert.deepEqual(
    describePaneMode({ previewEnabled: false, splitEnabled: false }),
    { showEditor: true, showPreview: false, splitEnabled: false },
  );
});
