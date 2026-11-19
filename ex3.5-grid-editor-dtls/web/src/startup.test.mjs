import test from "node:test";
import assert from "node:assert/strict";

import { relayHasAuthoritativeHistory, shouldApplySeed, shouldRecoverRelayHistory } from "./startup.js";

test("relayHasAuthoritativeHistory treats snapshots as authoritative", () => {
  assert.equal(relayHasAuthoritativeHistory({ snapshot_present: true, message_count: 0 }), true);
});

test("relayHasAuthoritativeHistory treats relay message history as authoritative", () => {
  assert.equal(relayHasAuthoritativeHistory({ snapshot_present: false, message_count: 3 }), true);
  assert.equal(relayHasAuthoritativeHistory({ snapshot_present: false, message_count: 0 }), false);
});

test("shouldApplySeed only applies to truly empty brand-new docs", () => {
  assert.equal(shouldApplySeed("# seeded", "", { snapshot_present: false, message_count: 0 }), true);
  assert.equal(shouldApplySeed("# seeded", "relay text", { snapshot_present: false, message_count: 0 }), false);
  assert.equal(shouldApplySeed("# seeded", "", { snapshot_present: true, message_count: 0 }), false);
  assert.equal(shouldApplySeed("# seeded", "", { snapshot_present: false, message_count: 4 }), false);
});

test("shouldRecoverRelayHistory only when relay has history but browser text is empty", () => {
  assert.equal(shouldRecoverRelayHistory("", { snapshot_present: true, message_count: 0 }), true);
  assert.equal(shouldRecoverRelayHistory("", { snapshot_present: false, message_count: 4 }), true);
  assert.equal(shouldRecoverRelayHistory("shared text", { snapshot_present: false, message_count: 4 }), false);
  assert.equal(shouldRecoverRelayHistory("", { snapshot_present: false, message_count: 0 }), false);
});
