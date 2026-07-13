import { Compartment, EditorSelection, EditorState } from "@codemirror/state";
import { EditorView, keymap, lineNumbers } from "@codemirror/view";
import { defaultKeymap, history, historyKeymap } from "@codemirror/commands";
import { markdown } from "@codemirror/lang-markdown";
import * as cmView from "@codemirror/view";
import * as cmState from "@codemirror/state";
import { wrapSelectedText } from "./formatting.js";

export function createEditor(parent, awareness, participantID, onLocalUpdate, onSelectionChange) {
  injectCursorStyles();
  let applyingRemote = false;
  const lineNumbersCompartment = new Compartment();

  // Intent: Move the browser embodiment to a real CodeMirror surface so CRDT
  // text convergence and remote cursor rendering share the same editor model.
  // Source: DI-zegov; DI-favok; DI-vasul
  const extensions = [
    lineNumbersCompartment.of(lineNumbers()),
    history(),
    keymap.of([...defaultKeymap, ...historyKeymap]),
    markdown(),
    EditorView.lineWrapping,
    EditorView.updateListener.of((update) => {
      if (update.docChanged && !applyingRemote) {
        onLocalUpdate(update);
      }
      if (update.selectionSet) {
        const selection = update.state.selection.main;
        onSelectionChange(selection.anchor, selection.head);
      }
    }),
    ...createRemoteCursorExtensions(cmView, cmState, awareness, participantID),
  ];

  const state = EditorState.create({
    doc: "",
    extensions,
  });

  const view = new EditorView({
    state,
    parent,
  });

  return {
    view,
    setText(text) {
      const current = view.state.doc.toString();
      if (current === text) {
        return;
      }
      const selection = view.state.selection.main;
      applyingRemote = true;
      view.dispatch({
        changes: { from: 0, to: current.length, insert: text },
        selection: {
          anchor: Math.min(selection.anchor, text.length),
          head: Math.min(selection.head, text.length),
        },
      });
      applyingRemote = false;
    },
    focus() {
      view.focus();
    },
    getText() {
      return view.state.doc.toString();
    },
    getCursorLine() {
      return view.state.doc.lineAt(view.state.selection.main.head).number;
    },
    setLineNumbers(enabled) {
      view.dispatch({
        effects: lineNumbersCompartment.reconfigure(enabled ? lineNumbers() : []),
      });
    },
    findNext(query, options = {}) {
      if (!query) {
        return false;
      }
      const doc = view.state.doc.toString();
      const current = view.state.selection.main;
      const match = findMatch(doc, query, current.to, options) || findMatch(doc, query, 0, options);
      if (!match) {
        return false;
      }
      view.dispatch({
        selection: EditorSelection.range(match.from, match.to),
        scrollIntoView: true,
      });
      view.focus();
      return true;
    },
    replaceAll(query, replacement, options = {}) {
      if (!query) {
        return 0;
      }
      const source = view.state.doc.toString();
      const next = replaceMatches(source, query, replacement, options);
      if (next.count === 0) {
        return 0;
      }
      view.dispatch({
        changes: { from: 0, to: source.length, insert: next.text },
      });
      view.focus();
      return next.count;
    },
    goToLine(lineNumber) {
      const total = view.state.doc.lines;
      const line = view.state.doc.line(Math.max(1, Math.min(total, lineNumber)));
      view.dispatch({
        selection: EditorSelection.cursor(line.from),
        scrollIntoView: true,
      });
      view.focus();
    },
    insertAtCursor(text) {
      const selection = view.state.selection.main;
      view.dispatch({
        changes: { from: selection.from, to: selection.to, insert: text },
        selection: EditorSelection.cursor(selection.from + text.length),
        scrollIntoView: true,
      });
      view.focus();
    },
    wrapSelection(prefix, suffix) {
      const selection = view.state.selection.main;
      const next = wrapSelectedText(view.state.doc.toString(), selection.from, selection.to, prefix, suffix);
      view.dispatch({
        changes: { from: selection.from, to: selection.to, insert: next.text.slice(selection.from, next.text.length - (view.state.doc.length - selection.to)) },
        selection: EditorSelection.range(next.selectionFrom, next.selectionTo),
        scrollIntoView: true,
      });
      view.focus();
    },
    destroy() {
      view.destroy();
    },
  };
}

