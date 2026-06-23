import { AlignLeft, AlignCenter, AlignRight } from "lucide-react";

import RichTextButton from "./RichTextButton";
import { runCommand, isActive } from "./richTextCommands";

const ALIGNS = [
  { command: "justifyLeft", icon: AlignLeft, labelKey: "rte.align.left" },
  { command: "justifyCenter", icon: AlignCenter, labelKey: "rte.align.center" },
  { command: "justifyRight", icon: AlignRight, labelKey: "rte.align.right" },
];

function RichTextAlignControls({ editor, t }) {
  return (
    <div className="rte-group">
      {ALIGNS.map((align) => (
        <RichTextButton
          key={align.command}
          active={isActive(align.command)}
          label={t(align.labelKey)}
          onClick={() => runCommand(editor, align.command)}
        >
          <align.icon size={16} strokeWidth={1.9} aria-hidden="true" />
        </RichTextButton>
      ))}
    </div>
  );
}

export default RichTextAlignControls;
