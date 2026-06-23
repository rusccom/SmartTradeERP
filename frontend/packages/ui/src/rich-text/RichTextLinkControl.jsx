import { useState } from "react";
import { Link2 } from "lucide-react";

import RichTextButton from "./RichTextButton";

function RichTextLinkControl({ editor, t }) {
  const [open, setOpen] = useState(false);
  const [href, setHref] = useState("");
  const toggle = () => openPopover({ editor, setHref, setOpen, open });
  return (
    <div className="rte-link">
      <RichTextButton label={t("rte.link.label")} active={editor.isActive("link") || open} onClick={toggle}>
        <Link2 size={16} strokeWidth={1.9} aria-hidden="true" />
      </RichTextButton>
      {open && renderPopover({ editor, href, setHref, setOpen, t })}
    </div>
  );
}

function openPopover({ editor, setHref, setOpen, open }) {
  setHref(editor.getAttributes("link").href || "");
  setOpen(!open);
}

function renderPopover({ editor, href, setHref, setOpen, t }) {
  return (
    <div className="rte-link-popover" role="dialog">
      <input
        className="rte-link-input"
        type="url"
        autoFocus
        placeholder={t("rte.link.placeholder")}
        value={href}
        onChange={(event) => setHref(event.target.value)}
        onKeyDown={(event) => onKeyDown(event, { editor, href, setOpen })}
      />
      {renderActions({ editor, href, setOpen, t })}
    </div>
  );
}

function renderActions({ editor, href, setOpen, t }) {
  return (
    <div className="rte-link-actions">
      <button type="button" className="rte-link-apply" onMouseDown={(e) => e.preventDefault()} onClick={() => applyLink(editor, href, setOpen)}>{t("rte.link.apply")}</button>
      <button type="button" className="rte-link-remove" onMouseDown={(e) => e.preventDefault()} onClick={() => removeLink(editor, setOpen)}>{t("rte.link.remove")}</button>
    </div>
  );
}

function onKeyDown(event, params) {
  if (event.key === "Enter") {
    event.preventDefault();
    applyLink(params.editor, params.href, params.setOpen);
  }
}

function applyLink(editor, href, setOpen) {
  const value = href.trim();
  if (!value) return removeLink(editor, setOpen);
  editor.chain().focus().extendMarkRange("link").setLink({ href: value }).run();
  setOpen(false);
}

function removeLink(editor, setOpen) {
  editor.chain().focus().extendMarkRange("link").unsetLink().run();
  setOpen(false);
}

export default RichTextLinkControl;
