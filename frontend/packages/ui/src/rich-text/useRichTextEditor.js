import { useEffect, useRef } from "react";
import { useEditor } from "@tiptap/react";

import { buildRichTextExtensions } from "./richTextExtensions";

export function useRichTextEditor({ initialContent, documentKey, onHtmlChange }) {
  const prevKey = useRef(documentKey);
  const editor = useEditor({
    extensions: buildRichTextExtensions(),
    content: initialContent || "",
    immediatelyRender: false,
    onUpdate: ({ editor: instance }) => onHtmlChange(instance.getHTML()),
  });
  useRehydrate({ editor, documentKey, initialContent, prevKey });
  return editor;
}

function useRehydrate({ editor, documentKey, initialContent, prevKey }) {
  useEffect(() => {
    if (!editor || prevKey.current === documentKey) return;
    prevKey.current = documentKey;
    editor.commands.setContent(initialContent || "", { emitUpdate: false });
  }, [editor, documentKey, initialContent]);
}
