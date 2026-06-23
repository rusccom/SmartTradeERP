import { useState } from "react";
import { Baseline } from "lucide-react";

import RichTextButton from "./RichTextButton";
import { runCommand } from "./richTextCommands";

const PALETTE = [
  "#1a1a1a", "#5c5f62", "#b91c1c", "#c2410c",
  "#a16207", "#15803d", "#0e7490", "#2563eb",
  "#6d28d9", "#be185d",
];

function RichTextColorControl({ editor, t }) {
  const [open, setOpen] = useState(false);
  const apply = (color) => { runCommand(editor, "foreColor", color); setOpen(false); };
  return (
    <div className="rte-color">
      <RichTextButton label={t("rte.color.label")} active={open} onClick={() => setOpen((value) => !value)}>
        <Baseline size={16} strokeWidth={1.9} aria-hidden="true" />
      </RichTextButton>
      {open && renderMenu(apply, t)}
    </div>
  );
}

function renderMenu(apply, t) {
  return (
    <div className="rte-color-menu" role="menu">
      <div className="rte-color-grid">
        {PALETTE.map((color) => renderSwatch(color, apply, t))}
      </div>
    </div>
  );
}

function renderSwatch(color, apply, t) {
  return (
    <button
      key={color}
      type="button"
      className="rte-color-swatch"
      style={{ background: color }}
      title={color}
      aria-label={t("rte.color.swatch")}
      onMouseDown={(event) => event.preventDefault()}
      onClick={() => apply(color)}
    />
  );
}

export default RichTextColorControl;
