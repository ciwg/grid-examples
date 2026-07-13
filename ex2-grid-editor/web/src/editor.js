import { EditorState } from "@codemirror/state";
import { EditorView, keymap, lineNumbers } from "@codemirror/view";
import { defaultKeymap, history, historyKeymap } from "@codemirror/commands";
import { markdown } from "@codemirror/lang-markdown";
import * as cmView from "@codemirror/view";
import * as cmState from "@codemirror/state";

export function createEditor(parent, awareness, participantID, onLocalUpdate, onSelectionChange) {
  injectCursorStyles();
  let applyingRemote = false;

  // Intent: Move the browser embodiment to a real CodeMirror surface so CRDT
  // text convergence and remote cursor rendering share the same editor model.
  // Source: DI-zegov
  const extensions = [
    lineNumbers(),
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
    constructor(name, color, clientID) {
      super();
      this.name = name || "User";
      this.color = normalizeColor(color);
      this.clientID = clientID || "";
    }

    toDOM() {
      const wrapper = document.createElement("span");
      wrapper.className = "grid-remote-cursor";
      wrapper.style.setProperty("--grid-peer-color", this.color);

      const label = document.createElement("span");
      label.className = "grid-remote-cursor-label";
      label.textContent = this.name;
      label.style.backgroundColor = this.color;
      wrapper.appendChild(label);
      return wrapper;
    }

    eq(other) {
      return this.clientID === other.clientID && this.name === other.name && this.color === other.color;
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
          widget: new CursorWidget(user.name, color, id),
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
  `;
  document.head.appendChild(style);
}

function normalizeColor(value) {
  return typeof value === "string" && /^#[0-9a-fA-F]{6}$/.test(value) ? value : "#999999";
}
