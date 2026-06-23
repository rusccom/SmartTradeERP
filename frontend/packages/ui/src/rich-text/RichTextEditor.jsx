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
  const editor = useRichTextEditor({ initialContent, documentKey, onHtmlChange });
  const [source, setSource] = useState(null);
  if (!editor) return null;
  const sourceOpen = source !== null;
  return (
    <div className="rte">
      <RichTextToolbar editor={editor} htmlOpen={sourceOpen} imageDisabled={imageDisabled} onRequestImage={onRequestImage} onToggleHtml={() => setSource(toggleSource(editor, source))} t={t} />
      {!sourceOpen && <RichTextTableControls editor={editor} t={t} />}
      {renderBody({ editor, source, setSource, t })}
    </div>
  );
}

// toggleSource flips between the rendered editor and the HTML source. Opening
// seeds the draft from the editor; closing applies the edited draft back.
function toggleSource(editor, source) {
  if (source === null) return editor.getHTML();
  editor.commands.setContent(source, { emitUpdate: true });
  return null;
}

function renderBody({ editor, source, setSource, t }) {
  if (source !== null) return <RichTextHtmlView value={source} onChange={setSource} t={t} />;
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
