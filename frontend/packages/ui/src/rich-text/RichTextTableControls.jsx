import { useEditorState } from "@tiptap/react";

import RichTextButton from "./RichTextButton";

const ACTIONS = [
  { command: "addColumnBefore", labelKey: "rte.table.addColBefore", text: "+Col←" },
  { command: "addColumnAfter", labelKey: "rte.table.addColAfter", text: "+Col→" },
  { command: "deleteColumn", labelKey: "rte.table.delCol", text: "-Col" },
  { command: "addRowBefore", labelKey: "rte.table.addRowBefore", text: "+Row↑" },
  { command: "addRowAfter", labelKey: "rte.table.addRowAfter", text: "+Row↓" },
  { command: "deleteRow", labelKey: "rte.table.delRow", text: "-Row" },
  { command: "deleteTable", labelKey: "rte.table.delete", text: "✕" },
];

function RichTextTableControls({ editor, t }) {
  const inTable = useEditorState({ editor, selector: ({ editor: e }) => Boolean(e) && e.isActive("table") });
  if (!inTable) return null;
  return (
    <div className="rte-toolbar rte-table-bar">
      {ACTIONS.map((action) => (
        <RichTextButton
          key={action.command}
          label={t(action.labelKey)}
          onClick={() => editor.chain().focus()[action.command]().run()}
        >
          <span className="rte-btn-text">{action.text}</span>
        </RichTextButton>
      ))}
    </div>
  );
}

export default RichTextTableControls;
