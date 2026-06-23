// Thin wrapper over document.execCommand for the contenteditable surface. Every
// command focuses the editable element first, runs the command, then fires the
// onInput sync so the React value tracks the new innerHTML.
export function runCommand(editor, command, value) {
  const el = editor?.ref?.current;
  if (!el) return;
  el.focus();
  // Emit inline styles (style="color:…") instead of presentational tags like
  // <font>, which the server sanitizer strips. Keeps color/bold round-tripping.
  document.execCommand("styleWithCSS", false, true);
  document.execCommand(command, false, value);
  editor.onInput();
}

export function setBlock(editor, tag) {
  runCommand(editor, "formatBlock", tag);
}

export function insertImageUrl(editor, url) {
  if (url) runCommand(editor, "insertImage", url);
}

export function createLink(editor, url) {
  if (url) runCommand(editor, "createLink", url);
}

export function isActive(command) {
  try {
    return document.queryCommandState(command);
  } catch {
    return false;
  }
}
