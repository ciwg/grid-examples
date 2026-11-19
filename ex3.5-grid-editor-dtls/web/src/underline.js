// Intent: Keep underline parsing in a pure helper so the browser can render
// `<u>...</u>` visually in CodeMirror while preserving the exact saved bytes
// and export behavior. Source: DI-naruv
export function findUnderlineRanges(text) {
  const source = String(text || "");
  const ranges = [];
  const stack = [];
  let index = 0;
  while (index < source.length) {
    if (source.startsWith("<u>", index)) {
      stack.push(index);
      index += 3;
      continue;
    }
    if (source.startsWith("</u>", index)) {
      const openFrom = stack.pop();
      if (openFrom != null) {
        const openTo = openFrom + 3;
        const closeFrom = index;
        const closeTo = index + 4;
        if (openTo <= closeFrom) {
          ranges.push({
            openFrom,
            openTo,
            contentFrom: openTo,
            contentTo: closeFrom,
            closeFrom,
            closeTo,
          });
        }
      }
      index += 4;
      continue;
    }
    index += 1;
  }
  return ranges;
}
