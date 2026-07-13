import { EditorState } from "@codemirror/state";
import { EditorView, keymap, lineNumbers } from "@codemirror/view";
import { defaultKeymap, history, historyKeymap } from "@codemirror/commands";
import { markdown } from "@codemirror/lang-markdown";
import { remoteCursorPlugin, injectStyles } from "@collab-editor/awareness";
import * as cmView from "@codemirror/view";
import * as cmState from "@codemirror/state";

export function createEditor(parent, awareness, participantID, onLocalUpdate, onSelectionChange) {
  injectStyles();
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
    ...remoteCursorPlugin(cmView, cmState, awareness, participantID),
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
