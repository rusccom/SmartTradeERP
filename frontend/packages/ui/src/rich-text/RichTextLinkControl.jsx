import { useState } from "react";
import { Link2 } from "lucide-react";

import RichTextButton from "./RichTextButton";
import { createLink } from "./richTextCommands";

function RichTextLinkControl({ editor, t }) {
  const [open, setOpen] = useState(false);
  const [href, setHref] = useState("");
  const apply = () => { createLink(editor, href.trim()); setOpen(false); setHref(""); };
  return (
    <div className="rte-link">
      <RichTextButton label={t("rte.link.label")} active={open} onClick={() => setOpen((value) => !value)}>
        <Link2 size={16} strokeWidth={1.9} aria-hidden="true" />
      </RichTextButton>
      {open && renderPopover({ href, setHref, apply, t })}
    </div>
  );
}

function renderPopover({ href, setHref, apply, t }) {
  return (
    <div className="rte-link-popover" role="dialog">
      <input
        className="rte-link-input"
        type="url"
        autoFocus
        placeholder={t("rte.link.placeholder")}
        value={href}
        onChange={(event) => setHref(event.target.value)}
        onKeyDown={(event) => onKeyDown(event, apply)}
      />
      <div className="rte-link-actions">
        <button type="button" className="rte-link-apply" onMouseDown={(e) => e.preventDefault()} onClick={apply}>{t("rte.link.apply")}</button>
      </div>
    </div>
  );
}

function onKeyDown(event, apply) {
  if (event.key !== "Enter") return;
  event.preventDefault();
  apply();
}

export default RichTextLinkControl;
