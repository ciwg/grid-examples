// Intent: Keep Phase 4 metadata parsing and normalization in pure helpers so
// relay-backed catalog labels, favorites, archive state, and search results
// can be tested without driving the full browser UI shell. Source: DI-loruk;
// DI-sukip
export function parseMetadataList(raw) {
  const seen = new Set();
  const values = [];
  for (const entry of String(raw || "").split(",")) {
    const value = entry.trim();
    if (!value) {
      continue;
    }
    const key = value.toLowerCase();
    if (seen.has(key)) {
      continue;
    }
    seen.add(key);
    values.push(value);
  }
  return values;
}

export function formatMetadataList(values) {
  return Array.isArray(values) ? values.join(", ") : "";
}

export function normalizeMetadataRecord(documentID, record = {}) {
  return {
    offset: Number(record.offset || 0),
    envelope_cid: String(record.envelope_cid || ""),
    document_id: String(record.document_id || documentID || ""),
    author: String(record.author || ""),
    participant_id: String(record.participant_id || ""),
    title: String(record.title || ""),
    description: String(record.description || ""),
    summary: String(record.summary || ""),
    tags: Array.isArray(record.tags) ? record.tags.filter(Boolean) : [],
    collections: Array.isArray(record.collections) ? record.collections.filter(Boolean) : [],
    favorite: Boolean(record.favorite),
    archived: Boolean(record.archived),
    updated_at: String(record.updated_at || ""),
    embodiment: String(record.embodiment || ""),
    received_at: String(record.received_at || ""),
    lamport: Number(record.lamport || 0),
  };
}

export function metadataDisplayTitle(documentID, localTitle, record) {
  const relayTitle = String(record?.title || "").trim();
  if (relayTitle) {
    return relayTitle;
  }
  const local = String(localTitle || "").trim();
  if (local) {
    return local;
  }
  return `Document ${documentID}`;
}
