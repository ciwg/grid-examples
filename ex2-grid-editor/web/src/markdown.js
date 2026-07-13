export function renderMarkdown(source) {
  // Intent: Provide a repo-local markdown preview for Phase 2 workflow and
  // export surfaces without adding another heavy browser dependency before the
  // larger document workflow semantics are settled. Source: DI-dovoz
  const escaped = escapeHTML(source || "");
  const lines = escaped.split("\n");
  const blocks = [];
  let paragraph = [];
  let inCode = false;
  let codeLines = [];
  let inList = false;

  const flushParagraph = () => {
    if (paragraph.length === 0) {
      return;
    }
    blocks.push(`<p>${inlineFormat(paragraph.join(" "))}</p>`);
    paragraph = [];
  };

  const flushList = () => {
    if (inList) {
      blocks.push("</ul>");
      inList = false;
    }
  };

  for (const line of lines) {
    if (line.startsWith("```")) {
      flushParagraph();
      flushList();
      if (inCode) {
        blocks.push(`<pre><code>${codeLines.join("\n")}</code></pre>`);
        codeLines = [];
        inCode = false;
      } else {
        inCode = true;
      }
      continue;
    }

    if (inCode) {
      codeLines.push(line);
      continue;
    }

    if (/^#{1,6}\s/.test(line)) {
      flushParagraph();
      flushList();
      const level = line.match(/^#+/)[0].length;
      blocks.push(`<h${level}>${inlineFormat(line.slice(level + 1))}</h${level}>`);
      continue;
    }

    if (/^[-*]\s/.test(line)) {
      flushParagraph();
      if (!inList) {
        blocks.push("<ul>");
        inList = true;
      }
      blocks.push(`<li>${inlineFormat(line.slice(2))}</li>`);
      continue;
    }

    if (line.trim() === "") {
      flushParagraph();
      flushList();
      continue;
    }

    paragraph.push(line.trim());
  }

  flushParagraph();
  flushList();

  if (inCode) {
    blocks.push(`<pre><code>${codeLines.join("\n")}</code></pre>`);
  }

  return blocks.join("\n");
}

export function extractHeadings(source) {
  const headings = [];
  const lines = (source || "").split("\n");
  for (let index = 0; index < lines.length; index += 1) {
    const line = lines[index];
    const match = line.match(/^(#{1,6})\s+(.*)$/);
    if (!match) {
      continue;
    }
    headings.push({
      level: match[1].length,
      text: match[2].trim(),
      line: index + 1,
    });
  }
  return headings;
}

function inlineFormat(text) {
  return text
    .replace(/!\[([^\]]*)\]\(([^)]+)\)/g, '<span class="md-image">🖼 $1</span>')
    .replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2">$1</a>')
    .replace(/\*\*([^*]+)\*\*/g, "<strong>$1</strong>")
    .replace(/\*([^*]+)\*/g, "<em>$1</em>")
    .replace(/`([^`]+)`/g, "<code>$1</code>")
    .replace(/<u>(.*?)<\/u>/g, "<u>$1</u>");
}

function escapeHTML(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;");
}
