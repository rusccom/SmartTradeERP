import { useEditorState } from "@tiptap/react";

const OPTIONS = [
  { value: "paragraph", labelKey: "rte.format.paragraph" },
  { value: "h1", labelKey: "rte.format.h1" },
  { value: "h2", labelKey: "rte.format.h2" },
  { value: "h3", labelKey: "rte.format.h3" },
];

function RichTextFormatSelect({ editor, t }) {
  const value = useEditorState({ editor, selector: ({ editor: e }) => readBlockType(e) });
  return (
    <select
      className="rte-select"
      aria-label={t("rte.format.label")}
      value={value}
      onMouseDown={(event) => event.stopPropagation()}
      onChange={(event) => applyBlock(editor, event.target.value)}
    >
      {OPTIONS.map((option) => (
        <option key={option.value} value={option.value}>{t(option.labelKey)}</option>
      ))}
    </select>
  );
}

function readBlockType(editor) {
  if (!editor) return "paragraph";
  for (const level of [1, 2, 3]) {
    if (editor.isActive("heading", { level })) return `h${level}`;
  }
  return "paragraph";
}

function applyBlock(editor, value) {
  const chain = editor.chain().focus();
  if (value === "paragraph") chain.setParagraph().run();
  else chain.setHeading({ level: Number(value.slice(1)) }).run();
}

export default RichTextFormatSelect;
