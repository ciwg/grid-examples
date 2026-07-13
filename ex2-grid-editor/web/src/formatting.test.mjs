import test from "node:test";
import assert from "node:assert/strict";

import { wrapSelectedText } from "./formatting.js";

test("wrapSelectedText wraps selected text with markdown markers", () => {
  const next = wrapSelectedText("hello world", 6, 11, "**", "**");
  assert.equal(next.insert, "**world**");
  assert.equal(next.text, "hello **world**");
  assert.equal(next.selectionFrom, 8);
  assert.equal(next.selectionTo, 13);
});

test("wrapSelectedText inserts fallback text for an empty selection", () => {
  const next = wrapSelectedText("hello", 5, 5, "<u>", "</u>");
  assert.equal(next.text, "hello<u>text</u>");
  assert.equal(next.selectionFrom, 8);
  assert.equal(next.selectionTo, 12);
});

test("wrapSelectedText preserves combined bold and underline wrappers", () => {
  const bold = wrapSelectedText("word", 0, 4, "**", "**");
  const underlineInsideBold = wrapSelectedText(bold.text, bold.selectionFrom, bold.selectionTo, "<u>", "</u>");
  assert.equal(underlineInsideBold.text, "**<u>word</u>**");

  const underline = wrapSelectedText("word", 0, 4, "<u>", "</u>");
  const boldInsideUnderline = wrapSelectedText(underline.text, underline.selectionFrom, underline.selectionTo, "**", "**");
  assert.equal(boldInsideUnderline.text, "<u>**word**</u>");
});
