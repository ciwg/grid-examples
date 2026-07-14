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

  // Intent: Keep workflow and review metadata local to the browser while
  // giving the app one coherent registry for recent docs, timestamps,
  // bookmarks, snapshots, comments, and review history instead of scattering
  // those values across ad hoc storage keys. Source: DI-dovoz; DI-nuvif;
  // DI-safor; DI-lapek
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
        comments: [],
        activity: [],
        recentParticipants: [],
        savedVersions: [],
      };
    }
    return this.state.documents[documentID];
  }

  touchViewed(documentID) {
    const document = this.ensureDocument(documentID);
    document.lastViewedAt = nowISO();
    this.pushActivity(document, {
      type: "viewed",
      label: "Viewed document",
      at: document.lastViewedAt,
    });
    this.bumpRecent(documentID);
    this.save();
    return structuredClone(document);
  }

  touchEdited(documentID) {
    const document = this.ensureDocument(documentID);
    document.lastEditedAt = nowISO();
    this.pushActivity(document, {
      type: "edited",
      label: "Edited document",
      at: document.lastEditedAt,
    });
    this.bumpRecent(documentID);
    this.save();
    return structuredClone(document);
  }

  touchExported(documentID) {
    const document = this.ensureDocument(documentID);
    document.lastExportedAt = nowISO();
    this.pushActivity(document, {
      type: "exported",
      label: "Exported document",
      at: document.lastExportedAt,
    });
    this.save();
    return structuredClone(document);
  }

  updateTitle(documentID, title) {
    const document = this.ensureDocument(documentID);
    document.title = title?.trim() || defaultTitle(documentID);
    this.pushActivity(document, {
      type: "retitled",
      label: `Updated title to ${document.title}`,
      at: nowISO(),
    });
    this.save();
    return structuredClone(document);
  }

  setTitle(documentID, title) {
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
    this.pushActivity(document, {
      type: "bookmark",
      label: `Added bookmark ${bookmark.label}`,
      at: bookmark.createdAt || nowISO(),
    });
    this.save();
    return structuredClone(document);
  }

  addSnapshot(documentID, snapshot) {
    const document = this.ensureDocument(documentID);
    document.snapshots = [snapshot, ...document.snapshots].slice(0, 24);
    this.pushActivity(document, {
      type: "snapshot",
      label: `Published snapshot ${snapshot.title || "snapshot"}`,
      at: snapshot.createdAt || nowISO(),
    });
    this.save();
    return structuredClone(document);
  }

  addComment(documentID, comment) {
    const document = this.ensureDocument(documentID);
    document.comments = [comment, ...document.comments.filter((value) => value.id !== comment.id)].slice(0, 64);
    this.pushActivity(document, {
      type: "comment",
      label: `Added comment by ${comment.authorName}`,
      at: comment.createdAt || nowISO(),
    });
    this.save();
    return structuredClone(document);
  }

  updateComment(documentID, commentID, updater) {
    const document = this.ensureDocument(documentID);
    document.comments = document.comments.map((comment) => comment.id === commentID ? updater(structuredClone(comment)) : comment);
    this.save();
    return structuredClone(document);
  }

  addReaction(documentID, commentID, reaction) {
    const document = this.ensureDocument(documentID);
    document.comments = document.comments.map((comment) => {
      if (comment.id !== commentID) {
        return comment;
      }
      return {
        ...comment,
        reactions: [reaction, ...(comment.reactions || [])].slice(0, 16),
      };
    });
    this.pushActivity(document, {
      type: "reaction",
      label: `Reacted ${reaction.emoji} to a comment`,
      at: reaction.createdAt || nowISO(),
    });
    this.save();
    return structuredClone(document);
  }

  toggleCommentResolved(documentID, commentID, resolvedBy) {
    const document = this.ensureDocument(documentID);
    let resolved = false;
    document.comments = document.comments.map((comment) => {
      if (comment.id !== commentID) {
        return comment;
      }
      resolved = !comment.resolved;
      return {
        ...comment,
        resolved,
        resolvedAt: resolved ? nowISO() : "",
        resolvedBy: resolved ? resolvedBy : "",
      };
    });
    this.pushActivity(document, {
      type: resolved ? "resolved" : "reopened",
      label: resolved ? "Resolved comment" : "Reopened comment",
      at: nowISO(),
    });
    this.save();
    return structuredClone(document);
  }

  noteParticipant(documentID, participant) {
    const document = this.ensureDocument(documentID);
    const next = {
      participantID: participant.participantID,
      name: participant.name || participant.participantID,
      color: participant.color || "#999999",
      lastSeenAt: participant.lastSeenAt || nowISO(),
      lastEditedAt: participant.lastEditedAt || "",
      lastViewedAt: participant.lastViewedAt || "",
    };
    document.recentParticipants = [
      next,
      ...document.recentParticipants.filter((value) => value.participantID !== next.participantID),
    ].slice(0, 24);
    this.save();
    return structuredClone(document);
  }

  addSavedVersion(documentID, version) {
    const document = this.ensureDocument(documentID);
    document.savedVersions = [version, ...document.savedVersions.filter((value) => value.id !== version.id)].slice(0, 24);
    this.pushActivity(document, {
      type: "version",
      label: `Saved version ${version.name}`,
      at: version.createdAt || nowISO(),
    });
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
      comments: [],
      activity: [],
      recentParticipants: [],
      savedVersions: [],
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

  listComments(documentID) {
    return structuredClone(this.ensureDocument(documentID).comments || []);
  }

  listActivity(documentID) {
    return structuredClone(this.ensureDocument(documentID).activity || []);
  }

  listRecentParticipants(documentID) {
    return structuredClone(this.ensureDocument(documentID).recentParticipants || []);
  }

  listSavedVersions(documentID) {
    return structuredClone(this.ensureDocument(documentID).savedVersions || []);
  }

  bumpRecent(documentID) {
    this.state.recent = [documentID, ...this.state.recent.filter((value) => value !== documentID)].slice(0, 12);
  }

  pushActivity(document, event) {
    document.activity = [event, ...(document.activity || [])].slice(0, 80);
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
      comments: Array.isArray(value.comments) ? value.comments : [],
      activity: Array.isArray(value.activity) ? value.activity : [],
      recentParticipants: Array.isArray(value.recentParticipants) ? value.recentParticipants : [],
      savedVersions: Array.isArray(value.savedVersions)
        ? value.savedVersions.map((version) => ({
          ...version,
          content: version.content || "",
          replicaBase64: version.replicaBase64 || "",
        }))
        : [],
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
