// Intent: Keep review/history helpers pure so the Phase 3 browser surfaces can
// be tested without the relay or DOM, while preserving a seam for later
// PromiseGrid-native review models. Source: DI-safor; DI-lapek
export function summarizeDocument(text) {
  const lines = String(text || "").split("\n").map((line) => line.trim()).filter(Boolean);
  if (lines.length === 0) {
    return "No summary yet.";
  }
  const headings = lines.filter((line) => line.startsWith("#")).slice(0, 2).map((line) => line.replace(/^#+\s*/, ""));
  const prose = lines.filter((line) => !line.startsWith("#")).slice(0, 3);
  return [...headings, ...prose].join(" ").slice(0, 240) || "No summary yet.";
}

export function extractMentions(text, participantIndex = new Map()) {
  const matches = [...String(text || "").matchAll(/(^|\s)@([a-zA-Z0-9._-]+)/g)];
  return matches.map((match) => {
    const label = match[2];
    const participantID = participantIndex.get(label.toLowerCase()) || "";
    return {
      label,
      participantID,
    };
  });
}

export function inferVersionName(title, text) {
  const heading = String(text || "").match(/^#\s+(.+)$/m)?.[1]?.trim();
  if (heading) {
    return heading.slice(0, 60);
  }
  return `${title || "Document"} version`;
}
