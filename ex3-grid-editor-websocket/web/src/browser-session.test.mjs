import test from "node:test";
import assert from "node:assert/strict";

import { dismissWelcome, getOrCreateParticipantID, isWelcomeDismissed } from "./browser-session.js";

test("getOrCreateParticipantID reuses the safe session storage value", () => {
  const first = getOrCreateParticipantID(() => "abc123");
  const second = getOrCreateParticipantID(() => "different");
  assert.equal(first, "browser-abc123");
  assert.equal(second, "browser-abc123");
});

test("dismissWelcome persists through safe local storage", () => {
  dismissWelcome();
  assert.equal(isWelcomeDismissed(), true);
});
