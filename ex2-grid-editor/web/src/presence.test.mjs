import assert from "node:assert/strict";
import test from "node:test";

import { cancelPresenceRefresh, presenceState, schedulePresenceRefresh } from "./presence.js";

function fakeScheduler(now) {
  const scheduled = [];
  const cleared = [];
  return {
    now: () => now,
    setTimeout: (callback, delay) => {
      const timer = { callback, delay };
      scheduled.push(timer);
      return timer;
    },
    clearTimeout: (timer) => {
      cleared.push(timer);
    },
    scheduled,
    cleared,
  };
}

test("presenceState follows the normal live, stale, offline, and removal boundaries", () => {
  const observedAt = "2026-08-14T12:00:00Z";
  const start = Date.parse(observedAt);

  assert.equal(presenceState(observedAt, "normal", start + 60_000), "live");
  assert.equal(presenceState(observedAt, "normal", start + 60_001), "stale");
  assert.equal(presenceState(observedAt, "normal", start + 300_001), "offline");
  assert.equal(presenceState(observedAt, "normal", start + 900_001), "gone");
});

test("schedulePresenceRefresh wakes at the next local lifecycle boundary", () => {
  const observedAt = "2026-08-14T12:00:00Z";
  const scheduler = fakeScheduler(Date.parse(observedAt));
  const onRefresh = () => {};

  const timer = schedulePresenceRefresh([{ lastSeenAt: observedAt }], "normal", onRefresh, scheduler);

  assert.equal(timer, scheduler.scheduled[0]);
  assert.equal(scheduler.scheduled[0].callback, onRefresh);
  assert.equal(scheduler.scheduled[0].delay, 60_001);
});

test("schedulePresenceRefresh selects the nearest peer boundary and skips unusable observations", () => {
  const observedAt = "2026-08-14T12:00:00Z";
  const scheduler = fakeScheduler(Date.parse(observedAt) + 60_001);

  schedulePresenceRefresh([
    { lastSeenAt: "not-a-time" },
    { lastSeenAt: observedAt },
    { lastSeenAt: "2026-08-14T12:00:30Z" },
  ], "normal", () => {}, scheduler);

  assert.equal(scheduler.scheduled[0].delay, 30_000);
});

test("cancelPresenceRefresh clears the local timer and absent peers schedule nothing", () => {
  const scheduler = fakeScheduler(Date.parse("2026-08-14T12:00:00Z"));

  assert.equal(schedulePresenceRefresh([], "normal", () => {}, scheduler), null);
  const timer = { id: "presence" };
  cancelPresenceRefresh(timer, scheduler);

  assert.deepEqual(scheduler.cleared, [timer]);
});
