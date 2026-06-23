import { useState } from "react";

import { useEditableHtml } from "./useEditableHtml";
import RichTextToolbar from "./RichTextToolbar";
import RichTextHtmlView from "./RichTextHtmlView";
import "./rich-text.css";

// Raw-HTML contenteditable editor (Shopify body_html style): the editable
// element's innerHTML IS the value, pasted HTML is preserved verbatim, and the
// </> toggle swaps the rendered surface for an editable HTML source view.
function RichTextEditor({ initialContent, documentKey, onHtmlChange, onRequestImage, imageDisabled, t }) {
  const editor = useEditableHtml({ initialContent, documentKey, onHtmlChange });
  const [source, setSource] = useState(null);
  const sourceOpen = source !== null;
  const onToggleHtml = () => setSource(toggleSource(editor, source));
  return (
    <div className="rte">
      <RichTextToolbar editor={editor} htmlOpen={sourceOpen} imageDisabled={imageDisabled} onRequestImage={onRequestImage} onToggleHtml={onToggleHtml} t={t} />
      {renderBody({ editor, source, setSource, onHtmlChange, t })}
    </div>
  );
}

// toggleSource flips between the rendered editor and the HTML source. Opening
// seeds the draft from the live innerHTML; closing applies the edited draft
// back into the contenteditable (via seed) and reports the new value.
function toggleSource(editor, source) {
  if (source === null) {
    const el = editor.ref.current;
    return el ? el.innerHTML : "";
  }
  editor.seed(source);
  return null;
}

function renderBody({ editor, source, setSource, onHtmlChange, t }) {
  if (source !== null) return <RichTextHtmlView value={source} onChange={(html) => applySource(html, setSource, onHtmlChange)} t={t} />;
  return (
    <div
      className="rte-content"
      contentEditable
      suppressContentEditableWarning
      ref={editor.attach}
      onInput={editor.onInput}
    />
  );
}

// applySource keeps the form value live while the HTML source view is open, so
// editing in </> and clicking Save (without toggling back) never loses edits.
function applySource(html, setSource, onHtmlChange) {
  setSource(html);
  onHtmlChange(html);
}

export default RichTextEditor;
