// Intent: Keep preview and split toolbar semantics in a pure helper so the
// browser UI can guarantee distinct visible outcomes for both buttons.
// Source: DI-zosuf
export function nextPaneState(current, action) {
  const state = {
    previewEnabled: Boolean(current?.previewEnabled),
    splitEnabled: Boolean(current?.splitEnabled),
  };
  if (action === "preview") {
    return {
      previewEnabled: !state.previewEnabled || state.splitEnabled,
      splitEnabled: false,
    };
  }
  if (action === "split") {
    return {
      previewEnabled: true,
      splitEnabled: !state.splitEnabled,
    };
  }
  return state;
}

// Intent: Keep preview-pane visibility rules explicit so preview-only mode
// actually hides the editor pane instead of relying on layout side effects.
// Source: DI-zosuf
export function describePaneMode(state) {
  const previewEnabled = Boolean(state?.previewEnabled);
  const splitEnabled = previewEnabled && Boolean(state?.splitEnabled);
  return {
    showEditor: !previewEnabled || splitEnabled,
    showPreview: previewEnabled,
    splitEnabled,
  };
}
