import { useState } from "react";
import { Baseline } from "lucide-react";

import RichTextButton from "./RichTextButton";

const PALETTE = [
  "#1a1a1a", "#5c5f62", "#b91c1c", "#c2410c",
  "#a16207", "#15803d", "#0e7490", "#2563eb",
  "#6d28d9", "#be185d",
];

function RichTextColorControl({ editor, t }) {
  const [open, setOpen] = useState(false);
  return (
    <div className="rte-color">
      <RichTextButton label={t("rte.color.label")} active={open} onClick={() => setOpen((value) => !value)}>
        <Baseline size={16} strokeWidth={1.9} aria-hidden="true" />
      </RichTextButton>
      {open && renderMenu(editor, t, () => setOpen(false))}
    </div>
  );
}

function renderMenu(editor, t, close) {
  return (
    <div className="rte-color-menu" role="menu">
      <div className="rte-color-grid">
        {PALETTE.map((color) => renderSwatch(editor, color, t, close))}
      </div>
      <button type="button" className="rte-color-clear" onMouseDown={(event) => event.preventDefault()} onClick={() => applyClear(editor, close)}>
        {t("rte.color.clear")}
      </button>
    </div>
  );
}

function renderSwatch(editor, color, t, close) {
  return (
    <button
      key={color}
      type="button"
      className="rte-color-swatch"
      style={{ background: color }}
      title={color}
      aria-label={t("rte.color.swatch")}
      onMouseDown={(event) => event.preventDefault()}
      onClick={() => applyColor(editor, color, close)}
    />
  );
}

function applyColor(editor, color, close) {
  editor.chain().focus().setColor(color).run();
  close();
}

function applyClear(editor, close) {
  editor.chain().focus().unsetColor().run();
  close();
}

export default RichTextColorControl;
