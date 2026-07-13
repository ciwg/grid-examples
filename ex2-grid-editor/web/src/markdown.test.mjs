import test from "node:test";
import assert from "node:assert/strict";

import { extractHeadings, renderMarkdown } from "./markdown.js";

test("renderMarkdown formats headings, emphasis, and code blocks", () => {
  const html = renderMarkdown("# Title\n\nSome **bold** text and `code`.\n\n```js\nconst x = 1;\n```");
  assert.match(html, /<h1>Title<\/h1>/);
  assert.match(html, /<strong>bold<\/strong>/);
  assert.match(html, /<code>code<\/code>/);
  assert.match(html, /<pre><code>const x = 1;/);
});

test("extractHeadings returns line-aware heading metadata", () => {
  const headings = extractHeadings("# Top\n\n## Deep\ntext");
  assert.deepEqual(headings, [
    { level: 1, text: "Top", line: 1 },
    { level: 2, text: "Deep", line: 3 },
  ]);
});

test("renderMarkdown supports checklists and tables", () => {
  const html = renderMarkdown("- [x] Done\n- [ ] Todo\n\n| Name | Value |\n| --- | --- |\n| A | B |");
  assert.match(html, /type="checkbox" disabled checked/);
  assert.match(html, /<table>/);
  assert.match(html, /<th>Name<\/th>/);
  assert.match(html, /<td>B<\/td>/);
});
