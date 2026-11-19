import test from "node:test";
import assert from "node:assert/strict";

import { safeLocalStorage, safeSessionStorage } from "./safe-storage.js";

test("safeSessionStorage falls back to memory when session storage is unavailable", () => {
  const originalWindow = globalThis.window;
  globalThis.window = {
    sessionStorage: {
      getItem() {
        throw new Error("blocked");
      },
      setItem() {
        throw new Error("blocked");
      },
      removeItem() {
        throw new Error("blocked");
      },
    },
  };
  try {
    safeSessionStorage.setItem("participant", "browser-a");
    assert.equal(safeSessionStorage.getItem("participant"), "browser-a");
    safeSessionStorage.removeItem("participant");
    assert.equal(safeSessionStorage.getItem("participant"), null);
  } finally {
    globalThis.window = originalWindow;
  }
});

test("safeLocalStorage uses browser storage when available", () => {
  const data = new Map();
  const originalWindow = globalThis.window;
  globalThis.window = {
    localStorage: {
      getItem(key) {
        return data.has(key) ? data.get(key) : null;
      },
      setItem(key, value) {
        data.set(key, String(value));
      },
      removeItem(key) {
        data.delete(key);
      },
    },
  };
  try {
    safeLocalStorage.setItem("welcome", "true");
    assert.equal(safeLocalStorage.getItem("welcome"), "true");
    assert.equal(data.get("welcome"), "true");
  } finally {
    globalThis.window = originalWindow;
  }
});
