export function relayHasAuthoritativeHistory(state) {
  if (!state || typeof state !== "object") {
    return false;
  }
  if (state.snapshot_present) {
    return true;
  }
  return Number(state.message_count || 0) > 0;
}

export function shouldApplySeed(seed, relayText, state) {
  // Intent: Only let browser-local seed content populate brand-new documents.
  // Once the relay has snapshot or message history, the relay must be treated
  // as authoritative to avoid opening the wrong doc contents in the browser.
  // Source: DI-ramuv; DI-lumek; DI-gafit
  if (!seed) {
    return false;
  }
  if (relayHasAuthoritativeHistory(state)) {
    return false;
  }
  return relayText === "";
}
