import test from "node:test";
import assert from "node:assert/strict";

import { resolveFormattingSelection, wrapSelectedText } from "./formatting.js";

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

test("wrapSelectedText keeps both wrappers after repeated formatting", () => {
  const first = wrapSelectedText("word", 0, 4, "**", "**");
  const second = wrapSelectedText(first.text, first.selectionFrom, first.selectionTo, "<u>", "</u>");
  assert.match(second.text, /\*\*<u>word<\/u>\*\*|<u>\*\*word\*\*<\/u>/);
  assert.match(second.text, /\*\*/);
  assert.match(second.text, /<u>.*<\/u>/);
});

test("resolveFormattingSelection prefers the current non-empty selection", () => {
  assert.deepEqual(
    resolveFormattingSelection({ from: 2, to: 6 }, { from: 0, to: 4 }),
    { from: 2, to: 6 },
  );
});

test("resolveFormattingSelection falls back to the previous non-empty selection", () => {
  assert.deepEqual(
    resolveFormattingSelection({ from: 5, to: 5 }, { from: 1, to: 4 }),
    { from: 1, to: 4 },
  );
});
