import test from "node:test";
import assert from "node:assert/strict";

import { presenceState } from "./presence.js";

test("presenceState uses demo-friendly aging windows", (t) => {
  t.mock.timers.enable({ apis: ["Date"], now: new Date("2026-07-20T20:17:05Z") });
  assert.equal(presenceState("2026-07-20T20:16:30Z", "demo"), "live");
  assert.equal(presenceState("2026-07-20T20:10:00Z", "demo"), "stale");
  assert.equal(presenceState("2026-07-20T19:50:00Z", "demo"), "offline");
  assert.equal(presenceState("2026-07-20T19:40:00Z", "demo"), "gone");
});

test("presenceState uses shorter normal-profile windows", (t) => {
  t.mock.timers.enable({ apis: ["Date"], now: new Date("2026-07-20T20:17:05Z") });
  assert.equal(presenceState("2026-07-20T20:16:30Z", "normal"), "live");
  assert.equal(presenceState("2026-07-20T20:15:30Z", "normal"), "stale");
  assert.equal(presenceState("2026-07-20T20:10:30Z", "normal"), "offline");
  assert.equal(presenceState("2026-07-20T20:00:00Z", "normal"), "gone");
});
