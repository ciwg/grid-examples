import test from "node:test";
import assert from "node:assert/strict";

import { normalizeBrowserColor } from "./color.js";

test("normalizeBrowserColor keeps valid colors and uppercases them", () => {
  assert.equal(normalizeBrowserColor("#1d6fd6"), "#1D6FD6");
});

test("normalizeBrowserColor falls back for invalid input", () => {
  assert.equal(normalizeBrowserColor("blue"), "#1D6FD6");
  assert.equal(normalizeBrowserColor(""), "#1D6FD6");
});
