// Intent: Keep the PromiseGrid flow badge, trace caption, and protocol-label
// styling rules testable outside the browser bootstrap so the demo surface can
// be verified without spinning up the whole page. Source: DI-holoz; DI-dogub
export function formatTransportSummary(syncMode, awarenessMode) {
  return `browser sync: ${syncMode || "-"} · awareness: ${awarenessMode || "-"} · path: relay`;
}

// Intent: Keep the on-screen trace caption aligned with actual relay-observed
// message availability so the page explains whether viewers are seeing live
// traffic or an empty document. Source: DI-dogub
export function traceCaption(documentID, entryCount) {
  if (entryCount > 0) {
    return `Live relay-observed PromiseGrid traffic for ${documentID}. Click a message for decoded payload and raw CBOR base64.`;
  }
  return `No relay traffic yet for ${documentID}. Start typing to watch signed messages flow.`;
}

// Intent: Keep protocol-specific trace styling stable as the page distinguishes
// document sync, awareness, metadata, and publish traffic in the same list.
// Source: DI-dogub
export function traceProtocolClass(protocolName) {
  if (protocolName === "live-awareness") {
    return "awareness";
  }
  if (protocolName === "document-metadata") {
    return "metadata";
  }
  if (protocolName === "publish-document") {
    return "publish";
  }
  return "document";
}