function createRemoteCursorExtensions(cmView, cmState, awareness, clientID) {
  const { Decoration, ViewPlugin, EditorView, WidgetType } = cmView;
  const { StateEffect, StateField } = cmState;
  const setRemoteCursors = StateEffect.define();

  class CursorWidget extends WidgetType {
    constructor(name, color, clientID, typing) {
      super();
      this.name = name || "User";
      this.color = normalizeColor(color);
      this.clientID = clientID || "";
      this.typing = Boolean(typing);
    }

    toDOM() {
      const wrapper = document.createElement("span");
      wrapper.className = "grid-remote-cursor";
      if (this.typing) {
        wrapper.classList.add("typing");
      }
      wrapper.style.setProperty("--grid-peer-color", this.color);

      const label = document.createElement("span");
      label.className = "grid-remote-cursor-label";
      label.textContent = this.name;
      label.style.backgroundColor = this.color;
      wrapper.appendChild(label);
      return wrapper;
    }

    eq(other) {
      return this.clientID === other.clientID && this.name === other.name && this.color === other.color && this.typing === other.typing;
    }
  }

  const remoteCursorField = StateField.define({
    create() {
      return Decoration.none;
    },
    update(decorations, transaction) {
      for (const effect of transaction.effects) {
        if (effect.is(setRemoteCursors)) {
          return effect.value;
        }
      }
      return decorations.map(transaction.changes);
    },
    provide: (field) => EditorView.decorations.from(field),
  });

  const plugin = ViewPlugin.fromClass(class {
    constructor(view) {
      this.view = view;
      this.updateDecorations = this.updateDecorations.bind(this);
      awareness.on("change", this.updateDecorations);
      this.updateDecorations();
    }

    update(update) {
      if (update.docChanged || update.selectionSet) {
        this.updateDecorations();
      }
    }

    updateDecorations() {
      const docLength = this.view.state.doc.length;
      const decorations = [];
      const states = awareness.getStates();
      states.forEach((state, id) => {
        if (id === clientID) {
          return;
        }
        const user = state.user;
        const selection = state.selection;
        if (!user || !selection || typeof selection.anchor !== "number") {
          return;
        }
        const color = normalizeColor(user.color);
        const anchor = Math.max(0, Math.min(selection.anchor, docLength));
        decorations.push(Decoration.widget({
          widget: new CursorWidget(user.name, color, id, state.typing),
          side: -1,
        }).range(anchor));

        if (typeof selection.head === "number" && selection.head !== selection.anchor) {
          const head = Math.max(0, Math.min(selection.head, docLength));
          const from = Math.min(anchor, head);
          const to = Math.max(anchor, head);
          if (from < to) {
            decorations.push(Decoration.mark({
              class: "grid-remote-selection",
              attributes: { style: `background: ${color}33` },
            }).range(from, to));
          }
        }
      });
      this.view.dispatch({ effects: setRemoteCursors.of(Decoration.set(decorations, true)) });
    }

    destroy() {
      awareness.off("change", this.updateDecorations);
    }
  });

  return [remoteCursorField, plugin];
}

function injectCursorStyles() {
  if (document.getElementById("grid-editor-remote-cursor-styles")) {
    return;
  }
  const style = document.createElement("style");
  style.id = "grid-editor-remote-cursor-styles";
  style.textContent = `
    .grid-remote-cursor {
      position: relative;
      display: inline-block;
      width: 0;
      height: 1.25em;
      margin-left: -1px;
      margin-right: -1px;
      border-left: 2px solid var(--grid-peer-color, #999999);
      pointer-events: none;
      vertical-align: text-bottom;
    }

    .grid-remote-cursor.typing {
      animation: grid-editor-typing-caret 0.8s steps(1) infinite;
    }

    .grid-remote-cursor-label {
      position: absolute;
      top: -1.5em;
      left: -1px;
      font-size: 10px;
      font-family: sans-serif;
      font-weight: 600;
      line-height: 1.2;
      white-space: nowrap;
      color: white;
      padding: 1px 4px;
      border-radius: 3px 3px 3px 0;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.18);
    }

    .grid-remote-selection {
      border-radius: 2px;
    }

    @keyframes grid-editor-typing-caret {
      0% { opacity: 1; }
      50% { opacity: 0.2; }
      100% { opacity: 1; }
    }
  `;
  document.head.appendChild(style);
}

function normalizeColor(value) {
  return typeof value === "string" && /^#[0-9a-fA-F]{6}$/.test(value) ? value : "#999999";
}

function replaceMatches(source, query, replacement, options) {
  if (options.regex) {
    const flags = options.caseSensitive ? "g" : "gi";
    const expression = new RegExp(query, flags);
    let count = 0;
    return {
      text: source.replace(expression, (...args) => {
        count += 1;
        return typeof replacement === "function" ? replacement(...args) : replacement;
      }),
      count,
    };
  }
  const haystack = options.caseSensitive ? source : source.toLowerCase();
  const needle = options.caseSensitive ? query : query.toLowerCase();
  let count = 0;
  let cursor = 0;
  let output = "";
  while (cursor < source.length) {
    const index = haystack.indexOf(needle, cursor);
    if (index === -1) {
      output += source.slice(cursor);
      break;
    }
    output += source.slice(cursor, index);
    output += replacement;
    cursor = index + query.length;
    count += 1;
  }
  return { text: output, count };
}

function findMatch(source, query, start, options) {
  if (options.regex) {
    const flags = options.caseSensitive ? "g" : "gi";
    const expression = new RegExp(query, flags);
    expression.lastIndex = start;
    const match = expression.exec(source);
    if (!match) {
      return null;
    }
    return { from: match.index, to: match.index + match[0].length };
  }
  const haystack = options.caseSensitive ? source : source.toLowerCase();
  const needle = options.caseSensitive ? query : query.toLowerCase();
  const index = haystack.indexOf(needle, start);
  if (index === -1) {
    return null;
  }
  return { from: index, to: index + query.length };
}
