import { useState } from "react";

function RichTextHtmlView({ editor, t, onClose }) {
  const [value, setValue] = useState(() => editor.getHTML());
  return (
    <div className="rte-html">
      <textarea
        className="rte-html-area"
        spellCheck={false}
        aria-label={t("rte.html.label")}
        value={value}
        onChange={(event) => setValue(event.target.value)}
      />
      <div className="rte-html-actions">
        <button type="button" className="rte-html-apply" onClick={() => apply(editor, value, onClose)}>{t("rte.html.apply")}</button>
        <button type="button" className="rte-html-cancel" onClick={onClose}>{t("rte.html.cancel")}</button>
      </div>
    </div>
  );
}

function apply(editor, value, onClose) {
  editor.commands.setContent(value, { emitUpdate: true });
  onClose();
}

export default RichTextHtmlView;
