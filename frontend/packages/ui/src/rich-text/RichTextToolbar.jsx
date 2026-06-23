import { useEffect, useReducer } from "react";
import { Code2 } from "lucide-react";

import RichTextButton from "./RichTextButton";
import RichTextFormatSelect from "./RichTextFormatSelect";
import RichTextMarkControls from "./RichTextMarkControls";
import RichTextColorControl from "./RichTextColorControl";
import RichTextAlignControls from "./RichTextAlignControls";
import RichTextListControls from "./RichTextListControls";
import RichTextLinkControl from "./RichTextLinkControl";
import RichTextInsertControls from "./RichTextInsertControls";

function RichTextToolbar({ editor, htmlOpen, imageDisabled, onRequestImage, onToggleHtml, t }) {
  useSelectionRefresh(htmlOpen);
  return (
    <div className="rte-toolbar">
      {!htmlOpen && renderControls({ editor, imageDisabled, onRequestImage, t })}
      <RichTextButton label={t("rte.html.label")} active={htmlOpen} onClick={onToggleHtml}>
        <Code2 size={16} strokeWidth={1.9} aria-hidden="true" />
      </RichTextButton>
    </div>
  );
}

// Re-render the toolbar on caret/selection moves so queryCommandState-based
// active highlights stay in sync. Skipped while the HTML source view is open.
function useSelectionRefresh(htmlOpen) {
  const [, refresh] = useReducer((value) => value + 1, 0);
  useEffect(() => {
    if (htmlOpen) return undefined;
    document.addEventListener("selectionchange", refresh);
    return () => document.removeEventListener("selectionchange", refresh);
  }, [htmlOpen]);
}

function renderControls({ editor, imageDisabled, onRequestImage, t }) {
  return (
    <>
      <RichTextFormatSelect editor={editor} t={t} />
      <RichTextMarkControls editor={editor} t={t} />
      <RichTextColorControl editor={editor} t={t} />
      <RichTextAlignControls editor={editor} t={t} />
      <RichTextListControls editor={editor} t={t} />
      <RichTextLinkControl editor={editor} t={t} />
      <RichTextInsertControls editor={editor} imageDisabled={imageDisabled} onRequestImage={onRequestImage} t={t} />
    </>
  );
}

export default RichTextToolbar;
