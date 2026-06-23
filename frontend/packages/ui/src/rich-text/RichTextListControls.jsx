import { useEditorState } from "@tiptap/react";
import { List, ListOrdered, Quote } from "lucide-react";

import RichTextButton from "./RichTextButton";

const ITEMS = [
  { name: "bulletList", icon: List, labelKey: "rte.list.bullet", command: "toggleBulletList" },
  { name: "orderedList", icon: ListOrdered, labelKey: "rte.list.ordered", command: "toggleOrderedList" },
  { name: "blockquote", icon: Quote, labelKey: "rte.list.blockquote", command: "toggleBlockquote" },
];

function RichTextListControls({ editor, t }) {
  const active = useEditorState({ editor, selector: ({ editor: e }) => readActiveItems(e) });
  return (
    <div className="rte-group">
      {ITEMS.map((item) => (
        <RichTextButton
          key={item.name}
          active={active[item.name]}
          label={t(item.labelKey)}
          onClick={() => editor.chain().focus()[item.command]().run()}
        >
          <item.icon size={16} strokeWidth={1.9} aria-hidden="true" />
        </RichTextButton>
      ))}
    </div>
  );
}

function readActiveItems(editor) {
  if (!editor) return {};
  return ITEMS.reduce((acc, item) => ({ ...acc, [item.name]: editor.isActive(item.name) }), {});
}

export default RichTextListControls;
