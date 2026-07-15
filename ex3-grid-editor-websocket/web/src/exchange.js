import { summarizeDocument } from "./review.js";

// Intent: Keep Phase 4 exchange and export rules in pure helpers so the
// browser can preserve raw markdown bytes and stable publish/import behavior
// under direct tests instead of only through the full UI shell.
// Source: DI-tavul; DI-gosaf
export function buildPublishSource(currentText, currentReplicaBase64, currentTitle, savedVersions, requestedValue) {
  const current = {
    sourceKind: "current",
    sourceVersionID: "",
    sourceVersionName: "",
    title: currentTitle,
    summary: summarizeDocument(currentText),
    text: currentText,
    replicaBase64: currentReplicaBase64,
  };
  if (!Array.isArray(savedVersions) || savedVersions.length === 0) {
    return current;
  }
  const normalized = String(requestedValue || "").trim();
  if (!normalized || normalized.toLowerCase() === "current") {
    return current;
  }
  const version = savedVersions.find((value) => value.name === normalized || value.id === normalized);
  if (!version || !version.content || !version.replicaBase64) {
    return null;
  }
  return {
    sourceKind: "saved_version",
    sourceVersionID: version.id,
    sourceVersionName: version.name,
    title: version.name,
    summary: version.summary || summarizeDocument(version.content),
    text: version.content,
    replicaBase64: version.replicaBase64,
  };
}

export function parsePublishedURL(raw, origin) {
  try {
    const url = new URL(String(raw || "").trim(), origin);
    const parts = url.pathname.split("/").filter(Boolean);
    const publishedIndex = parts.indexOf("published");
    if (publishedIndex === -1 || !parts[publishedIndex + 1]) {
      return null;
    }
    return {
      origin: url.origin,
      envelopeCID: parts[publishedIndex + 1],
    };
  } catch {
    return null;
  }
}

export function buildExportArtifact(format, title, text, replicaBytes, renderMarkdown) {
  const safeTitle = title || "document";
  const rawText = text || "";
  if (format === "html") {
    return {
      extension: "html",
      mime: "text/html;charset=utf-8",
      body: wrapHTML(safeTitle, renderMarkdown(rawText)),
    };
  }
  if (format === "text") {
    return {
      extension: "txt",
      mime: "text/plain;charset=utf-8",
      body: rawText,
    };
  }
  if (format === "automerge") {
    return {
      extension: "automerge",
      mime: "application/octet-stream",
      body: replicaBytes,
    };
  }
  return {
    extension: "md",
    mime: "text/markdown;charset=utf-8",
    body: rawText,
  };
}

function wrapHTML(title, body) {
  return `<!doctype html><html><head><meta charset="utf-8"><title>${escapeHTML(title)}</title></head><body>${body}</body></html>`;
}

function escapeHTML(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;");
}
