import { safeLocalStorage, safeSessionStorage } from "./safe-storage.js";

const PARTICIPANT_KEY = "grid-editor-participant-id";
const WELCOME_KEY = "grid-editor-dismissed-welcome";

export function getOrCreateParticipantID(randomUUID = crypto.randomUUID.bind(crypto)) {
  const existing = safeSessionStorage.getItem(PARTICIPANT_KEY);
  if (existing) {
    return existing;
  }
  const created = `browser-${randomUUID()}`;
  safeSessionStorage.setItem(PARTICIPANT_KEY, created);
  return created;
}

export function isWelcomeDismissed() {
  return safeLocalStorage.getItem(WELCOME_KEY) === "true";
}

export function dismissWelcome() {
  safeLocalStorage.setItem(WELCOME_KEY, "true");
}
