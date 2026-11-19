import test from "node:test";
import assert from "node:assert/strict";

import { findUnderlineRanges } from "./underline.js";

test("findUnderlineRanges finds a simple underline pair", () => {
  const ranges = findUnderlineRanges("hello <u>world</u>");
  assert.deepEqual(ranges, [{
    openFrom: 6,
    openTo: 9,
    contentFrom: 9,
    contentTo: 14,
    closeFrom: 14,
    closeTo: 18,
  }]);
});

test("findUnderlineRanges keeps multiple underline pairs", () => {
  const ranges = findUnderlineRanges("<u>one</u> and <u>two</u>");
  assert.equal(ranges.length, 2);
  assert.equal(ranges[0].contentFrom, 3);
  assert.equal(ranges[0].contentTo, 6);
  assert.equal(ranges[1].contentFrom, 18);
  assert.equal(ranges[1].contentTo, 21);
});

test("findUnderlineRanges ignores unmatched tags", () => {
  assert.deepEqual(findUnderlineRanges("<u>broken"), []);
  assert.deepEqual(findUnderlineRanges("</u>broken"), []);
});

test("findUnderlineRanges supports bold around underline content", () => {
  const [range] = findUnderlineRanges("**<u>word</u>**");
  assert.equal(range.contentFrom, 5);
  assert.equal(range.contentTo, 9);
});
