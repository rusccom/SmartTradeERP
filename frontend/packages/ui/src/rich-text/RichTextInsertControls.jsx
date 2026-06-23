import { ImagePlus } from "lucide-react";

import RichTextButton from "./RichTextButton";
import { insertImageUrl } from "./richTextCommands";

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
    </div>
  );
}

async function requestImage(editor, onRequestImage) {
  if (!onRequestImage) return;
  const url = await onRequestImage();
  insertImageUrl(editor, url);
}

export default RichTextInsertControls;
