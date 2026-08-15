// Intent: Render awareness using the approved demo/normal lifecycle windows
// so the main peer list answers "who is here now?" while keeping presence an
// embodiment-local interpretation rather than a relay membership decision.
// Source: DI-mivor; DI-vasul; DI-dizut; DI-dazin.
export function presenceState(lastSeenAt, profile, now = Date.now()) {
  if (!lastSeenAt) {
    return "live";
  }
  const ageMs = now - new Date(lastSeenAt).getTime();
  const thresholds = profile === "normal"
    ? { live: 60_000, stale: 5 * 60_000, offline: 15 * 60_000 }
    : { live: 5 * 60_000, stale: 15 * 60_000, offline: 30 * 60_000 };
  if (ageMs <= thresholds.live) {
    return "live";
  }
  if (ageMs <= thresholds.stale) {
    return "stale";
  }
  if (ageMs <= thresholds.offline) {
    return "offline";
  }
  return "gone";
}

// Intent: Wake only when a local display state can change, keeping the
// lifecycle accurate without polling or producing relay traffic. Source:
// DI-dizut; DI-dazin.
export function schedulePresenceRefresh(peers, profile, onRefresh, scheduler = {
  now: () => Date.now(),
  setTimeout: window.setTimeout.bind(window),
  clearTimeout: window.clearTimeout.bind(window),
}) {
  const now = scheduler.now();
  const thresholds = profile === "normal"
    ? [60_000, 5 * 60_000, 15 * 60_000]
    : [5 * 60_000, 15 * 60_000, 30 * 60_000];
  let nextRefreshAt = Infinity;
  for (const peer of peers) {
    const observedAt = new Date(peer.lastSeenAt).getTime();
    if (!Number.isFinite(observedAt)) {
      continue;
    }
    for (const threshold of thresholds) {
      const boundary = observedAt + threshold + 1;
      if (boundary > now && boundary < nextRefreshAt) {
        nextRefreshAt = boundary;
      }
    }
  }
  if (!Number.isFinite(nextRefreshAt)) {
    return null;
  }
  return scheduler.setTimeout(onRefresh, Math.max(1, nextRefreshAt - now));
}

// Intent: Ensure an embodiment never retains a timer for a document or
// awareness projection it no longer presents. Source: DI-dizut; DI-dazin.
export function cancelPresenceRefresh(presenceRefreshTimer, scheduler = {
  now: () => Date.now(),
  setTimeout: window.setTimeout.bind(window),
  clearTimeout: window.clearTimeout.bind(window),
}) {
  if (presenceRefreshTimer !== null) {
    scheduler.clearTimeout(presenceRefreshTimer);
  }
}
