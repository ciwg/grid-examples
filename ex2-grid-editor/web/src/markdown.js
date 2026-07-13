export function renderMarkdown(source) {
  // Intent: Provide a repo-local markdown preview for Phase 2 workflow and
  // Phase 3 review/content surfaces without adding another heavy browser
  // dependency before the larger document workflow semantics are settled.
  // Source: DI-dovoz; DI-safor
  const escaped = escapeHTML(source || "");
  const lines = escaped.split("\n");
  const blocks = [];
  let paragraph = [];
  let inCode = false;
  let codeLines = [];
  let inList = false;
  let tableLines = [];

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

  const flushTable = () => {
    if (tableLines.length === 0) {
      return;
    }
    blocks.push(renderTable(tableLines));
    tableLines = [];
  };

  for (const line of lines) {
    if (line.startsWith("```")) {
      flushParagraph();
      flushList();
      flushTable();
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
      flushTable();
      const level = line.match(/^#+/)[0].length;
      blocks.push(`<h${level}>${inlineFormat(line.slice(level + 1))}</h${level}>`);
      continue;
    }

    if (/^\|.+\|$/.test(line.trim())) {
      flushParagraph();
      flushList();
      tableLines.push(line.trim());
      continue;
    }

    flushTable();

    const checklistMatch = line.match(/^[-*]\s+\[([ xX])\]\s+(.*)$/);
    if (checklistMatch) {
      flushParagraph();
      if (!inList) {
        blocks.push("<ul>");
        inList = true;
      }
      const checked = checklistMatch[1].toLowerCase() === "x" ? " checked" : "";
      blocks.push(`<li><label><input type="checkbox" disabled${checked}> ${inlineFormat(checklistMatch[2])}</label></li>`);
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
      flushTable();
      continue;
    }

    paragraph.push(line.trim());
  }

  flushParagraph();
  flushList();
  flushTable();

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
    .replace(/!\[([^\]]*)\]\(([^)]+)\)/g, (_match, alt, url) => renderMedia(alt, url))
    .replace(/\[([^\]]+)\]\(([^)]+)\)/g, (_match, label, url) => renderLink(label, url))
    .replace(/\*\*([^*]+)\*\*/g, "<strong>$1</strong>")
    .replace(/\*([^*]+)\*/g, "<em>$1</em>")
    .replace(/`([^`]+)`/g, "<code>$1</code>")
    .replace(/<u>(.*?)<\/u>/g, "<u>$1</u>");
}

function renderLink(label, url) {
  const safeURL = sanitizeURL(url);
  return `<a href="${safeURL}" target="_blank" rel="noreferrer">${label}</a>`;
}

function renderMedia(alt, url) {
  const safeURL = sanitizeURL(url);
  if (/\.(png|jpe?g|gif|webp|svg)$/i.test(safeURL) || safeURL.startsWith("data:image/")) {
    return `<figure><img src="${safeURL}" alt="${alt}"><figcaption>${alt}</figcaption></figure>`;
  }
  if (/\.(mp4|webm|ogg)$/i.test(safeURL)) {
    return `<figure><video controls src="${safeURL}"></video><figcaption>${alt}</figcaption></figure>`;
  }
  return `<span class="md-image">🖼 ${alt}</span>`;
}

function renderTable(lines) {
  const rows = lines.map((line) => line.slice(1, -1).split("|").map((cell) => cell.trim()));
  const [header, ...body] = rows;
  const headHTML = `<tr>${header.map((cell) => `<th>${inlineFormat(cell)}</th>`).join("")}</tr>`;
  const bodyHTML = body
    .filter((row) => !row.every((cell) => /^:?-{3,}:?$/.test(cell)))
    .map((row) => `<tr>${row.map((cell) => `<td>${inlineFormat(cell)}</td>`).join("")}</tr>`)
    .join("");
  return `<table><thead>${headHTML}</thead><tbody>${bodyHTML}</tbody></table>`;
}

function sanitizeURL(url) {
  const value = String(url || "");
  if (value.startsWith("javascript:")) {
    return "#";
  }
  return value;
}

function escapeHTML(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;");
}
