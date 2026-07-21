import { safeLocalStorage } from "./safe-storage.js";

const STORAGE_KEY = "grid-editor-phase1-preferences";

const DEFAULT_SHORTCUTS = {
  search: "Mod-F",
  bold: "Mod-B",
  italic: "Mod-I",
  underline: "Mod-Shift-U",
  settings: "Mod-,",
  help: "F1",
};

const DEFAULTS = {
  displayName: "Browser User",
  color: "#1d6fd6",
  theme: "paper",
  lineNumbers: true,
  fontSize: 16,
  dyslexiaMode: false,
  profile: "demo",
  shortcuts: DEFAULT_SHORTCUTS,
};

export class PreferencesStore {
  constructor(storage = safeLocalStorage) {
    this.storage = storage;
    this.value = this.load();
  }

  load() {
    try {
      const raw = this.storage.getItem(STORAGE_KEY);
      if (!raw) {
        return structuredClone(DEFAULTS);
      }
      return normalizePreferences(JSON.parse(raw));
    } catch {
      return structuredClone(DEFAULTS);
    }
  }

  get() {
    return structuredClone(this.value);
  }

  update(patch) {
    this.value = normalizePreferences({
      ...this.value,
      ...patch,
      shortcuts: {
        ...this.value.shortcuts,
        ...(patch.shortcuts || {}),
      },
    });
    this.storage.setItem(STORAGE_KEY, JSON.stringify(this.value));
    return this.get();
  }
}

export function defaultShortcuts() {
  return { ...DEFAULT_SHORTCUTS };
}

export function formatShortcut(shortcut) {
  return (shortcut || "")
    .replaceAll("Mod", navigator.platform.includes("Mac") ? "Cmd" : "Ctrl")
    .replaceAll("-", "+");
}

function normalizePreferences(value) {
  const next = {
    ...DEFAULTS,
    ...(value || {}),
    shortcuts: {
      ...DEFAULT_SHORTCUTS,
      ...((value && value.shortcuts) || {}),
    },
  };
  next.displayName = String(next.displayName || DEFAULTS.displayName).slice(0, 80);
  next.color = /^#[0-9a-fA-F]{6}$/.test(next.color) ? next.color : DEFAULTS.color;
  next.theme = ["paper", "night"].includes(next.theme) ? next.theme : DEFAULTS.theme;
  next.lineNumbers = Boolean(next.lineNumbers);
  next.dyslexiaMode = Boolean(next.dyslexiaMode);
  next.profile = ["demo", "normal"].includes(next.profile) ? next.profile : DEFAULTS.profile;
  next.fontSize = Math.max(14, Math.min(24, Number(next.fontSize) || DEFAULTS.fontSize));
  return next;
}
