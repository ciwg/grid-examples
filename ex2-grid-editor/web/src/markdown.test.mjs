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

test("renderMarkdown preserves underline tags and safe links", () => {
  const html = renderMarkdown("<u>underlined</u> [ok](https://example.com) [bad](javascript:alert(1))");
  assert.match(html, /<u>underlined<\/u>/);
  assert.match(html, /href="https:\/\/example.com"/);
  assert.match(html, /href="#"/);
});

test("renderMarkdown handles bold and underline in either nesting order", () => {
  const boldOuter = renderMarkdown("**<u>word</u>**");
  assert.match(boldOuter, /<strong><u>word<\/u><\/strong>/);

  const underlineOuter = renderMarkdown("<u>**word**</u>");
  assert.match(underlineOuter, /<u><strong>word<\/strong><\/u>/);
});

test("renderMarkdown renders image references", () => {
  const html = renderMarkdown("![Alt](https://example.com/demo.png)");
  assert.match(html, /<img src="https:\/\/example.com\/demo.png" alt="Alt">/);
});

test("extractHeadings ignores non-heading markdown lines", () => {
  const headings = extractHeadings("plain\n- list\n### Deep");
  assert.deepEqual(headings, [
    { level: 3, text: "Deep", line: 3 },
  ]);
});

test("renderMarkdown preserves fenced code content", () => {
  const html = renderMarkdown("```txt\n**not bold**\n```");
  assert.match(html, /<pre><code>\*\*not bold\*\*<\/code><\/pre>/);
});

test("renderMarkdown escapes raw html outside allowed underline tags", () => {
  const html = renderMarkdown("<script>alert(1)</script>");
  assert.doesNotMatch(html, /<script>/);
  assert.match(html, /&lt;script&gt;alert\(1\)&lt;\/script&gt;/);
});
