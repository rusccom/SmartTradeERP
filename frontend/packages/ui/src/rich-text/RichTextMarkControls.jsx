import { useEditorState } from "@tiptap/react";
import { Bold, Italic, Underline, Strikethrough, Code } from "lucide-react";

import RichTextButton from "./RichTextButton";

const MARKS = [
  { name: "bold", icon: Bold, labelKey: "rte.mark.bold" },
  { name: "italic", icon: Italic, labelKey: "rte.mark.italic" },
  { name: "underline", icon: Underline, labelKey: "rte.mark.underline" },
  { name: "strike", icon: Strikethrough, labelKey: "rte.mark.strike" },
  { name: "code", icon: Code, labelKey: "rte.mark.code" },
];

function RichTextMarkControls({ editor, t }) {
  const active = useEditorState({ editor, selector: ({ editor: e }) => readActiveMarks(e) });
  return (
    <div className="rte-group">
      {MARKS.map((mark) => (
        <RichTextButton
          key={mark.name}
          active={active[mark.name]}
          label={t(mark.labelKey)}
          onClick={() => toggleMark(editor, mark.name)}
        >
          <mark.icon size={16} strokeWidth={1.9} aria-hidden="true" />
        </RichTextButton>
      ))}
    </div>
  );
}

function readActiveMarks(editor) {
  if (!editor) return {};
  return MARKS.reduce((acc, mark) => ({ ...acc, [mark.name]: editor.isActive(mark.name) }), {});
}

function toggleMark(editor, name) {
  editor.chain().focus().toggleMark(name).run();
}

export default RichTextMarkControls;
