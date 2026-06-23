import { useState } from "react";
import { Youtube } from "lucide-react";

import RichTextButton from "./RichTextButton";

function RichTextVideoControl({ editor, t }) {
  const [open, setOpen] = useState(false);
  const [src, setSrc] = useState("");
  const close = () => { setOpen(false); setSrc(""); };
  return (
    <div className="rte-video">
      <RichTextButton label={t("rte.insert.video")} active={open} onClick={() => setOpen((value) => !value)}>
        <Youtube size={16} strokeWidth={1.9} aria-hidden="true" />
      </RichTextButton>
      {open && renderPopover({ editor, src, setSrc, close, t })}
    </div>
  );
}

function renderPopover({ editor, src, setSrc, close, t }) {
  return (
    <div className="rte-link-popover" role="dialog">
      <input
        className="rte-link-input"
        type="url"
        autoFocus
        placeholder={t("rte.insert.videoPlaceholder")}
        value={src}
        onChange={(event) => setSrc(event.target.value)}
        onKeyDown={(event) => onKeyDown(event, { editor, src, close })}
      />
      <div className="rte-link-actions">
        <button type="button" className="rte-link-apply" onMouseDown={(e) => e.preventDefault()} onClick={() => applyVideo(editor, src, close)}>{t("rte.link.apply")}</button>
      </div>
    </div>
  );
}

function onKeyDown(event, params) {
  if (event.key !== "Enter") return;
  event.preventDefault();
  applyVideo(params.editor, params.src, params.close);
}

function applyVideo(editor, src, close) {
  const value = src.trim();
  if (value) editor.commands.setYoutubeVideo({ src: value });
  close();
}

export default RichTextVideoControl;
