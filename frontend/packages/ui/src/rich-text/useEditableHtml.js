import { useCallback, useEffect, useRef } from "react";

// useEditableHtml owns the contenteditable element. The element's innerHTML IS
// the value: pasted HTML is preserved verbatim. We seed innerHTML on mount and
// only re-seed when documentKey changes, so typing never clobbers the caret.
export function useEditableHtml({ initialContent, documentKey, onHtmlChange }) {
  const ref = useRef(null);
  const prevKey = useRef(null);
  // Holds the html to write the next time the editable node mounts (initial
  // mount, documentKey change, or returning from the HTML source view).
  const pending = useRef(initialContent);
  if (prevKey.current !== documentKey) {
    prevKey.current = documentKey;
    pending.current = initialContent;
  }

  const attach = useCallback((node) => {
    ref.current = node;
    if (node && pending.current !== null) {
      node.innerHTML = pending.current || "";
      pending.current = null;
    }
  }, []);

  useEffect(() => {
    if (ref.current && pending.current !== null) {
      ref.current.innerHTML = pending.current || "";
      pending.current = null;
    }
  }, [documentKey]);

  const onInput = useCallback(() => {
    if (ref.current) onHtmlChange(ref.current.innerHTML);
  }, [onHtmlChange]);

  const seed = useCallback((html) => {
    pending.current = html;
    if (ref.current) {
      ref.current.innerHTML = html || "";
      pending.current = null;
    }
    onHtmlChange(html);
  }, [onHtmlChange]);

  return { ref, attach, onInput, seed };
}
