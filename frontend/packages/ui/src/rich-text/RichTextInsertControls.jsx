import { ImagePlus, Table as TableIcon } from "lucide-react";

import RichTextButton from "./RichTextButton";
import RichTextVideoControl from "./RichTextVideoControl";

function RichTextInsertControls({ editor, imageDisabled, onRequestImage, t }) {
  return (
    <div className="rte-group">
      <RichTextButton
        label={imageDisabled ? t("products.form.mediaAfterSave") : t("rte.insert.image")}
        disabled={imageDisabled}
        onClick={() => requestImage(editor, onRequestImage)}
      >
        <ImagePlus size={16} strokeWidth={1.9} aria-hidden="true" />
      </RichTextButton>
      <RichTextVideoControl editor={editor} t={t} />
      <RichTextButton label={t("rte.insert.table")} onClick={() => insertTable(editor)}>
        <TableIcon size={16} strokeWidth={1.9} aria-hidden="true" />
      </RichTextButton>
    </div>
  );
}

async function requestImage(editor, onRequestImage) {
  if (!onRequestImage) return;
  const url = await onRequestImage();
  if (url) editor.chain().focus().setImage({ src: url }).run();
}

function insertTable(editor) {
  editor.chain().focus().insertTable({ rows: 3, cols: 3, withHeaderRow: true }).run();
}

export default RichTextInsertControls;
