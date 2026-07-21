export function presenceState(lastSeenAt, profile) {
  if (!lastSeenAt) {
    return "live";
  }
  const ageMs = Date.now() - new Date(lastSeenAt).getTime();
  // Intent: Render awareness using the approved demo/normal lifecycle windows
  // so the main peer list answers "who is here now?" while still giving
  // demos enough time before a peer is dimmed or removed. Source: DI-mivor;
  // DI-vasul
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
