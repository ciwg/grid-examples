const STORAGE_KEY = "grid-editor-phase2-documents";

export class DocumentRegistry {
  constructor(storage = window.localStorage) {
    this.storage = storage;
    this.state = this.load();
  }

  load() {
    try {
      const raw = this.storage.getItem(STORAGE_KEY);
      if (!raw) {
        return { documents: {}, recent: [], openTabs: [] };
      }
      return normalizeState(JSON.parse(raw));
    } catch {
      return { documents: {}, recent: [], openTabs: [] };
    }
  }

  save() {
    this.storage.setItem(STORAGE_KEY, JSON.stringify(this.state));
  }

  snapshot() {
    return structuredClone(this.state);
  }

  // Intent: Keep Phase 2 workflow metadata local to the browser while giving
  // the app one coherent registry for recent docs, timestamps, bookmarks, and
  // snapshots instead of scattering those values across ad hoc storage keys.
  // Source: DI-dovoz; DI-nuvif
  ensureDocument(documentID) {
    if (!this.state.documents[documentID]) {
      this.state.documents[documentID] = {
        documentID,
        title: defaultTitle(documentID),
        createdAt: nowISO(),
        lastViewedAt: nowISO(),
        lastEditedAt: "",
        lastExportedAt: "",
        bookmarks: [],
        snapshots: [],
      };
    }
    return this.state.documents[documentID];
  }

  touchViewed(documentID) {
    const document = this.ensureDocument(documentID);
    document.lastViewedAt = nowISO();
    this.bumpRecent(documentID);
    this.save();
    return structuredClone(document);
  }

  touchEdited(documentID) {
    const document = this.ensureDocument(documentID);
    document.lastEditedAt = nowISO();
    this.bumpRecent(documentID);
    this.save();
    return structuredClone(document);
  }

  touchExported(documentID) {
    const document = this.ensureDocument(documentID);
    document.lastExportedAt = nowISO();
    this.save();
    return structuredClone(document);
  }

  updateTitle(documentID, title) {
    const document = this.ensureDocument(documentID);
    document.title = title?.trim() || defaultTitle(documentID);
    this.save();
    return structuredClone(document);
  }

  openTab(documentID) {
    this.ensureDocument(documentID);
    this.state.openTabs = [documentID, ...this.state.openTabs.filter((value) => value !== documentID)].slice(0, 8);
    this.bumpRecent(documentID);
    this.save();
  }

  closeTab(documentID) {
    this.state.openTabs = this.state.openTabs.filter((value) => value !== documentID);
    this.save();
  }

  addBookmark(documentID, bookmark) {
    const document = this.ensureDocument(documentID);
    document.bookmarks = [bookmark, ...document.bookmarks.filter((value) => value.id !== bookmark.id)].slice(0, 24);
    this.save();
    return structuredClone(document);
  }

  addSnapshot(documentID, snapshot) {
    const document = this.ensureDocument(documentID);
    document.snapshots = [snapshot, ...document.snapshots].slice(0, 24);
    this.save();
    return structuredClone(document);
  }

  duplicateDocument(documentID, nextDocumentID, content) {
    const source = this.ensureDocument(documentID);
    this.state.documents[nextDocumentID] = {
      ...structuredClone(source),
      documentID: nextDocumentID,
      title: `${source.title} copy`,
      createdAt: nowISO(),
      lastViewedAt: nowISO(),
      lastEditedAt: nowISO(),
      lastExportedAt: "",
      snapshots: [],
      bookmarks: [],
      seedContent: content,
    };
    this.openTab(nextDocumentID);
    this.save();
    return structuredClone(this.state.documents[nextDocumentID]);
  }

  registerSeedContent(documentID, content) {
    const document = this.ensureDocument(documentID);
    document.seedContent = content;
    this.save();
  }

  seedContent(documentID) {
    return this.ensureDocument(documentID).seedContent || "";
  }

  listRecent() {
    return this.state.recent.map((documentID) => structuredClone(this.ensureDocument(documentID)));
  }

  listOpenTabs() {
    return this.state.openTabs.map((documentID) => structuredClone(this.ensureDocument(documentID)));
  }

  get(documentID) {
    return structuredClone(this.ensureDocument(documentID));
  }

  bumpRecent(documentID) {
    this.state.recent = [documentID, ...this.state.recent.filter((value) => value !== documentID)].slice(0, 12);
  }
}

export function templateCatalog() {
  return [
    {
      id: "blank",
      label: "Blank Doc",
      content: "",
    },
    {
      id: "notes",
      label: "Meeting Notes",
      content: "# Meeting notes\n\n## Agenda\n\n- \n\n## Notes\n\n- \n\n## Follow-up\n\n- ",
    },
    {
      id: "checklist",
      label: "Checklist",
      content: "# Checklist\n\n- [ ] Item one\n- [ ] Item two\n- [ ] Item three",
    },
    {
      id: "demo",
      label: "Demo Sample",
      content: "# grid-editor demo\n\nThis is a shared sample document.\n\n## Try this\n\n- type from another browser\n- move your cursor\n- toggle preview\n- export a snapshot",
    },
  ];
}

export function generateDemoText() {
  return [
    "# Generated demo document",
    "",
    "## Summary",
    "",
    "This is a generated test document for Phase 2 workflow surfaces.",
    "",
    "## Checklist",
    "",
    "- [ ] review layout",
    "- [ ] open preview",
    "- [ ] export markdown",
    "- [ ] add bookmark",
  ].join("\n");
}

export function normalizeState(raw) {
  const state = raw && typeof raw === "object" ? raw : {};
  const documents = {};
  for (const [documentID, value] of Object.entries(state.documents || {})) {
    documents[documentID] = {
      documentID,
      title: value.title || defaultTitle(documentID),
      createdAt: value.createdAt || nowISO(),
      lastViewedAt: value.lastViewedAt || "",
      lastEditedAt: value.lastEditedAt || "",
      lastExportedAt: value.lastExportedAt || "",
      bookmarks: Array.isArray(value.bookmarks) ? value.bookmarks : [],
      snapshots: Array.isArray(value.snapshots) ? value.snapshots : [],
      seedContent: value.seedContent || "",
    };
  }
  return {
    documents,
    recent: Array.isArray(state.recent) ? state.recent.filter((value) => documents[value]) : [],
    openTabs: Array.isArray(state.openTabs) ? state.openTabs.filter((value) => documents[value]) : [],
  };
}

function defaultTitle(documentID) {
  return `Document ${documentID}`;
}

function nowISO() {
  return new Date().toISOString();
}
