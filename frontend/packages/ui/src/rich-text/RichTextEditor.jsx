import { useState } from "react";
import { EditorContent } from "@tiptap/react";
import { BubbleMenu } from "@tiptap/react/menus";

import { useRichTextEditor } from "./useRichTextEditor";
import RichTextToolbar from "./RichTextToolbar";
import RichTextMarkControls from "./RichTextMarkControls";
import RichTextTableControls from "./RichTextTableControls";
import RichTextHtmlView from "./RichTextHtmlView";
import "./rich-text.css";

function RichTextEditor({ initialContent, documentKey, onHtmlChange, onRequestImage, imageDisabled, t }) {
  const [htmlOpen, setHtmlOpen] = useState(false);
  const editor = useRichTextEditor({ initialContent, documentKey, onHtmlChange });
  if (!editor) return null;
  return (
    <div className="rte">
      <RichTextToolbar editor={editor} htmlOpen={htmlOpen} imageDisabled={imageDisabled} onRequestImage={onRequestImage} onToggleHtml={() => setHtmlOpen((value) => !value)} t={t} />
      <RichTextTableControls editor={editor} t={t} />
      {renderBody({ editor, htmlOpen, setHtmlOpen, t })}
    </div>
  );
}

function renderBody({ editor, htmlOpen, setHtmlOpen, t }) {
  if (htmlOpen) return <RichTextHtmlView editor={editor} t={t} onClose={() => setHtmlOpen(false)} />;
  return (
    <>
      <BubbleMenu editor={editor} options={{ placement: "top", offset: 6 }}>
        <div className="rte-bubble"><RichTextMarkControls editor={editor} t={t} /></div>
      </BubbleMenu>
      <EditorContent className="rte-content" editor={editor} />
    </>
  );
}

export default RichTextEditor;
