function createMemoryStorage() {
  const data = new Map();
  return {
    getItem(key) {
      return data.has(key) ? data.get(key) : null;
    },
    setItem(key, value) {
      data.set(key, String(value));
    },
    removeItem(key) {
      data.delete(key);
    },
  };
}

function buildSafeStorage(getStorage) {
  const memory = createMemoryStorage();

  return {
    getItem(key) {
      try {
        const storage = getStorage();
        return storage ? storage.getItem(key) : memory.getItem(key);
      } catch {
        return memory.getItem(key);
      }
    },
    setItem(key, value) {
      try {
        const storage = getStorage();
        if (storage) {
          storage.setItem(key, value);
          return;
        }
      } catch {
        // Intent: Let private/incognito or policy-restricted browsers fall
        // back to in-memory workflow/session state so relay-backed document
        // sync still works even when local storage APIs reject reads or
        // writes. Source: DI-ribaf
      }
      memory.setItem(key, value);
    },
    removeItem(key) {
      try {
        const storage = getStorage();
        if (storage) {
          storage.removeItem(key);
          return;
        }
      } catch {
        // Intent: Keep cleanup non-fatal when browser storage is restricted
        // so temporary local state cannot break the shared-document path.
        // Source: DI-ribaf
      }
      memory.removeItem(key);
    },
  };
}

export const safeLocalStorage = buildSafeStorage(() => window.localStorage);
export const safeSessionStorage = buildSafeStorage(() => window.sessionStorage);
