import { List, ListOrdered } from "lucide-react";

import RichTextButton from "./RichTextButton";
import { runCommand, isActive } from "./richTextCommands";

const ITEMS = [
  { command: "insertUnorderedList", icon: List, labelKey: "rte.list.bullet" },
  { command: "insertOrderedList", icon: ListOrdered, labelKey: "rte.list.ordered" },
];

function RichTextListControls({ editor, t }) {
  return (
    <div className="rte-group">
      {ITEMS.map((item) => (
        <RichTextButton
          key={item.command}
          active={isActive(item.command)}
          label={t(item.labelKey)}
          onClick={() => runCommand(editor, item.command)}
        >
          <item.icon size={16} strokeWidth={1.9} aria-hidden="true" />
        </RichTextButton>
      ))}
    </div>
  );
}

export default RichTextListControls;
