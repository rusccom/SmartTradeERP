import { useEditorState } from "@tiptap/react";
import { AlignLeft, AlignCenter, AlignRight, AlignJustify } from "lucide-react";

import RichTextButton from "./RichTextButton";

const ALIGNS = [
  { value: "left", icon: AlignLeft, labelKey: "rte.align.left" },
  { value: "center", icon: AlignCenter, labelKey: "rte.align.center" },
  { value: "right", icon: AlignRight, labelKey: "rte.align.right" },
  { value: "justify", icon: AlignJustify, labelKey: "rte.align.justify" },
];

function RichTextAlignControls({ editor, t }) {
  const active = useEditorState({ editor, selector: ({ editor: e }) => readActiveAligns(e) });
  return (
    <div className="rte-group">
      {ALIGNS.map((align) => (
        <RichTextButton
          key={align.value}
          active={active[align.value]}
          label={t(align.labelKey)}
          onClick={() => editor.chain().focus().setTextAlign(align.value).run()}
        >
          <align.icon size={16} strokeWidth={1.9} aria-hidden="true" />
        </RichTextButton>
      ))}
    </div>
  );
}

function readActiveAligns(editor) {
  if (!editor) return {};
  return ALIGNS.reduce((acc, align) => ({ ...acc, [align.value]: editor.isActive({ textAlign: align.value }) }), {});
}

export default RichTextAlignControls;
