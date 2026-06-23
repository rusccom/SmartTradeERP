import { Bold, Italic, Underline, Strikethrough } from "lucide-react";

import RichTextButton from "./RichTextButton";
import { runCommand, isActive } from "./richTextCommands";

const MARKS = [
  { command: "bold", icon: Bold, labelKey: "rte.mark.bold" },
  { command: "italic", icon: Italic, labelKey: "rte.mark.italic" },
  { command: "underline", icon: Underline, labelKey: "rte.mark.underline" },
  { command: "strikeThrough", icon: Strikethrough, labelKey: "rte.mark.strike" },
];

function RichTextMarkControls({ editor, t }) {
  return (
    <div className="rte-group">
      {MARKS.map((mark) => (
        <RichTextButton
          key={mark.command}
          active={isActive(mark.command)}
          label={t(mark.labelKey)}
          onClick={() => runCommand(editor, mark.command)}
        >
          <mark.icon size={16} strokeWidth={1.9} aria-hidden="true" />
        </RichTextButton>
      ))}
    </div>
  );
}

export default RichTextMarkControls;
