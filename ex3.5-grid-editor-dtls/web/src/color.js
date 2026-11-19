// Intent: Keep browser color preview formatting in a pure helper so the
// Chrome visibility fix can be tested without depending on the DOM or the
// native color input widget. Source: DI-pafob
export function normalizeBrowserColor(value) {
  const candidate = String(value || "").trim();
  if (/^#[0-9a-fA-F]{6}$/.test(candidate)) {
    return candidate.toUpperCase();
  }
  return "#1D6FD6";
}
