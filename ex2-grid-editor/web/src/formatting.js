// Keep formatting wrapper behavior in a pure helper so browser toolbar actions
// and Node tests exercise the same selection logic.
export function wrapSelectedText(text, from, to, prefix, suffix) {
  const selected = text.slice(from, to);
  const fallback = selected || "text";
  const insert = `${prefix}${selected || fallback}${suffix}`;
  const nextText = `${text.slice(0, from)}${insert}${text.slice(to)}`;
  return {
    text: nextText,
    selectionFrom: from + prefix.length,
    selectionTo: from + prefix.length + fallback.length,
  };
}
